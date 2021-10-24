package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type userModel struct {
	Username         string `json:"username"`
	PineappleOnPizza bool   `json:"pineapple_on_pizza"`
}

type createUserRequest struct {
	userModel
}

type createUserResponse struct {
	UserGuid string `json:"user_guid"`
}

type createUser struct {
	Datastore    db.Datastore        `spiderweb:"resource=datastore"`
	RequestBody  *createUserRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createUserResponse `spiderweb:"response,mime=application/json"`
}

func (self *createUser) Handle(ctx context.Context) (int, error) {
	userAttributes := user.Attributes{
		PineappleOnPizza:  self.RequestBody.PineappleOnPizza,
		Username:          self.RequestBody.Username,
		MailCarryCapacity: user.DefaultCarryCapacity,
	}

	newUser := user.NewUser(userAttributes)

	if err := self.Datastore.CreateUser(ctx, newUser); err != nil {
		if errors.Is(err, db.ErrUserGuidExists) {
			return http.StatusConflict, errors.Propagate(icCreateUserUserGuidConflict, err)
		} else if errors.Is(err, db.ErrUsernameExists) {
			return http.StatusConflict, errors.Propagate(icCreateUserUsernameConflict, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icCreateUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icCreateUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	log.Debug(ctx, "created user: %+v", newUser)

	self.ResponseBody.UserGuid = string(newUser.UserGuid)

	return http.StatusCreated, nil
}
