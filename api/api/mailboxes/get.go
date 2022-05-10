package mailboxes

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/db"
)

type getMailboxResponse struct {
	mime.Json

	mailboxModel
}

// FIXME: Must only allow mailbox owner to get mailbox details.
type getMailbox struct {
	pathParams
	Datastore db.Datastore
	body.Response[getMailboxResponse]
}

func (self *getMailbox) Handle(ctx context.Context) (int, error) {
	foundMailbox, err := self.Datastore.GetMailbox(ctx, self.MailboxAddress)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	self.ResponseBody.Address = foundMailbox.Address
	self.ResponseBody.Owner = foundMailbox.UserGuid
	self.ResponseBody.Location = geoLocation{
		Longitude: foundMailbox.Location.Longitude,
		Latitude:  foundMailbox.Location.Latitude,
	}

	return http.StatusOK, nil
}
