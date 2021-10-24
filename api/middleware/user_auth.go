package middleware

import (
	"bytes"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"

	"github.com/wspowell/spiderweb/http"
)

var (
	ErrUserUnauthorized = errors.New("auth-user-1", "user unauthorized")
)

type UserAuth struct {
	AuthenticatedUser user.User
	datastore         db.Datastore
}

func NewUserAuth(datastore db.Datastore) UserAuth {
	return UserAuth{
		datastore: datastore,
	}
}

func (self UserAuth) Auth(ctx context.Context, VisitAllHeaders func(func(key, value []byte))) (int, error) {
	authHeader := []byte(http.HeaderAuthorization)

	status := http.StatusOK
	var err error
	VisitAllHeaders(func(key, value []byte) {
		if err == nil && bytes.EqualFold(authHeader, key) {
			foundUser, dbErr := self.datastore.GetUser(ctx, user.Guid(value))
			if dbErr != nil {
				status = http.StatusUnauthorized
				err = ErrUserUnauthorized
				return
			}

			self.AuthenticatedUser = *foundUser
		}

	})

	return status, err
}
