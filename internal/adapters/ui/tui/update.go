package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles user input and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "w":
			// Copy selection and quit
			m.CopySelection()
			return m, tea.Quit

		case "p":
			// Toggle preview mode
			m.PreviewMode = !m.PreviewMode
			if m.PreviewMode {
				m.LoadPreviewContent()
			}

		case "L", "l":
			// Move focus to right panel if there are selections or preview mode
			if m.PreviewMode || len(m.Selector.GetSelection()) > 0 {
				m.FocusRight = true
			}

		case "H", "h":
			// Move focus to left panel
			m.FocusRight = false

		case "up", "k":
			if m.FocusRight {
				if m.PreviewMode {
					// Scroll up in preview
					if m.PreviewScroll > 0 {
						m.PreviewScroll--
					}
				} else {
					// Scroll up in right panel
					if m.RightScroll > 0 {
						m.RightScroll--
					}
				}
			} else {
				// Move cursor up in left panel
				m.MoveCursorUp()
				if m.PreviewMode {
					m.LoadPreviewContent()
				}
			}

		case "down", "j":
			if m.FocusRight {
				if m.PreviewMode {
					// Scroll down in preview
					lines := strings.Split(m.PreviewContent, "\n")
					if m.PreviewScroll < len(lines)-1 {
						m.PreviewScroll++
					}
				} else {
					// Scroll down in right panel
					selectedNodes := m.GetAllSelectedNodes()
					if m.RightScroll < len(selectedNodes)-1 {
						m.RightScroll++
					}
				}
			} else {
				// Move cursor down in left panel
				m.MoveCursorDown()
				if m.PreviewMode {
					m.LoadPreviewContent()
				}
			}

		case "K":
			if !m.FocusRight {
				// Move to previous directory
				m.MoveToPreviousDirectory()
				if m.PreviewMode {
					m.LoadPreviewContent()
				}
			}

		case "J":
			if !m.FocusRight {
				// Move to next directory
				m.MoveToNextDirectory()
				if m.PreviewMode {
					m.LoadPreviewContent()
				}
			}

		case "enter":
			if !m.FocusRight {
				// Toggle expand for directories, select for files
				if m.Cursor.IsDir {
					m.ToggleExpand()
				} else {
					m.ToggleSelect()
				}
			}

		case "space":
			if !m.FocusRight {
				// Toggle selection
				m.ToggleSelect()
			}
		}
	}

	return m, nil
}
