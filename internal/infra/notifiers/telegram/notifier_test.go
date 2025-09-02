package telegram_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
	"github.com/DanilaKorobkov/defi-monitoring/internal/infra/notifiers/telegram"
	mocks "github.com/DanilaKorobkov/defi-monitoring/mocks/internal_/infra/notifiers/telegram"
	"github.com/DanilaKorobkov/defi-monitoring/test/generators"
)

type notifierSuite struct {
	suite.Suite
}

func TestNotifier(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(notifierSuite))
}

func (s *notifierSuite) TestFake() {
	type TestCase struct {
		name                string
		makePosition        func() domain.LiquidityPoolPosition
		expectedMessageText string
	}

	testCases := []TestCase{
		{
			name:                "In range",
			makePosition:        makePosition,
			expectedMessageText: inRangePositionText,
		},
		{
			name: "Out of range lower",
			makePosition: func() domain.LiquidityPoolPosition {
				position := makePosition()
				position.CurrentTick = position.TickLower - 1
				return position
			},
			expectedMessageText: outOfRangeLowerText,
		},
		{
			name: "Out of range upper",
			makePosition: func() domain.LiquidityPoolPosition {
				position := makePosition()
				position.CurrentTick = position.TickUpper + 1
				return position
			},
			expectedMessageText: outOfRangeUpperText,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			ctx := context.Background()

			subject := generators.NewSubjectGenerator().Slim().Result()
			position := testCase.makePosition()

			expectedMessage := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: subject.TelegramUserID,
				},
				DisableWebPagePreview: true,
				ParseMode:             tgbotapi.ModeHTML,
				Text:                  strings.TrimSpace(testCase.expectedMessageText),
			}
			tgBot := mocks.NewTgBotApi(s.T())
			tgBot.EXPECT().
				Send(expectedMessage).
				Return(tgbotapi.Message{}, nil).
				Once()
			notifier := telegram.NewNotifier(tgBot)
			err := notifier.NotifyLiquidityPoolPositions(ctx, subject, position)
			s.Require().NoError(err)
		})
	}
}

func makePosition() domain.LiquidityPoolPosition {
	return domain.LiquidityPoolPosition{
		Chain:        domain.ChainBase,
		Dex:          domain.DexUniswapV3,
		PositionLink: "https://google.com",
		Token0: domain.Token{
			Name:     "WETH",
			Decimals: 18,
		},
		Token1: domain.Token{
			Name:     "USDC",
			Decimals: 6,
		},
		TickLower:   -192660,
		CurrentTick: -191000,
		TickUpper:   -190940,
	}
}

const inRangePositionText = `
<b>Statuses:</b> ✅

<b>Status: ✅</b>
<b>Chain:</b> Base
<b>Dex:</b> Uniswap V3
<b>Position:</b> <a href="https://google.com">link</a>
<b>Proportion:</b> WETH (3,49%) : USDC (96,51%)
<b>Range low price:</b> 1 WETH = 4298,34 USDC
<b>Range up price:</b> 1 WETH = 5105,00 USDC
<b>Current price:</b> 1 WETH = 5074,46 USDC
`

const outOfRangeLowerText = `
<b>Statuses:</b> ❌

<b>Status: ❌</b>
<b>Chain:</b> Base
<b>Dex:</b> Uniswap V3
<b>Position:</b> <a href="https://google.com">link</a>
<b>Proportion:</b> WETH (100,00%) : USDC (0,00%)
<b>Range low price:</b> 1 WETH = 4298,34 USDC
<b>Range up price:</b> 1 WETH = 5105,00 USDC
<b>Current price:</b> 1 WETH = 4297,91 USDC
`

const outOfRangeUpperText = `
<b>Statuses:</b> ❌

<b>Status: ❌</b>
<b>Chain:</b> Base
<b>Dex:</b> Uniswap V3
<b>Position:</b> <a href="https://google.com">link</a>
<b>Proportion:</b> WETH (0,00%) : USDC (100,00%)
<b>Range low price:</b> 1 WETH = 4298,34 USDC
<b>Range up price:</b> 1 WETH = 5105,00 USDC
<b>Current price:</b> 1 WETH = 5105,51 USDC
`
