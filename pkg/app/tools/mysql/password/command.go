package password

import (
	"context"
	"crypto/rand"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/cli"
)

const (
	passwordChars    = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rootPasswordFile = "./secrets/mysql_root_password"
	userPasswordFile = "./secrets/mysql_user_password"
)

type command struct {
}

func NewCommand() *cobra.Command {
	c := &command{}

	cmd := &cobra.Command{
		Use:   "generate-password",
		Short: "Generate MySQL password",
		RunE:  cli.WithContext(c.execute),
	}
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rootPassword, err := makePassword(16)
	if err != nil {
		slog.Error("failed to generate MySQL root password", err)
		return err
	}
	userPassword, err := makePassword(16)
	if err != nil {
		slog.Error("failed to generate MySQL user password", err)
		return err
	}

	if err := os.WriteFile(filepath.Join(pwd, rootPasswordFile), []byte(rootPassword), 0644); err != nil {
		slog.Error("failed to write MySQL root password to the file", err)
		return err
	}

	if err := os.WriteFile(filepath.Join(pwd, userPasswordFile), []byte(userPassword), 0644); err != nil {
		slog.Error("failed to write MySQL user password to the file", err)
		return err
	}

	slog.Info("Completed generating the root and user passwords")
	return nil
}

func makePassword(digit int) (string, error) {
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(passwordChars[int(v)%len(passwordChars)])
	}
	return result, nil
}
