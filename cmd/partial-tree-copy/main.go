package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/makinzm/partial-tree-copy/internal/app"
)

func main() {
	webMode := flag.Bool("web", false, "Launch browser-based GUI instead of TUI")
	webPort := flag.Int("port", 8080, "Port for the web UI server")
	flag.Parse()

	// Create and initialize the application
	application, err := app.NewApplication(*webMode, *webPort)
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
