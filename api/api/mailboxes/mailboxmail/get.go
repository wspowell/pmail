package mailboxmail

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
)

type getMailboxMailResponse struct {
	mime.Json

	Mail []mailItem `json:"mail"`
}

type mailItem struct {
	MailGuid string `json:"mailGuid"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type getMailboxMail struct {
	auth.User
	pathParams
	Datastore db.Datastore
	body.Response[getMailboxMailResponse]
}

func (self *getMailboxMail) Handle(ctx context.Context) (int, error) {
	foundMail, err := self.Datastore.GetMailboxMail(ctx, self.MailboxAddress)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	self.ResponseBody.Mail = make([]mailItem, len(foundMail))
	for mailIndex := range foundMail {

		self.ResponseBody.Mail[mailIndex] = mailItem{
			MailGuid: foundMail[mailIndex].Guid,
			From:     foundMail[mailIndex].From,
			To:       foundMail[mailIndex].To,
		}
	}

	return http.StatusOK, nil
}
