package users

import (
	"net/http"
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type userModel struct {
	Username         string `json:"username"`
	PineappleOnPizza bool   `json:"pineappleOnPizza"`
}

type createUserRequest struct {
	userModel
	Password string `json:"password"`
}

type createUserResponse struct {
	UserGuid string `json:"userGuid"`
}

type createUser struct {
	Datastore    db.Datastore        `spiderweb:"resource=datastore"`
	RequestBody  *createUserRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *createUserResponse `spiderweb:"response,mime=application/json"`
}

func (self *createUser) Handle(ctx context.Context) (int, error) {
	// FIXME: Move to validation.
	if self.RequestBody.Username == "" {
		return httpstatus.UnprocessableEntity, errors.Propagate(icCreateUserUsernameBlank, errInvalidUsername)
	}
	if self.RequestBody.Password == "" {
		return httpstatus.UnprocessableEntity, errors.Propagate(icCreateUserPasswordBlank, errInvalidPassword)
	}

	userAttributes := user.Attributes{
		PineappleOnPizza:  self.RequestBody.PineappleOnPizza,
		Username:          self.RequestBody.Username,
		MailCarryCapacity: user.DefaultCarryCapacity,
		CreatedOn:         time.Now().UTC(),
	}

	newUser := user.NewUser(userAttributes)

	password, err := auth.Password(ctx, self.RequestBody.Password)
	if err != nil {
		return http.StatusInternalServerError, errors.Propagate(icCreateUserPasswordError, err)
	}

	if err := self.Datastore.CreateUser(ctx, newUser, password); err != nil {
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
