package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type updateUserRequest struct {
	PineappleOnPizza bool `json:"pineappleOnPizza"`
}

type updateUser struct {
	UserGuid    string             `spiderweb:"path=user_guid"`
	Datastore   db.Datastore       `spiderweb:"resource=datastore"`
	RequestBody *updateUserRequest `spiderweb:"request,mime=application/json"`
}

func (self *updateUser) Handle(ctx context.Context) (int, error) {
	updateUser, err := self.Datastore.GetUser(ctx, user.Guid(self.UserGuid))
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return http.StatusNotFound, errors.Propagate(icUpdateUserGetUserNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icUpdateUserGetUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icUpdateUserGetUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	updateUser.Attributes = user.Attributes{
		Username:          updateUser.Username,
		PineappleOnPizza:  self.RequestBody.PineappleOnPizza,
		MailCarryCapacity: updateUser.MailCarryCapacity,
	}

	if err := self.Datastore.UpdateUser(ctx, *updateUser); err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return http.StatusNotFound, errors.Propagate(icUpdateUserGetUserNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icUpdateUserUpdateUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icUpdateUserUpdateUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	log.Debug(ctx, "updated user: %+v", updateUser)

	return http.StatusOK, nil
}
