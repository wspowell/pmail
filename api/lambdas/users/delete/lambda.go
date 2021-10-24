package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/users"
	"github.com/wspowell/spiderweb/server/lambda"
)

func main() {
	lambda.New(api.Config(), users.RouteDelete).Start()
}
