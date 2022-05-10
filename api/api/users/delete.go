package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/resources/db"
)

type deleteUser struct {
	pathParams
	Datastore db.Datastore
}

func (self *deleteUser) Handle(ctx context.Context) (int, error) {
	if err := self.Datastore.DeleteUser(ctx, self.UserGuid); err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			// Success, but do not log.
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	} else {
		log.Debug(ctx, "deleted user: %s", self.UserGuid)
	}

	return httpstatus.NoContent, nil
}
