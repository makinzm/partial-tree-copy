package tui

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
)

// buildTreeView constructs the tree view (left panel)
func (m *Model) buildTreeView(maxLines int) string {
	// Get all visible nodes
	visibleNodes := m.GetVisibleNodes()

	// Find cursor position
	cursorIdx := -1
	for i, node := range visibleNodes {
		if node == m.Cursor {
			cursorIdx = i
			break
		}
	}

	// Determine display range
	startIdx := 0
	endIdx := len(visibleNodes)

	if cursorIdx >= 0 && len(visibleNodes) > maxLines {
		// Center the cursor in view
		halfHeight := maxLines / 2

		if cursorIdx > halfHeight {
			startIdx = cursorIdx - halfHeight
		}

		if startIdx+maxLines > len(visibleNodes) {
			startIdx = len(visibleNodes) - maxLines
		}

		if startIdx < 0 {
			startIdx = 0
		}

		endIdx = startIdx + maxLines
		if endIdx > len(visibleNodes) {
			endIdx = len(visibleNodes)
		}
	}

	// Build the view string
	var s strings.Builder

	// Add breadcrumb path at top
	s.WriteString("Path: " + m.formatBreadcrumbs() + "\n\n")

	// Indicate if there are hidden nodes above
	if startIdx > 0 {
		s.WriteString("...\n")
	}

	// Render visible nodes
	for i := startIdx; i < endIdx; i++ {
		node := visibleNodes[i]
		level := m.GetNodeLevel(node)
		s.WriteString(m.renderSingleNode(node, level, node == m.Cursor, !m.FocusRight))
	}

	// Indicate if there are hidden nodes below
	if endIdx < len(visibleNodes) {
		s.WriteString("...\n")
	}

	return s.String()
}

// buildSelectionView constructs the selection view (right panel)
func (m *Model) buildSelectionView(maxLines int) string {
	var s strings.Builder

	// Add title with selection count
	s.WriteString("Selected Files (" + strconv.Itoa(len(m.Selector.GetSelection())) + "):\n\n")

	// Show message if no files are selected
	if len(m.Selector.GetSelection()) == 0 {
		s.WriteString("No files selected\n")
		return s.String()
	}

	// Get all selected nodes
	selectedNodes := m.GetAllSelectedNodes()

	// Determine display range based on scroll position
	startIdx := m.RightScroll
	if len(selectedNodes) > 0 && startIdx >= len(selectedNodes) {
		startIdx = len(selectedNodes) - 1
	}

	// Calculate visible range
	visibleCount := maxLines - 2 // Subtract title and empty line
	endIdx := startIdx + visibleCount
	if endIdx > len(selectedNodes) {
		endIdx = len(selectedNodes)
	}

	// Indicate if there are hidden nodes above
	if startIdx > 0 {
		s.WriteString("...\n")
	}

	// Render selected files
	for i := 0; i < len(selectedNodes); i++ {
		// Skip nodes outside visible range
		if i < startIdx || i >= endIdx {
			continue
		}

		node := selectedNodes[i]

		// Create relative path
		relPath := m.getRelativePath(node.Path, m.Root.Path)

		// Format line with index
		numStr := strconv.Itoa(i + 1)
		padding := " "
		if i+1 >= 10 { // Adjust padding for double-digit numbers
			padding = ""
		}

		line := numStr + "." + padding + relPath

		// Highlight current scroll position if right panel is focused
		if m.FocusRight && i == m.RightScroll {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("> " + line)
		} else {
			line = "  " + line
		}

		s.WriteString(line + "\n")
	}

	// Indicate if there are hidden nodes below
	if endIdx < len(selectedNodes) {
		s.WriteString("...\n")
	}

	return s.String()
}

// renderSingleNode renders a single node for the tree view
func (m *Model) renderSingleNode(node *entities.FileNode, level int, isCursor bool, isFocused bool) string {
	prefix := strings.Repeat("  ", level)
	line := prefix

	// Add cursor indicator with proper coloring
	if isCursor {
		if isFocused {
			// Color the cursor if this panel has focus
			line += lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("> ")
		} else {
			line += "> "
		}
	} else {
		line += "  "
	}

	// Render node based on type and state
	if node.IsDir {
		if node.Expanded {
			line += "📂 " + filepath.Base(node.Path) + "\n"
		} else {
			line += "📁 " + filepath.Base(node.Path) + "\n"
		}
	} else {
		if node.Selected {
			line += "[✓] " + filepath.Base(node.Path) + "\n"
		} else {
			line += "[ ] " + filepath.Base(node.Path) + "\n"
		}
	}

	return line
}

// formatBreadcrumbs creates a breadcrumb navigation string
func (m *Model) formatBreadcrumbs() string {
	breadcrumbs := m.GetBreadcrumbs()
	if len(breadcrumbs) == 0 {
		return ""
	}

	var result strings.Builder

	for i, node := range breadcrumbs {
		// Special handling for root node
		if i == 0 {
			result.WriteString("/")
			continue
		}

		// Highlight current (last) node
		if i == len(breadcrumbs)-1 {
			result.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Render(node.Name))
		} else {
			result.WriteString(node.Name)
		}

		// Add separator between nodes
		if i < len(breadcrumbs)-1 {
			result.WriteString("/")
		}
	}

	return result.String()
}

// getRelativePath returns the path relative to a base path
func (m *Model) getRelativePath(path, rootPath string) string {
	// If path starts with root path, return relative path
	if strings.HasPrefix(path, rootPath) {
		relPath, err := filepath.Rel(rootPath, path)
		if err == nil {
			return relPath
		}
	}

	// Default to filename if relative path can't be determined
	return filepath.Base(path)
}
