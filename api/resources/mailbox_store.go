package resources

import (
	"context"

	"github.com/wspowell/errors"
)

// Mailbox store errors
var (
	ErrMailboxStoreFailure            = errors.New("resources-mailboxstore-1", "internal mailbox store failure")
	ErrHomeMailboxExists              = errors.New("resources-mailboxstore-2", "user already has home mailbox")
	ErrMailboxLabelExists             = errors.New("resources-mailboxstore-3", "mailbox label already exists")
	ErrMailboxNotFound                = errors.New("resources-mailboxstore-5", "mailbox not found")
	ErrInsufficientMailboxPermissions = errors.New("resources-mailboxstore-9", "invalid mailbox permissions")
)

type Latitude float32
type Longitude float32

type GeoCoordinate struct {
	Lat Latitude
	Lng Longitude
}

// FIXME: Ideally userId should be a uint64.
type MailboxId uint32

type Mailbox struct {
	Id MailboxId
	MailboxAttributes
}

type MailboxAttributes struct {
	// Label of the mailbox to show the user.
	// Must be globally unique.
	Label string

	// Owner of the mailbox
	// 0 - No owner
	Owner UserId

	// Location in the world.
	Location GeoCoordinate

	// Capacity of mail in the mailbox.
	// 0 - Pickup only
	Capacity uint32
}

type MailboxStore interface {
	// CreateMailbox and place into the world.
	// Errors:
	//   * ErrHomeMailboxExists
	//   * ErrMailboxLabelExists
	//   * ErrMailboxStoreFailure
	CreateMailbox(ctx context.Context, attributes MailboxAttributes) (MailboxId, error)

	// GetMailbox using the mailbox ID.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrMailboxStoreFailure
	GetMailbox(ctx context.Context, mailboxId MailboxId) (Mailbox, error)

	// GetUserMailbox using the user ID.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrMailboxStoreFailure
	GetUserMailbox(ctx context.Context, userId UserId) (Mailbox, error)

	// GetNearbyMailboxes using coordinates.
	// Errors:
	//   * ErrMailboxStoreFailure
	GetNearbyMailboxes(ctx context.Context, location GeoCoordinate, radiusMeters uint32) ([]Mailbox, error)

	// DropOffMail in a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInsufficientMailboxPermissions
	//   * ErrMailboxStoreFailure
	DropOffMail(ctx context.Context, userId UserId, mailboxId MailboxId, mailIds []MailId) error

	// PickUpMail from a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInsufficientMailboxPermissions
	//   * ErrMailboxStoreFailure
	PickUpMail(ctx context.Context, userId UserId, mailboxId MailboxId) ([]MailId, error)
}
