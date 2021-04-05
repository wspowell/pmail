package server

import (
	"time"

	"github.com/wspowell/logging"
	"github.com/wspowell/pmail/api"
	"github.com/wspowell/spiderweb"
)

func New() *spiderweb.Server {
	serverConfig := &spiderweb.ServerConfig{
		LogConfig:    logging.NewConfig(logging.LevelDebug, map[string]interface{}{}),
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := spiderweb.NewServer(serverConfig)
	api.Routes(server)

	return server
}
