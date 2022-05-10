package mail

import "github.com/wspowell/errors"

var (
	errUncaughtDbError  = errors.New("uncaught database error")
	errMailNotFound     = errors.New("mail not found")
	errInvalidRecipient = errors.New("invalid recipient")
	errEmptyContents    = errors.New("mail contents empty")
)
