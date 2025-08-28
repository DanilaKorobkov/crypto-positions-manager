package domain

import "context"

type UniswapProvider interface {
	// GetPositionsWithLiquidity returns uniswap positions that have liquidity.
	GetPositionsWithLiquidity(ctx context.Context, wallet string) ([]UniswapV3Position, error)
}

type Notifier interface {
	// Notify sends message to subject.
	Notify(ctx context.Context, user Subject, message string) error
}
