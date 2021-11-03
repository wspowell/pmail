package mailboxmail

import (
	"net/http"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"

	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/snailmail/resources/models/mail"
	"github.com/wspowell/snailmail/resources/models/mailbox"
	"github.com/wspowell/snailmail/resources/models/user"
)

type createMailboxResponse struct {
	DroppedOffMail []string `json:"dropped_off_mail_guids"`
	PickedUpMail   []string `json:"picked_up_mail_guids"`
}

type exchangeMail struct {
	AuthorizedUser *middleware.UserAuth   `spiderweb:"auth"`
	MailboxGuid    string                 `spiderweb:"path=mailbox_guid"`
	Datastore      db.Datastore           `spiderweb:"resource=datastore"`
	ResponseBody   *createMailboxResponse `spiderweb:"response,mime=application/json"`
}

func (self *exchangeMail) Handle(ctx context.Context) (int, error) {
	foundMailbox, err := self.Datastore.GetMailbox(ctx, mailbox.Guid(self.MailboxGuid))
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, errors.Propagate(icExchangeMailGetMailboxNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icExchangeMailGetMailboxDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icExchangeMailGetMailboxUnknownDbError, err, errUncaughtDbError)
		}
	}

	droppedOffMail, err := self.Datastore.DropOffMail(ctx, user.Guid(self.AuthorizedUser.UserGuid), foundMailbox.MailboxGuid)
	if err != nil {
		if errors.Is(err, db.ErrMailboxNotFound) {
			return http.StatusNotFound, errors.Propagate(icExchangeMailDropOffMailNotFound, err)
		} else if errors.Is(err, db.ErrInternalFailure) {
			return http.StatusInternalServerError, errors.Propagate(icExchangeMailDropOffMailDbError, err)
		} else {
			return http.StatusInternalServerError, errors.Convert(icExchangeMailDropOffMailUnknownDbError, err, errUncaughtDbError)
		}
	}

	self.ResponseBody.DroppedOffMail = mail.ToStrings(droppedOffMail)

	if foundMailbox.Owner == "" || foundMailbox.Owner == user.Guid(self.AuthorizedUser.UserGuid) {
		pickedUpMail, err := self.Datastore.PickUpMail(ctx, user.Guid(self.AuthorizedUser.UserGuid), foundMailbox.MailboxGuid)
		if err != nil {
			if errors.Is(err, db.ErrMailboxNotFound) {
				return http.StatusNotFound, errors.Propagate(icExchangeMailPickUpMailNotFound, err)
			} else if errors.Is(err, db.ErrInternalFailure) {
				return http.StatusInternalServerError, errors.Propagate(icExchangeMailPickUpMailDbError, err)
			} else {
				return http.StatusInternalServerError, errors.Convert(icExchangeMailPickUpMailUnknownDbError, err, errUncaughtDbError)
			}
		}

		self.ResponseBody.PickedUpMail = mail.ToStrings(pickedUpMail)
	}

	return http.StatusOK, nil
}
