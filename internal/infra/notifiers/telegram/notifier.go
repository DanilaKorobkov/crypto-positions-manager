package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

type Notifier struct {
	telegramBot *tgbotapi.BotAPI
}

func NewNotifier(telegramBot *tgbotapi.BotAPI) *Notifier {
	return &Notifier{
		telegramBot: telegramBot,
	}
}

func (n *Notifier) NotifyLiquidityPoolPositionInRange(
	ctx context.Context,
	subject domain.Subject,
	position domain.LiquidityPoolPosition,
) error {
	return n.sendMessage(ctx, subject, buildActivePositionMessage(position))
}

func (n *Notifier) NotifyLiquidityPoolPositionOutOfRange(
	ctx context.Context,
	subject domain.Subject,
	position domain.LiquidityPoolPosition,
) error {
	return n.sendMessage(ctx, subject, buildOutOfRangeMessage(position))
}

func (n *Notifier) sendMessage(_ context.Context, subject domain.Subject, text string) error {
	message := tgbotapi.NewMessage(subject.TelegramUserID, text)
	message.ParseMode = tgbotapi.ModeMarkdownV2
	message.DisableWebPagePreview = true

	_, err := n.telegramBot.Send(message)
	if err != nil {
		return fmt.Errorf("telegram.Send: %w", err)
	}

	return nil
}

func buildOutOfRangeMessage(position domain.LiquidityPoolPosition) string {
	return fmt.Sprintf(
		"❌ %s %s: [position](%s) %s : %s out of range",
		position.Chain,
		position.Dex,
		position.PositionLink,
		position.Token0.Name,
		position.Token1.Name,
	)
}

func buildActivePositionMessage(position domain.LiquidityPoolPosition) string {
	token0, token1 := position.GetTokensPercentage()

	return fmt.Sprintf(
		`✅ %s %s [position](%s) in range\. %s\(%s%%\) : %s\(%s%%\)\.`,
		position.Chain,
		position.Dex,
		position.PositionLink,
		position.Token0.Name,
		formatAndEscape(token0),
		position.Token1.Name,
		formatAndEscape(token1),
	)
}

func formatAndEscape(value float64) string {
	cut := fmt.Sprintf("%.2f", value)
	return strings.Replace(cut, ".", ",", 1)
}
