package mailbox

import (
	"net/http"

	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server"
)

func Routes(server *server.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/users/{user_id}/mailbox", &createMailbox{})
	server.Handle(config, http.MethodGet, "/users/{user_id}/mailbox", &checkMailbox{})
}
