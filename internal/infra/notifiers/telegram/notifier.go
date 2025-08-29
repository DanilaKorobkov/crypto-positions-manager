package telegram

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/samber/lo"

	_ "embed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

//go:embed templates/dex_lp_position.html
var dexLPTemplate string

type Notifier struct {
	telegramBot *tgbotapi.BotAPI
}

func NewNotifier(telegramBot *tgbotapi.BotAPI) *Notifier {
	return &Notifier{
		telegramBot: telegramBot,
	}
}

func (n *Notifier) NotifyLiquidityPoolPositions(
	_ context.Context,
	subject domain.Subject,
	positions ...domain.LiquidityPoolPosition,
) error {
	messageText, err := makeMessageText(positions)
	if err != nil {
		return fmt.Errorf("makeMessageText: %w", err)
	}

	message := tgbotapi.NewMessage(subject.TelegramUserID, messageText)
	message.ParseMode = tgbotapi.ModeHTML
	message.DisableWebPagePreview = true

	_, err = n.telegramBot.Send(message)
	if err != nil {
		return fmt.Errorf("telegram.Send: %w", err)
	}

	return nil
}

func makeMessageText(positions []domain.LiquidityPoolPosition) (string, error) {
	statuses := convertToAnotherSlice(positions, getStatus)
	data := renderInfo{
		Statuses:  strings.Join(statuses, " "),
		Positions: convertToAnotherSlice(positions, makePositionRenderInfo),
	}

	message, err := renderMessage(data)
	if err != nil {
		return "", fmt.Errorf("renderMessage: %w", err)
	}

	return message, nil
}

func getStatus(position domain.LiquidityPoolPosition) string {
	if position.IsInRange() {
		return "✅"
	}
	return "❌"
}

func makePositionRenderInfo(position domain.LiquidityPoolPosition) positionRenderInfo {
	token0, token1 := position.GetTokensPercentage()

	return positionRenderInfo{
		Status:        "✅",
		Chain:         string(position.Chain),
		Dex:           string(position.Dex),
		PositionLink:  position.PositionLink,
		Token0:        position.Token0.Name,
		Token0Percent: formatAndEscape(token0),
		Token1:        position.Token1.Name,
		Token1Percent: formatAndEscape(token1),
		LowPrice:      formatAndEscape(position.GetLowerPrice()),
		UpPrice:       formatAndEscape(position.GetUpperPrice()),
		CurrentPrice:  formatAndEscape(position.GetCurrentPrice()),
	}
}

func renderMessage(data renderInfo) (string, error) {
	tmpl, err := template.New("telegramMsg").Parse(dexLPTemplate)
	if err != nil {
		return "", fmt.Errorf("template.New: %w", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("template.New: %w", err)
	}

	return buf.String(), nil
}

func formatAndEscape(value float64) string {
	cut := fmt.Sprintf("%.2f", value)
	return strings.Replace(cut, ".", ",", 1)
}

func convertToAnotherSlice[T any, R any](items []T, cast func(T) R) []R {
	return lo.Map(items, func(item T, _ int) R {
		return cast(item)
	})
}

type renderInfo struct {
	Statuses  string
	Positions []positionRenderInfo
}

type positionRenderInfo struct {
	Status        string
	Chain         string
	Dex           string
	PositionLink  string
	Token0        string
	Token0Percent string
	Token1        string
	Token1Percent string
	LowPrice      string
	UpPrice       string
	CurrentPrice  string
}
