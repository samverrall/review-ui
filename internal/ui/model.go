package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samverrall/review-ui/internal/diff"
	"github.com/samverrall/review-ui/internal/git"
)

type model struct {
	changedFiles   []string            // All changed files
	currentIndex   int                 // Current file index
	diffs          map[string]string   // Cached formatted diffs
	viewport       viewport.Model      // Scrollable viewport
	ready          bool                // Terminal size known
	width          int                 // Terminal width
	height         int                 // Terminal height
	err            error               // Error state
	cursorLine     int                 // Current cursor line position
	commentInput   textinput.Model     // Text input for comments
	commentMode    bool                // Whether we're in comment input mode
	comments       map[string][]string // Comments by "filename:lineNumber" or "filename:startLine-endLine"
	commentLine    int                 // Line number where comment is being added (or start of range)
	commentEndLine int                 // End line of comment range (-1 for single line)
	selectionMode  bool                // Whether we're in visual selection mode
	selectionStart int                 // Start line of selection
	statusMessage  string              // Status message to display to user
	fileListMode   bool                // Whether we're in file list selection mode
	fileListCursor int                 // Current cursor position in file list
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

	// Initialize comment input
	ti := textinput.New()
	ti.Placeholder = "Enter your comment..."
	ti.CharLimit = 200
	ti.Width = 80

	m := model{
		changedFiles: files,
		currentIndex: 0,
		diffs:        make(map[string]string),
		viewport:     viewport.New(0, 0),
		commentInput: ti,
		commentMode:  false,
		comments:     make(map[string][]string),
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

// Width returns the current terminal width
func (m model) Width() int {
	return m.width
}

// getCommentKey generates a unique key for storing comments
func (m *model) getCommentKey(lineNumber int) string {
	if m.currentIndex < 0 || m.currentIndex >= len(m.changedFiles) {
		return ""
	}
	return fmt.Sprintf("%s:%d", m.changedFiles[m.currentIndex], lineNumber)
}

// getCommentKeyForRange generates a unique key for storing range comments
func (m *model) getCommentKeyForRange(startLine, endLine int) string {
	if m.currentIndex < 0 || m.currentIndex >= len(m.changedFiles) {
		return ""
	}
	// Ensure start <= end
	if startLine > endLine {
		startLine, endLine = endLine, startLine
	}
	return fmt.Sprintf("%s:%d-%d", m.changedFiles[m.currentIndex], startLine, endLine)
}

// getSelectionRange returns the start and end lines of the current selection (ordered)
func (m *model) getSelectionRange() (int, int) {
	start, end := m.selectionStart, m.cursorLine
	if start > end {
		return end, start
	}
	return start, end
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
		m.cursorLine = 0
		m.selectionMode = false
		return nil
	}

	// Fetch diff from git
	rawDiff, err := git.GetFileDiff(filename)
	if err != nil {
		return fmt.Errorf("failed to load diff for %s: %w", filename, err)
	}

	// Format the diff with colors
	formattedDiff := diff.FormatDiff(m.width, rawDiff)
	m.diffs[filename] = formattedDiff

	// Update viewport content
	m.viewport.SetContent(formattedDiff)
	m.viewport.GotoTop()
	m.cursorLine = 0
	m.selectionMode = false

	return nil
}

// exportComments formats all comments for export
func (m *model) exportComments() string {
	if len(m.comments) == 0 {
		return "No comments to export."
	}

	var builder strings.Builder
	builder.WriteString("# Code Review Comments\n")
	builder.WriteString(fmt.Sprintf("# Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Sort keys for consistent output
	keys := make([]string, 0, len(m.comments))
	for key := range m.comments {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Group by file
	currentFile := ""
	for _, key := range keys {
		comments := m.comments[key]

		// Parse the key: "filename:lineNumber" or "filename:start-end"
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			continue
		}

		filename := parts[0]
		location := parts[1]

		// Add file header if we've moved to a new file
		if filename != currentFile {
			if currentFile != "" {
				builder.WriteString("\n")
			}
			builder.WriteString(fmt.Sprintf("## File: %s\n\n", filename))
			currentFile = filename
		}

		// Format the location
		var locationStr string
		if strings.Contains(location, "-") {
			// Range comment
			var start, end int
			if _, err := fmt.Sscanf(location, "%d-%d", &start, &end); err == nil {
				locationStr = fmt.Sprintf("Lines %d-%d", start+1, end+1)
			}
		} else {
			// Single line comment
			var lineNum int
			if _, err := fmt.Sscanf(location, "%d", &lineNum); err == nil {
				locationStr = fmt.Sprintf("Line %d", lineNum+1)
			}
		}

		// Add comments for this location
		builder.WriteString(fmt.Sprintf("### %s\n", locationStr))
		for _, comment := range comments {
			builder.WriteString(fmt.Sprintf("- %s\n", comment))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// saveCommentsToFile saves all comments to a file
func (m *model) saveCommentsToFile() error {
	content := m.exportComments()
	if content == "No comments to export." {
		return fmt.Errorf("no comments to save")
	}

	filename := fmt.Sprintf("code-review-comments-%s.md", time.Now().Format("20060102-150405"))

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	m.statusMessage = fmt.Sprintf("💾 Saved to %s", filename)
	return nil
}

// copyCommentsToClipboard copies all comments to the clipboard
func (m *model) copyCommentsToClipboard() error {
	content := m.exportComments()
	if content == "No comments to export." {
		return fmt.Errorf("no comments to copy")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	m.statusMessage = "📋 Copied to clipboard"
	return nil
}
