package app

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLog(cnf *Config) (*zap.Logger, error) {
	var loggerConfig zap.Config

	switch cnf.Environment {
	case EnvProduction:
		loggerConfig = zap.NewProductionConfig()
	case EnvDevelopment:
		loggerConfig = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("InitLog failed: unknown environment '%s'", cnf.Environment)
	}

	if len(cnf.Logger.Dir) > 0 {
		fileName := fmt.Sprintf("%s/%s.log", cnf.Logger.Dir, cnf.AppName)
		loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, fileName)
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("cannot initialize zap logger: %v", err)
	}
	logger = logger.Named(cnf.AppName)
	logger.Info("application start", zap.String("appName", cnf.AppName))

	return logger, nil
}
