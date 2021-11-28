package authorize

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

func RouteAuthorizeUser() route.Route { return route.Post("/authorize/user", &authUser{}) }

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteAuthorizeUser())
}
