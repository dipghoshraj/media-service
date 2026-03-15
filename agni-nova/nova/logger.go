package nova

import (
	"log/slog"
	"os"
)

// Logger is the shared structured JSON logger for the nova package.
// All log output includes a "component" field set to "agni-nova" for
// easy filtering in log aggregators.
var Logger = slog.New(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
).With("component", "agni-nova")
