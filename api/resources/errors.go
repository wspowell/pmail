package resources

import "github.com/wspowell/errors"

const (
	icHomeMailboxExists = "resources-mailboxstore-2"
	icMailboxNotFound   = "resources-mailboxstore-3"
)

// Mailbox errors
var (
	ErrHomeMailboxExists = errors.New(icHomeMailboxExists, "user already has home mailbox")
	ErrMailboxNotFound   = errors.New(icMailboxNotFound, "mailbox not found")
)

const (
	icMailNotDelivered = "resources-mailstore-1"
	icUserNotRecipient = "resources-mailstore-2"
)

// Mail errors
var (
	ErrMailNotDelivered = errors.New(icMailNotDelivered, "mail not delivered")
	ErrUserNotRecipient = errors.New(icUserNotRecipient, "user not recipient")
)
