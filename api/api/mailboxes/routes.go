package mailboxes

import (
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/request"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteCreate(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Post, "/mailboxes", handler.NewHandle(createMailbox{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}
func RouteGet(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Get, "/mailboxes/{mailboxAddress}", handler.NewHandle(getMailbox{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	//server.Handle(RouteCreate(apiResources))
	//server.Handle(RouteGet(apiResources))
}

type pathParams struct {
	MailboxAddress string
}

func (self *pathParams) PathParameters() []request.Parameter {
	return []request.Parameter{
		request.NewParam("mailboxAddress", &self.MailboxAddress),
	}
}
