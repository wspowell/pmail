package mail

import (
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/user"
)

type mailResponse struct {
	MailGuid string `json:"mail_guid"`

	From     string `json:"from"`
	To       string `json:"to"`
	Contents string `json:"contents"`

	SentOn      time.Time `json:"sent_on"`
	DeliveredOn time.Time `json:"delivered_on"`
	OpenedOn    time.Time `json:"opened_on"`
}

type getMailResponse struct {
	mailResponse
}

type openMail struct {
	AuthorizedUser *middleware.UserAuth `spiderweb:"auth"`
	MailGuid       string               `spiderweb:"path=mail_guid"`
	Datastore      db.Datastore         `spiderweb:"resource=datastore"`
	ResponseBody   *getMailResponse     `spiderweb:"response,mime=application/json"`
}

func (self *openMail) Handle(ctx context.Context) (int, error) {
	foundMail, err := self.Datastore.GetMail(ctx, mail.Guid(self.MailGuid))
	if err != nil {
		if errors.Is(err, db.ErrMailNotFound) {
			return httpstatus.InternalServerError, errors.Propagate(icGetMailGetMailMailNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, errors.Propagate(icGetMailGetMailDbError, err)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icGetMailGetMailUnknownDbError, err, errUncaughtDbError)
		}
	}

	// User must be the recipient and must be delivered in order to open the mail.
	if foundMail.CanOpen(user.Guid(self.AuthorizedUser.UserGuid)) && foundMail.IsDelivered() {
		// Check if this is the first time opening the mail.
		if !foundMail.IsOpened() {
			foundMail.OpenedOn = time.Now().UTC()
			err := self.Datastore.OpenMail(ctx, mail.Guid(self.MailGuid), foundMail.OpenedOn)
			if err != nil {
				if errors.Is(err, db.ErrMailNotFound) {
					return httpstatus.InternalServerError, errors.Propagate(icGetMailOpenMailMailNotFound, err)
				} else if errors.Is(err, db.ErrInternalFailure) {
					return httpstatus.InternalServerError, errors.Propagate(icGetMailOpenMailDbError, err)
				} else {
					return httpstatus.InternalServerError, errors.Convert(icGetMailOpenMailUnknownDbError, err, errUncaughtDbError)
				}
			}
		}

		self.ResponseBody.MailGuid = string(foundMail.MailGuid)
		self.ResponseBody.From = string(foundMail.From)
		self.ResponseBody.To = string(foundMail.To)
		self.ResponseBody.Contents = foundMail.Contents
		self.ResponseBody.SentOn = foundMail.SentOn
		self.ResponseBody.DeliveredOn = foundMail.DeliveredOn
		self.ResponseBody.OpenedOn = foundMail.OpenedOn

		return httpstatus.OK, nil
	}

	return httpstatus.NotFound, errors.Propagate(icGetMailUserNotRecipient, errMailNotFound)
}
