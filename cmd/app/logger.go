package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// setupLogger configures slog based on debug mode
func setupLogger(debug bool) *slog.Logger {
	if !debug {
		// Discard all logs when not in debug mode
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	// Write logs to debug.log file
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to open debug.log: %v\n", err)
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(logFile, opts)
	return slog.New(handler)
}
