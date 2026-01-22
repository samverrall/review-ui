package ui

import (
	"fmt"
	"strings"
)

// renderWithCursor highlights the cursor line, selection, and displays comments
func (m model) renderWithCursor() string {
	// Get all lines from the viewport's total content
	lines := strings.Split(m.viewport.View(), "\n")

	// Get selection range if in selection mode
	var selStart, selEnd int
	if m.selectionMode {
		selStart, selEnd = m.getSelectionRange()
	}

	// Build output with cursor/selection highlighting and comments
	var result []string
	processedRanges := make(map[string]bool)

	for i, line := range lines {
		// Calculate actual line number in the diff
		actualLineNumber := m.viewport.YOffset + i

		// Calculate cursor index in visible area
		cursorIndex := m.cursorLine - m.viewport.YOffset

		// Apply selection highlighting if in selection mode
		if m.selectionMode && actualLineNumber >= selStart && actualLineNumber <= selEnd {
			// Set width to fill the entire terminal width for consistency
			line = selectionStyle.Width(m.width).Render(line)
		} else if cursorIndex >= 0 && cursorIndex < len(lines) && i == cursorIndex {
			// Highlight cursor line if not in selection and add cursor indicator
			// Only highlight if cursor is actually visible

			// Create content with cursor indicator
			content := "▶ " + line

			// Apply cursor line style with full width
			// The key is to set the width BEFORE rendering so it fills the entire line
			line = cursorLineStyle.Width(m.width).Render(content)
		}

		result = append(result, line)

		// Check for single line comments
		key := m.getCommentKey(actualLineNumber)
		if comments, exists := m.comments[key]; exists {
			for _, comment := range comments {
				commentLine := commentStyle.Render(fmt.Sprintf("💬 %s", comment))
				result = append(result, commentLine)
			}
		}

		// Check for range comments that end at this line
		// We need to check all possible ranges that could include this line
		currentFile := ""
		if m.currentIndex >= 0 && m.currentIndex < len(m.changedFiles) {
			currentFile = m.changedFiles[m.currentIndex]
		}
		for rangeKey, comments := range m.comments {
			// Parse range keys (format: "filename:start-end")
			if strings.Contains(rangeKey, "-") && !processedRanges[rangeKey] {
				// Check if this comment belongs to the current file
				if !strings.HasPrefix(rangeKey, currentFile+":") {
					continue
				}
				// Extract the range from the key
				parts := strings.Split(rangeKey, ":")
				if len(parts) >= 2 {
					rangePart := parts[len(parts)-1]
					var start, end int
					if _, err := fmt.Sscanf(rangePart, "%d-%d", &start, &end); err == nil {
						// This is a range comment
						// Show it after the last line of the range
						if actualLineNumber == end {
							processedRanges[rangeKey] = true
							for _, comment := range comments {
								commentLine := commentStyle.Render(fmt.Sprintf("💬 [lines %d-%d] %s", start+1, end+1, comment))
								result = append(result, commentLine)
							}
						}
					}
				}
			}
		}
	}

	return strings.Join(result, "\n")
}

// renderFileList renders the file selection list
func (m model) renderFileList() string {
	var b strings.Builder

	// Header
	headerText := fmt.Sprintf("📂 Select File (%d files)", len(m.changedFiles))
	header := headerStyle.Render(headerText)
	if m.width > 0 {
		header = headerStyle.Width(m.width - 8).Render(headerText) // Account for modal padding
	}
	b.WriteString(header)
	b.WriteString("\n\n")

	// File list
	for i, file := range m.changedFiles {
		if i == m.fileListCursor {
			// Highlight the current selection
			line := fileListSelectedStyle.Render(fmt.Sprintf("  %s", file))
			b.WriteString(line)
		} else {
			line := fileListItemStyle.Render(fmt.Sprintf("  %s", file))
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	helpText := "↑↓ navigate | ↵ select | esc cancel"
	footer := footerStyle.Render(helpText)
	b.WriteString(footer)

	// Wrap in modal container
	return modalContainer.Render(b.String())
}

// View renders the current state of the model
func (m model) View() string {
	// Handle error state
	if m.err != nil {
		return modalContainer.Render(errorStyle.Render(fmt.Sprintf("❌ Error: %v\n\nPress q to quit.", m.err)))
	}

	// Handle no changes state
	if len(m.changedFiles) == 0 {
		return modalContainer.Render(infoStyle.Render("ℹ️  No unstaged changes found.\n\nPress q to quit."))
	}

	// Handle not ready state (terminal size not yet known)
	if !m.ready {
		return modalContainer.Render("Initializing...")
	}

	// Handle file list mode
	if m.fileListMode {
		return m.renderFileList()
	}

	// Build main view with modal-style centering
	var b strings.Builder

	// Header: File counter and name (prominent)
	currentFile := m.changedFiles[m.currentIndex]
	headerText := fmt.Sprintf("📄 File %d/%d: %s", m.currentIndex+1, len(m.changedFiles), currentFile)
	header := headerStyle.Width(m.width).Render(headerText)
	b.WriteString(header)
	b.WriteString("\n")

	// Viewport: Diff content with cursor highlighting
	b.WriteString(m.renderWithCursor())
	b.WriteString("\n")

	// Comment input area (if in comment mode)
	if m.commentMode {
		var commentPrompt string
		if m.commentEndLine >= 0 && m.commentEndLine != m.commentLine {
			// Range comment
			commentPrompt = fmt.Sprintf("💬 Adding comment to lines %d-%d:", m.commentLine+1, m.commentEndLine+1)
		} else {
			// Single line comment
			commentPrompt = fmt.Sprintf("💬 Adding comment to line %d:", m.commentLine+1)
		}
		inputArea := commentInputStyle.Render(
			fmt.Sprintf("%s\n%s", commentPrompt, m.commentInput.View()),
		)
		b.WriteString(inputArea)
		b.WriteString("\n")
	}

	// Status message (if present)
	if m.statusMessage != "" {
		statusLine := statusStyle.Render(m.statusMessage)
		b.WriteString(statusLine)
		b.WriteString("\n")
	}

	// Footer: Help text
	helpText := "tab files | n next | p prev | jk move | v select | c comment | s save | y copy | q quit"
	if m.commentMode {
		helpText = "↵ save | esc cancel"
	} else if m.selectionMode {
		helpText = "↑↓ extend selection | 💬 comment selection | v/esc exit selection"
	}
	footer := footerStyle.Width(m.width).Render(helpText)
	b.WriteString(footer)

	// Wrap everything in modal container for centered appearance
	content := b.String()
	return content
}
