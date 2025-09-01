package migrators

import (
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// Postgres driver.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// File system driver for embed.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MakePostgresMigrator(url string, embed fs.FS) (*migrate.Migrate, error) {
	return MakePostgresMigratorWithPath(url, embed, "sql")
}

func MakePostgresMigratorWithPath(url string, embed fs.FS, path string) (*migrate.Migrate, error) {
	dataDriver, err := iofs.New(embed, path)
	if err != nil {
		return nil, fmt.Errorf("load embed migration files: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", dataDriver, url)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return migrator, nil
}
