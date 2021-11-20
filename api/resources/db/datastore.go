package db

import (
	"time"

	"github.com/wspowell/context"

	"github.com/wspowell/snailmail/resources/models/geo"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

type Datastore interface {
	// CreateUser with given attributes.
	// Errors:
	//   * ErrUserGuidExists
	//   * ErrUsernameExists
	//   * ErrInternalFailure
	CreateUser(ctx context.Context, newUser user.User, password string) error

	// GetUser using user GUID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	GetUser(ctx context.Context, userGuid user.Guid) (*user.User, error)

	// AuthUser using user GUID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	AuthUser(ctx context.Context, username string, password string) (*user.User, error)

	// DeleteUser using user GUID.
	// Errors:
	//   * ErrInternalFailure
	DeleteUser(ctx context.Context, userGuid user.Guid) error

	// UpdateUser with new user.
	// Errors:
	//   * ErrUsernameExists
	//   * ErrInternalFailure
	UpdateUser(ctx context.Context, updatedUser user.User) error

	// CreateMail to send.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	CreateMail(ctx context.Context, newMail mail.Mail) error

	// OpenMail for viewing.
	// Errors:
	//   * ErrMailNotDelivered
	//   * ErrMailNotFound
	//   * ErrInternalFailure
	GetMail(ctx context.Context, mailGuid mail.Guid) (*mail.Mail, error)

	// GetUserMail for viewing all mail delivered to user and picked up from mailbox.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	GetUserMail(ctx context.Context, userGuid user.Guid) ([]mail.Mail, error)

	// DeleteMail permanently.
	// Errors:
	//   * ErrInternalFailure
	DeleteMail(ctx context.Context, mailGuid mail.Guid) error

	// CreateMailbox and place into the world.
	// Errors:
	//   * ErrMailboxGuidExists
	//   * ErrUserMailboxExists
	//   * ErrMailboxLabelExists
	//   * ErrInternalFailure
	CreateMailbox(ctx context.Context, newMailbox mailbox.Mailbox) error

	// GetMailbox using the mailbox GUID.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetMailbox(ctx context.Context, mailboxGuid mailbox.Guid) (*mailbox.Mailbox, error)

	// GetMailbox using the mailbox label.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetMailboxByLabel(ctx context.Context, mailboxLabel string) (*mailbox.Mailbox, error)

	// DeleteMailbox using the mailbox GUID.
	// Errors:
	//   * ErrInternalFailure
	DeleteMailbox(ctx context.Context, mailboxGuid mailbox.Guid) error

	// GetUserMailbox using the user GUID.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetUserMailbox(ctx context.Context, userGuid user.Guid) (*mailbox.Mailbox, error)

	// GetNearbyMailboxes using coordinates.
	// Errors:
	//   * ErrInternalFailure
	GetNearbyMailboxes(ctx context.Context, location geo.Coordinate, radiusMeters uint32) ([]mailbox.Mailbox, error)

	// GetMailboxMail stored in the mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetMailboxMail(ctx context.Context, mailboxGuid mailbox.Guid) ([]mail.Mail, error)

	// DropOffMail in a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	DropOffMail(ctx context.Context, carrierGuid user.Guid, mailboxGuid mailbox.Guid) ([]mail.Guid, error)

	// PickUpMail from a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	PickUpMail(ctx context.Context, carrierGuid user.Guid, mailboxGuid mailbox.Guid) ([]mail.Guid, error)

	// OpenMail for the first time.
	// Errors:
	//   * ErrMailNotFound
	//   * ErrInternalFailure
	OpenMail(ctx context.Context, mailGuid mail.Guid, openedAt time.Time) error
}
