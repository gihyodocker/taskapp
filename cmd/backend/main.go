package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/backend/config"
	"github.com/gihyodocker/taskapp/pkg/app/backend/server"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	c := cli.NewCLI("taskapp-backend", "The backend API application of taskapp")
	c.AddCommands(
		server.NewCommand(),
		config.NewCommand(),
	)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
