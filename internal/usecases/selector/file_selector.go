package selector

import (
	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
)

// FileSelector handles the selection of files in the tree
type FileSelector struct {
	selection map[string]*entities.FileNode
}

// NewFileSelector creates a new FileSelector
func NewFileSelector() *FileSelector {
	return &FileSelector{
		selection: make(map[string]*entities.FileNode),
	}
}

// ToggleSelect toggles the selection state of a file node
func (fs *FileSelector) ToggleSelect(node *entities.FileNode) {
	if !node.IsDir {
		node.Selected = !node.Selected
		if node.Selected {
			fs.selection[node.Path] = node
		} else {
			delete(fs.selection, node.Path)
		}
	}
}

// GetSelection returns the current selection map
func (fs *FileSelector) GetSelection() map[string]*entities.FileNode {
	return fs.selection
}

// GetSelectedNodes returns all selected nodes in a slice
func (fs *FileSelector) GetSelectedNodes() []*entities.FileNode {
	var selectedNodes []*entities.FileNode

	for _, node := range fs.selection {
		selectedNodes = append(selectedNodes, node)
	}

	// Sort nodes by path for consistent ordering
	fs.sortNodesByPath(selectedNodes)

	return selectedNodes
}

// GetSelectedNodesInTreeOrder returns selected nodes in the order they appear in the tree
func (fs *FileSelector) GetSelectedNodesInTreeOrder(visibleNodes []*entities.FileNode) []*entities.FileNode {
	var selectedNodes []*entities.FileNode

	for _, node := range visibleNodes {
		if node.Selected {
			selectedNodes = append(selectedNodes, node)
		}
	}

	return selectedNodes
}

// sortNodesByPath sorts nodes by their path
func (fs *FileSelector) sortNodesByPath(nodes []*entities.FileNode) {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].Path > nodes[j].Path {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}
