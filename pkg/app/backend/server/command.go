package server

import (
	"context"
	"fmt"
	"github.com/gihyodocker/taskapp/pkg/config"
	"time"

	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/server"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

type command struct {
	port        int
	gracePeriod time.Duration
	configFile  string
}

func NewCommand() *cobra.Command {
	c := &command{
		port:        8180,
		gracePeriod: 5 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start up the backend server",
		RunE:  cli.WithContext(c.execute),
	}
	cmd.Flags().IntVar(&c.port, "port", c.port, "The port number used to run HTTP server.")
	cmd.Flags().DurationVar(&c.gracePeriod, "grace-period", c.gracePeriod, "How long to wait for graceful shutdown.")
	cmd.Flags().StringVar(&c.configFile, "config-file", c.configFile, "The path to the config file.")
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	appConfig, err := config.LoadConfigFile(c.configFile)
	if err != nil {
		slog.Error("failed to load api configuration",
			slog.String("config-file", c.configFile),
			err,
		)
		return err
	}
	// TODO Use appConfig
	fmt.Println(appConfig)

	options := []server.Option{
		server.WithGracePeriod(c.gracePeriod),
	}
	httpServer := server.NewHTTPServer(c.port, options...)
	group.Go(func() error {
		return httpServer.Serve(ctx)
	})

	if err := group.Wait(); err != nil {
		slog.Error("failed while running", err)
		return err
	}
	return nil
}
