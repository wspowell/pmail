package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/db"
)

type getUserResponse struct {
	mime.Json

	MailboxAddress string `json:"mailboxAddress,omitempty"`
}

type getUser struct {
	pathParams
	Datastore db.Datastore
	body.Response[getUserResponse]
}

func (self *getUser) Handle(ctx context.Context) (int, error) {
	// foundUser, err := self.Datastore.GetUser(ctx, user.Guid(self.UserGuid))
	// if err != nil {
	// 	if errors.Is(err, db.ErrUserNotFound) {
	// 		return http.StatusNotFound, err
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return http.StatusInternalServerError, err
	// 	} else {
	// 		return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// }

	// if userMailbox, err := self.Datastore.GetUserMailbox(ctx, foundUser.UserGuid); err != nil {
	// 	if errors.Is(err, db.ErrMailboxNotFound) {
	// 		// No mailbox assigned to user.
	// 		log.Debug(ctx, "user does not have a mailbox")
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return http.StatusInternalServerError, err
	// 	} else {
	// 		return http.StatusInternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// } else {
	// 	self.ResponseBody.MailboxAddress = userMailbox.Address
	// }

	// self.ResponseBody.userModel.Username = foundUser.Username
	// self.ResponseBody.userModel.PineappleOnPizza = foundUser.Attributes.PineappleOnPizza

	// log.Debug(ctx, "got user: %+v", foundUser)

	return http.StatusOK, nil
}
