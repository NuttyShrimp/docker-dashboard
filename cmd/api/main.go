package main

import (
	"fmt"

	"github.com/nuttyshrimp/docker-dashboard/internal/server"
	"github.com/nuttyshrimp/docker-dashboard/pkg/config"
	"github.com/nuttyshrimp/docker-dashboard/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	loggerFile := config.GetString("logger.file")
	loggerLevelStr := config.GetString("logger.level")
	var loggerLevel *zapcore.Level
	if loggerLevelStr != "" {
		loggerLevelTmp, err := zapcore.ParseLevel(loggerLevelStr)
		if err != nil {
			panic(fmt.Errorf("invalid logger level %s | %w", loggerLevelStr, err))
		}
		loggerLevel = &loggerLevelTmp
	}

	zapLogger, err := logger.New(logger.Config{
		Console: true,
		File:    loggerFile,
		Level:   loggerLevel,
	})
	if err != nil {
		panic(fmt.Errorf("zap logger initialization failed: %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

	api := server.New()

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}
