package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	SimpleConfig = `
database:
  host: localhost
  username: root
  password: password
  dbname: taskapp
  maxIdleConns: 5
  maxOpenConns: 10
  connMaxLifetime: 1h
`
)

func TestLoadConfigFile(t *testing.T) {
	dir := t.TempDir()
	configFilePath := dir + "/config.yaml"
	if err := os.WriteFile(configFilePath, []byte(SimpleConfig), 0644); err != nil {
		t.Error("failed to write config file", err)
	}

	actual, err := LoadConfigFile(configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, actual.Database)
	assert.Equal(t, "localhost", actual.Database.Host)
	assert.Equal(t, "root", actual.Database.Username)
	assert.Equal(t, "password", actual.Database.Password)
	assert.Equal(t, "taskapp", actual.Database.DBName)
	assert.Equal(t, 5, actual.Database.MaxIdleConns)
	assert.Equal(t, 10, actual.Database.MaxOpenConns)
	assert.Equal(t, 1*time.Hour, actual.Database.ConnMaxLifetime)
}
