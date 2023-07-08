package server

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
	"github.com/gihyodocker/taskapp/pkg/app/web/page"
	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/server"
)

type command struct {
	port        int
	apiAddress  string
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
	cmd.Flags().StringVar(&c.apiAddress, "api-address", c.apiAddress, "The api API address.")
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

	// Pages
	indexPage := page.NewIndex(taskCli)
	deletePage := page.NewDelete(taskCli)
	updatePage := page.NewUpdate(taskCli)
	createPage := page.NewCreate(taskCli)

	httpServer := server.NewHTTPServer(c.port, options...)
	httpServer.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// endpoints
	httpServer.Post("/tasks/{id}/update/complete", updatePage.Complete)
	httpServer.Get("/tasks/{id}/update", updatePage.Input)
	httpServer.Post("/tasks/{id}/delete/complete", deletePage.Complete)
	httpServer.Get("/tasks/{id}/delete", deletePage.Confirm)
	httpServer.Post("/tasks/create/complete", createPage.Complete)
	httpServer.Get("/tasks/create", createPage.Input)
	httpServer.Get("/", indexPage.Index)
	group.Go(func() error {
		return httpServer.Serve(ctx)
	})

	if err := group.Wait(); err != nil {
		slog.Error("failed while running", err)
		return err
	}
	return nil
}
