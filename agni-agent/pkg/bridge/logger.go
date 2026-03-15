package bridge

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger is the shared structured JSON logger for agni-agent.
// All log lines carry a "component" field for easy filtering.
var Logger = slog.New(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
).With("component", "agni-agent")

// FormatError builds a prefixed error string, kept for compatibility.
func FormatError(message string, err error) error {
	if err == nil {
		return fmt.Errorf("[Agni Agent] %s ", message)
	}
	return fmt.Errorf("[Agni Agent] %s -- %v", message, err)
}
