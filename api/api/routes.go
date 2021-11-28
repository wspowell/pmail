package api

import (
	"net/http"
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"

	"github.com/wspowell/snailmail/api/authorize"
	"github.com/wspowell/snailmail/api/mail"
	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/geo"
	mailm "github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

func Config() *endpoint.Config {
	ctx := context.Local()
	ctx = log.WithContext(ctx, log.NewConfig(log.LevelDebug))

	//datastore := db.NewInMemory()
	datastore := db.NewMySql()

	var err error

	if err := datastore.Connect(); err != nil {
		log.Fatal(ctx, "%v", err)
	}
	if err := datastore.Migrate(); err != nil {
		log.Fatal(ctx, "%v", err)
	}

	nearbyMailboxes, err := datastore.GetNearbyMailboxes(ctx, geo.Coordinate{
		Lat: 33.09387418,
		Lng: -96.90730757,
	}, 1000.0)

	log.Info(ctx, "%+v", nearbyMailboxes)

	if err != nil {
		panic(err)
	}

	// Setup test data.

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
		Address: "AAAAAAAAAAAA",
		Attributes: mailbox.Attributes{
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
		Address: "BBBBBBBBBBBB",
		Attributes: mailbox.Attributes{
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

	err = datastore.CreateMail(ctx, mailm.Mail{
		MailGuid: "mail-123",
		Attributes: mailm.Attributes{
			From:     "cba-321",
			To:       "abc-123",
			Contents: "Hello there!",
		},
		SentOn:      time.Now().UTC().Add(-3 * time.Hour),
		DeliveredOn: time.Now().UTC().Add(-2 * time.Hour),
	})
	if err != nil {
		panic(err)
	}

	err = datastore.CreateMail(ctx, mailm.Mail{
		MailGuid: "mail-123-2",
		Attributes: mailm.Attributes{
			From:     "cba-321",
			To:       "abc-123",
			Contents: "Greetings!",
		},
		SentOn:      time.Now().UTC().Add(-6 * time.Hour),
		DeliveredOn: time.Now().UTC().Add(-5 * time.Hour),
		OpenedOn:    time.Now().UTC().Add(-4 * time.Hour),
	})
	if err != nil {
		panic(err)
	}

	err = datastore.CreateMail(ctx, mailm.Mail{
		MailGuid: "mail-123-3",
		Attributes: mailm.Attributes{
			From:     "cba-321",
			To:       "abc-123",
			Contents: "Why wont you answer?",
		},
		SentOn:      time.Now().UTC().Add(-6 * time.Hour),
		DeliveredOn: time.Now().UTC().Add(-5 * time.Hour),
		OpenedOn:    time.Now().UTC().Add(-4 * time.Hour),
	})
	if err != nil {
		panic(err)
	}

	_, err = datastore.DropOffMail(ctx, user.Guid("cba-321"), "AAAAAAAAAAAA")
	if err != nil {
		panic(err)
	}

	_, err = datastore.PickUpMail(ctx, user.Guid("abc-123"), "AAAAAAAAAAAA")
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
