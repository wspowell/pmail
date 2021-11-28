package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/authorize"
	"github.com/wspowell/spiderweb/server/lambda"
)

func main() {
	lambda.New(api.Config(), authorize.RouteAuthorizeUser()).Start()
}
