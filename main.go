package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/backend"
	"github.com/gihyodocker/taskapp/pkg/app/web"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	// TODO: Initialize slog
	c := cli.NewCLI("taskapp", "Simple task management application")
	c.AddCommands(backend.Command)
	c.AddCommands(web.Command)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
