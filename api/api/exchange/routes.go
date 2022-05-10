package exchange

import (
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/request"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteExchange(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Put, "/exchangemail", handler.NewHandle(exchangeMail{User: apiResources.UserAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	server.Handle(RouteExchange(apiResources))
}

type pathParams struct {
	MailGuid string
}

func (self *pathParams) PathParameters() []request.Parameter {
	return []request.Parameter{
		request.NewParam("mailGuid", &self.MailGuid),
	}
}
