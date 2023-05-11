package mysql

import (
	"github.com/spf13/cobra"

	"github.com/gihyodocker/taskapp/pkg/app/mysql/password"
)

var Command = &cobra.Command{
	Use:   "mysql",
	Short: "Command line utilities for MySQL",
	Long:  "Command line utilities fir MySQL for operating in taskapp",
}

func init() {
	Command.AddCommand(password.NewCommand())
}
