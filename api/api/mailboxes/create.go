package mailboxes

import (
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

type mailboxModel struct {
	Label    string      `json:"label"`
	Location geoLocation `json:"location"`
	Capacity uint32      `json:"capacity"`
	Owner    string      `json:"owner"`
}

type geoLocation struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type createMailboxRequest struct {
	mailboxModel
}

type createMailboxResponse struct {
	MailboxGuid string `json:"mailbox_guid"`
}

// FIXME: This should only allow admins to create non-owned mailboxes.
type createMailbox struct {
	Datastore    db.Datastore           `spiderweb:"resource=datastore"`
	RequestBody  *createMailboxRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *createMailbox) Handle(ctx context.Context) (int, error) {
	if self.RequestBody.Owner != "" {
		// Ensure user exists.
		if _, err := self.Datastore.GetUser(ctx, user.Guid(self.RequestBody.Owner)); err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				return httpstatus.NotFound, errors.Propagate(icCreateMailboxUserNotFound, err)
			} else if errors.Is(err, db.ErrInternalFailure) {
				return httpstatus.InternalServerError, errors.Propagate(icCreateMailboxGetUserDbError, err)
			} else {
				return httpstatus.InternalServerError, errors.Convert(icCreateMailboxGetUserUnknownDbError, err, errUncaughtDbError)
			}
		}
	}

	attributes := mailbox.Attributes{
		Label: self.RequestBody.Label,
		Owner: user.Guid(self.RequestBody.Owner),
		Location: geo.Coordinate{
			Lat: geo.Latitude(self.RequestBody.Location.Latitude),
			Lng: geo.Longitude(self.RequestBody.Location.Longitude),
		},
		Capacity: self.RequestBody.Capacity,
	}
	newMailbox := mailbox.NewMailbox(attributes)

	if err := self.Datastore.CreateMailbox(ctx, newMailbox); err != nil {
		if errors.Is(err, db.ErrMailboxGuidExists) {
			return httpstatus.Conflict, errors.Propagate(icCreateMailboxGuidConflict, err)
		} else if errors.Is(err, db.ErrUserMailboxExists) {
			return httpstatus.Conflict, errors.Propagate(icCreateMailboxUserMailboxConflict, err)
		} else if errors.Is(err, db.ErrMailboxLabelExists) {
			return httpstatus.Conflict, errors.Propagate(icCreateMailboxLabelConflict, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, errors.Propagate(icCreateMailboxDbError, err)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icCreateMailboxUnknownDbError, err, errUncaughtDbError)
		}
	}

	log.Debug(ctx, "created mailbox: %+v", newMailbox)

	self.ResponseBody.MailboxGuid = string(newMailbox.MailboxGuid)

	return httpstatus.Created, nil
}
