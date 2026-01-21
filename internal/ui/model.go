package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samverrall/text-editor/internal/diff"
	"github.com/samverrall/text-editor/internal/git"
)

type model struct {
	changedFiles []string          // All changed files
	currentIndex int               // Current file index
	diffs        map[string]string // Cached formatted diffs
	viewport     viewport.Model    // Scrollable viewport
	ready        bool              // Terminal size known
	width        int               // Terminal width
	height       int               // Terminal height
	err          error             // Error state
}

// New creates and initializes a new model
func New() (model, error) {
	// Check if we're in a git repository
	isRepo, err := git.IsGitRepo()
	if err != nil {
		return model{}, fmt.Errorf("failed to check git repository: %w", err)
	}
	if !isRepo {
		return model{}, fmt.Errorf("not a git repository")
	}

	// Get changed files
	files, err := git.GetChangedFiles()
	if err != nil {
		return model{}, fmt.Errorf("failed to get changed files: %w", err)
	}

	m := model{
		changedFiles: files,
		currentIndex: 0,
		diffs:        make(map[string]string),
		viewport:     viewport.New(0, 0),
	}

	// Load first diff if we have files
	if len(files) > 0 {
		if err := m.loadDiff(0); err != nil {
			return model{}, err
		}
	}

	return m, nil
}

// Init initializes the model (required by Bubbletea)
func (m model) Init() tea.Cmd {
	return nil
}

// loadDiff loads and caches the diff for the file at the given index
func (m *model) loadDiff(index int) error {
	if index < 0 || index >= len(m.changedFiles) {
		return nil
	}

	filename := m.changedFiles[index]

	// Check if already cached
	if _, exists := m.diffs[filename]; exists {
		m.viewport.SetContent(m.diffs[filename])
		return nil
	}

	// Fetch diff from git
	rawDiff, err := git.GetFileDiff(filename)
	if err != nil {
		return fmt.Errorf("failed to load diff for %s: %w", filename, err)
	}

	// Format the diff with colors
	formattedDiff := diff.FormatDiff(rawDiff)
	m.diffs[filename] = formattedDiff

	// Update viewport content
	m.viewport.SetContent(formattedDiff)
	m.viewport.GotoTop()

	return nil
}
