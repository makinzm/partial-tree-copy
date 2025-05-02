package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/makinzm/partial-tree-copy/internal/adapters/ui/tui"
	"github.com/makinzm/partial-tree-copy/internal/usecases/copier"
	"github.com/makinzm/partial-tree-copy/internal/usecases/navigator"
	"github.com/makinzm/partial-tree-copy/internal/usecases/selector"
)

// UIPresenter is responsible for handling the presentation layer
type UIPresenter struct {
	navigator *navigator.FileNavigator
	selector  *selector.FileSelector
	copier    *copier.FileCopier
}

// NewUIPresenter creates a new UIPresenter
func NewUIPresenter(
	navigator *navigator.FileNavigator,
	selector *selector.FileSelector,
	copier *copier.FileCopier,
) *UIPresenter {
	return &UIPresenter{
		navigator: navigator,
		selector:  selector,
		copier:    copier,
	}
}

// StartUI starts the terminal UI
func (p *UIPresenter) StartUI() error {
	// Create the TUI model
	model, err := tui.NewModel(
		p.navigator,
		p.selector,
		p.copier,
		20, // Maximum visible rows
	)
	if err != nil {
		return fmt.Errorf("failed to create UI model: %w", err)
	}

	// Initialize BubbleTea program
	program := tea.NewProgram(*model)

	// Start the program
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("error running UI: %w", err)
	}

	return nil
}

// HandleError presents an error message and exits the application
func (p *UIPresenter) HandleError(err error, message string) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
