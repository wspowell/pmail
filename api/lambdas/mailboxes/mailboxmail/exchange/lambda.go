package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/spiderweb/server/lambda"
)

func main() {
	lambda.New(api.Config(), mailboxmail.RouteExchange).Start()
}
