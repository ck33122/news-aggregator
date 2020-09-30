package app

import (
	"io/ioutil"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

type Environment string

const (
	EnvDevelopment Environment = "Development"
	EnvProduction  Environment = "Production"
)

type LoggerConfig struct {
	Dir string
}

type DatabaseConfig struct {
	User     string
	Password string
	Database string
	Address  string
}

type ApiConfig struct {
	Listen string
}

type ImportChannelConfig struct {
	Address string
	Id      uuid.UUID
}

type Config struct {
	AppName     string
	Environment Environment
	Logger      LoggerConfig
	Database    DatabaseConfig
	Api         ApiConfig
	Import      []ImportChannelConfig
}

func NewConfig(configPath, appName string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	config.AppName = appName
	return &config, nil
}
