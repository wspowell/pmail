package api

import (
	"net/http"
	"time"

	"github.com/wspowell/snailmail/api/authorize"
	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"

	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
)

func Config() *endpoint.Config {
	datastore := db.NewInMemory()

	signingKey, err := auth.GetSigningKey()
	if err != nil {
		// FIXME: Need to add Context() to restful.Server
		log.NewLog(log.NewConfig(log.LevelError)).Fatal("failed to get jwt signing key: %s", err)
	}

	jwtAuth := auth.NewJwt(signingKey)
	middleware.JwtAuth = jwtAuth

	return &endpoint.Config{
		//Auther: userAuth,
		//ErrorHandler: error_handlers.ErrorJsonWithCodeResponse{},
		LogConfig: log.NewConfig(log.LevelDebug),
		//MimeTypeHandlers: endpoint.NewMimeTypeHandlers(),
		//RequestValidator:  validators.NoopRequest{},
		//ResponseValidator: validators.NoopResponse{},

		Resources: map[string]interface{}{
			"datastore": datastore,
			"jwt":       jwtAuth,
		},
		Timeout: 30 * time.Second,
	}
}

func Routes(server *restful.Server) {
	config := Config()

	server.HandleNotFound(config, &noRoute{})
	authorize.Routes(server, config)
	users.Routes(server, config)
	mailboxes.Routes(server, config)
	mailboxmail.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx context.Context) (int, error) {
	return http.StatusNotFound, nil
}
