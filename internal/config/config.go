package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	HTTP_port     int    `env:"HTTP_PORT" envDefault:"3000"`
	IsProd        bool   `env:"IS_PROD" envDefault:"false"`
	Salt          string `env:"SALT" envDefault:"b9PDPbt4"`
	Secret        string `env:"SECRET" envDefault:"secret"`
	Level_str     string `env:"LEVEL" envDefault:"debug"`
	Level         zapcore.Level
	DSN           string `env:"DSN" envDefault:"https://b90ccf5ce9514de88ab2166dd3696827@o4503908933566464.ingest.sentry.io/4503908956766208"`
	ProbesPort    string `env:"PROBES_PORT" envDefault:"3030"`
	ServiceName   string `env:"SERVICE_NAME" envDefault:"auth"`
	JaegerAddress string `env:"JAEGER_ADDRESS" envDefault:"http://jaeger-instance-collector.observability:14268/api/traces"`
	JaegerPort    string `env:"JAEGER_PORT" envDefault:"9000"`
	DB_url        string `env:"DB_URL" envDefault:"mongodb://localhost:27017/"`
	DB_used       bool   `env:"DB_USED" envDefault:"false"`
}

var config Config = Config{}

func GetLevel(level string) zapcore.Level {
	switch level {
	case "debug", "dbg", "all":
		return zapcore.DebugLevel
	case "info", "inf", "": // make the zero value useful
		return zapcore.InfoLevel
	case "error", "err":
		return zapcore.ErrorLevel
	}
	return zapcore.DebugLevel
}

func GetConfig() (*Config, error) {
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("read logger configuration failed: %w", err)
	}
	config.Level = GetLevel(config.Level_str)
	return &config, nil
}
