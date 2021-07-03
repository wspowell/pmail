package resources

import (
	"context"
)

type CarrierStore interface {
	// DropoffMail in a Mailbox.
	DropoffMail(ctx context.Context, userId UserId, mailboxId MailboxId, mailIds []MailId) error
	// PickupMail from a Mailbox.
	// Returns the Mail collected.
	PickupMail(ctx context.Context, userId UserId, mailboxId MailboxId) ([]MailId, error)
}
