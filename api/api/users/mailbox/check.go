package mailbox

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/resources/db"
)

type checkMailboxResponse struct {
	Mail []*resources.Mail `json:"mail"`
}

type checkMailbox struct {
	UserId       uint32                `spiderweb:"path=user_id"`
	Mails        resources.MailStore   `spiderweb:"resource=mailstore"`
	Mailboxes    db.Datastore          `spiderweb:"resource=datastore"`
	ResponseBody *checkMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *checkMailbox) Handle(ctx context.Context) (int, error) {
	mailbox, err := self.Mailboxes.GetMailboxByUserId(ctx, self.UserId)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(icCheckMailboxLookupError, err)
	}

	collectedMail, err := self.Mails.CollectMail(ctx, mailbox.Id)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(icCheckMailboxCollectError, err)
	}

	readMail := make([]*resources.Mail, 0, len(collectedMail))
	for _, m := range collectedMail {
		read, err := self.Mails.ReadMail(ctx, m)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(icCheckMailboxReadError, err)
		}

		readMail = append(readMail, read)
	}

	self.ResponseBody.Mail = readMail

	return http.StatusOK, nil
}
