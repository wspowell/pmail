package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models"
)

type createUserRequest struct {
	mime.Json

	PublicKey string      `json:"publicKey"`
	Location  geoLocation `json:"mailboxLocation"`
}

type geoLocation struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type createUserResponse struct {
	mime.Json

	UserGuid       string `json:"userGuid"`
	SignedKey      string `json:"signedKey"`
	MailboxAddress string `json:"mailboxAddress"`
}

type createUser struct {
	Datastore db.Datastore
	body.Request[createUserRequest]
	body.Response[createUserResponse]
}

func (self *createUser) Handle(ctx context.Context) (int, error) {
	// FIXME: Move to validation.
	if self.RequestBody.PublicKey == "" {
		return httpstatus.UnprocessableEntity, errMissingPublicKey
	}

	if self.RequestBody.Location.Longitude == 0 || self.RequestBody.Location.Latitude == 0 {
		return httpstatus.UnprocessableEntity, errMissingMailboxLocation
	}

	newUser := models.CreateUser(self.RequestBody.PublicKey, models.NewCoordinate(self.RequestBody.Location.Longitude, self.RequestBody.Location.Latitude))

	if err := self.Datastore.CreateUser(ctx, newUser); err != nil {
		if errors.Is(err, db.ErrAddressExists) {
			return http.StatusConflict, err
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, err
		} else {
			return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	self.ResponseBody.UserGuid = newUser.Guid
	self.ResponseBody.SignedKey = auth.Sign(newUser.PublicKey, newUser.Signature)
	self.ResponseBody.MailboxAddress = newUser.Mailbox.Address

	return http.StatusCreated, nil
}
