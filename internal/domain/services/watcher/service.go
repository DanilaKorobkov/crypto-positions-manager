package watcher

import (
	"context"
	"log/slog"
	"time"

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
	logger := service.logger.With(slog.String("wallet", subject.Wallets[0]))

	positions, err := service.liquidityPoolPositions.GetPositionsWithLiquidity(
		ctx,
		subject.Wallets[0],
	)
	if err != nil {
		logger.Error("GetPositionsWithLiquidity", slog.String("err", err.Error()))
		return
	}

	if len(positions) == 0 {
		logger.Info("no positions found")
		return
	}

	err = service.notifier.NotifyLiquidityPoolPositions(ctx, subject, positions...)
	if err != nil {
		logger.Error("NotifyLiquidityPoolPositions", slog.String("err", err.Error()))
		return
	}
}
