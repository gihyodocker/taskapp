package cli

import "github.com/spf13/cobra"

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
