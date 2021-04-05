package resources

import "github.com/wspowell/errors"

const (
	icUserNotFound     = "resources-userstore-1"
	icUsernameConflict = "resources-userstore-2"
)

var (
	ErrUserNotFound     = errors.New(icUserNotFound, "user id not found")
	ErrUsernameConflict = errors.New(icUsernameConflict, "username already exists")
)

const (
	icMailboxOwned      = "resources-mailboxstore-1"
	icHomeMailboxExists = "resources-mailboxstore-2"
)

var (
	ErrMailboxOwned      = errors.New(icMailboxOwned, "mailbox already owned")
	ErrHomeMailboxExists = errors.New(icHomeMailboxExists, "user already has home mailbox")
)
