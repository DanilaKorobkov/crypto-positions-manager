package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hasura/go-graphql-client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DanilaKorobkov/crypto-positions-manager/internal"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain/services/watcher"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/infra/notifiers/telegram"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/infra/positions_providers/uniswap_v3_base"
)

type Config struct {
	TelegramBotToken            string        `env:"TELEGRAM_BOT_TOKEN,required,unset"`
	ErrorReceiverTelegramUserID int64         `env:"ERROR_RECEIVER_TELEGRAM_USER_ID,required,unset"`
	WatchTelegramUserID         int64         `env:"WATCH_TELEGRAM_USER_ID,required,unset"`
	WatchWallet                 string        `env:"WATCH_WALLET,required,unset"`
	SubgraphID                  string        `env:"SUBGRAPH_ID,required"`
	TheGraphToken               string        `env:"THE_GRAPH_TOKEN,required,unset"`
	CheckInterval               time.Duration `env:"CHECK_INTERVAL,required"`
}

//nolint:funlen,maintidx // How to make better?
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		fatal(slog.Default(), fmt.Errorf("env.Parse: %w", err))
	}

	loggerConfig := internal.LoggerConfig{
		TelegramBotToken:       config.TelegramBotToken,
		TelegramReceiverUserID: config.ErrorReceiverTelegramUserID,
	}
	logger := internal.NewLogger(loggerConfig)

	telegramBot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		fatal(slog.Default(), fmt.Errorf("tgbotapi.NewBotAPI: %w", err))
	}

	telegramNotifier := telegram.NewNotifier(telegramBot)

	url := "https://gateway.thegraph.com/api/subgraphs/id/" + config.SubgraphID
	client := graphql.NewClient(url, http.DefaultClient).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer "+config.TheGraphToken)
		})
	uniswap := uniswap_v3_base.NewProvider(client)

	watcherConfig := watcher.ServiceConfig{
		Provider:      uniswap,
		Notifier:      telegramNotifier,
		CheckInterval: config.CheckInterval,
		Logger:        logger,
	}
	watcherService := watcher.NewService(watcherConfig)

	watchForUser := domain.User{
		TelegramUserID: config.WatchTelegramUserID,
		Wallet:         config.WatchWallet,
	}

	logger.Info("starting watcher")
	watcherService.StartWatching(ctx, watchForUser)
	logger.Info("watcher finished")
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("-", slog.String("err", err.Error()))
	os.Exit(-1) //nolint:revive // It's easier
}
