package users

import (
	"net/http"

	"github.com/wspowell/spiderweb/endpoint"
	"github.com/wspowell/spiderweb/lambda"
	"github.com/wspowell/spiderweb/server"
)

type definition struct {
	method  string
	path    string
	handler endpoint.Handler
}

var (
	create = definition{http.MethodPost, "/users", &createUser{}}
	get    = definition{http.MethodGet, "/users/{id}", &getUser{}}
	update = definition{http.MethodPut, "/users/{id}", &updateUser{}}
	delete = definition{http.MethodDelete, "/users/{id}", &deleteUser{}}
)

func Routes(server *server.Server, config *endpoint.Config) {
	server.Handle(config, create.method, create.path, create.handler)
	server.Handle(config, get.method, get.path, get.handler)
	server.Handle(config, update.method, update.path, update.handler)
	server.Handle(config, delete.method, delete.path, delete.handler)
}

func LambdaCreate(config *endpoint.Config) *lambda.Lambda {
	return lambda.New(config, create.path, create.handler)
}

func LambdaGet(config *endpoint.Config) *lambda.Lambda {
	return lambda.New(config, get.path, get.handler)
}

func LambdaUpdate(config *endpoint.Config) *lambda.Lambda {
	return lambda.New(config, update.path, update.handler)
}

func LambdaDelete(config *endpoint.Config) *lambda.Lambda {
	return lambda.New(config, delete.path, delete.handler)
}
