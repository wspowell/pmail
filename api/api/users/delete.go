package users

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type deleteUser struct {
	UserId uint32              `spiderweb:"path=id"`
	Users  resources.UserStore `spiderweb:"resource=userstore"`
}

func (self *deleteUser) Handle(ctx *endpoint.Context) (int, error) {
	if err := self.Users.DeleteUser(self.UserId); err != nil {
		ctx.Error(icDeleteUserError, "failed to delete user: %#v", err)
		return http.StatusInternalServerError, errors.Wrap(icDeleteUserError, err)
	}

	ctx.Debug("deleted user: %d", self.UserId)

	return http.StatusOK, nil
}
