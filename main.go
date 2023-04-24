package main

import (
	"log"

	"github.com/gihyodocker/todoapp/pkg/app/backend"
	"github.com/gihyodocker/todoapp/pkg/app/web"
	"github.com/gihyodocker/todoapp/pkg/cli"
)

func main() {
	// TODO: Initialize slog
	c := cli.NewCLI("todoapp", "TODO application")
	c.AddCommands(backend.Command)
	c.AddCommands(web.Command)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
