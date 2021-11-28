package mailboxmail

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"

	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/user"
)

type createMailboxResponse struct {
	DroppedOffMailGuids []string `json:"droppedOffMailGuids"`
	PickedUpMailGuids   []string `json:"pickedUpMailGuids"`
}

type exchangeMail struct {
	AuthorizedUser *middleware.UserAuth   `spiderweb:"auth"`
	MailboxAddress string                 `spiderweb:"path=mailbox_address"`
	Datastore      db.Datastore           `spiderweb:"resource=datastore"`
	ResponseBody   *createMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *exchangeMail) Handle(ctx context.Context) (int, error) {
	foundMailbox, err := self.Datastore.GetMailbox(ctx, self.MailboxAddress)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, errors.Propagate(icExchangeMailGetMailboxNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icExchangeMailGetMailboxDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icExchangeMailGetMailboxUnknownDbError, err, errUncaughtDbError)
		}
	}

	droppedOffMail, err := self.Datastore.DropOffMail(ctx, user.Guid(self.AuthorizedUser.UserGuid), foundMailbox.Address)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, errors.Propagate(icExchangeMailDropOffMailNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icExchangeMailDropOffMailDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icExchangeMailDropOffMailUnknownDbError, err, errUncaughtDbError)
		}
	}

	self.ResponseBody.DroppedOffMailGuids = mail.ToStrings(droppedOffMail)

	if foundMailbox.Owner == "" || foundMailbox.Owner == user.Guid(self.AuthorizedUser.UserGuid) {
		pickedUpMail, err := self.Datastore.PickUpMail(ctx, user.Guid(self.AuthorizedUser.UserGuid), foundMailbox.Address)
		if err != nil {
			if errors.Is(err, db.ErrMailboxNotFound) {
				return http.StatusNotFound, errors.Propagate(icExchangeMailPickUpMailNotFound, err)
			} else if errors.Is(err, db.ErrInternalFailure) {
				return http.StatusInternalServerError, errors.Propagate(icExchangeMailPickUpMailDbError, err)
			} else {
				return http.StatusInternalServerError, errors.Convert(icExchangeMailPickUpMailUnknownDbError, err, errUncaughtDbError)
			}
		}

		self.ResponseBody.PickedUpMailGuids = mail.ToStrings(pickedUpMail)
	}

	return http.StatusOK, nil
}
