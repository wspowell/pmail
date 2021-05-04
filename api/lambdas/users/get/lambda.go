package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/users"
)

func main() {
	users.LambdaGet(api.Config()).Start()
}
