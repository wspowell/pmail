package authorize

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError = errors.New("snailmail-auth-1", "uncaught database error")
	errJwtError        = errors.New("snailmail-auth-1", "failure signing jwt")
)

const (
	icAuthUserUserNotFound   = "snailmail-auth-1"
	icAuthUserDbError        = "snailmail-auth-2"
	icAuthUserUnknownDbError = "snailmail-auth-3"
	icAuthUserJwtError       = "snailmail-auth-4"
)
