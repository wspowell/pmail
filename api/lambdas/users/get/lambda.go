package main

import (
	"github.com/wspowell/spiderweb/server/lambda"

	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/users"
)

func main() {
	lambda.New(api.Config(), users.RouteGet()).Start()
}
