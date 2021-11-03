package users

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

var (
	RouteCreate = route.Post("/users", &createUser{})
	RouteGet    = route.Get("/users/{user_guid}", &getUser{})
	RouteUpdate = route.Put("/users/{user_guid}", &updateUser{})
	RouteDelete = route.Delete("/users/{user_guid}", &deleteUser{})
)

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate)
	server.Handle(config, RouteGet)
	server.Handle(config, RouteUpdate)
	server.Handle(config, RouteDelete)
}
