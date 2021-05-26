package main

import (
	"github.com/wspowell/snailmail/server"
)

func main() {
	server.New().Listen()
}
