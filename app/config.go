package app

import (
	"flag"

	ini "gopkg.in/ini.v1"
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
type Config struct {
	Environment Environment
	Logger      LoggerConfig
	Database    DatabaseConfig
	Api         ApiConfig
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
	return ini.MapTo(&config, configPath)
}
