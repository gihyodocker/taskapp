package cmd

import (
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start up some application server",
	Long:  "Start up some application server, and listen for incoming requests",
}
