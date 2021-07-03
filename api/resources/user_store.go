package resources

import (
	"context"

	"github.com/wspowell/errors"
)

const (
	icCreateUserUsernameConflict = "resources-userstore-1"
	icCreateUserFailed           = "resources-userstore-2"
	icGetUserUserNotFound        = "resources-userstore-3"
	icUpdateUserUserNotFound     = "resources-userstore-4"
)

var (
	ErrCreateUserErrorUsernameConflict = errors.New(icCreateUserUsernameConflict, "username already exists")
	ErrCreateUserErrorCreateFailure    = errors.New(icCreateUserFailed, "failed to create")
)

var (
	ErrGetUserErrorUserNotFound = errors.New(icGetUserUserNotFound, "user id not found")
)

var (
	ErrUpdateUserErrorUserNotFound = errors.New(icUpdateUserUserNotFound, "user id not found")
)

// FIXME: Ideally userId should be a uint64.
type UserId uint32

type User struct {
	Id UserId
	UserAttributes
}

type UserAttributes struct {
	Username         string
	PineappleOnPizza bool
}

type UserStore interface {
	CreateUser(ctx context.Context, attributes UserAttributes) (User, error)
	GetUser(ctx context.Context, userId UserId) (User, error)
	DeleteUser(ctx context.Context, userId UserId) error
	UpdateUser(ctx context.Context, userId UserId, newAttributes UserAttributes) error
}
