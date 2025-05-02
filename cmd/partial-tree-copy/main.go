package main

import (
	"fmt"
	"os"

	"github.com/makinzm/partial-tree-copy/internal/app"
)

func main() {
	// Create and initialize the application
	application, err := app.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Run the application
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
