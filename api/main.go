package main

import (
	"github.com/wspowell/snailmail/server"
)

func main() {
	snailmail := server.New()
	snailmail.Listen()
}
