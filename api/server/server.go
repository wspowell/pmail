package server

import (
	"time"

	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/server/restful"

	"github.com/wspowell/snailmail/api"
)

func New() *restful.Server {
	serverConfig := &restful.ServerConfig{
		LogConfig:    log.NewConfig(log.LevelDebug),
		Host:         "0.0.0.0",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		EnablePprof:  false,
	}

	server := restful.NewServer(serverConfig)
	api.Routes(server)

	return server
}
