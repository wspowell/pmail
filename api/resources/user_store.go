package resources

type UserAttributes struct {
	Username string
}

type UserStore interface {
	CreateUser(attributes UserAttributes) (uint, error)
	DeleteUser(userId uint) error
	UpdateUser(newAttributes UserAttributes) error
}
