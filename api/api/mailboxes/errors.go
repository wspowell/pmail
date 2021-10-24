package mailboxes

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError = errors.New("snailmail-mailbox-1", "uncaught database error")
)

const (
	icCreateMailboxUserNotFound          = "snailmail-mailbox-5"
	icCreateMailboxGetUserDbError        = "snailmail-mailbox-5"
	icCreateMailboxGetUserUnknownDbError = "snailmail-mailbox-5"
	icCreateMailboxGuidConflict          = "snailmail-mailbox-5"
	icCreateMailboxUserMailboxConflict   = "snailmail-mailbox-5"
	icCreateMailboxLabelConflict         = "snailmail-mailbox-5"
	icCreateMailboxDbError               = "snailmail-mailbox-5"
	icCreateMailboxUnknownDbError        = "snailmail-mailbox-5"
)

const (
	icGetMailboxNotFound            = "snailmail-mailbox-5"
	icGetMailboxDbError             = "snailmail-mailbox-5"
	icGetMailboxUnknownDbError      = "snailmail-mailbox-5"
	icGetMailboxMailMailboxNotFound = "snailmail-mailbox-5"
	icGetMailboxMailDbError         = "snailmail-mailbox-5"
	icGetMailboxMailUnknownDbError  = "snailmail-mailbox-5"
)
