package mailboxes

import (
	"net/http"

	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodGet, "/mailboxes/{id}", &getMailbox{})
}
