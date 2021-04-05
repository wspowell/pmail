package db

import (
	"math/rand"

	"github.com/wspowell/pmail/resources"

	"github.com/wspowell/errors"
)

var _ resources.UserStore = (*Users)(nil)

type Users struct {
	userIdToUser     map[uint32]resources.User
	usernameToUserId map[string]uint32
}

func NewUsers() *Users {
	return &Users{
		userIdToUser:     map[uint32]resources.User{},
		usernameToUserId: map[string]uint32{},
	}
}

func (self *Users) userIdExists(userId uint32) bool {
	_, exists := self.userIdToUser[userId]
	return exists
}

func (self *Users) usernameExists(username string) bool {
	_, exists := self.usernameToUserId[username]
	return exists
}

func (self *Users) CreateUser(username string, attributes resources.UserAttributes) (uint32, error) {
	if self.usernameExists(username) {
		return 0, errors.Wrap(icCreateUserUsernameConflict, resources.ErrUsernameConflict)
	}

	userId := rand.Uint32()

	self.userIdToUser[userId] = resources.User{
		Username:   username,
		Attributes: attributes,
	}
	self.usernameToUserId[username] = userId

	return userId, nil
}

func (self *Users) GetUser(userId uint32) (*resources.User, error) {
	if user, exists := self.userIdToUser[userId]; exists {
		return &user, nil
	}

	return nil, errors.Wrap(icGetUserUserNotFound, resources.ErrUserNotFound)
}

func (self *Users) DeleteUser(userId uint32) error {
	if user, exists := self.userIdToUser[userId]; exists {
		delete(self.userIdToUser, userId)
		delete(self.usernameToUserId, user.Username)
	}
	return nil
}

func (self *Users) UpdateUser(userId uint32, newAttributes resources.UserAttributes) error {
	if user, exists := self.userIdToUser[userId]; exists {
		user.Attributes = newAttributes

		self.userIdToUser[userId] = user
		return nil
	}

	return errors.New(icUpdateUserUserNotFound, "user id not found: %v", userId)
}
