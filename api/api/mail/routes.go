package mail

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

var (
	RouteCreate = route.Post("/mail", &createMail{})
	RouteList   = route.Get("/mail", &listMail{})
	RouteOpen   = route.Get("/mail/{mail_guid}", &openMail{})
)

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate)
	server.Handle(config, RouteList)
	server.Handle(config, RouteOpen)
}
