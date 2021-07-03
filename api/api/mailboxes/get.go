package mailboxes

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
)

type mailboxModel struct {
	Location geoLocation `json:"location"`
}

type geoLocation struct {
	Latitude  resources.Latitude  `json:"latitude"`
	Longitude resources.Longitude `json:"longitude"`
}

type getMailboxResponse struct {
	mailboxModel
}

type getMailbox struct {
	MailboxId    uint32                 `spiderweb:"path=id"`
	Mailboxes    resources.MailboxStore `spiderweb:"resource=mailboxstore"`
	ResponseBody *getMailboxResponse    `spiderweb:"response,mime=application/json"`
}

func (self *getMailbox) Handle(ctx context.Context) (int, error) {
	mailbox, err := self.Mailboxes.GetMailboxById(ctx, self.MailboxId)
	if err != nil {
		if errors.Is(err, resources.ErrorMailboxNotFound) {
			return http.StatusNotFound, errors.Wrap(icGetMailboxNotFound, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icGetMailboxError, err)
	}

	self.ResponseBody.Location.Latitude = mailbox.Attributes.Location.Lat
	self.ResponseBody.Location.Longitude = mailbox.Attributes.Location.Lng

	return http.StatusOK, nil
}
