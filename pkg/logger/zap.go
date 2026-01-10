// Package logger initiates a zap logger
package logger

import (
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
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

	logger, err := cfg.Build(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		return nil, err
	}

	env := config.GetDefaultString("app.env", "development")
	logger = logger.With(zap.String("env", env))

	return logger, nil
}
