package mailboxes

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type mailboxModel struct {
	Latitude    resources.Latitude  `json:"latitude"`
	Longitude   resources.Longitude `json:"longitude"`
	MailboxType uint8               `json:"mailbox_type"`
}

type createMailboxRequest struct {
	mailboxModel
}

type createMailboxResponse struct {
	MailboxId uint32 `json:"mailbox_id"`
}

type createMailbox struct {
	Users        resources.MailboxStore `spiderweb:"resource=mailboxstore"`
	RequestBody  *createMailboxRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *createMailbox) Handle(ctx *endpoint.Context) (int, error) {
	mailboxAttributes := resources.MailboxAttributes{
		Type: resources.MailboxType(self.RequestBody.MailboxType),
		Location: resources.GeoCoordinate{
			Lat: self.RequestBody.Latitude,
			Lng: self.RequestBody.Longitude,
		},
	}

	mailboxId, err := self.Users.CreateMailbox(mailboxAttributes)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(icCreateMailboxError, err)
	}

	ctx.Debug("created mailbox: %d", mailboxId)

	self.ResponseBody.MailboxId = mailboxId

	return http.StatusCreated, nil
}
