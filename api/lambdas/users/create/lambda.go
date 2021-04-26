package main

import (
	"github.com/wspowell/pmail/api"
	"github.com/wspowell/pmail/api/users"
)

func main() {
	users.LambdaCreate(api.Config()).Start()
}
