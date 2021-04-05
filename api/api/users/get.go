package users

import (
	"net/http"

	"github.com/wspowell/errors"
	"github.com/wspowell/pmail/resources"
	"github.com/wspowell/spiderweb/endpoint"
)

type getUserResponse struct {
	userModel
}

type User struct {
	Username   string
	Attributes UserAttributes
}

type UserAttributes struct {
	PineappleOnPizza bool
}

type getUser struct {
	UserId       uint32              `spiderweb:"path=id"`
	Users        resources.UserStore `spiderweb:"resource=userstore"`
	ResponseBody *getUserResponse    `spiderweb:"response,mime=application/json"`
}

func (self *getUser) Handle(ctx *endpoint.Context) (int, error) {
	user, err := self.Users.GetUser(self.UserId)
	if err != nil {
		if errors.Is(err, resources.ErrUserNotFound) {
			return http.StatusNotFound, errors.Wrap(icGetUserUserNotFound, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icGetUserGetUserError, err)
	}

	ctx.Debug("deleted user: %d", self.UserId)

	self.ResponseBody.userModel.Username = user.Username
	self.ResponseBody.userModel.PineappleOnPizza = user.Attributes.PineappleOnPizza

	return http.StatusOK, nil
}
