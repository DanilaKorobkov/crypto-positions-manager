package watcher

import (
	"context"
	"fmt"
	"log/slog"
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
		"❌[position](https://app.uniswap.org/positions/v3/base/%s) %s:%s is out of range",
		inactivePosition.ID,
		inactivePosition.Token0.Name,
		inactivePosition.Token1.Name,
	)
}

func buildActivePositionMessage(activePosition domain.UniswapV3Position) string {
	return fmt.Sprintf(
		"✅[position](https://app.uniswap.org/positions/v3/base/%s) %s:%s in range: %d%%",
		activePosition.ID,
		activePosition.Token0.Name,
		activePosition.Token1.Name,
		int(calculateTickPercentagePosition(activePosition)),
	)
}

func calculateTickPercentagePosition(position domain.UniswapV3Position) float64 {
	return float64(position.CurrentTick-position.TickLower) / float64(position.TickUpper-position.TickLower) * 100
}
