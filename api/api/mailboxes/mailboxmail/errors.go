package mailboxmail

import "github.com/wspowell/errors"

var (
	errUncaughtDbError = errors.New("snailmail-mailbox-1", "uncaught database error")

	ErrInsufficientMailboxPermissions = errors.New("snailmail-mailboxmail-1", "insufficient mailbox permissions")
)

const (
	icExchangeMailGetMailboxNotFound        = "snailmail-mailbox-mail-1"
	icExchangeMailGetMailboxDbError         = "snailmail-mailbox-mail-1"
	icExchangeMailGetMailboxUnknownDbError  = "snailmail-mailbox-mail-1"
	icExchangeMailDropOffMailNotFound       = "snailmail-mailbox-mail-1"
	icExchangeMailDropOffMailDbError        = "snailmail-mailbox-mail-1"
	icExchangeMailDropOffMailUnknownDbError = "snailmail-mailbox-mail-1"
	icExchangeMailPickUpMailNotFound        = "snailmail-mailbox-mail-1"
	icExchangeMailPickUpMailDbError         = "snailmail-mailbox-mail-1"
	icExchangeMailPickUpMailUnknownDbError  = "snailmail-mailbox-mail-1"
)
