package users

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type userModel struct {
	Username         string `json:"username"`
	PineappleOnPizza bool   `json:"pineapple_on_pizza"`
	MailboxId        uint32 `json:"mailbox_id,omitempty"`
}

type createUserRequest struct {
	userModel
}

type createUserResponse struct {
	UserId uint32 `json:"user_id"`
}

type createUser struct {
	Users        resources.UserStore `spiderweb:"resource=userstore"`
	RequestBody  *createUserRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createUserResponse `spiderweb:"response,mime=application/json"`
}

func (self *createUser) Handle(ctx *endpoint.Context) (int, error) {
	userAttributes := resources.UserAttributes{
		PineappleOnPizza: self.RequestBody.PineappleOnPizza,
	}

	userId, err := self.Users.CreateUser(self.RequestBody.Username, userAttributes)
	if err != nil {
		if errors.Is(err, resources.ErrUsernameConflict) {
			return http.StatusConflict, errors.Wrap(icCreateUserUsernameConflict, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icCreateUserError, err)
	}

	ctx.Debug("created user: %d", userId)

	self.ResponseBody.UserId = userId

	return http.StatusCreated, nil
}
