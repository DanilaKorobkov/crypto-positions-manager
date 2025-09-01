package migrators

import (
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	// Postgres driver.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// File system driver for embed.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func MakePostgresMigrator(url string, fs fs.FS) (*migrate.Migrate, error) {
	return MakePostgresMigratorWithPath(url, fs, "sql")
}

func MakePostgresMigratorWithPath(url string, fs fs.FS, path string) (*migrate.Migrate, error) {
	dataDriver, err := iofs.New(fs, path)
	if err != nil {
		return nil, fmt.Errorf("load embed migration files: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", dataDriver, url)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return migrator, nil
}
