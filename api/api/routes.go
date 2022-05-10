package api

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/server/restful"

	"github.com/wspowell/snailmail/api/authorize"
	"github.com/wspowell/snailmail/api/exchange"
	"github.com/wspowell/snailmail/api/mail"
	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models"
)

func testSetup(ctx context.Context, datastore db.Datastore) {
	var err error

	// nearbyMailboxes, err := datastore.GetNearbyMailboxes(ctx, models.Coordinate{
	// 	Lat: 33.09387418,
	// 	Lng: -96.90730757,
	// }, 1000.0)

	// log.Info(ctx, "%+v", nearbyMailboxes)

	if err != nil {
		panic(err)
	}

	// Setup test data.

	if true {
		err = datastore.CreateUser(ctx,
			models.CreateUser("public key 1", models.NewCoordinate(-96.9073063043077, 33.09389393115629)))
		if err != nil {
			panic(err)
		}

		// err = datastore.CreateMailbox(ctx, mailbox.Mailbox{
		// 	Address: "AAAAAAAAAAAA",
		// 	Attributes: mailbox.Attributes{
		// 		Owner: "abc-123",
		// 		Location: geo.Coordinate{
		// 			Lat: 88.8,
		// 			Lng: 99.9,
		// 		},
		// 		Capacity: 20,
		// 	},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		err = datastore.CreateUser(ctx, models.CreateUser("public key 2", models.NewCoordinate(-96.91355023160014, 33.08658597452532)))
		if err != nil {
			panic(err)
		}

		// err = datastore.CreateMailbox(ctx, mailbox.Mailbox{
		// 	Address: "BBBBBBBBBBBB",
		// 	Attributes: mailbox.Attributes{
		// 		Owner: "cba-321",
		// 		Location: geo.Coordinate{
		// 			Lat: 11.1,
		// 			Lng: 22.2,
		// 		},
		// 		Capacity: 20,
		// 	},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = datastore.CreateMail(ctx, mailm.Mail{
		// 	MailGuid: "mail-123",
		// 	Attributes: mailm.Attributes{
		// 		From:     "cba-321",
		// 		To:       "abc-123",
		// 		Contents: "Hello there!",
		// 	},
		// 	SentOn:      time.Now().UTC().Add(-3 * time.Hour),
		// 	DeliveredOn: time.Now().UTC().Add(-2 * time.Hour),
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = datastore.CreateMail(ctx, mailm.Mail{
		// 	MailGuid: "mail-123-2",
		// 	Attributes: mailm.Attributes{
		// 		From:     "cba-321",
		// 		To:       "abc-123",
		// 		Contents: "Greetings!",
		// 	},
		// 	SentOn:      time.Now().UTC().Add(-6 * time.Hour),
		// 	DeliveredOn: time.Now().UTC().Add(-5 * time.Hour),
		// 	OpenedOn:    time.Now().UTC().Add(-4 * time.Hour),
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = datastore.CreateMail(ctx, mailm.Mail{
		// 	MailGuid: "mail-123-3",
		// 	Attributes: mailm.Attributes{
		// 		From:     "cba-321",
		// 		To:       "abc-123",
		// 		Contents: "Why wont you answer?",
		// 	},
		// 	SentOn:      time.Now().UTC().Add(-6 * time.Hour),
		// 	DeliveredOn: time.Now().UTC().Add(-5 * time.Hour),
		// 	OpenedOn:    time.Now().UTC().Add(-4 * time.Hour),
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// _, err = datastore.DropOffMail(ctx, "cba-321", "AAAAAAAAAAAA")
		// if err != nil {
		// 	log.Error(ctx, "failed to drop off mail: %+v", err)
		// 	panic(err)
		// }

		// _, err = datastore.PickUpMail(ctx, "abc-123", "AAAAAAAAAAAA")
		// if err != nil {
		// 	log.Error(ctx, "failed to pick up mail: %+v", err)
		// 	panic(err)
		// }
	}
}

func Routes(server *restful.Server) {
	apiResources := resources.Load()

	//testSetup(ctx, datastore)

	server.HandleNotFound(handler.NewHandle(noRoute{}))
	authorize.Routes(server, apiResources)
	users.Routes(server, apiResources)
	mail.Routes(server, apiResources)
	mailboxes.Routes(server, apiResources)
	mailboxmail.Routes(server, apiResources)
	exchange.Routes(server, apiResources)
}

type noRoute struct{}

func (self *noRoute) Handle(ctx context.Context) (int, error) {
	return http.StatusNotFound, nil
}
