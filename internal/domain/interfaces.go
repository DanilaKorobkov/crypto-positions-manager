package domain

import "context"

type UniswapProvider interface {
	// GetPositionsWithLiquidity returns uniswap positions that hav liquidity.
	GetPositionsWithLiquidity(ctx context.Context, walletAddress string) ([]UniswapV3Position, error)
}

type Notifier interface {
	// Notify sends message to user.
	Notify(ctx context.Context, user User, message string) error
}
