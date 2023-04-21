package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start up the Web server",
	Long:  "Start up the Web server, and listen for incoming requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Start up the Web server")
		return nil
	},
}
