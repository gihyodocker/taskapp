package backend

import (
	"github.com/spf13/cobra"

	"github.com/gihyodocker/taskapp/pkg/app/backend/config"
	"github.com/gihyodocker/taskapp/pkg/app/backend/server"
)

var Command = &cobra.Command{
	Use:   "backend",
	Short: "Start up the backend application",
	Long:  "Start up the backend application, and listen for incoming requests",
}

func init() {
	Command.AddCommand(server.NewCommand())
	Command.AddCommand(config.NewCommand())
}
