package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/db"
)

type updateUserRequest struct {
	mime.Json

	PineappleOnPizza bool `json:"pineappleOnPizza"`
}

type updateUser struct {
	pathParams
	Datastore db.Datastore
	body.Request[updateUserRequest]
}

func (self *updateUser) Handle(ctx context.Context) (int, error) {
	// updateUser, err := self.Datastore.GetUser(ctx, user.Guid(self.UserGuid))
	// if err != nil {
	// 	if errors.Is(err, db.ErrUserNotFound) {
	// 		return http.StatusNotFound, err
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return http.StatusInternalServerError, err
	// 	} else {
	// 		return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// }

	// updateUser.Attributes = user.Attributes{
	// 	Username:          updateUser.Username,
	// 	PineappleOnPizza:  self.RequestBody.PineappleOnPizza,
	// 	MailCarryCapacity: updateUser.MailCarryCapacity,
	// }

	// if err := self.Datastore.UpdateUser(ctx, *updateUser); err != nil {
	// 	if errors.Is(err, db.ErrUserNotFound) {
	// 		return http.StatusNotFound, err
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return http.StatusInternalServerError, err
	// 	} else {
	// 		return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// }

	// log.Debug(ctx, "updated user: %+v", updateUser)

	return http.StatusOK, nil
}
