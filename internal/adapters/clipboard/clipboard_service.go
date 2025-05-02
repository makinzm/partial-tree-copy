package clipboard

import (
	"github.com/atotto/clipboard"
)

// ClipboardService provides access to system clipboard operations
type ClipboardService struct{}

// NewClipboardService creates a new ClipboardService
func NewClipboardService() *ClipboardService {
	return &ClipboardService{}
}

// WriteToClipboard writes the given content to the system clipboard
func (cs *ClipboardService) WriteToClipboard(content string) error {
	return clipboard.WriteAll(content)
}
