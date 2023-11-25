package server

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
	"github.com/gihyodocker/taskapp/pkg/app/web/handler"
	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/server"
)

type command struct {
	port        int
	apiAddress  string
	assetsDir   string
	gracePeriod time.Duration
}

func NewCommand() *cobra.Command {
	c := &command{
		port:        8280,
		apiAddress:  "http://127.0.0.1:8180",
		gracePeriod: 5 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start up the web server",
		RunE:  cli.WithContext(c.execute),
	}
	cmd.Flags().IntVar(&c.port, "port", c.port, "The port number used to run HTTP api.")
	cmd.Flags().StringVar(&c.apiAddress, "api-address", c.apiAddress, "The API address.")
	cmd.Flags().StringVar(&c.assetsDir, "assets-dir", c.assetsDir, "The path to the assets directory.")
	cmd.Flags().DurationVar(&c.gracePeriod, "grace-period", c.gracePeriod, "How long to wait for graceful shutdown.")
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	options := []server.Option{
		server.WithGracePeriod(c.gracePeriod),
	}

	// HTTP clients
	taskCli := client.NewTask(c.apiAddress)

	// Handlers
	indexHandler := handler.NewIndex(taskCli)
	deleteHandler := handler.NewDelete(taskCli)
	updateHandler := handler.NewUpdate(taskCli)
	createHandler := handler.NewCreate(taskCli)

	httpServer := server.NewHTTPServer(c.port, options...)
	// Health check
	httpServer.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Static files
	if c.assetsDir != "" {
		httpServer.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir(c.assetsDir))))
	}

	// Task application endpoints
	httpServer.Post("/tasks/{id}/update/complete", updateHandler.Complete)
	httpServer.Get("/tasks/{id}/update", updateHandler.Input)
	httpServer.Post("/tasks/{id}/delete/complete", deleteHandler.Complete)
	httpServer.Get("/tasks/{id}/delete", deleteHandler.Confirm)
	httpServer.Post("/tasks/create/complete", createHandler.Complete)
	httpServer.Get("/tasks/create", createHandler.Input)
	httpServer.Get("/", indexHandler.Index)
	group.Go(func() error {
		return httpServer.Serve(ctx)
	})

	if err := group.Wait(); err != nil {
		slog.Error("failed while running", err)
		return err
	}
	return nil
}
