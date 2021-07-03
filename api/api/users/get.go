package users

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
)

type getUserResponse struct {
	userModel
}

type getUser struct {
	UserId       uint32                 `spiderweb:"path=id"`
	Users        resources.UserStore    `spiderweb:"resource=userstore"`
	MailboxStore resources.MailboxStore `spiderweb:"resource=mailboxstore"`
	ResponseBody *getUserResponse       `spiderweb:"response,mime=application/json,etag"`
}

func (self *getUser) Handle(ctx context.Context) (int, error) {
	user, err := self.Users.GetUser(ctx, self.UserId)
	if err != nil {
		if errors.Is(err, resources.ErrUserNotFound) {
			return http.StatusNotFound, errors.Wrap(icGetUserUserNotFound, err)
		}
		return http.StatusInternalServerError, errors.Wrap(icGetUserGetUserError, err)
	}

	mailbox, err := self.MailboxStore.GetMailboxByUserId(ctx, self.UserId)
	if err != nil {
		if errors.Is(err, resources.ErrorMailboxNotFound) {
			// Ignore.
		} else {
			return http.StatusInternalServerError, errors.Wrap(icGetUserMailboxLookupError, err)
		}
	}

	if mailbox != nil {
		self.ResponseBody.MailboxId = mailbox.Id
	}

	self.ResponseBody.userModel.Username = user.Username
	self.ResponseBody.userModel.PineappleOnPizza = user.Attributes.PineappleOnPizza

	return http.StatusOK, nil
}
