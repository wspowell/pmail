package mailboxes

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

func RouteCreate() route.Route { return route.Post("/mailboxes", &createMailbox{}) }
func RouteGet() route.Route    { return route.Get("/mailboxes/{mailbox_address}", &getMailbox{}) }

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate())
	server.Handle(config, RouteGet())
}
