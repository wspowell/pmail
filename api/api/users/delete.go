package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type deleteUser struct {
	UserGuid  string       `spiderweb:"path=user_guid"`
	Datastore db.Datastore `spiderweb:"resource=datastore"`
}

func (self *deleteUser) Handle(ctx context.Context) (int, error) {
	if err := self.Datastore.DeleteUser(ctx, user.Guid(self.UserGuid)); err != nil {
		if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icDeleteUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icDeleteUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	log.Debug(ctx, "deleted user: %s", self.UserGuid)

	return http.StatusOK, nil
}
