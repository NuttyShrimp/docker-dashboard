// Package logger initiates a zap logger
package logger

import (
	"errors"
	"fmt"
	"os"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/nuttyshrimp/docker-dashboard/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	File    string
	Console bool
	Level   *zapcore.Level
}

// LocalLogger is the raw logger before zapsentry core is attached.
// This is useful for logging errors locally without sending duplicates to Sentry.
var LocalLogger *zap.Logger

func New(logCfg Config) (*zap.Logger, error) {
	if logCfg.File != "" {
		// nolint:gosec // I want to set my own permissions
		err := os.Mkdir("logs", 0o755)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("create logs directory %w", err)
		}
	}

	outputPaths := []string{}
	errorOutputPaths := []string{}
	if logCfg.File != "" {
		outputPaths = append(outputPaths, fmt.Sprintf("logs/%s.log", logCfg.File))
		errorOutputPaths = append(errorOutputPaths, fmt.Sprintf("logs/%s.log", logCfg.File))
	}
	if logCfg.Console {
		outputPaths = append(outputPaths, "stdout")
		errorOutputPaths = append(errorOutputPaths, "stderr")
	}

	if len(outputPaths) == 0 {
		return nil, errors.New("no output paths specified")
	}

	var cfg zap.Config

	if config.IsDev() {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		level := zap.DebugLevel
		if logCfg.Level != nil {
			level = *logCfg.Level
		}
		cfg.Level.SetLevel(level)
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		level := zap.WarnLevel
		if logCfg.Level != nil {
			level = *logCfg.Level
		}
		cfg.Level.SetLevel(level)
	}

	cfg.OutputPaths = outputPaths
	cfg.ErrorOutputPaths = errorOutputPaths
	cfg.DisableStacktrace = true

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	LocalLogger = logger

	if dsn := config.GetString("app.dsn"); dsn != "" {
		zapCfg := zapsentry.Configuration{
			EnableBreadcrumbs: true,
			BreadcrumbLevel:   zapcore.InfoLevel,
			Level:             zapcore.ErrorLevel,
			Tags: map[string]string{
				"component": config.GetString("app.name"),
			},
		}

		core, err := zapsentry.NewCore(zapCfg, zapsentry.NewSentryClientFromClient(sentry.CurrentHub().Client()))
		if err != nil {
			return nil, fmt.Errorf("failed to initialize zap core %w", err)
		}

		logger = zapsentry.AttachCoreToLogger(core, logger)
	}

	return logger, nil
}
