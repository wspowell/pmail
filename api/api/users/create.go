package users

import (
	"net/http"

	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/resources/db"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
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
	Users        db.Datastore        `spiderweb:"resource=datastore"`
	RequestBody  *createUserRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createUserResponse `spiderweb:"response,mime=application/json"`
}

func (self *createUser) Handle(ctx context.Context) (int, error) {
	userAttributes := resources.UserAttributes{
		PineappleOnPizza: self.RequestBody.PineappleOnPizza,
	}

	user, err := self.Users.CreateUser(ctx, self.RequestBody.Username, userAttributes)
	if err != nil {
		switch err {
		case resources.CreateUserErrorUsernameConflict:
			return http.StatusConflict, errors.Wrap(icCreateUserUsernameConflict, err)
		case resources.CreateUserErrorCreateFailure:
			return http.StatusInternalServerError, errors.Wrap(icCreateUserError, err)
		}
	}

	log.Debug(ctx, "created user: %d", user.Id)

	self.ResponseBody.UserId = user.Id

	return http.StatusCreated, nil
}
