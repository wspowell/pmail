package mailbox

import (
	"github.com/google/uuid"
	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/user"
)

type Guid string

type Mailbox struct {
	MailboxGuid Guid
	Attributes
}

type Attributes struct {
	// Label of the mailbox to show the user.
	// Must be globally unique.
	Label string

	// Owner of the mailbox
	// If owner is blank, then the mailbox is a public exchange mailbox.
	Owner user.Guid

	// Location in the world.
	Location geo.Coordinate

	// Capacity of mail in the mailbox.
	// 0 - Pickup only
	Capacity uint32
}

func NewMailbox(attributes Attributes) Mailbox {
	return Mailbox{
		MailboxGuid: Guid(uuid.New().String()),
		Attributes:  attributes,
	}
}

func (self Mailbox) IsDropoff() bool {
	return self.Capacity != 0
}

func (self Mailbox) IsPublic() bool {
	return self.Owner != ""
}

func (self Mailbox) IsNearby(location geo.Coordinate, radiusMeters uint32) bool {
	// TODO
	return true
}
