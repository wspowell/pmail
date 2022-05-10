package exchange

import (
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models"
)

type exchangeMailRequest struct {
	mime.Json

	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type exchangeMail struct {
	auth.User
	Datastore db.Datastore
	body.Request[exchangeMailRequest]
}

func (self *exchangeMail) Handle(ctx context.Context) (int, error) {
	if self.RequestBody.Longitude == 0 || self.RequestBody.Latitude == 0 {
		return httpstatus.UnprocessableEntity, errMissingLocation
	}

	location := models.NewCoordinate(self.RequestBody.Longitude, self.RequestBody.Latitude)
	err := self.Datastore.ExchangeMail(ctx, self.AuthorizedUser.Guid, location)
	if err != nil {
		if errors.Is(err, db.ErrInternalFailure) {
			return httpstatus.InternalServerError, err
		} else {
			return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
		}
	}

	return httpstatus.NoContent, nil
}
