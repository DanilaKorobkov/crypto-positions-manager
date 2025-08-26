package domain

import "context"

type UniswapProvider interface {
	// GetPositionsWithLiquidity returns uniswap positions that hav liquidity.
	GetPositionsWithLiquidity(ctx context.Context, walletAddress string) ([]UniswapV3Position, error)
}

type Notifier interface {
	// Notify sends notify to user.
	Notify(ctx context.Context, notify Notify) error
}
