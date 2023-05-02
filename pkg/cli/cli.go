package cli

import (
	"fmt"

	"github.com/gihyodocker/taskapp/pkg/version"
	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd *cobra.Command
}

func NewCLI(name, desc string) *CLI {
	c := &CLI{
		rootCmd: &cobra.Command{
			Use:   name,
			Short: desc,
		},
	}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the information of current binary",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Get())
		},
	}
	c.rootCmd.AddCommand(versionCmd)
	// TODO: Set global flags?
	return c
}

func (c *CLI) AddCommands(commands ...*cobra.Command) {
	for _, cmd := range commands {
		c.rootCmd.AddCommand(cmd)
	}
}

func (c *CLI) Execute() error {
	return c.rootCmd.Execute()
}
