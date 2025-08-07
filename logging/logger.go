package logging

import (
	"log/slog"
	"os"
)

func parseLogLevel(s string) slog.Level {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	if err != nil {
		slog.Error("Invalid log level, defaulting to WARN", "level", s, "error", err)
		return slog.LevelWarn
	}
	return level
}

func SetupLogger(p string) *slog.Logger {
	level := os.Getenv("L337_LOG_LEVEL")
	if level == "" {
		level = "WARN"
	}
	logLevel := parseLogLevel(level)
	options := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	}
	handler := slog.NewTextHandler(os.Stdout, &options)
	logger := slog.New(handler).With("package", p)
	return logger
}
