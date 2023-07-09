package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/tools/cmd/mysql"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	c := cli.NewCLI("taskapp-tools", "The utility tools of taskapp")
	c.AddCommands(
		mysql.NewCommand(),
	)
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
