package resources

import "context"

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
	Label string
	// Owner of the mailbox
	// 0 - No owner
	Owner    UserId
	Location GeoCoordinate
	// Capacity of mail in the mailbox.
	// 0 - Pickup only
	Capacity uint32
}

type MailboxStore interface {
	CreateMailbox(ctx context.Context, attributes MailboxAttributes) (MailboxId, error)
	GetMailbox(ctx context.Context, mailboxId MailboxId) (Mailbox, error)
	GetUserMailbox(ctx context.Context, userId UserId) (Mailbox, error)
	GetNearbyMailboxes(ctx context.Context, location GeoCoordinate, radius float32) ([]Mailbox, error)
}
