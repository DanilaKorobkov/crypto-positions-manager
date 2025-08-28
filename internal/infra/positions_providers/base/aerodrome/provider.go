package aerodrome

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hasura/go-graphql-client"
	"github.com/samber/lo"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

type ProviderTheGraph struct {
	client *graphql.Client
}

func NewProviderTheGraph(client *graphql.Client) *ProviderTheGraph {
	return &ProviderTheGraph{
		client: client,
	}
}

func (*ProviderTheGraph) GetName() string {
	return "Base Aerodrome"
}

func (provider *ProviderTheGraph) GetPositionsWithLiquidity(
	ctx context.Context,
	wallet string,
) ([]domain.LiquidityPoolPosition, error) {
	var unclosedPosition unclosedPositionsQuery

	variables := map[string]any{
		"wallet": wallet,
	}

	err := provider.client.Query(ctx, &unclosedPosition, variables)
	if err != nil {
		return nil, fmt.Errorf("graphql.Query: %w", err)
	}

	return convertToDomain(unclosedPosition.Positions), nil
}

func convertToDomain(unclosedPositions []position) []domain.LiquidityPoolPosition {
	if len(unclosedPositions) == 0 {
		return nil
	}

	return lo.Map(unclosedPositions, func(pos position, _ int) domain.LiquidityPoolPosition {
		return domain.LiquidityPoolPosition{
			Chain:        domain.ChainBase,
			Dex:          domain.DexAerodrome,
			PositionLink: "https://aerodrome.finance/dash",
			Token0: domain.Token{
				Name: pos.Pool.Token0.Symbol,
			},
			Token1: domain.Token{
				Name: pos.Pool.Token1.Symbol,
			},
			CurrentTick: mustConvertToInt(pos.Pool.Tick),
			TickLower:   mustConvertToInt(pos.TickLower.TickIdx),
			TickUpper:   mustConvertToInt(pos.TickUpper.TickIdx),
		}
	})
}

func mustConvertToInt(value string) int {
	integer, err := strconv.Atoi(value)
	if err != nil {
		message := "mustConvertToInt: " + value
		panic(message)
	}

	return integer
}

type unclosedPositionsQuery struct {
	Positions []position `graphql:"positions(where: {owner: $wallet, liquidity_gt: 0})"`
}

type position struct {
	ID        string
	TickLower tick
	TickUpper tick
	Pool      pool
}

type pool struct {
	Tick   string
	Token0 token
	Token1 token
}

type token struct {
	Symbol string
}

type tick struct {
	TickIdx string
}
