package server

import (
	"context"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/todoapp/pkg/cli"
)

type server struct {
	port int
	// TODO: graceful shutdown
}

func NewCommand() *cobra.Command {
	s := &server{
		port: 8280,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start up the web server",
		RunE:  cli.WithContext(s.execute),
	}
	// TODO Set flags
	return cmd
}

func (s *server) execute(ctx context.Context) error {
	slog.Info("This is web server.")
	return nil
}
