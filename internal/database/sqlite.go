package database

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type SqliteConfig struct {
	Name string `mapstructure:"name"`
}

func NewSqlite(config *SqliteConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", config.Name)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
