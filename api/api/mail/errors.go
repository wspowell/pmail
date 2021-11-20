package mail

import "github.com/wspowell/errors"

var (
	errUncaughtDbError  = errors.New("snailmail-mail-1", "uncaught database error")
	errMailNotFound     = errors.New("snailmail-mail-2", "mail not found")
	errInvalidRecipient = errors.New("snailmail-mail-3", "invalid recipient")
	errEmptyContents    = errors.New("snailmail-mail-4", "mail contents empty")
)

const (
	icCreateMailCreateMailDbError                = "snailmail-mail-1"
	icCreateMailCreateMailUnknownDbError         = "snailmail-mail-1"
	icCreateMailEmptyContents                    = "snailmail-mail-1"
	icCreateMailGetMailboxByLabelMailboxNotFound = "snailmail-mail-1"
	icCreateMailGetMailboxByLabelDbError         = "snailmail-mail-1"
	icCreateMailGetMailboxByLabelUnknownDbError  = "snailmail-mail-1"
	icCreateMailInvalidRecipient                 = "snailmail-mail-1"
	icCreateMailCreateMailboxMailboxNotFound     = "snailmail-mail-1"
)

const (
	icGetMailGetMailMailNotFound    = "snailmail-mail-1"
	icGetMailGetMailDbError         = "snailmail-mail-1"
	icGetMailGetMailUnknownDbError  = "snailmail-mail-1"
	icGetMailUserNotRecipient       = "snailmail-mail-1"
	icGetMailOpenMailMailNotFound   = "snailmail-mail-1"
	icGetMailOpenMailDbError        = "snailmail-mail-1"
	icGetMailOpenMailUnknownDbError = "snailmail-mail-1"
)

const (
	icListMailGetUserMailUserNotFound   = "snailmail-mail-1"
	icListMailGetUserMailDbError        = "snailmail-mail-1"
	icListMailGetUserMailUnknownDbError = "snailmail-mail-1"
)
