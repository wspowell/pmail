package users

import (
	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
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
	userRouteConfig := *config
	userRouteConfig.Auther = middleware.NewUserAuth(config.Resources["datastore"].(db.Datastore))

	server.Handle(&userRouteConfig, RouteCreate)
	server.Handle(&userRouteConfig, RouteGet)
	server.Handle(&userRouteConfig, RouteUpdate)
	server.Handle(&userRouteConfig, RouteDelete)
}
