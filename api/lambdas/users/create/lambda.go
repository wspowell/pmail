package main

import (
	"github.com/wspowell/snailmail/api"
	"github.com/wspowell/snailmail/api/users"
)

func main() {
	users.LambdaCreate(api.Config()).Start()
}
