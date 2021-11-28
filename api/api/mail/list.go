package mail

import (
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type listMailResponse struct {
	Mail []mailResponse `json:"mail"`
}

type listMail struct {
	AuthorizedUser *middleware.UserAuth `spiderweb:"auth"`
	Datastore      db.Datastore         `spiderweb:"resource=datastore"`
	ResponseBody   *listMailResponse    `spiderweb:"response,mime=application/json"`
}

func (self *listMail) Handle(ctx context.Context) (int, error) {
	userMail, err := self.Datastore.GetUserMail(ctx, user.Guid(self.AuthorizedUser.UserGuid))
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return httpstatus.InternalServerError, errors.Propagate(icListMailGetUserMailUserNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, errors.Propagate(icListMailGetUserMailDbError, err)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icListMailGetUserMailUnknownDbError, err, errUncaughtDbError)
		}
	}

	self.ResponseBody.Mail = []mailResponse{}
	for index := range userMail {
		mailResponseItem := mailResponse{
			MailGuid:    string(userMail[index].MailGuid),
			From:        string(userMail[index].From),
			To:          string(userMail[index].To),
			Contents:    userMail[index].Contents,
			SentOn:      userMail[index].SentOn,
			DeliveredOn: userMail[index].DeliveredOn,
			OpenedOn:    userMail[index].OpenedOn,
		}

		self.ResponseBody.Mail = append(self.ResponseBody.Mail, mailResponseItem)
	}

	return httpstatus.OK, nil
}
