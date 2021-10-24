package api

import (
	"net/http"
	"time"

	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/resources/db"

	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
)

func Config() *endpoint.Config {
	datastore := db.NewInMemory()

	return &endpoint.Config{
		//Auther: userAuth,
		//ErrorHandler: error_handlers.ErrorJsonWithCodeResponse{},
		LogConfig: log.NewConfig(log.LevelDebug),
		//MimeTypeHandlers: endpoint.NewMimeTypeHandlers(),
		//RequestValidator:  validators.NoopRequest{},
		//ResponseValidator: validators.NoopResponse{},
		Resources: map[string]interface{}{
			"datastore": datastore,
		},
		Timeout: 30 * time.Second,
	}
}

func Routes(server *restful.Server) {
	config := Config()

	server.HandleNotFound(config, &noRoute{})
	users.Routes(server, config)
	mailboxes.Routes(server, config)
	mailboxmail.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx context.Context) (int, error) {
	return http.StatusNotFound, nil
}
