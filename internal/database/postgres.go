package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgConfig struct {
	Database string `mapstructure:"database"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func NewPostgres(config *PgConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		"postgres", config.User, config.Password, config.Host, config.Port, config.Database)
	db, err := sqlx.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
