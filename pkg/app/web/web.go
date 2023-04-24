package web

import (
	"github.com/spf13/cobra"

	"github.com/gihyodocker/todoapp/pkg/app/web/server"
)

var Command = &cobra.Command{
	Use:   "web",
	Short: "Start up the web application",
	Long:  "Start up the web application, and listen for incoming requests",
}

func init() {
	Command.AddCommand(server.NewCommand())
}
