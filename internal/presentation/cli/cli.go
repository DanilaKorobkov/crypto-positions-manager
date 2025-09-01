package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v3"

	"github.com/DanilaKorobkov/defi-monitoring/pkg/migrators"
)

type Config struct {
	DBCommandConfig DBCommandConfig
}

type DBCommandConfig struct {
	MigratorCommandConfig MigratorCommandConfig
}

type MigratorCommandConfig struct {
	PostgresURLEnvName string
	Migrations         fs.FS
	Logger             *slog.Logger
}

func New(config Config) *cli.Command {
	return &cli.Command{
		Commands: []*cli.Command{
			newDBCommands(config.DBCommandConfig),
		},
	}
}

func newDBCommands(config DBCommandConfig) *cli.Command {
	return &cli.Command{
		Name: "db",
		Commands: []*cli.Command{
			newMigrateCommand(config.MigratorCommandConfig),
		},
	}
}

//nolint:maintidx // How to simplify?
func newMigrateCommand(config MigratorCommandConfig) *cli.Command {
	var pgURL string

	return &cli.Command{
		Name:  "migrations",
		Usage: "migrate database to latest version",
		Flags: []cli.Flag{
			makeToURLFlag(&pgURL, config.PostgresURLEnvName),
		},
		Action: func(context.Context, *cli.Command) error {
			migrator, err := migrators.MakePostgresMigratorWithPath(pgURL, config.Migrations, "sql")
			if err != nil {
				message := fmt.Sprintf("migrators.MakePostgresMigratorWithPath: %s", err)
				return cli.Exit(message, -1)
			}
			err = migrateAction(migrator)
			if err != nil {
				return cli.Exit(err.Error(), -1)
			}
			err = migrateCheck(migrator, config.Logger)
			if err != nil {
				return cli.Exit(err.Error(), -1)
			}
			return nil
		},
	}
}

func makeToURLFlag(destination *string, envName string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "url",
		Usage:       "Postgres DSN e.g: postgres://user:password@localhost:5432/dbname?sslmode=disable",
		Required:    true,
		Destination: destination,
		Sources:     cli.EnvVars(strings.ToUpper(envName)),
	}
}

func migrateCheck(migrator *migrate.Migrate, logger *slog.Logger) error {
	resultVersion, isDirty, err := migrator.Version()
	if err != nil {
		return fmt.Errorf("check migrations: %w", err)
	}

	logger.With("version", resultVersion).
		With("isDirty", isDirty).
		Info("migrator version")

	return nil
}

func migrateAction(migrator *migrate.Migrate) error {
	err := migrator.Up()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return fmt.Errorf("migrate to latest: %w", err)
}
