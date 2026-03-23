package app

import (
	"github.com/makinzm/partial-tree-copy/internal/adapters/repositories"
	"github.com/makinzm/partial-tree-copy/internal/adapters/ui"
	"github.com/makinzm/partial-tree-copy/internal/adapters/ui/web"
	"github.com/makinzm/partial-tree-copy/internal/usecases/copier"
	"github.com/makinzm/partial-tree-copy/internal/usecases/navigator"
	"github.com/makinzm/partial-tree-copy/internal/usecases/selector"
)

// Application is the main application struct that wires everything together
type Application struct {
	presenter *ui.UIPresenter
	webMode   bool
	webPort   int
}

// NewApplication creates and initializes a new Application
func NewApplication(webMode bool, webPort int) (*Application, error) {
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
		webMode:   webMode,
		webPort:   webPort,
	}, nil
}

// Run starts the application
func (app *Application) Run() error {
	if app.webMode {
		return web.StartServer(".", app.webPort)
	}
	return app.presenter.StartUI()
}
