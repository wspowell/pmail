package mailboxes

import (
	"net/http"

	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server/restful"
)

func Routes(server *restful.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodGet, "/mailboxes/{id}", &getMailbox{})
}
