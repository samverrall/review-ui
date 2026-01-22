package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all incoming messages and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize
		// Account for header (4 lines), footer (3 lines), modal padding (2 lines), and buffer (2 lines)
		headerHeight := 4
		footerHeight := 3
		modalPaddingHeight := 2
		bufferHeight := 2
		verticalMarginHeight := headerHeight + footerHeight + modalPaddingHeight + bufferHeight

		if !m.ready {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Handle comment input mode separately
		if m.commentMode {
			switch msg.String() {
			case "enter":
				// Save comment
				commentText := m.commentInput.Value()
				if commentText != "" {
					var key string
					// Check if this is a range comment or single line comment
					if m.commentEndLine >= 0 && m.commentEndLine != m.commentLine {
						// Range comment
						key = m.getCommentKeyForRange(m.commentLine, m.commentEndLine)
					} else {
						// Single line comment
						key = m.getCommentKey(m.commentLine)
					}
					if key != "" {
						m.comments[key] = append(m.comments[key], commentText)
					}
				}
				// Exit comment mode
				m.commentMode = false
				m.commentInput.Reset()
				m.commentEndLine = -1
				return m, nil

			case "esc":
				// Cancel comment
				m.commentMode = false
				m.commentInput.Reset()
				return m, nil

			default:
				// Pass keys to text input
				m.commentInput, cmd = m.commentInput.Update(msg)
				return m, cmd
			}
		}

		// File list mode handlers
		if m.fileListMode {
			switch msg.String() {
			case "enter":
				// Select the file at fileListCursor
				m.currentIndex = m.fileListCursor
				if err := m.loadDiff(m.currentIndex); err != nil {
					m.err = err
				}
				m.fileListMode = false
				return m, nil

			case "esc":
				// Exit file list mode
				m.fileListMode = false
				return m, nil

			case "j":
				// Move cursor down in file list
				if m.fileListCursor < len(m.changedFiles)-1 {
					m.fileListCursor++
				}
				return m, nil

			case "k":
				// Move cursor up in file list
				if m.fileListCursor > 0 {
					m.fileListCursor--
				}
				return m, nil
			}
		}

		// Normal mode key handlers
		switch msg.String() {
		case "q", "ctrl+c":
			// Quit the application
			return m, tea.Quit

		case "tab":
			// Enter file list mode
			m.fileListMode = true
			m.fileListCursor = m.currentIndex
			return m, nil

		case "c":
			// Open comment input at current cursor line or selection
			m.commentMode = true
			if m.selectionMode {
				// Get the selection range
				start, end := m.getSelectionRange()
				m.commentLine = start
				m.commentEndLine = end
				// Exit selection mode after starting comment
				m.selectionMode = false
			} else {
				// Single line comment
				m.commentLine = m.cursorLine
				m.commentEndLine = -1
			}
			m.commentInput.Focus()
			return m, textinput.Blink

		case "v":
			// Toggle visual selection mode
			if !m.selectionMode {
				// Enter selection mode - set selection start to current cursor
				m.selectionMode = true
				m.selectionStart = m.cursorLine
			} else {
				// Exit selection mode
				m.selectionMode = false
			}
			return m, nil

		case "esc":
			// Exit selection mode if active
			if m.selectionMode {
				m.selectionMode = false
				return m, nil
			}

		case "s":
			// Save comments to file
			m.statusMessage = "" // Clear previous status
			if err := m.saveCommentsToFile(); err != nil {
				m.statusMessage = fmt.Sprintf("✗ Error: %v", err)
			}
			return m, nil

		case "y":
			// Copy comments to clipboard
			m.statusMessage = "" // Clear previous status
			if err := m.copyCommentsToClipboard(); err != nil {
				m.statusMessage = fmt.Sprintf("✗ Error: %v", err)
			}
			return m, nil

		case "n":
			// Next file
			m.statusMessage = "" // Clear status message
			if len(m.changedFiles) > 0 {
				m.currentIndex = (m.currentIndex + 1) % len(m.changedFiles)
				if err := m.loadDiff(m.currentIndex); err != nil {
					m.err = err
				}
			}
			return m, nil

		case "p":
			// Previous file
			m.statusMessage = "" // Clear status message
			if len(m.changedFiles) > 0 {
				m.currentIndex = (m.currentIndex - 1 + len(m.changedFiles)) % len(m.changedFiles)
				if err := m.loadDiff(m.currentIndex); err != nil {
					m.err = err
				}
			}
			return m, nil

		case "j":
			// Move cursor down
			totalLines := m.viewport.TotalLineCount()
			if m.cursorLine < totalLines-1 {
				m.cursorLine++
				// Auto-scroll viewport if cursor goes below visible area
				if m.cursorLine >= m.viewport.YOffset+m.viewport.Height {
					m.viewport.ScrollDown(1)
				}
			}

		case "k":
			// Move cursor up
			if m.cursorLine > 0 {
				m.cursorLine--
				// Auto-scroll viewport if cursor goes above visible area
				if m.cursorLine < m.viewport.YOffset {
					m.viewport.ScrollUp(1)
				}
			}

		default:
			// Pass other keys to viewport for scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}
