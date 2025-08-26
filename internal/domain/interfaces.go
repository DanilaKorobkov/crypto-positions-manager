package domain

import "context"

type UniswapProvider interface {
	GetOpenPositions(ctx context.Context, walletAddress string) ([]UniswapV3Position, error)
}

type Notifier interface {
	Notify(ctx context.Context, notify Notify) error
}
