package domain

import (
	"math"
)

type (
	Chain string
	Dex   string
)

const (
	ChainBase Chain = "Base"

	DexUniswapV3 Dex = "Uniswap V3"
	DexAerodrome Dex = "Aerodrome"

	tickBase float64 = 1.0001
)

type Token struct {
	Name     string
	Decimals int
}

type Subject struct {
	TelegramUserID int64
	Wallet         string
}

type LiquidityPoolPosition struct {
	Chain        Chain
	Dex          Dex
	PositionLink string
	Token0       Token
	Token1       Token
	CurrentTick  int
	TickLower    int
	TickUpper    int
}

func (p LiquidityPoolPosition) GetCurrentPrice() float64 {
	return p.tickToPrice(p.CurrentTick)
}

func (p LiquidityPoolPosition) GetLowerPrice() float64 {
	return p.tickToPrice(p.TickLower)
}

func (p LiquidityPoolPosition) GetTokensPercentage() (token0, token1 float64) {
	if p.CurrentTick <= p.TickLower {
		return 100, 0
	}

	if p.CurrentTick >= p.TickUpper {
		return 0, 100
	}

	token1 = float64(p.CurrentTick-p.TickLower) / float64(p.TickUpper-p.TickLower) * 100
	token0 = 100 - token1

	return token0, token1
}

func (p LiquidityPoolPosition) GetUpperPrice() float64 {
	return p.tickToPrice(p.TickUpper)
}

func (p LiquidityPoolPosition) IsInRange() bool {
	return p.TickLower <= p.CurrentTick && p.CurrentTick <= p.TickUpper
}

func (p LiquidityPoolPosition) tickToPrice(tick int) float64 {
	decimal0 := p.Token0.Decimals
	decimal1 := p.Token1.Decimals
	return math.Pow(tickBase, float64(tick)) * math.Pow(10, math.Abs(float64(decimal0-decimal1)))
}
