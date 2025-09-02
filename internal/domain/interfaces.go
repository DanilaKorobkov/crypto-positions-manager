package domain

import "context"

type LiquidityPoolPositionsProvider interface {
	// GetName returns driver information.
	GetName() string
	// GetPositionsWithLiquidity get wallet positions with liquidity.
	GetPositionsWithLiquidity(ctx context.Context, wallet string) ([]LiquidityPoolPosition, error)
}

type Notifier interface {
	// NotifyLiquidityPoolPositions notify subject the positions status and info about.
	NotifyLiquidityPoolPositions(ctx context.Context, subject Subject, positions ...LiquidityPoolPosition) error
}

type SubjectsRepository interface {
	// Add subject and override if already exists.
	Add(ctx context.Context, subject Subject) error
	// GetAll returns all stored subjects.
	GetAll(ctx context.Context) ([]Subject, error)
}
