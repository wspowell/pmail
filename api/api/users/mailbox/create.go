package mailbox

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type mailboxModel struct {
	Location geoLocation `json:"location"`
}

type geoLocation struct {
	Latitude  resources.Latitude  `json:"latitude"`
	Longitude resources.Longitude `json:"longitude"`
}

type createMailboxRequest struct {
	mailboxModel
}

type createMailboxResponse struct {
	MailboxId uint32 `json:"mailbox_id"`
}

type createMailbox struct {
	UserId       uint32                 `spiderweb:"path=user_id"`
	Mailboxes    resources.MailboxStore `spiderweb:"resource=mailboxstore"`
	RequestBody  *createMailboxRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *createMailbox) Handle(ctx *endpoint.Context) (int, error) {
	mailbox, err := self.Mailboxes.GetMailboxByUserId(self.UserId)
	if err != nil {
		// This error is expected.
		if !errors.Is(err, resources.ErrorMailboxNotFound) {
			return http.StatusInternalServerError, errors.Wrap(icCreateMailboxExistenceCheckError, err)
		}
	}

	if mailbox != nil {
		return http.StatusConflict, errors.Wrap(icCreateMailboxMailboxExists, err)
	}

	mailboxAttributes := resources.MailboxAttributes{
		Location: resources.GeoCoordinate{
			Lat: self.RequestBody.Location.Latitude,
			Lng: self.RequestBody.Location.Longitude,
		},
	}

	mailboxId, err := self.Mailboxes.CreateMailbox(self.UserId, mailboxAttributes)
	if err != nil {
		if errors.Is(err, resources.ErrHomeMailboxExists) {
			return http.StatusConflict, errors.Wrap(icCreateMailboxMailboxExistsOnCreate, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icCreateMailboxError, err)
	}

	log.Debug(ctx, "created mailbox: %d", mailboxId)

	self.ResponseBody.MailboxId = mailboxId

	return http.StatusCreated, nil
}
