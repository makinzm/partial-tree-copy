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

	// Build the selection view (right panel)
	rightView := m.buildSelectionView(maxLines)

	// Set fixed width for the tree view
	treeViewWidth := 50 // Adjust as needed

	// Apply styles to panels
	leftStyle := lipgloss.NewStyle().Width(treeViewWidth)
	rightStyle := lipgloss.NewStyle()

	// Join panels horizontally
	combinedView := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftView),
		rightStyle.Render(rightView))

	// Add help text at bottom
	helpText := "\nHow to use\n" +
		"Press 'w'/Ctrl+'c' to quit and copy, 'Space' to select file, 'Enter' to expand/collapse dir\n" +
		"Navigation: 'h'/'l' to switch panels, 'j'/'k' to move up/down, 'J'/'K' to jump between directories"

	return combinedView + helpText
}
