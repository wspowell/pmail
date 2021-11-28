package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	DefaultCarryCapacity = 10
)

type Guid string

type User struct {
	Attributes
	UserGuid Guid
}

func NewUser(attributes Attributes) User {
	return User{
		UserGuid:   Guid(uuid.New().String()),
		Attributes: attributes,
	}
}

type Attributes struct {
	// MailCarryCapacity limits the amount of mail a single user may carry.
	MailCarryCapacity uint32

	// Username to identity the user.
	// Must be globally unique.
	Username string

	// CreatedOn timestamp.
	CreatedOn time.Time

	// PineappleOnPizza is always true, duh.
	PineappleOnPizza bool
}
