package mailboxmail

import (
	"github.com/wspowell/snailmail/middleware"
	"github.com/wspowell/snailmail/resources/db"
	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
	"github.com/wspowell/spiderweb/server/route"
)

var (
	RouteExchange = route.Get("/mailboxes/{mailbox_guid}/mail", &exchangeMail{})
)

func Routes(server *restful.Server, config *endpoint.Config) {
	userRouteConfig := *config
	userRouteConfig.Auther = middleware.NewUserAuth(config.Resources["datastore"].(db.Datastore))

	server.Handle(&userRouteConfig, RouteExchange)
}
