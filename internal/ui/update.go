package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all incoming messages and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize
		headerHeight := 1
		footerHeight := 1
		verticalMarginHeight := headerHeight + footerHeight

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
		switch msg.String() {
		case "q", "ctrl+c":
			// Quit the application
			return m, tea.Quit

		case "n":
			// Next file
			if len(m.changedFiles) > 0 {
				m.currentIndex = (m.currentIndex + 1) % len(m.changedFiles)
				if err := m.loadDiff(m.currentIndex); err != nil {
					m.err = err
				}
			}

		case "p":
			// Previous file
			if len(m.changedFiles) > 0 {
				m.currentIndex = (m.currentIndex - 1 + len(m.changedFiles)) % len(m.changedFiles)
				if err := m.loadDiff(m.currentIndex); err != nil {
					m.err = err
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
