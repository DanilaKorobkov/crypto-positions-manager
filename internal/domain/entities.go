package domain

type Token struct {
	Name string
}

type Notify struct {
	Message string
}

type UniswapV3Position struct {
	ID          string
	Token0      Token
	Token1      Token
	TickLower   int
	TickUpper   int
	CurrentTick int
}

// IsActive check if the position in range.
func (position UniswapV3Position) IsActive() bool {
	return position.TickLower <= position.CurrentTick && position.CurrentTick <= position.TickUpper
}
