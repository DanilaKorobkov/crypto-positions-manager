package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
)

type Notifier struct {
	telegramBot *tgbotapi.BotAPI
}

func NewNotifier(telegramBot *tgbotapi.BotAPI) *Notifier {
	return &Notifier{
		telegramBot: telegramBot,
	}
}

func (n *Notifier) Notify(_ context.Context, subject domain.Subject, text string) error {
	message := tgbotapi.NewMessage(subject.TelegramUserID, text)
	message.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.telegramBot.Send(message)
	if err != nil {
		return fmt.Errorf("telegram.Send: %w", err)
	}

	return nil
}
