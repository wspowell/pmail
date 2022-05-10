package authorize

import (
	"net/http"
	"time"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models"
)

type jwtRequest struct {
	mime.Json

	UserGuid  string `json:"userGuid"`
	SignedKey string `json:"signedKey"`
}

type jwtResponse struct {
	mime.Json

	JwtToken string `json:"jwtToken"`
}

type UserTokener interface {
	UserToken(authUser *models.User, expiresAt time.Time) (string, error)
}

type authUser struct {
	JwtAuth   UserTokener
	Datastore db.Datastore
	body.Request[jwtRequest]
	body.Response[jwtResponse]
}

func (self *authUser) Handle(ctx context.Context) (int, error) {
	foundUser, err := self.Datastore.GetUser(ctx, self.RequestBody.UserGuid)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return httpstatus.NotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	if self.RequestBody.SignedKey != auth.Sign(foundUser.PublicKey, foundUser.Signature) {
		// Invalid credentials.
		// Treat this case as if the user was not found.
		return httpstatus.NotFound, db.ErrUserNotFound
	}

	userMailbox, err := self.Datastore.GetUserMailbox(ctx, self.RequestBody.UserGuid)
	if err != nil {
		if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.NotFound, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	foundUser.Mailbox = *userMailbox

	jwtToken, err := self.JwtAuth.UserToken(foundUser, time.Now().UTC().Add(time.Hour))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, errJwtError)
	}

	self.ResponseBody.JwtToken = jwtToken

	return httpstatus.OK, nil
}
