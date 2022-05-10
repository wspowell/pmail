package mailboxmail

import "github.com/wspowell/errors"

var (
	errUncaughtDbError = errors.New("uncaught database error")

	errInsufficientMailboxPermissions = errors.New("insufficient mailbox permissions")
)
