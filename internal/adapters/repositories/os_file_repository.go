package repositories

import (
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// OSDirEntry implements the DirEntry interface using os.DirEntry
type OSDirEntry struct {
	entry os.DirEntry
}

// Name returns the name of the directory entry
func (ode *OSDirEntry) Name() string {
	return ode.entry.Name()
}

// IsDir reports whether the entry describes a directory
func (ode *OSDirEntry) IsDir() bool {
	return ode.entry.IsDir()
}

// OSFileRepository is a file repository implementation using OS file operations
type OSFileRepository struct{}

// NewOSFileRepository creates a new OSFileRepository
func NewOSFileRepository() *OSFileRepository {
	return &OSFileRepository{}
}

// GetCurrentDirectory returns the current working directory
func (r *OSFileRepository) GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// ReadDirectory reads a directory and returns its entries
func (r *OSFileRepository) ReadDirectory(path string) ([]repositories.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []repositories.DirEntry
	for _, entry := range entries {
		result = append(result, &OSDirEntry{entry: entry})
	}

	return result, nil
}

// ReadFile reads the content of a file
func (r *OSFileRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetRelativePath returns the path of target relative to base
func (r *OSFileRepository) GetRelativePath(target, base string) (string, error) {
	return filepath.Rel(base, target)
}

// WriteToClipboard writes content to the system clipboard
func (r *OSFileRepository) WriteToClipboard(content string) error {
	return clipboard.WriteAll(content)
}
