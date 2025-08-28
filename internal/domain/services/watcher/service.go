package watcher

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
	"github.com/DanilaKorobkov/crypto-positions-manager/pkg/tickers"
)

type ServiceConfig struct {
	Provider      domain.UniswapProvider
	Notifier      domain.Notifier
	CheckInterval time.Duration
	Logger        *slog.Logger
}

type Service struct {
	provider      domain.UniswapProvider
	notifier      domain.Notifier
	checkInterval time.Duration
	logger        *slog.Logger
}

func NewService(config ServiceConfig) *Service {
	return &Service{
		provider:      config.Provider,
		notifier:      config.Notifier,
		checkInterval: config.CheckInterval,
		logger:        config.Logger,
	}
}

func (service *Service) StartWatching(ctx context.Context, subject domain.Subject) {
	ticker := tickers.NewTickerChanWithInitial(service.checkInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			service.checkPositions(context.WithoutCancel(ctx), subject)
		}
	}
}

func (service *Service) checkPositions(ctx context.Context, subject domain.Subject) {
	logger := service.logger.With(slog.String("wallet", subject.Wallet))

	positions, err := service.provider.GetPositionsWithLiquidity(ctx, subject.Wallet)
	if err != nil {
		logger.Error("GetPositionsWithLiquidity", slog.String("err", err.Error()))
		return
	}

	if len(positions) == 0 {
		logger.Info("no positions found")
		return
	}

	err = service.processPositions(ctx, subject, positions)
	if err != nil {
		logger.Error("processPositions", slog.String("err", err.Error()))
		return
	}
}

func (service *Service) processPosition(
	ctx context.Context,
	subject domain.Subject,
	position domain.UniswapV3Position,
) error {
	message := buildOutOfRangeMessage(position)
	if position.IsActive() {
		message = buildActivePositionMessage(position)
	}

	err := service.notifier.Notify(ctx, subject, message)
	if err != nil {
		return fmt.Errorf("notifier.Send: %w", err)
	}

	return nil
}

func (service *Service) processPositions(
	ctx context.Context,
	subject domain.Subject,
	positions []domain.UniswapV3Position,
) error {
	if len(positions) == 0 {
		return nil
	}

	executor := pool.New().WithContext(ctx).WithFirstError().WithCancelOnError()
	for _, position := range positions {
		executor.Go(func(ctx context.Context) error {
			return service.processPosition(ctx, subject, position)
		})
	}

	err := executor.Wait()
	if err != nil {
		return fmt.Errorf("executor.Wait: %w", err)
	}

	return nil
}

func buildOutOfRangeMessage(inactivePosition domain.UniswapV3Position) string {
	return fmt.Sprintf(
		"❌ [Position](https://app.uniswap.org/positions/v3/base/%s) %s : %s out of range",
		inactivePosition.ID,
		inactivePosition.Token0.Name,
		inactivePosition.Token1.Name,
	)
}

func buildActivePositionMessage(activePosition domain.UniswapV3Position) string {
	token0, token1 := calculateTokensPercentage(activePosition)

	return fmt.Sprintf(
		`✅ [Position](https://app.uniswap.org/positions/v3/base/%s) in range\. %s\(%s%%\) : %s\(%s%%\)\.`,
		activePosition.ID,
		activePosition.Token0.Name,
		formatAndEscape(token0),
		activePosition.Token1.Name,
		formatAndEscape(token1),
	)
}

func calculateTokensPercentage(position domain.UniswapV3Position) (token0, token1 float64) {
	token1 = float64(position.CurrentTick-position.TickLower) / float64(position.TickUpper-position.TickLower) * 100
	token0 = 100 - token1
	return token0, token1
}

func formatAndEscape(value float64) string {
	cut := fmt.Sprintf("%.2f", value)
	return strings.Replace(cut, ".", ",", 1)
}
