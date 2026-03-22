package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/usecases/copier"
	"github.com/makinzm/partial-tree-copy/internal/usecases/navigator"
	"github.com/makinzm/partial-tree-copy/internal/usecases/selector"
)

// Model represents the state of the file tree viewer
type Model struct {
	Root           *entities.FileNode // Root node of the file tree
	Cursor         *entities.FileNode // Current position of the cursor in the tree
	MaxVisibleRows int                // Maximum number of visible rows in the tree view
	FocusRight     bool               // Indicates if the right pane is focused
	RightScroll    int                // Scroll position of the right pane

	// Use cases
	Navigator *navigator.FileNavigator
	Selector  *selector.FileSelector
	Copier    *copier.FileCopier
}

// NewModel creates a new Model with the given use cases and settings
func NewModel(
	navigator *navigator.FileNavigator,
	selector *selector.FileSelector,
	copier *copier.FileCopier,
	maxVisibleRows int,
) (*Model, error) {
	// Build the root node
	rootNode, err := navigator.BuildRootNode()
	if err != nil {
		return nil, err
	}

	return &Model{
		Root:           rootNode,
		Cursor:         rootNode,
		MaxVisibleRows: maxVisibleRows,
		FocusRight:     false,
		RightScroll:    0,
		Navigator:      navigator,
		Selector:       selector,
		Copier:         copier,
	}, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}
