package main

import (
	"kw_tool/cmd/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("The following error happened while starting up application, err = %s", err)
	}
}
