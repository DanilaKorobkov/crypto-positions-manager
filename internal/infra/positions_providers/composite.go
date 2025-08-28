package positions_providers

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
)

type Composite struct {
	impls []domain.LiquidityPoolPositionsProvider
}

func NewComposite(impls ...domain.LiquidityPoolPositionsProvider) *Composite {
	return &Composite{
		impls: impls,
	}
}

func (*Composite) GetName() string {
	return "Composite"
}

func (c *Composite) GetPositionsWithLiquidity(
	ctx context.Context,
	wallet string,
) ([]domain.LiquidityPoolPosition, error) {
	p := pool.NewWithResults[[]domain.LiquidityPoolPosition]().WithContext(ctx).WithFirstError().WithCancelOnError()
	for _, impl := range c.impls {
		p.Go(func(ctx context.Context) ([]domain.LiquidityPoolPosition, error) {
			positions, err := impl.GetPositionsWithLiquidity(ctx, wallet)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", impl.GetName(), err)
			}
			return positions, nil
		})
	}

	pos, err := p.Wait()
	if err != nil {
		return nil, fmt.Errorf("p.Wait: %w", err)
	}

	return lo.Flatten(pos), nil
}
