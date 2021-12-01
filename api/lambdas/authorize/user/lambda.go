package main

import (
	"github.com/wspowell/spiderweb/server/lambda"

	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/authorize"
)

func main() {
	lambda.New(api.Config(), authorize.RouteAuthorizeUser()).Start()
}
