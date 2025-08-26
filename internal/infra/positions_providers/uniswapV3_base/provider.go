package uniswapV3_base

import (
	"context"
	"fmt"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
	"github.com/hasura/go-graphql-client"
	"github.com/samber/lo"
	"strconv"
)

type Provider struct {
	client *graphql.Client
}

func NewProvider(client *graphql.Client) *Provider {
	return &Provider{
		client: client,
	}
}

func (w *Provider) GetOpenPositions(
	ctx context.Context,
	walletAddress string,
) ([]domain.UniswapV3Position, error) {
	var unclosedPosition unclosedPositionsQuery
	variables := map[string]any{
		"walletAddress": graphql.String(walletAddress),
	}
	err := w.client.Query(ctx, &unclosedPosition, variables)
	if err != nil {
		return nil, fmt.Errorf("uniswap v3 query: %w", err)
	}
	return convertToDomain(unclosedPosition.Positions), nil
}

func convertToDomain(unclosedPositions []position) []domain.UniswapV3Position {
	if len(unclosedPositions) == 0 {
		return nil
	}
	return lo.Map(unclosedPositions, func(pos position, _ int) domain.UniswapV3Position {
		return domain.UniswapV3Position{
			ID:          pos.ID,
			TickLower:   pos.TickLower,
			TickUpper:   pos.TickUpper,
			CurrentTick: mustConvertToInt(pos.Pool.Tick),
			Token0: domain.Token{
				Name: pos.Pool.Token0.Symbol,
			},
			Token1: domain.Token{
				Name: pos.Pool.Token1.Symbol,
			},
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
	Positions []position `graphql:"positions(where: {owner: $walletAddress, liquidity_gt: 0})"`
}

type position struct {
	ID        string
	TickLower int
	TickUpper int
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
