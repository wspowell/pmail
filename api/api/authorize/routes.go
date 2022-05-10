package authorize

import (
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteAuthorizeUser(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Post, "/authorize/user", handler.NewHandle(authUser{JwtAuth: apiResources.UserAuth.JwtAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func RouteAuthorizeSession(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Post, "/authorize/session", handler.NewHandle(authSession{User: apiResources.UserAuth, Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	server.Handle(RouteAuthorizeUser(apiResources))
	//server.Handle(RouteAuthorizeSession(apiResources))
}
