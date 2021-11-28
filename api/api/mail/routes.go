package mail

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

func RouteCreate() route.Route { return route.Post("/mail", &createMail{}) }
func RouteList() route.Route   { return route.Get("/mail", &listMail{}) }
func RouteOpen() route.Route   { return route.Get("/mail/{mail_guid}", &openMail{}) }

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate())
	server.Handle(config, RouteList())
	server.Handle(config, RouteOpen())
}
