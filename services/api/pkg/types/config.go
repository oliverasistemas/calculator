package types

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr     string   `envconfig:"ADDR" default:":4000"`
	LogLevel LogLevel `envconfig:"LOG_LEVEL" default:"info"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load() // optional .env file

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}

// LogLevel wraps slog.Level so envconfig can decode it from a string.
type LogLevel slog.Level

func (l *LogLevel) Decode(value string) error {
	var level slog.Level
	if err := level.UnmarshalText([]byte(value)); err != nil {
		return fmt.Errorf("unsupported log level %q: %w", value, err)
	}
	*l = LogLevel(level)
	return nil
}
