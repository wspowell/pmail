package users

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError = errors.New("snailmail-users-1", "uncaught database error")
	errInvalidUsername = errors.New("snailmail-users-2", "invalid username")
	errInvalidPassword = errors.New("snailmail-users-3", "invalid password")
)

const (
	icCreateUserUserGuidConflict = "snailmail-users-1"
	icCreateUserUsernameConflict = "snailmail-users-2"
	icCreateUserDbError          = "snailmail-users-3"
	icCreateUserUnknownDbError   = "snailmail-users-4"
	icCreateUserUsernameBlank    = "snailmail-users-5"
	icCreateUserPasswordBlank    = "snailmail-users-6"
	icCreateUserPasswordError    = "snailmail-users-7"
)

const (
	icGetUserUserNotFound          = "snailmail-users-8"
	icGetUserDbError               = "snailmail-users-9"
	icGetUserUnknownDbError        = "snailmail-users-10"
	icGetUserMailboxDbError        = "snailmail-users-11"
	icGetUserMailboxUnknownDbError = "snailmail-users-12"
)

const (
	icUpdateUserGetUserNotFound          = "snailmail-users-13"
	icUpdateUserGetUserDbError           = "snailmail-users-14"
	icUpdateUserGetUserUnknownDbError    = "snailmail-users-15"
	icUpdateUserUpdateUserDbError        = "snailmail-users-16"
	icUpdateUserUpdateUserUnknownDbError = "snailmail-users-17"
)

const (
	icDeleteUserDbError        = "snailmail-users-18"
	icDeleteUserUnknownDbError = "snailmail-users-19"
)
