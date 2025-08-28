package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/rs/zerolog"

	slogtelegram "github.com/platx/slog-telegram"
	slogmulti "github.com/samber/slog-multi"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type LoggerConfig struct {
	TelegramBotToken       string
	TelegramReceiverUserID int64
}

func NewLogger(config LoggerConfig) *slog.Logger {
	compositeHandler := slogmulti.Router().
		Add(makeZeroLogHandler(), matchLogToZeroLog).
		Add(makeTelegramHandler(config), matchLogToTelegram).
		Handler()
	return slog.New(compositeHandler)
}

func makeTelegramHandler(config LoggerConfig) slog.Handler {
	options := slogtelegram.HandlerOptions{
		Formatter: slogtelegram.FormatterOptions{},
		Sender: slogtelegram.SenderOptions{
			Token:  config.TelegramBotToken,
			ChatID: config.TelegramReceiverUserID,
		},
	}

	return slogtelegram.NewHandler(options)
}

func makeZeroLogHandler() slog.Handler {
	zeroLogger := zerolog.New(os.Stdout)

	slogConfig := slogzerolog.Option{
		Level:  slog.LevelDebug,
		Logger: &zeroLogger,
	}

	return slogConfig.NewZerologHandler()
}

func matchLogToZeroLog(context.Context, slog.Record) bool {
	return true
}

//nolint:gocritic // It's the library contract
func matchLogToTelegram(_ context.Context, record slog.Record) bool {
	return record.Level == slog.LevelError
}
