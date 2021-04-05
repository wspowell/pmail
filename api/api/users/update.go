package users

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type updateUserRequest struct {
	PineappleOnPizza bool `json:"pineapple_on_pizza"`
}

type updateUser struct {
	UserId      uint32              `spiderweb:"path=id"`
	Users       resources.UserStore `spiderweb:"resource=userstore"`
	RequestBody *updateUserRequest  `spiderweb:"request,mime=application/json"`
}

func (self *updateUser) Handle(ctx *endpoint.Context) (int, error) {
	userAttributes := resources.UserAttributes{
		PineappleOnPizza: self.RequestBody.PineappleOnPizza,
	}

	if err := self.Users.UpdateUser(self.UserId, userAttributes); err != nil {
		if errors.Is(err, resources.ErrUserNotFound) {
			return http.StatusNotFound, errors.Wrap(icUpdateUserUserNotFound, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icUpdateUserError, err)
	}

	ctx.Debug("updated user: %d", self.UserId)

	return http.StatusOK, nil
}
