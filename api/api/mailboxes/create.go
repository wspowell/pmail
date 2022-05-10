package mailboxes

import (
	"github.com/wspowell/context"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/db"
)

type mailboxModel struct {
	Address  string      `json:"address"`
	Location geoLocation `json:"location"`
	Owner    string      `json:"owner"`
}

type geoLocation struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type createMailboxRequest struct {
	mime.Json

	mailboxModel
}

type createMailboxResponse struct {
	mime.Json

	MailboxAddress string `json:"mailboxAddress"`
}

// FIXME: This should only allow admins to create non-owned mailboxes.
type createMailbox struct {
	Datastore db.Datastore
	body.Request[createMailboxRequest]
	body.Response[createMailboxResponse]
}

func (self *createMailbox) Handle(ctx context.Context) (int, error) {
	// if self.RequestBody.Owner != "" {
	// 	// Ensure user exists.
	// 	if _, err := self.Datastore.GetUser(ctx, user.Guid(self.RequestBody.Owner)); err != nil {
	// 		if errors.Is(err, db.ErrUserNotFound) {
	// 			return httpstatus.NotFound, err
	// 		} else if errors.Is(err, db.ErrInternalFailure) {
	// 			return httpstatus.InternalServerError, err
	// 		} else {
	// 			return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 		}
	// 	}
	// }

	// attributes := mailbox.Attributes{
	// 	Owner: user.Guid(self.RequestBody.Owner),
	// 	Location: geo.Coordinate{
	// 		Lat: geo.Latitude(self.RequestBody.Location.Latitude),
	// 		Lng: geo.Longitude(self.RequestBody.Location.Longitude),
	// 	},
	// 	Capacity: self.RequestBody.Capacity,
	// }
	// newMailbox := mailbox.NewMailbox(attributes)

	// if err := self.Datastore.CreateMailbox(ctx, newMailbox); err != nil {
	// 	if errors.Is(err, db.ErrMailboxAddressExists) {
	// 		return httpstatus.Conflict, err
	// 	} else if errors.Is(err, db.ErrUserMailboxExists) {
	// 		return httpstatus.Conflict, err
	// 	} else if errors.Is(err, db.ErrMailboxLabelExists) {
	// 		return httpstatus.Conflict, err
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return httpstatus.InternalServerError, err
	// 	} else {
	// 		return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// }

	// log.Debug(ctx, "created mailbox: %+v", newMailbox)

	// self.ResponseBody.MailboxAddress = newMailbox.Address

	return httpstatus.Created, nil
}
