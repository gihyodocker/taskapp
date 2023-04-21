package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start up the API server",
	Long:  "Start up the API server, and listen for incoming requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Start up the API server")
		return nil
	},
}
