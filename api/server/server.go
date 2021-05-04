package server

import (
	"time"

	"github.com/wspowell/log"
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/spiderweb/server"
)

func New() *server.Server {
	serverConfig := &server.Config{
		LogConfig:    log.NewConfig(log.LevelDebug),
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		EnablePprof:  false,
	}

	server := server.New(serverConfig)
	api.Routes(server)

	return server
}
