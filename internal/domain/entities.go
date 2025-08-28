package domain

type (
	Chain string
	Dex   string
)

const (
	ChainBase Chain = "Base"

	DexUniswapV3 Dex = "Uniswap V3"
	DexAerodrome Dex = "Aerodrome"
)

type Token struct {
	Name string
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

func (p LiquidityPoolPosition) GetTokensPercentage() (token0, token1 float64) {
	token1 = float64(p.CurrentTick-p.TickLower) / float64(p.TickUpper-p.TickLower) * 100
	token0 = 100 - token1
	return token0, token1
}

func (p LiquidityPoolPosition) IsInRange() bool {
	return p.TickLower <= p.CurrentTick && p.CurrentTick <= p.TickUpper
}
