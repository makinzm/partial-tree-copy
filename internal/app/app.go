package app

import (
	"github.com/makinzm/partial-tree-copy/internal/adapters/repositories"
	"github.com/makinzm/partial-tree-copy/internal/adapters/ui"
	"github.com/makinzm/partial-tree-copy/internal/usecases/copier"
	"github.com/makinzm/partial-tree-copy/internal/usecases/navigator"
	"github.com/makinzm/partial-tree-copy/internal/usecases/selector"
)

// Application is the main application struct that wires everything together
type Application struct {
	presenter *ui.UIPresenter
}

// NewApplication creates and initializes a new Application
func NewApplication() (*Application, error) {
	// Initialize repository
	fileRepo := repositories.NewOSFileRepository()

	// Initialize use cases
	fileNavigator := navigator.NewFileNavigator(fileRepo)
	fileSelector := selector.NewFileSelector()
	fileCopier := copier.NewFileCopier(fileRepo)

	// Initialize UI presenter
	presenter := ui.NewUIPresenter(fileNavigator, fileSelector, fileCopier)

	return &Application{
		presenter: presenter,
	}, nil
}

// Run starts the application
func (app *Application) Run() error {
	return app.presenter.StartUI()
}
