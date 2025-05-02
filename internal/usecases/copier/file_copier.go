package copier

import (
	"strings"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// FileCopier handles copying selected files to clipboard
type FileCopier struct {
	repo repositories.FileRepository
}

// NewFileCopier creates a new FileCopier
func NewFileCopier(repo repositories.FileRepository) *FileCopier {
	return &FileCopier{
		repo: repo,
	}
}

// CopySelectionToClipboard copies all selected files to clipboard
func (fc *FileCopier) CopySelectionToClipboard(selection map[string]*entities.FileNode) error {
	currentDir, err := fc.repo.GetCurrentDirectory()
	if err != nil {
		return err
	}

	var builder strings.Builder
	for _, node := range selection {
		// Get path relative to current directory
		relativePath, err := fc.repo.GetRelativePath(node.Path, currentDir)
		if err != nil {
			continue
		}

		// Read file content
		content, err := fc.repo.ReadFile(node.Path)
		if err != nil {
			continue
		}

		// Add to clipboard content
		builder.WriteString("★★ The contents of " + relativePath + " is below.\n")
		builder.Write(content)
		builder.WriteString("\n\n")
	}

	// Write to clipboard
	return fc.repo.WriteToClipboard(builder.String())
}
