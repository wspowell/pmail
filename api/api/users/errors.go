package users

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError = errors.New("snailmail-users-1", "uncaught database error")
	errInvalidUsername = errors.New("snailmail-users-2", "invalid username")
	errInvalidPassword = errors.New("snailmail-users-2", "invalid password")
)

const (
	icCreateUserUserGuidConflict = "snailmail-users-2"
	icCreateUserUsernameConflict = "snailmail-users-3"
	icCreateUserDbError          = "snailmail-users-1"
	icCreateUserUnknownDbError   = "snailmail-users-4"
	icCreateUserUsernameBlank    = "snailmail-users-4"
	icCreateUserPasswordBlank    = "snailmail-users-4"
)

const (
	icGetUserUserNotFound          = "snailmail-users-4"
	icGetUserDbError               = "snailmail-users-6"
	icGetUserUnknownDbError        = "snailmail-users-6"
	icGetUserMailboxDbError        = "snailmail-users-6"
	icGetUserMailboxUnknownDbError = "snailmail-users-6"
)

const (
	icUpdateUserGetUserNotFound          = "snailmail-users-7"
	icUpdateUserGetUserDbError           = "snailmail-users-7"
	icUpdateUserGetUserUnknownDbError    = "snailmail-users-7"
	icUpdateUserUpdateUserDbError        = "snailmail-users-7"
	icUpdateUserUpdateUserUnknownDbError = "snailmail-users-7"
)

const (
	icDeleteUserDbError        = "snailmail-users-9"
	icDeleteUserUnknownDbError = "snailmail-users-9"
)
