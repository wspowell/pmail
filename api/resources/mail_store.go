package resources

import (
	"context"
	"time"
)

type MailId uint32

type Mail struct {
	Id MailId

	// Metadata
	SentOn      time.Time
	DeliveredOn time.Time
	OpenedOn    time.Time
}

type MailAttributes struct {
	From     UserId
	To       UserId
	Contents string
}

type MailStore interface {
	CreateMail(ctx context.Context, attributes MailAttributes) (MailId, error)

	// OpenMail for viewing.
	// Errors:
	//   * ErrMailNotDelivered
	//   * ErrUserNotRecipient
	OpenMail(ctx context.Context, mailId MailId) (Mail, error)
}
