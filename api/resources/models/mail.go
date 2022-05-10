package models

import (
	"time"

	"github.com/google/uuid"
)

type Mail struct {
	Guid string

	FromGuid           string
	FromMailboxAddress string
	ToMailboxAddress   string
	ToGuid             string

	Contents MailContents

	// Metadata
	SentOn      time.Time
	DeliveredOn time.Time
	OpenedOn    time.Time
}

type MailContents struct {
	From string
	To   string
	Body string
}

type MailPreview struct {
	Guid string
	From string
	To   string
}

func CreateMail(fromGuid string, fromMailboxAddress string, toMailboxAddress string, contents MailContents) Mail {
	return Mail{
		Guid:               uuid.New().String(),
		FromGuid:           fromGuid,
		FromMailboxAddress: fromMailboxAddress,
		ToMailboxAddress:   toMailboxAddress,
		Contents:           contents,
		SentOn:             time.Now().UTC(),
	}
}

func (self Mail) IsDelivered() bool {
	return !self.DeliveredOn.IsZero()
}

func (self Mail) IsOpened() bool {
	return !self.OpenedOn.IsZero()
}
