// Package logger initiates a zap logger
package logger

import (
	"fmt"
	"os"

	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	File    string
	Console bool
}

func New(logCfg Config) (*zap.Logger, error) {
	err := os.Mkdir("logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("create logs directory %w", err)
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
		return nil, fmt.Errorf("no output paths specified")
	}

	var cfg zap.Config

	if config.IsDev() {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zap.WarnLevel)
	}

	cfg.OutputPaths = outputPaths
	cfg.ErrorOutputPaths = errorOutputPaths

	logger, err := cfg.Build(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		return nil, err
	}

	env := config.GetDefaultString("app.env", "development")
	logger = logger.With(zap.String("env", env))

	return logger, nil
}
