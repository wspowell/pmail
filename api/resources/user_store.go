package resources

type User struct {
	Username   string
	Attributes UserAttributes
}

type UserAttributes struct {
	PineappleOnPizza bool
}

// FIXME: Ideally userId should be a uint64.
type UserStore interface {
	CreateUser(username string, attributes UserAttributes) (uint32, error)
	GetUser(userId uint32) (*User, error)
	DeleteUser(userId uint32) error
	UpdateUser(userId uint32, newAttributes UserAttributes) error
}
