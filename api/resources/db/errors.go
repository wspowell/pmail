package db

import "github.com/wspowell/errors"

var (
	ErrInternalFailure = errors.New("internal users failure")
)

var (
	ErrOneTimePasswordNotFound = errors.New("one time password not found")
)

// Users
var (
	ErrAddressExists = errors.New("address already exists")
	ErrUserNotFound  = errors.New("user not found")
)

// Mail
var (
	ErrMailGuidExists   = errors.New("mail guid already exists")
	ErrInvalidRecipient = errors.New("invalid recipient")
	ErrEmptyMail        = errors.New("mail contents empty")
	ErrMailNotFound     = errors.New("mail not found")
)

// Mailboxes
var (
	ErrMailboxAddressExists = errors.New("mailbox guid already exists")
	ErrUserMailboxExists    = errors.New("user mailbox already exists")
	ErrMailboxFull          = errors.New("mailbox full")
	ErrMailboxLabelExists   = errors.New("mailbox label already exists")
	ErrMailboxNotFound      = errors.New("mailbox not found")
)
