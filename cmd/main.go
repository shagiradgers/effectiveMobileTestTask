package main

import (
	"effectiveMobileTestTask/internal/database"
	"effectiveMobileTestTask/internal/server"
)

type Config struct {
	ServerConfig   server.Config   `mapstructure:"server"`
	DatabaseConfig database.Config `mapstructure:"database"`
}

type ProjectConfig Config

func init() {

}

func main() {

}
