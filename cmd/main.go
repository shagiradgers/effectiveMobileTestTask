package main

import (
	"effectiveMobileTestTask/internal/database"
	"effectiveMobileTestTask/internal/handler"
	"effectiveMobileTestTask/internal/server"
	"effectiveMobileTestTask/internal/service"
	"effectiveMobileTestTask/internal/store"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerConfig   *server.Config         `mapstructure:"server"`
	DatabaseConfig *database.SqliteConfig `mapstructure:"database"`
}

var ProjectConfig Config

func init() {
	ProjectConfig = Config{
		ServerConfig: &server.Config{
			Host: "localhost",
			Port: 8081,
		},
		DatabaseConfig: &database.SqliteConfig{
			Name: "db",
		},
	}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	db, err := database.NewSqlite(ProjectConfig.DatabaseConfig)
	if err != nil {
		panic(err)
	}

	userStore, err := store.NewUserStore(db)
	if err != nil {
		panic(err)
	}

	userService := service.New(logrus.New(), userStore)
	h := handler.New(userService)

	s := server.New(h.Engine, ProjectConfig.ServerConfig)
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}
