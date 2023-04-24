package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

type Runner func(ctx context.Context) error

func WithContext(runner Runner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			select {
			case s := <-ch:
				slog.Info("stopping due to signal", slog.String("signal", s.String()))
				cancel()
			case <-ctx.Done():
			}
		}()

		slog.Info(fmt.Sprintf("running application by %s command", cmd.Name()))
		return runner(ctx)
	}
}
