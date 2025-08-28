package watcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
	"github.com/DanilaKorobkov/defi-monitoring/pkg/tickers"
)

type ServiceConfig struct {
	LiquidityPoolPositions domain.LiquidityPoolPositionsProvider
	Notifier               domain.Notifier
	CheckInterval          time.Duration
	Logger                 *slog.Logger
}

type Service struct {
	liquidityPoolPositions domain.LiquidityPoolPositionsProvider
	notifier               domain.Notifier
	checkInterval          time.Duration
	logger                 *slog.Logger
}

func NewService(config ServiceConfig) *Service {
	return &Service{
		liquidityPoolPositions: config.LiquidityPoolPositions,
		notifier:               config.Notifier,
		checkInterval:          config.CheckInterval,
		logger:                 config.Logger,
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

	positions, err := service.liquidityPoolPositions.GetPositionsWithLiquidity(ctx, subject.Wallet)
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
	position domain.LiquidityPoolPosition,
) error {
	notify := service.notifier.NotifyLiquidityPoolPositionInRange
	if !position.IsInRange() {
		notify = service.notifier.NotifyLiquidityPoolPositionOutOfRange
	}

	err := notify(ctx, subject, position)
	if err != nil {
		return fmt.Errorf("notifier.Send: %w", err)
	}

	return nil
}

func (service *Service) processPositions(
	ctx context.Context,
	subject domain.Subject,
	positions []domain.LiquidityPoolPosition,
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
