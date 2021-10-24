package user

import (
	"github.com/google/uuid"
)

const (
	DefaultCarryCapacity = 10
)

type Guid string

type User struct {
	UserGuid Guid
	Attributes
}

func NewUser(attributes Attributes) User {
	return User{
		UserGuid:   Guid(uuid.New().String()),
		Attributes: attributes,
	}
}

type Attributes struct {
	// Username to identity the user.
	// Must be globally unique.
	Username string

	// PineappleOnPizza is always true, duh.
	PineappleOnPizza bool

	// MailCarryCapacity limits the amount of mail a single user may carry.
	MailCarryCapacity uint32
}
