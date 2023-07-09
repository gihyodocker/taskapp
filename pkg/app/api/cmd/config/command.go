package config

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/config"
)

type command struct {
	outputFile string

	database *config.Database
}

func NewCommand() *cobra.Command {
	c := &command{
		outputFile: "api-config-local.yaml",
		database: &config.Database{
			Host:            "127.0.0.1",
			Username:        "taskapp_user",
			DBName:          "taskapp",
			MaxIdleConns:    5,
			MaxOpenConns:    10,
			ConnMaxLifetime: 1 * time.Hour,
		},
	}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Generate the api configuration file",
		RunE:  cli.WithContext(c.execute),
	}
	cmd.Flags().StringVar(&c.outputFile, "output-file", c.outputFile, "The config file output filename")
	cmd.Flags().StringVar(&c.database.Host, "database-host", c.database.Host, "The MySQL host address")
	cmd.Flags().StringVar(&c.database.Username, "database-username", c.database.Username, "The MySQL user name")
	cmd.Flags().StringVar(&c.database.Password, "database-password", c.database.Password, "The MySQL user password")
	cmd.Flags().StringVar(&c.database.DBName, "database-dbname", c.database.DBName, "The dbname name of taskapp")
	cmd.Flags().IntVar(&c.database.MaxIdleConns, "database-max-idle-conns", c.database.MaxIdleConns, "the maximum number of idle connections")
	cmd.Flags().IntVar(&c.database.MaxOpenConns, "database-max-open-conns", c.database.MaxOpenConns, "The maximum number of open connections")
	cmd.Flags().DurationVar(&c.database.ConnMaxLifetime, "database-conn-max-lifetime", c.database.ConnMaxLifetime, "The maximum amount of time a connection")

	cmd.MarkFlagRequired("database-password")
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	conf := config.Application{
		Database: c.database,
	}

	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(pwd, c.outputFile)
	if err := os.WriteFile(filepath.Join(pwd, c.outputFile), data, 0644); err != nil {
		slog.Error("failed to write api config file", err)
		return err
	}

	slog.Info("Completed generating the api config file.", slog.String("outputPath", outputPath))
	return nil
}
