package users

import (
	"time"

	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/spiderweb/handler"
	"github.com/wspowell/spiderweb/httpmethod"
	"github.com/wspowell/spiderweb/request"
	"github.com/wspowell/spiderweb/server/restful"
)

func RouteCreate(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Post, "/users", handler.NewHandle(createUser{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}
func RouteGet(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Get, "/users/{userGuid}", handler.NewHandle(getUser{Datastore: apiResources.Datastore}).
		WithETag(30 * time.Second).
		WithLogConfig(apiResources.LogConfig)
}
func RouteUpdate(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Put, "/users/{userGuid}", handler.NewHandle(updateUser{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}
func RouteDelete(apiResources *resources.Resources) (string, string, *handler.Handle) {
	return httpmethod.Delete, "/users/{userGuid}", handler.NewHandle(deleteUser{Datastore: apiResources.Datastore}).
		WithLogConfig(apiResources.LogConfig)
}

func Routes(server *restful.Server, apiResources *resources.Resources) {
	server.Handle(RouteCreate(apiResources))
	server.Handle(RouteGet(apiResources))
	//server.Handle(RouteUpdate(apiResources))
	server.Handle(RouteDelete(apiResources))
}

type pathParams struct {
	UserGuid string
}

func (self *pathParams) PathParameters() []request.Parameter {
	return []request.Parameter{
		request.NewParam("userGuid", &self.UserGuid),
	}
}
