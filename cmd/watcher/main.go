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

	"github.com/DanilaKorobkov/defi-monitoring/internal"
	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
	"github.com/DanilaKorobkov/defi-monitoring/internal/domain/services/watcher"
	"github.com/DanilaKorobkov/defi-monitoring/internal/infra/notifiers/telegram"
	"github.com/DanilaKorobkov/defi-monitoring/internal/infra/positions_providers"
	"github.com/DanilaKorobkov/defi-monitoring/internal/infra/positions_providers/base/aerodrome"
	"github.com/DanilaKorobkov/defi-monitoring/internal/infra/positions_providers/base/uniswap_v3"
)

type Config struct {
	TelegramBotToken            string        `env:"TELEGRAM_BOT_TOKEN,required,unset"`
	ErrorReceiverTelegramUserID int64         `env:"ERROR_RECEIVER_TELEGRAM_USER_ID,required,unset"`
	SubjectTelegramUserID       int64         `env:"SUBJECT_TELEGRAM_USER_ID,required,unset"`
	SubjectWallet               string        `env:"SUBJECT_WALLET,required,unset"`
	BaseUniswapV3SubgraphID     string        `env:"BASE_UNISWAP_V3_SUBGRAPH_ID,required"`
	BaseAerodromeSubgraphID     string        `env:"BASE_AERODROME_SUBGRAPH_ID,required"`
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

	baseUniswapV3 := makeBaseUniswapV3Provider(config)
	baseAerodrome := makeBaseAerodromeProvider(config)
	lp := positions_providers.NewComposite(baseUniswapV3, baseAerodrome)

	watcherConfig := watcher.ServiceConfig{
		LiquidityPoolPositions: lp,
		Notifier:               telegramNotifier,
		CheckInterval:          config.CheckInterval,
		Logger:                 logger,
	}
	watcherService := watcher.NewService(watcherConfig)

	subject := domain.Subject{
		TelegramUserID: config.SubjectTelegramUserID,
		Wallet:         config.SubjectWallet,
	}

	logger.Info("starting watcher")
	watcherService.StartWatching(ctx, subject)
	logger.Info("watcher finished")
}

func makeBaseUniswapV3Provider(config Config) *uniswap_v3.ProviderTheGraph {
	url := "https://gateway.thegraph.com/api/subgraphs/id/" + config.BaseUniswapV3SubgraphID
	setAuth := func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+config.TheGraphToken)
	}
	client := graphql.NewClient(url, http.DefaultClient).WithRequestModifier(setAuth)

	return uniswap_v3.NewProviderTheGraph(client)
}

func makeBaseAerodromeProvider(config Config) *aerodrome.ProviderTheGraph {
	url := "https://gateway.thegraph.com/api/subgraphs/id/" + config.BaseAerodromeSubgraphID
	setAuth := func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+config.TheGraphToken)
	}
	client := graphql.NewClient(url, http.DefaultClient).WithRequestModifier(setAuth)

	return aerodrome.NewProviderTheGraph(client)
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("-", slog.String("err", err.Error()))
	os.Exit(-1) //nolint:revive // It's easier
}
