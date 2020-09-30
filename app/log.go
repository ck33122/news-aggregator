package app

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

// GetLog returns static logger instance.
// If log.Init() was not called, returns null pointer.
// Shouldn't be called inside init().
func GetLog() *zap.Logger {
	return logger
}

// InitLog creates logger instance based on config.
// After this function call, GetLog() will return configured instance.
// If failed, returns error.
// name - name of application which will be used when logging.
func InitLog() error {
	// TODO kibana/graylog

	var loggerConfig zap.Config

	switch config.Environment {
	case EnvProduction:
		loggerConfig = zap.NewProductionConfig()
	case EnvDevelopment:
		loggerConfig = zap.NewDevelopmentConfig()
	default:
		return fmt.Errorf("InitLog failed: unknown environment '%s'", config.Environment)
	}

	if len(config.Logger.Dir) > 0 {
		fileName := fmt.Sprintf("%s/%s.log", config.Logger.Dir, appName)
		loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, fileName)
	}

	var err error
	logger, err = loggerConfig.Build()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
		return err
	}
	logger = logger.Named(appName)
	logger.Info("application start", zap.String("appName", appName))

	return nil
}

// DestroyLog flushes log buffers.
func DestroyLog() {
	_ = logger.Sync() // sync error ignored, see https://github.com/uber-go/zap/issues/328
}
