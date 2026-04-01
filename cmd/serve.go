package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/z2-cli/internal/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web API server",
	Long:  "Start an HTTP server that exposes the z2-cli data as a REST API for the web frontend.",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf(":%d", port)

		server := api.NewServer(addr, FrontendFS)
		return server.Start()
	},
}

func init() {
	serveCmd.Flags().Int("port", 8080, "port to listen on")
	rootCmd.AddCommand(serveCmd)
}
