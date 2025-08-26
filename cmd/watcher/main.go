package main

import (
	"context"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/domain/services/watcher"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/infra/notifiers/telegram"
	"github.com/DanilaKorobkov/crypto-positions-manager/internal/infra/positions_providers/uniswapV3_base"
	"github.com/caarlos0/env/v11"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	TelegramBotToken string        `env:"TELEGRAM_BOT_TOKEN,required,unset"`
	TelegramUserID   int64         `env:"TELEGRAM_USER_ID,required,unset"`
	WalletAddress    string        `env:"WALLET_ADDRESS,required,unset"`
	SubgraphID       string        `env:"SUBGRAPH_ID,required"`
	TheGraphToken    string        `env:"THE_GRAPH_TOKEN,required,unset"`
	CheckInterval    time.Duration `env:"CHECK_INTERVAL,required"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create telegram bot")
	}

	telegramNotifier := telegram.NewNotifier(telegram.NotifierConfig{
		Telegram: bot,
		UserID:   config.TelegramUserID,
	})

	url := "https://gateway.thegraph.com/api/subgraphs/id/" + config.SubgraphID
	client := graphql.NewClient(url, http.DefaultClient).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer "+config.TheGraphToken)
		})
	uniswap := uniswapV3_base.NewProvider(client)

	watcherConfig := watcher.ServiceConfig{
		Provider:      uniswap,
		Notifier:      telegramNotifier,
		CheckInterval: config.CheckInterval,
		Logger:        logger,
	}
	w := watcher.NewService(watcherConfig)
	logger.Info().Msg("starting watcher")
	w.StartWatching(ctx, config.WalletAddress)
	logger.Info().Msg("watcher finished")
}
