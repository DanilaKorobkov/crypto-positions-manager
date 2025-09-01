package main

import (
	"context"
	"github.com/DanilaKorobkov/defi-monitoring/internal/presentation/cli"
	"github.com/DanilaKorobkov/defi-monitoring/migrations"
	"log/slog"
	"os"
)

func main() {
	logger := slog.Default()
	logger.Info(os.Getenv("POSTGRES_URL"))
	config := cli.Config{
		DBCommandConfig: cli.DBCommandConfig{
			MigratorCommandConfig: cli.MigratorCommandConfig{
				PostgresURLEnvName: "POSTGRES_URL",
				Migrations:         migrations.GetMigrationsFS(),
				Logger:             logger,
			},
		},
	}
	command := cli.New(config)
	err := command.Run(context.Background(), os.Args)
	if err != nil {
		fatal(logger, err)
	}
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("-", slog.String("err", err.Error()))
	os.Exit(-1) //nolint:revive // It's easier
}
