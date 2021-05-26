package api

import (
	"net/http"
	"time"

	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/api/users/mailbox"
	"github.com/wspowell/snailmail/resources/db"

	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
)

func Config() *endpoint.Config {
	return &endpoint.Config{
		//Auther:       auth.Noop{},
		//ErrorHandler: error_handlers.ErrorJsonWithCodeResponse{},
		LogConfig: log.NewConfig(log.LevelDebug),
		//MimeTypeHandlers: endpoint.NewMimeTypeHandlers(),
		//RequestValidator:  validators.No opRequest{},
		//ResponseValidator: validators.NoopResponse{},
		Resources: map[string]interface{}{
			"userstore":    db.NewUsers(),
			"mailboxstore": db.NewMailboxes(),
			"mailstore":    db.NewMails(),
		},
		Timeout: 30 * time.Second,
	}
}

func Routes(server *restful.Server) {
	config := Config()

	server.HandleNotFound(config, &noRoute{})
	users.Routes(server, config)
	mailbox.Routes(server, config)
	mailboxes.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx context.Context) (int, error) {
	return http.StatusNotFound, nil
}
