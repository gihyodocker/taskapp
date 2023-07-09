package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/web/cmd/server"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	c := cli.NewCLI("taskapp-web", "The web application of taskapp")
	c.AddCommands(
		server.NewCommand(),
	)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
