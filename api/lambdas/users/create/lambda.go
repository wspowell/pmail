package main

import (
	"github.com/wspowell/spiderweb/server/lambda"

	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/snailmail/resources"
)

func main() {
	apiResources := resources.Load()

	_, path, handle := users.RouteCreate(apiResources)
	handler := lambda.New(path, handle)
	handler.Start()
}
