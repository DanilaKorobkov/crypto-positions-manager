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

type SubjectRepository interface {
	Add(ctx context.Context, subject Subject) error
	Get(ctx context.Context, id string) (Subject, error)
}
