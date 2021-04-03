package users

import (
	"net/http"

	"github.com/wspowell/spiderweb"
	"github.com/wspowell/spiderweb/endpoint"
)

func Routes(server *spiderweb.Server, config *endpoint.Config) {
	server.Handle(config, http.MethodPost, "/users", &createUser{})
	server.Handle(config, http.MethodPatch, "/users/{id}", &updateUser{})
	server.Handle(config, http.MethodDelete, "/users/{id}", &deleteUser{})
}

type updateUser struct{}

func (self *updateUser) Handle(ctx *endpoint.Context) (int, error) {
	return http.StatusNotImplemented, nil
}

type deleteUser struct{}

func (self *deleteUser) Handle(ctx *endpoint.Context) (int, error) {
	return http.StatusNotImplemented, nil
}
