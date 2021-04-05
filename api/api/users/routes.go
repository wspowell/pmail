package users

import (
	"net/http"

	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/users", &createUser{})
	server.Handle(config, http.MethodGet, "/users/{id}", &getUser{})
	server.Handle(config, http.MethodPut, "/users/{id}", &updateUser{})
	server.Handle(config, http.MethodDelete, "/users/{id}", &deleteUser{})
}
