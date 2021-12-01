package db

import "github.com/wspowell/errors"

var (
	ErrInternalFailure = errors.New("db-0", "internal users failure")
)

// Users
var (
	ErrUserGuidExists = errors.New("db-user-1", "user guid already exists")
	ErrUsernameExists = errors.New("db-user-2", "username already exists")
	ErrUserNotFound   = errors.New("db-user-3", "user not found")
)

// Mail
var (
	ErrMailGuidExists   = errors.New("db-mail-1", "mail guid already exists")
	ErrInvalidRecipient = errors.New("db-mail-2", "invalid recipient")
	ErrEmptyMail        = errors.New("db-mail-3", "mail contents empty")
	ErrMailNotFound     = errors.New("db-mail-4", "mail not found")
)

// Mailboxes
var (
	ErrMailboxAddressExists = errors.New("db-mailbox-1", "mailbox guid already exists")
	ErrUserMailboxExists    = errors.New("db-mailbox-2", "user mailbox already exists")
	ErrMailboxFull          = errors.New("db-mailbox-3", "mailbox full")
	ErrMailboxLabelExists   = errors.New("db-mailbox-4", "mailbox label already exists")
	ErrMailboxNotFound      = errors.New("db-mailbox-5", "mailbox not found")
)

const (
	icAuthUserUserNotFound = "inmemory-auth-1"
)

// In Memory User Database.
const (
	icCreateUserGuidConflict     = "inmemory-users-1"
	icCreateUserUsernameConflict = "inmemory-users-2"
	icUpdateUserUserNotFound     = "inmemory-users-3"
	icGetUserUserNotFound        = "inmemory-users-4"
)

// In Memory Mail Database.
const (
	icCreateMailGuidConflict = "inmemory-mail-1"
	icGetMailGuidNotFound    = "inmemory-mail-2"
	icOpenMailGuidNotFound   = "inmemory-mail-3"
)

// In Memory Mailbox Database.
const (
	icCreateMailboxAddressConflict      = "inmemory-mailbox-1"
	icGetMailboxAddressNotFound         = "inmemory-mailbox-2"
	icGetUserMailboxUserMailboxNotFound = "inmemory-mailbox-3"
	icGetUserMailboxAddressNotFound     = "inmemory-mailbox-4"
	icGetMailboxMailMailboxNotFound     = "inmemory-mailbox-5"
	icDropOffMailMailboxNotFound        = "inmemory-mailbox-6"
	icCreateMailboxUserMailboxConflict  = "inmemory-mailbox-7"
	icPickUpMailUserNotFound            = "inmemory-mailbox-8"
	icPickUpMailMailboxNotFound         = "inmemory-mailbox-9"
	icDropOffMailUserNotFound           = "inmemory-mailbox-10"
)

// MySql Database.
const (
	icConnectFailure = "mysql-connect-1"
)

const (
	icMigrateInitError        = "mysql-migrate-1"
	icMigrateUpError          = "mysql-migrate-2"
	icMigrateNewInstanceError = "mysql-migrate-3"
)

const (
	icCreateUserRandSeedError = "mysql-user-1"
)
