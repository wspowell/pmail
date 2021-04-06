package mailbox

import (
	"net/http"

	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/users/{user_id}/mailbox", &createMailbox{})
	server.Handle(config, http.MethodGet, "/users/{user_id}/mailbox", &checkMailbox{})
}
