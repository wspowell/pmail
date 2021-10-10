package resources

import (
	"context"
	"time"

	"github.com/wspowell/errors"
)

// Mail store errors
var (
	ErrMailStoreFailure = errors.New("resources-mailstore-1", "internal mail store failure")
	ErrInvalidRecipient = errors.New("resources-mailstore-2", "recipient is not valid")
	ErrEmptyMail        = errors.New("resources-mailstore-3", "mail contents is empty")
	ErrMailNotDelivered = errors.New("resources-mailstore-4", "mail not delivered")
	ErrMailNotFound     = errors.New("resources-mailstore-5", "mail not found")
)

type MailId uint32

type Mail struct {
	Id MailId
	MailAttributes

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
	// CreateMail to send.
	// Errors:
	//   * ErrInvalidRecipient
	//   * ErrEmptyMail
	//   * ErrMailStoreFailure
	CreateMail(ctx context.Context, attributes MailAttributes) (MailId, error)

	// OpenMail for viewing.
	// Errors:
	//   * ErrMailNotDelivered
	//   * ErrMailNotFound
	//   * ErrMailStoreFailure
	OpenMail(ctx context.Context, userId UserId, mailId MailId) (Mail, error)
}
