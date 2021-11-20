package api

import (
	"net/http"
	"time"

	"github.com/wspowell/snailmail/api/authorize"
	"github.com/wspowell/snailmail/api/mail"
	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"

	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
)

func Config() *endpoint.Config {
	datastore := db.NewInMemory()

	// Setup test data.
	var err error
	ctx := context.Local()

	err = datastore.CreateUser(ctx, user.User{
		UserGuid: user.Guid("abc-123"),
		Attributes: user.Attributes{
			Username:          "test1",
			PineappleOnPizza:  true,
			MailCarryCapacity: 20,
		},
	}, "a")
	if err != nil {
		panic(err)
	}

	err = datastore.CreateMailbox(ctx, mailbox.Mailbox{
		MailboxGuid: mailbox.Guid("mailbox-user1"),
		Attributes: mailbox.Attributes{
			Label: "ABCD",
			Owner: user.Guid("abc-123"),
			Location: geo.Coordinate{
				Lat: 88.8,
				Lng: 99.9,
			},
			Capacity: 20,
		},
	})
	if err != nil {
		panic(err)
	}

	err = datastore.CreateUser(ctx, user.User{
		UserGuid: user.Guid("cba-321"),
		Attributes: user.Attributes{
			Username:          "test2",
			PineappleOnPizza:  true,
			MailCarryCapacity: 20,
		},
	}, "a")
	if err != nil {
		panic(err)
	}

	err = datastore.CreateMailbox(ctx, mailbox.Mailbox{
		MailboxGuid: mailbox.Guid("mailbox-user2"),
		Attributes: mailbox.Attributes{
			Label: "EFGH",
			Owner: user.Guid("cba-321"),
			Location: geo.Coordinate{
				Lat: 11.1,
				Lng: 22.2,
			},
			Capacity: 20,
		},
	})
	if err != nil {
		panic(err)
	}

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
	mail.Routes(server, config)
	mailboxes.Routes(server, config)
	mailboxmail.Routes(server, config)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx context.Context) (int, error) {
	return http.StatusNotFound, nil
}
