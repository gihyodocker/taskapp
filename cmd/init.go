package cmd

func init() {
	serveCmd.AddCommand(apiCmd)
	serveCmd.AddCommand(webCmd)
	rootCmd.AddCommand(serveCmd)
}
