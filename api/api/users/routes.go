package users

import (
	"net/http"

	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/server"
)

func Routes(server *server.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/users", &createUser{})
	server.Handle(config, http.MethodGet, "/users/{id}", &getUser{})
	server.Handle(config, http.MethodPut, "/users/{id}", &updateUser{})
	server.Handle(config, http.MethodDelete, "/users/{id}", &deleteUser{})
}
