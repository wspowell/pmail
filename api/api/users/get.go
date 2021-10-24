package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/log"

	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/user"
)

type getUserResponse struct {
	userModel
	MailboxGuid string `json:"mailbox_guid,omitempty"`
}

type getUser struct {
	UserGuid     string           `spiderweb:"path=user_guid"`
	Datastore    db.Datastore     `spiderweb:"resource=datastore"`
	ResponseBody *getUserResponse `spiderweb:"response,mime=application/json,etag"`
}

func (self *getUser) Handle(ctx context.Context) (int, error) {
	foundUser, err := self.Datastore.GetUser(ctx, user.Guid(self.UserGuid))
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return http.StatusNotFound, errors.Propagate(icGetUserUserNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icGetUserDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icGetUserUnknownDbError, err, errUncaughtDbError)
		}
	}

	if userMailbox, err := self.Datastore.GetUserMailbox(ctx, foundUser.UserGuid); err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			// No mailbox assigned to user.
			log.Debug(ctx, "user does not have a mailbox")
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icGetUserMailboxDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icGetUserMailboxUnknownDbError, err, errUncaughtDbError)
		}
	} else {
		self.ResponseBody.MailboxGuid = string(userMailbox.MailboxGuid)
	}

	self.ResponseBody.userModel.Username = foundUser.Username
	self.ResponseBody.userModel.PineappleOnPizza = foundUser.Attributes.PineappleOnPizza

	log.Debug(ctx, "got user: %+v", foundUser)

	return http.StatusOK, nil
}
