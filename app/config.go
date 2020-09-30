package app

import (
	"flag"
	"io/ioutil"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

var (
	config Config
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
	Environment Environment
	Logger      LoggerConfig
	Database    DatabaseConfig
	Api         ApiConfig
	Import      []ImportChannelConfig
}

// GetConfig returns static config instance.
// If Parse() was not called before, returns empty instance.
func GetConfig() *Config {
	return &config
}

// InitConfig parses configuration file based on flags and fills static config instance.
// Also causes flags being parsed if it was not parsed before.
// Returns error if parsing failed.
// Closes application if parse flags failed or was given flag -help.
func InitConfig() error {
	if !flag.Parsed() {
		parseFlags()
	}
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, &config)
}
