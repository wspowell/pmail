package db

import (
	"github.com/wspowell/context"

	"github.com/wspowell/snailmail/resources/models"
)

type Datastore interface {
	// ExchangeMail for a user at a given location.
	// If there is an exchange at the location:
	//   * Any new mail the user has created is dropped off
	//   * Any mail the user is carrying is dropped off
	//   * Any mail in the exchange has a chance to be picked up
	//     * The same user may not pick up mail they dropped off within the last 24 hours
	// Errors:
	//   * ErrInternalFailure
	ExchangeMail(ctx context.Context, userGuid string, location models.Coordinate) error

	// CreateOneTimePassword for a websocket connection with the given value.
	// Errors:
	//   * ErrOneTimePasswordExists
	//   * ErrInternalFailure
	CreateOneTimePassword(ctx context.Context, connectionId string, oneTimePassword string) error

	// GetOneTimePassword with the given value.
	// Errors:
	//   * ErrOneTimePasswordNotFound
	//   * ErrInternalFailure
	GetOneTimePassword(ctx context.Context, oneTimePassword string) (string, error)

	// DeleteOneTimePassword for a websocket connection.
	// Errors:
	//   * ErrInternalFailure
	DeleteOneTimePassword(ctx context.Context, connectionId string) error

	// CreateUser with given attributes.
	// Errors:
	//   * ErrUserGuidExists
	//   * ErrUsernameExists
	//   * ErrInternalFailure
	CreateUser(ctx context.Context, newUser models.User) error

	// GetUser using user GUID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	GetUser(ctx context.Context, userGuid string) (*models.User, error)

	// DeleteUser using user GUID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	DeleteUser(ctx context.Context, userGuid string) error

	// UpdateUser with new user.
	// Errors:
	//   * ErrUsernameExists
	//   * ErrInternalFailure
	UpdateUser(ctx context.Context, updatedUser models.User) error

	// CreateMail to send.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	CreateMail(ctx context.Context, newMail models.Mail) error

	// OpenMail for viewing.
	// Errors:
	//   * ErrMailNotDelivered
	//   * ErrMailNotFound
	//   * ErrInternalFailure
	GetMail(ctx context.Context, mailGuid string) (*models.Mail, error)

	// GetUserMail for viewing all mail delivered to user and picked up from mailbox.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrInternalFailure
	GetUserMail(ctx context.Context, userGuid string) ([]models.Mail, error)

	// DeleteMail permanently.
	// Errors:
	//   * ErrInternalFailure
	DeleteMail(ctx context.Context, mailGuid string) error

	// CreateMailbox and place into the world.
	// Errors:
	//   * ErrMailboxAddressExists
	//   * ErrUserMailboxExists
	//   * ErrMailboxLabelExists
	//   * ErrInternalFailure
	CreateMailbox(ctx context.Context, userGuid string, newMailbox models.Mailbox) error

	// GetMailbox using the mailbox GUID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrMailboxNotFound
	//   * ErrInternalFailureGetUser
	GetMailbox(ctx context.Context, mailboxAddress string) (*models.Mailbox, error)

	// DeleteMailbox using the mailbox GUID.
	// Errors:
	//   * ErrInternalFailure
	DeleteMailbox(ctx context.Context, mailboxAddress string) error

	// GetUserMailbox using the user GUID.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetUserMailbox(ctx context.Context, userGuid string) (*models.Mailbox, error)

	// GetNearbyMailboxes using coordinates.
	// Errors:
	//   * ErrInternalFailure
	GetNearbyMailboxes(ctx context.Context, location models.Coordinate, radiusMeters float32) ([]models.Mailbox, error)

	// GetMailboxMail stored in the mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	GetMailboxMail(ctx context.Context, mailboxAddress string) ([]models.MailPreview, error)

	// DropOffMail in a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	DropOffMail(ctx context.Context, carrierGuid string, mailboxAddress string) ([]string, error)

	// PickUpMail from a Mailbox.
	// Errors:
	//   * ErrMailboxNotFound
	//   * ErrInternalFailure
	PickUpMail(ctx context.Context, carrierGuid string, mailboxAddress string) ([]string, error)

	// OpenMail for the first time.
	// Errors:
	//   * ErrMailNotFound
	//   * ErrInternalFailure
	OpenMail(ctx context.Context, mailGuid string) (*models.Mail, error)
}
