package mail

import (
	"github.com/wspowell/context"
	"github.com/wspowell/spiderweb/body"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/mime"

	"github.com/wspowell/snailmail/resources/auth"
	"github.com/wspowell/snailmail/resources/db"
)

type listMailResponse struct {
	mime.Json

	//Mail []mailResponse `json:"mail"`
}

type listMail struct {
	auth.User
	Datastore db.Datastore
	body.Response[listMailResponse]
}

func (self *listMail) Handle(ctx context.Context) (int, error) {
	// userMail, err := self.Datastore.GetUserMail(ctx, user.Guid(self.AuthorizedUser.UserGuid))
	// if err != nil {
	// 	if errors.Is(err, db.ErrUserNotFound) {
	// 		return httpstatus.InternalServerError, err
	// 	} else if errors.Is(err, db.ErrInternalFailure) {
	// 		return httpstatus.InternalServerError, err
	// 	} else {
	// 		return httpstatus.InternalServerError, errors.Wrap(err, errUncaughtDbError)
	// 	}
	// }

	// self.ResponseBody.Mail = []mailResponse{}
	// for index := range userMail {
	// 	mailResponseItem := mailResponse{
	// 		MailGuid:    string(userMail[index].MailGuid),
	// 		From:        string(userMail[index].From),
	// 		To:          string(userMail[index].To),
	// 		Contents:    userMail[index].Contents,
	// 		SentOn:      userMail[index].SentOn,
	// 		DeliveredOn: userMail[index].DeliveredOn,
	// 		OpenedOn:    userMail[index].OpenedOn,
	// 	}

	// 	self.ResponseBody.Mail = append(self.ResponseBody.Mail, mailResponseItem)
	// }

	return httpstatus.OK, nil
}
