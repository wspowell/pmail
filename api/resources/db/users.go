package db

import (
	"math/rand"

	"github.com/wspowell/pmail/resources"
)

type Users struct{}

func (self *Users) CreateUser(attributes resources.UserAttributes) (uint, error) {
	return uint(rand.Uint64()), nil
}

func (self *Users) DeleteUser(userId uint) error {
	panic("not implemented") // TODO: Implement
}

func (self *Users) UpdateUser(newAttributes resources.UserAttributes) error {
	panic("not implemented") // TODO: Implement
}
