package mailboxes

import (
	"net/http"

	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server"
)

func Routes(server *server.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodGet, "/mailboxes/{id}", &getMailbox{})
}
