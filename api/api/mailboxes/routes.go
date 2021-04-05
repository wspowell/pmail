package mailboxes

import (
	"net/http"

	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/mailboxes", &createMailbox{})
	// server.Handle(config, http.MethodGet, "/mailboxes/{id}", &getMailbox{})
	// server.Handle(config, http.MethodPut, "/mailboxes/{id}", &updateUser{})
	// server.Handle(config, http.MethodDelete, "/mailboxes/{id}", &deleteUser{})
}
