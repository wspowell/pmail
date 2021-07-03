package db

import (
	"context"
	"math/rand"

	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
)

var _ resources.UserStore = (*Users)(nil)

type Users struct {
	userIdToUser     map[uint32]*resources.User
	usernameToUserId map[string]uint32
}

func NewUsers() *Users {
	return &Users{
		userIdToUser:     map[uint32]*resources.User{},
		usernameToUserId: map[string]uint32{},
	}
}

func (self *Users) userIdExists(ctx context.Context, userId uint32) bool {
	_, exists := self.userIdToUser[userId]
	return exists
}

func (self *Users) usernameExists(ctx context.Context, username string) bool {
	_, exists := self.usernameToUserId[username]
	return exists
}

func (self *Users) CreateUser(ctx context.Context, username string, attributes resources.UserAttributes) (*resources.User, error) {
	if self.usernameExists(ctx, username) {
		return nil, errors.Wrap(icCreateUserUsernameConflict, resources.ErrUsernameConflict)
	}

	user := &resources.User{
		Id:         rand.Uint32(),
		Username:   username,
		Attributes: attributes,
	}

	self.userIdToUser[user.Id] = user
	self.usernameToUserId[username] = user.Id

	return user, nil
}

func (self *Users) GetUser(ctx context.Context, userId uint32) (*resources.User, error) {
	if user, exists := self.userIdToUser[userId]; exists {
		return user, nil
	}

	return nil, errors.Wrap(icGetUserUserNotFound, resources.ErrUserNotFound)
}

func (self *Users) DeleteUser(ctx context.Context, userId uint32) error {
	if user, exists := self.userIdToUser[userId]; exists {
		delete(self.userIdToUser, userId)
		delete(self.usernameToUserId, user.Username)
	}
	return nil
}

func (self *Users) UpdateUser(ctx context.Context, userId uint32, newAttributes resources.UserAttributes) error {
	if user, exists := self.userIdToUser[userId]; exists {
		user.Attributes = newAttributes

		self.userIdToUser[userId] = user
		return nil
	}

	return errors.Wrap(icUpdateUserUserNotFound, resources.ErrUserNotFound)
}
