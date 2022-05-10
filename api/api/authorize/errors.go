package authorize

import (
	"github.com/wspowell/errors"
)

var (
	errUncaughtDbError = errors.New("uncaught database error")
	errJwtError        = errors.New("failure signing jwt")
	errAwsError        = errors.New("failed contacting AWS")
	errNotFound        = errors.New("not found")
	errNoUserMailbox   = errors.New("user has no mailbox")
)
