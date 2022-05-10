package main

import (
	"github.com/wspowell/spiderweb/server/lambda"

	"github.com/wspowell/snailmail/api/mailboxes/mailboxmail"
	"github.com/wspowell/snailmail/resources"
)

func main() {
	apiResources := resources.Load()

	_, path, handle := mailboxmail.RouteGetMailboxMail(apiResources)
	handler := lambda.New(path, handle)
	handler.Start()
}
