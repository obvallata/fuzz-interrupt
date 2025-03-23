package main

import (
	"log"

	"diploma/chaosmachine/internal/action"
	"diploma/chaosmachine/internal/config"
	"diploma/chaosmachine/internal/interaction"
	"diploma/chaosmachine/internal/server"
)

func main() {
	// TODO: get config path from args or env
	conf, err := config.GetConfig("/Users/ddr/fuzz-interrupt/chaosmachine/config/config.yaml")
	if err != nil {
		log.Fatalf("get config: %s", err.Error())
	}

	clients, err := interaction.NewClients(conf.Clients)
	if err != nil {
		log.Fatalf("new clients: %s", err.Error())
	}
	defer clients.Close()

	a := action.NewAction(clients)
	server.Serve(conf.Server, a)

	a.BuildAutomaton()
}
