package mailboxmail

import (
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/request"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteGetMailboxMail(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Get, "/mailboxes/{mailboxAddress}/mail", handler.NewHandle(getMailboxMail{User: apiResources.UserAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	server.Handle(RouteGetMailboxMail(apiResources))
}

type pathParams struct {
	MailboxAddress string
}

func (self *pathParams) PathParameters() []request.Parameter {
	return []request.Parameter{
		request.NewParam("mailboxAddress", &self.MailboxAddress),
	}
}
