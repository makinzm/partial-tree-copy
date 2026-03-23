package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/makinzm/partial-tree-copy/internal/app"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: partial-tree-copy [options]\n\n")
		fmt.Fprintf(os.Stderr, "A CLI tool for selectively copying files from your project directory tree.\n\n")
		fmt.Fprintf(os.Stderr, "Modes:\n")
		fmt.Fprintf(os.Stderr, "  (default)  Terminal UI - navigate with keyboard, select files, copy to clipboard\n")
		fmt.Fprintf(os.Stderr, "  --web      Browser GUI - point-and-click file selection with content preview\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	webMode := flag.Bool("web", false, "Launch browser-based GUI instead of TUI")
	webPort := flag.Int("port", 8080, "Port for the web UI server (used with --web)")
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
