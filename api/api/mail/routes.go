package mail

import (
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/request"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteCreate(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Post, "/mail", handler.NewHandle(createMail{User: apiResources.UserAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}
func RouteList(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Get, "/mail", handler.NewHandle(listMail{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}
func RouteOpen(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Get, "/mail/{mailGuid}", handler.NewHandle(openMail{User: apiResources.UserAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	server.Handle(RouteCreate(apiResources))
	//server.Handle(RouteList(apiResources))
	server.Handle(RouteOpen(apiResources))
}

type pathParams struct {
	MailGuid string
}

func (self *pathParams) PathParameters() []request.Parameter {
	return []request.Parameter{
		request.NewParam("mailGuid", &self.MailGuid),
	}
}
