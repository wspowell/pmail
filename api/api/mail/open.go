package mail

import (
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
)

type openMailResponse struct {
	mime.Json

	MailGuid string `json:"mailGuid"`

	FromMailboxAddress string `json:"fromMailboxAddress"`
	ToMailboxAddress   string `json:"toMailboxAddress"`

	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`

	SentOn      time.Time `json:"sentOn"`
	DeliveredOn time.Time `json:"deliveredOn"`
	OpenedOn    time.Time `json:"openedOn"`
}

type openMail struct {
	auth.User
	pathParams
	Datastore db.Datastore
	body.Response[openMailResponse]
}

func (self *openMail) Handle(ctx context.Context) (int, error) {
	openedMail, err := self.Datastore.OpenMail(ctx, self.MailGuid)
	if err != nil {
		if errors.Is(err, db.ErrMailNotFound) {
			return httpstatus.NotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, err
		} else {
			return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	if openedMail.ToGuid != self.AuthorizedUser.Guid {
		return httpstatus.Forbidden, errInvalidRecipient
	}

	self.ResponseBody.MailGuid = openedMail.Guid
	self.ResponseBody.FromMailboxAddress = openedMail.FromMailboxAddress
	self.ResponseBody.ToMailboxAddress = openedMail.ToMailboxAddress
	self.ResponseBody.From = openedMail.Contents.From
	self.ResponseBody.To = openedMail.Contents.To
	self.ResponseBody.Body = openedMail.Contents.Body
	self.ResponseBody.SentOn = openedMail.SentOn
	self.ResponseBody.DeliveredOn = openedMail.DeliveredOn
	self.ResponseBody.OpenedOn = openedMail.OpenedOn

	return httpstatus.OK, nil
}
