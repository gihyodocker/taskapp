package mysql

import (
	"github.com/spf13/cobra"

	"github.com/gihyodocker/taskapp/pkg/app/tools/cmd/mysql/password"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "Command line utilities for MySQL",
		Long:  "Command line utilities fir MySQL for operating in taskapp",
	}
	cmd.AddCommand(password.NewCommand())

	return cmd
}
