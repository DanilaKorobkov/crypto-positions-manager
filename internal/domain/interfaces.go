package domain

import "context"

type LiquidityPoolPositionsProvider interface {
	// GetName returns driver information.
	GetName() string
	// GetPositionsWithLiquidity get wallet positions with liquidity.
	GetPositionsWithLiquidity(ctx context.Context, wallet string) ([]LiquidityPoolPosition, error)
}

type Notifier interface {
	// NotifyLiquidityPoolPositionInRange notify subject the position is active.
	NotifyLiquidityPoolPositionInRange(ctx context.Context, subject Subject, position LiquidityPoolPosition) error
	// NotifyLiquidityPoolPositionOutOfRange notify subject the position out of range.
	NotifyLiquidityPoolPositionOutOfRange(ctx context.Context, subject Subject, position LiquidityPoolPosition) error
}
