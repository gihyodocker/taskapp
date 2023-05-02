package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Application struct {
	Database *Database `yaml:"database"`
}

type Database struct {
	Host            string        `yaml:"host"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

func LoadConfigFile(configPath string) (*Application, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var appConfig Application
	if err := yaml.Unmarshal([]byte(data), &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
