package server

import (
	"io"
	"time"

	"github.com/wspowell/logging"
	"github.com/wspowell/pmail/api"
	"github.com/wspowell/spiderweb"
)

type NoopLogConfig struct {
	*logging.Config
}

func (self *NoopLogConfig) Out() io.Writer {
	return io.Discard
}

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
