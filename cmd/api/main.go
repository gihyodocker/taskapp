package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/api/cmd/config"
	"github.com/gihyodocker/taskapp/pkg/app/api/cmd/server"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	c := cli.NewCLI("taskapp-api", "The API application of taskapp")
	c.AddCommands(
		server.NewCommand(),
		config.NewCommand(),
	)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
