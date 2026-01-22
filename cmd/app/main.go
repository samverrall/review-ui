package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samverrall/review-ui/internal/ui"
)

func main() {
	// Parse command line flags
	debug := flag.Bool("debug", false, "enable debug logging to debug.log")
	flag.Parse()

	// Set up logger
	logger := setupLogger(*debug)
	if *debug {
		logger.Info("debug mode enabled")
	}

	// Create the model
	m, err := ui.NewWithLogger(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create the program with options
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support for scrolling
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
