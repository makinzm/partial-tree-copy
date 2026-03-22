package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the model as a string
func (m Model) View() string {
	// Number of lines to display in each panel
	maxLines := m.MaxVisibleRows - 4 // Reserve 4 rows for help text

	// Build the tree view (left panel)
	leftView := m.buildTreeView(maxLines)

	// Build the right panel (selection or preview)
	var rightView string
	if m.PreviewMode {
		rightView = m.buildPreviewView(maxLines)
	} else {
		rightView = m.buildSelectionView(maxLines)
	}

	// Set fixed width for the tree view (wider right panel in preview mode)
	treeViewWidth := 50 // Adjust as needed
	rightWidth := 40
	if m.PreviewMode {
		rightWidth = 80
	}

	// Apply styles to panels
	leftStyle := lipgloss.NewStyle().Width(treeViewWidth)
	rightStyle := lipgloss.NewStyle().Width(rightWidth)

	// Join panels horizontally
	combinedView := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftView),
		rightStyle.Render(rightView))

	// Add help text at bottom
	helpText := "\nHow to use\n" +
		"Press 'w'/Ctrl+'c' to quit and copy, 'Space' to select file, 'Enter' to expand/collapse dir\n" +
		"Navigation: 'h'/'l' to switch panels, 'j'/'k' to move up/down, 'J'/'K' to jump between directories\n" +
		"'p' to toggle file preview"

	return combinedView + helpText
}
