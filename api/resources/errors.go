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
	icHomeMailboxExists = "resources-mailboxstore-2"
	icMailboxNotFound   = "resources-mailboxstore-3"
)

var (
	ErrHomeMailboxExists = errors.New(icHomeMailboxExists, "user already has home mailbox")
	ErrorMailboxNotFound = errors.New(icMailboxNotFound, "mailbox not found")
)
