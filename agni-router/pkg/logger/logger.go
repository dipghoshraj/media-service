package logger

import (
	"log/slog"
	"os"
)

// Logger is the shared structured JSON logger for agni-router.
// All log lines carry a "component" field for easy filtering.
var Logger = slog.New(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
).With("component", "agni-router")
