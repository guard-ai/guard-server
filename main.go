package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/guard-ai/guard-server/server"
	config "github.com/guard-ai/guard-server/pkg"
)

func main() {
	e := server.CreateServer()
	address := fmt.Sprintf("%s:%s", config.ServerAddress, config.ServerPort)
	if err := e.Start(address); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
