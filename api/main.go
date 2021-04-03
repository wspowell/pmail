package main

import (
	"github.com/wspowell/pmail/server"
)

func main() {
	pmail := server.New()
	pmail.Listen()
}
