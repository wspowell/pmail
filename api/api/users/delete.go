package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/snailmail/resources"
)

type deleteUser struct {
	UserId uint32              `spiderweb:"path=id"`
	Users  resources.UserStore `spiderweb:"resource=userstore"`
}

func (self *deleteUser) Handle(ctx context.Context) (int, error) {
	if err := self.Users.DeleteUser(ctx, self.UserId); err != nil {
		log.Error(ctx, icDeleteUserError, "failed to delete user: %#v", err)
		return http.StatusInternalServerError, errors.Wrap(icDeleteUserError, err)
	}

	log.Debug(ctx, "deleted user: %d", self.UserId)

	return http.StatusOK, nil
}
