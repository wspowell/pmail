package api

import (
	"net/http"
	"time"

	"github.com/wspowell/pmail/api/mailboxes"
	"github.com/wspowell/pmail/api/users"
	"github.com/wspowell/pmail/api/users/mailbox"
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
		//RequestValidator:  validators.No opRequest{},
		//ResponseValidator: validators.NoopResponse{},
		Resources: map[string]interface{}{
			"userstore":    db.NewUsers(),
			"mailboxstore": db.NewMailboxes(),
			"mailstore":    db.NewMails(),
		},
		Timeout: 30 * time.Second,
	}

	server.HandleNotFound(config, &noRoute{})
	users.Routes(server, config)
	mailbox.Routes(server, config)
	mailboxes.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx *endpoint.Context) (int, error) {
	return http.StatusNotFound, nil
}
