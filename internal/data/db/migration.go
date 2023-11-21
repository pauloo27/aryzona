package db

import (
	"github.com/rubenv/sql-migrate"
)

var (
	migrations = migrate.FileMigrationSource{
		Dir: "migration",
	}
)

func (db *DBConn) Migrate() (int, error) {
	rawDB := db.Engine.DB().DB
	return migrate.Exec(rawDB, "postgres", migrations, migrate.Up)
}
