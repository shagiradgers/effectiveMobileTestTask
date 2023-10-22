package server

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func New(handler http.Handler, config *Config) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", config.Host, config.Port),
		Handler:           handler,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
	}
	return server
}
