package users

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

func RouteCreate() route.Route { return route.Post("/users", &createUser{}) }
func RouteGet() route.Route    { return route.Get("/users/{user_guid}", &getUser{}) }
func RouteUpdate() route.Route { return route.Put("/users/{user_guid}", &updateUser{}) }
func RouteDelete() route.Route { return route.Delete("/users/{user_guid}", &deleteUser{}) }

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate())
	server.Handle(config, RouteGet())
	server.Handle(config, RouteUpdate())
	server.Handle(config, RouteDelete())
}
