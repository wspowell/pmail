package authorize

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

var (
	RouteAuthorizeUser = route.Post("/authorize/user", &authUser{})
)

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteAuthorizeUser)
}
