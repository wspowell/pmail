package users

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/pmail/resources/db"
	"github.com/wspowell/spiderweb/endpoint"
)

type createUserRequest struct {
	Username string `json:"username"`
}

type createUserResponse struct {
	UserId uint `json:"user_id"`
}

type createUser struct {
	Users        *db.Users           `spiderweb:"resource=userstore"`
	RequestBody  *createUserRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createUserResponse `spiderweb:"response,mime=application/json"`
}

func (self *createUser) Handle(ctx *endpoint.Context) (int, error) {
	userAttributes := resources.UserAttributes{
		Username: self.RequestBody.Username,
	}

	userId, err := self.Users.CreateUser(userAttributes)
	if err != nil {
		ctx.Error(icCreateUserError, "failed to create user: %#v", err)
		return http.StatusInternalServerError, errors.Wrap(icCreateUserError, err)
	}

	ctx.Debug("created user: %d", userId)

	self.ResponseBody.UserId = userId

	return http.StatusCreated, nil
}
