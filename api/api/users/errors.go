package users

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError        = errors.New("uncaught database error")
	errMissingPublicKey       = errors.New("missing publicKey")
	errMissingMailboxLocation = errors.New("missing mailbox location")
)
