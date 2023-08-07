package logger

import (
	"auth/internal/config"
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = *zap.Logger

func initSentry(log *zap.Logger, sentryAddress, environment string, level zapcore.Level) *zap.Logger {
	if sentryAddress == "" {
		return log
	}
	cfg := zapsentry.Configuration{
		Level: level,
		Tags: map[string]string{
			"environment": environment,
			"app":         "auth",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryAddress))
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}
	return zapsentry.AttachCoreToLogger(core, log)
}

func New() (Logger, error) {
	cfg, _ := config.GetConfig()

	var zapCfg zap.Config
	if cfg.IsProd {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create logger failed: %w", err)
	}

	logger = initSentry(logger, cfg.DSN, "myenv", zapcore.Level(cfg.Level))
	defer func() {
		_ = logger.Sync()
	}()
	zap.ReplaceGlobals(logger)

	return logger, nil
}
