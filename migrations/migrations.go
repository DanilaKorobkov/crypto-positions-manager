package migrations

import (
	"embed"
	"io/fs"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

func GetMigrationsFS() fs.FS {
	return migrationsFS
}
