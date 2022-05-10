package mail

import (
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models"
)

type createMailRequest struct {
	mime.Json

	From             string `json:"from"`
	To               string `json:"to"`
	ToMailboxAddress string `json:"toMailboxAddress"`
	Body             string `json:"body"`
}

type createMail struct {
	auth.User
	Datastore db.Datastore
	body.Request[createMailRequest]
}

func (self *createMail) Handle(ctx context.Context) (int, error) {
	if self.RequestBody.Body == "" {
		return httpstatus.UnprocessableEntity, errEmptyContents
	}

	contents := models.MailContents{
		From: self.RequestBody.From,
		To:   self.RequestBody.To,
		Body: self.RequestBody.Body,
	}
	newMail := models.CreateMail(self.AuthorizedUser.Guid, self.AuthorizedUser.Mailbox.Address, self.RequestBody.ToMailboxAddress, contents)
	newMail.SentOn = time.Now().UTC()

	err := self.Datastore.CreateMail(ctx, newMail)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return httpstatus.NotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, err
		} else {
			return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	log.Debug(ctx, "created mail: %+v", newMail)

	return httpstatus.Created, nil
}
