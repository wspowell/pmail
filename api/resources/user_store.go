package resources

import (
	"context"

	"github.com/wspowell/errors"
)

// User store errors.
var (
	ErrUserStoreFailure = errors.New("resources-userstore-1", "internal user store failure")
	ErrUsernameExists   = errors.New("resources-userstore-2", "username already exists")
	ErrUserNotFound     = errors.New("resources-userstore-3", "user not found")
)

// FIXME: Ideally userId should be a uint64.
type UserId uint32

type User struct {
	Id UserId
	UserAttributes
}

type UserAttributes struct {
	// Username to identity the user.
	// Must be globally unique.
	Username string

	// PineappleOnPizza is always true, duh.
	PineappleOnPizza bool
}

type UserStore interface {
	// CreateUser with given attributes.
	// Errors:
	//   * ErrUsernameExists
	//   * ErrUserStoreFailure
	CreateUser(ctx context.Context, attributes UserAttributes) (User, error)

	// GetUser using user ID.
	// Errors:
	//   * ErrUserNotFound
	//   * ErrUserStoreFailure
	GetUser(ctx context.Context, userId UserId) (User, error)

	// DeleteUser using user ID.
	// Errors:
	//   * ErrUserStoreFailure
	DeleteUser(ctx context.Context, userId UserId) error

	// DeleteUser using user ID and new user attributes.
	// Errors:
	//   * ErrUsernameExists
	//   * ErrUserStoreFailure
	UpdateUser(ctx context.Context, userId UserId, newAttributes UserAttributes) error
}
