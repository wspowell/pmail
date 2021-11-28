package authorize

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type jwtRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jwtResponse struct {
	JwtToken string `json:"jwtToken"`
}

type userTokener interface {
	UserToken(authUser user.User) (string, error)
}

type authUser struct {
	JwtAuth      userTokener  `spiderweb:"resource=jwt"`
	Datastore    db.Datastore `spiderweb:"resource=datastore"`
	RequestBody  *jwtRequest  `spiderweb:"request,mime=application/json"`
	ResponseBody *jwtResponse `spiderweb:"response,mime=application/json"`
}

func (self *authUser) Handle(ctx context.Context) (int, error) {
	authedUser, err := self.Datastore.AuthUser(ctx, self.RequestBody.Username, self.RequestBody.Password)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return httpstatus.NotFound, errors.Propagate(icAuthUserUserNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icAuthUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icAuthUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	jwtToken, err := self.JwtAuth.UserToken(*authedUser)
	if err != nil {
		return http.StatusInternalServerError, errors.Convert(icAuthUserJwtError, err, errJwtError)
	}

	log.Debug(ctx, "authed user: %+v", authedUser)

	self.ResponseBody.JwtToken = jwtToken

	return httpstatus.OK, nil
}
