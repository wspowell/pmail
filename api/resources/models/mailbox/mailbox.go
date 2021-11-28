package mailbox

import (
	"math/rand"

	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/user"
)

type Mailbox struct {
	// Address of the mailbox to show the user.
	// Must be globally unique.
	// This is a code that can be used for others to send mail to you.
	Address string
	Attributes
}

type Attributes struct {
	// Owner of the mailbox
	// If owner is blank, then the mailbox is a public exchange mailbox.
	Owner user.Guid

	// Location in the world.
	Location geo.Coordinate

	// Capacity of mail in the mailbox.
	// 0 - Pickup only
	Capacity uint32
}

const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func NewMailbox(attributes Attributes) Mailbox {
	return Mailbox{
		// 599,555,620,984,320,000 permutations (36 character set, sets of 12)
		Address:    randStringBytes(12),
		Attributes: attributes,
	}
}

// FormatAddress as "AAAA-AAAA-AAAA"
func (self Mailbox) FormatAddress() string {
	formattedAddress := ""

	for i := 0; i < 4; i++ {
		formattedAddress += string(self.Address[i])
	}
	formattedAddress += "-"
	for i := 4; i < 8; i++ {
		formattedAddress += string(self.Address[i])
	}
	formattedAddress += "-"
	for i := 8; i < 12; i++ {
		formattedAddress += string(self.Address[i])
	}

	return formattedAddress
}

func (self Mailbox) IsDropoff() bool {
	return self.Capacity != 0
}

func (self Mailbox) IsPublic() bool {
	return self.Owner != ""
}

func (self Mailbox) IsNearby(location geo.Coordinate, radiusMeters float32) bool {
	// TODO
	return true
}
