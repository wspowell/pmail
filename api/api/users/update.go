package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/snailmail/resources"
)

type updateUserRequest struct {
	PineappleOnPizza bool `json:"pineapple_on_pizza"`
}

type updateUser struct {
	UserId      uint32              `spiderweb:"path=id"`
	Users       resources.UserStore `spiderweb:"resource=userstore"`
	RequestBody *updateUserRequest  `spiderweb:"request,mime=application/json"`
}

func (self *updateUser) Handle(ctx context.Context) (int, error) {
	userAttributes := resources.UserAttributes{
		PineappleOnPizza: self.RequestBody.PineappleOnPizza,
	}

	if err := self.Users.UpdateUser(self.UserId, userAttributes); err != nil {
		if errors.Is(err, resources.ErrUserNotFound) {
			return http.StatusNotFound, errors.Wrap(icUpdateUserUserNotFound, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icUpdateUserError, err)
	}

	log.Debug(ctx, "updated user: %d", self.UserId)

	return http.StatusOK, nil
}
