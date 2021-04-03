package api

import (
	"net/http"
	"time"

	"github.com/wspowell/pmail/api/users"
	"github.com/wspowell/pmail/resources/db"

	"github.com/wspowell/logging"
	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server) {
	config := &endpoint.Config{
		//Auther:       auth.Noop{},
		//ErrorHandler: error_handlers.ErrorJsonWithCodeResponse{},
		LogConfig: logging.NewConfig(logging.LevelDebug, map[string]interface{}{}),
		//MimeTypeHandlers: endpoint.NewMimeTypeHandlers(),
		//RequestValidator:  validators.NoopRequest{},
		//ResponseValidator: validators.NoopResponse{},
		Resources: map[string]interface{}{
			"userstore":    &db.Users{},
			"mailboxstore": &db.Mailboxes{},
			"mailstore":    &db.Mail{},
		},
		Timeout: 30 * time.Second,
	}

	server.HandleNotFound(config, &noRoute{})
	users.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx *endpoint.Context) (int, error) {
	return http.StatusNotFound, nil
}
