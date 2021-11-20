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

type createMailRequest struct {
	To       string `json:"to"`
	Contents string `json:"contents"`
}

type createMailResponse struct {
	MailGuid string `json:"mail_guid"`
}

type createMail struct {
	AuthorizedUser *middleware.UserAuth `spiderweb:"auth"`
	Datastore      db.Datastore         `spiderweb:"resource=datastore"`
	RequestBody    *createMailRequest   `spiderweb:"request,mime=application/json"`
	ResponseBody   *createMailResponse  `spiderweb:"response,mime=application/json"`
}

func (self *createMail) Handle(ctx context.Context) (int, error) {
	if self.RequestBody.Contents == "" {
		return httpstatus.UnprocessableEntity, errors.Propagate(icCreateMailEmptyContents, errEmptyContents)
	}

	// Check that the "to" mailbox code exists.
	toMailbox, err := self.Datastore.GetMailboxByLabel(ctx, self.RequestBody.To)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return httpstatus.NotFound, errors.Propagate(icCreateMailGetMailboxByLabelMailboxNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, errors.Propagate(icCreateMailGetMailboxByLabelDbError, err)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icCreateMailGetMailboxByLabelUnknownDbError, err, errUncaughtDbError)
		}
	}

	if toMailbox.Owner == user.Guid(self.AuthorizedUser.UserGuid) {
		return httpstatus.Conflict, errors.Propagate(icCreateMailInvalidRecipient, errInvalidRecipient)
	}

	mailAttributes := mail.Attributes{
		From:     user.Guid(self.AuthorizedUser.UserGuid),
		To:       toMailbox.Owner,
		Contents: self.RequestBody.Contents,
		Carrier:  user.Guid(self.AuthorizedUser.UserGuid),
	}
	newMail := mail.NewMail(mailAttributes)
	newMail.SentOn = time.Now().UTC()

	err = self.Datastore.CreateMail(ctx, newMail)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return httpstatus.NotFound, errors.Propagate(icCreateMailCreateMailboxMailboxNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, errors.Propagate(icCreateMailCreateMailDbError, err)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icCreateMailCreateMailUnknownDbError, err, errUncaughtDbError)
		}
	}

	self.ResponseBody.MailGuid = string(newMail.MailGuid)

	return httpstatus.Created, nil
}
