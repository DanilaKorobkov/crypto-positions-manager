package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
)

type ServiceConfig struct {
	Provider      domain.UniswapProvider
	Notifier      domain.Notifier
	CheckInterval time.Duration
	Logger        zerolog.Logger
}

type Service struct {
	provider      domain.UniswapProvider
	notifier      domain.Notifier
	checkInterval time.Duration
	logger        zerolog.Logger
}

func NewService(config ServiceConfig) *Service { //nolint:gocritic // Pointer is too much.
	return &Service{
		provider:      config.Provider,
		notifier:      config.Notifier,
		checkInterval: config.CheckInterval,
		logger:        config.Logger,
	}
}

func (service *Service) StartWatching(ctx context.Context, walletAddress string) {
	ticker := time.NewTicker(service.checkInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			service.checkPositions(context.WithoutCancel(ctx), walletAddress)
		}
	}
}

func (service *Service) checkPositions(ctx context.Context, walletAddress string) {
	logger := service.logger.With().Str("wallet", walletAddress).Logger()

	positions, err := service.provider.GetPositionsWithLiquidity(ctx, walletAddress)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	err = service.processPositions(ctx, positions)
	if err != nil {
		logger.Error().Err(err).Send()
	}
}

func (service *Service) processPosition(
	ctx context.Context,
	position domain.UniswapV3Position,
) error {
	notify := domain.Notify{
		Message: buildActivePositionMessage(position),
	}
	if !position.IsActive() {
		notify.Message = buildOutOfRangeMessage(position)
	}

	err := service.notifier.Notify(ctx, notify)
	if err != nil {
		return fmt.Errorf("notifier.Send: %w", err)
	}

	return nil
}

func (service *Service) processPositions(
	ctx context.Context,
	positions []domain.UniswapV3Position,
) error {
	if len(positions) == 0 {
		return nil
	}

	executor := pool.New().WithContext(ctx).WithFirstError().WithCancelOnError()
	for _, position := range positions {
		executor.Go(func(ctx context.Context) error {
			return service.processPosition(ctx, position)
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
