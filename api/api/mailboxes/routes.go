package mailboxes

import (
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

var (
	RouteCreate = route.Post("/mailboxes", &createMailbox{})
	RouteGet    = route.Get("/mailboxes/{mailbox_guid}", &getMailbox{})
)

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, RouteCreate)
	server.Handle(config, RouteGet)
}
