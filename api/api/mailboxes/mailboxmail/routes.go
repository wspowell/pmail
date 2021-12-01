package mailboxmail

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

func RouteExchange() route.Route {
	return route.Get("/mailboxes/{mailbox_address}/mail", &exchangeMail{})
}

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteExchange())
}
