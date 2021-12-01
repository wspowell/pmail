package mail

import (
	"time"

	"github.com/google/uuid"

	"github.com/wspowell/snailmail/resources/models/user"
)

type Guid string

func ToStrings(slice []Guid) []string {
	asStrings := make([]string, len(slice))
	for index := range slice {
		asStrings[index] = string(slice[index])
	}

	return asStrings
}

type Mail struct {
	MailGuid Guid
	Attributes

	// Metadata
	SentOn      time.Time
	DeliveredOn time.Time
	OpenedOn    time.Time
}

type Attributes struct {
	From     user.Guid
	To       user.Guid
	Carrier  user.Guid
	Contents string
}

func NewMail(attributes Attributes) Mail {
	return Mail{
		MailGuid:   Guid(uuid.New().String()),
		Attributes: attributes,
	}
}

func (self Mail) CanOpen(userGuid user.Guid) bool {
	return self.To == userGuid
}

func (self Mail) IsSent() bool {
	return !self.SentOn.IsZero()
}

func (self Mail) IsDelivered() bool {
	return !self.DeliveredOn.IsZero()
}

func (self Mail) IsOpened() bool {
	return !self.OpenedOn.IsZero()
}
