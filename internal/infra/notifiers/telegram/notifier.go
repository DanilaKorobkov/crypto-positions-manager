package telegram

import (
	"context"
	"fmt"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notifier struct {
	telegram *tgbotapi.BotAPI
	userID   int64
}

type NotifierConfig struct {
	Telegram *tgbotapi.BotAPI
	UserID   int64
}

func NewNotifier(config NotifierConfig) *Notifier {
	return &Notifier{
		telegram: config.Telegram,
		userID:   config.UserID,
	}
}

func (n *Notifier) Notify(_ context.Context, notify domain.Notify) error {
	msg := tgbotapi.NewMessage(n.userID, notify.Message)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.telegram.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram.Send: %w", err)
	}
	return nil
}
