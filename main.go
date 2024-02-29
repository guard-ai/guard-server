package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/guard-ai/guard-server/app/controllers"
	"github.com/guard-ai/guard-server/server"

	config "github.com/guard-ai/guard-server/pkg"
)

func main() {
	e := server.CreateServer()
	controller := controllers.New()
	server.RegisterRoutes(e, controller)

	address := fmt.Sprintf("%s:%s", config.Env().ServerAddress, config.Env().ServerPort)
	if err := e.Start(address); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
