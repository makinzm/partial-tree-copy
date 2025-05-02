package repositories

// FileRepository defines the interface for file system operations
type FileRepository interface {
	// GetCurrentDirectory returns the current working directory
	GetCurrentDirectory() (string, error)

	// ReadDirectory reads the contents of a directory at the given path
	// and returns a list of entries
	ReadDirectory(path string) ([]DirEntry, error)

	// ReadFile reads the content of a file at the given path
	ReadFile(path string) ([]byte, error)

	// GetRelativePath returns the path of target relative to base
	GetRelativePath(target, base string) (string, error)

	// WriteToClipboard writes the given content to the system clipboard
	WriteToClipboard(content string) error
}

// DirEntry represents an entry in a directory
type DirEntry interface {
	// Name returns the name of the directory entry
	Name() string

	// IsDir reports whether the entry describes a directory
	IsDir() bool
}
