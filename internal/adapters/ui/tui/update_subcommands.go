package tui

import (
	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
)

// GetVisibleNodes returns all currently visible nodes based on expansion state
func (m *Model) GetVisibleNodes() []*entities.FileNode {
	return m.Navigator.GetVisibleNodes(m.Root)
}

// GetAllSelectedNodes returns all selected nodes
func (m *Model) GetAllSelectedNodes() []*entities.FileNode {
	return m.Selector.GetSelectedNodes()
}

// MoveCursorUp moves the cursor up in the tree view
func (m *Model) MoveCursorUp() {
	visibleNodes := m.GetVisibleNodes()

	// Find current cursor index
	index := -1
	for i, node := range visibleNodes {
		if node == m.Cursor {
			index = i
			break
		}
	}

	// Move up if possible
	if index > 0 {
		m.Cursor = visibleNodes[index-1]
	}
}

// MoveCursorDown moves the cursor down in the tree view
func (m *Model) MoveCursorDown() {
	visibleNodes := m.GetVisibleNodes()

	// Find current cursor index
	index := -1
	for i, node := range visibleNodes {
		if node == m.Cursor {
			index = i
			break
		}
	}

	// Move down if possible
	if index >= 0 && index < len(visibleNodes)-1 {
		m.Cursor = visibleNodes[index+1]
	}
}

// ToggleExpand toggles expansion state of current directory
func (m *Model) ToggleExpand() {
	m.Navigator.ToggleExpand(m.Cursor)
}

// ToggleSelect toggles selection state of current file
func (m *Model) ToggleSelect() {
	m.Selector.ToggleSelect(m.Cursor)
}

// CopySelection copies all selected files to clipboard
func (m *Model) CopySelection() error {
	selection := m.Selector.GetSelection()
	return m.Copier.CopySelectionToClipboard(selection)
}

// MoveToPreviousDirectory moves to the previous directory in the tree
func (m *Model) MoveToPreviousDirectory() {
	visibleNodes := m.GetVisibleNodes()
	m.Cursor = m.Navigator.MoveToPreviousDirectory(visibleNodes, m.Cursor)
}

// MoveToNextDirectory moves to the next directory in the tree
func (m *Model) MoveToNextDirectory() {
	visibleNodes := m.GetVisibleNodes()
	m.Cursor = m.Navigator.MoveToNextDirectory(visibleNodes, m.Cursor)
}

// GetBreadcrumbs returns the path from root to current cursor
func (m *Model) GetBreadcrumbs() []*entities.FileNode {
	return m.Navigator.GetBreadcrumbs(m.Cursor)
}

// GetNodeLevel returns the depth level of a node
func (m *Model) GetNodeLevel(node *entities.FileNode) int {
	return m.Navigator.GetNodeLevel(node)
}
