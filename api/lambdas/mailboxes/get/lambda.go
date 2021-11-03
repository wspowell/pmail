package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/mailboxes"
	"github.com/wspowell/spiderweb/server/lambda"
)

func main() {
	lambda.New(api.Config(), mailboxes.RouteGet).Start()
}
