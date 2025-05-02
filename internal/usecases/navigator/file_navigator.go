package navigator

import (
	"path/filepath"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// FileNavigator handles the navigation through the file tree
type FileNavigator struct {
	repo repositories.FileRepository
}

// NewFileNavigator creates a new FileNavigator
func NewFileNavigator(repo repositories.FileRepository) *FileNavigator {
	return &FileNavigator{
		repo: repo,
	}
}

// BuildRootNode creates the root node for the file tree
func (fn *FileNavigator) BuildRootNode() (*entities.FileNode, error) {
	rootPath, err := fn.repo.GetCurrentDirectory()
	if err != nil {
		return nil, err
	}

	rootNode := entities.NewFileNode(rootPath, rootPath, true, nil)
	fn.BuildTree(rootNode)

	return rootNode, nil
}

// BuildTree populates the child nodes of the given fileNode
func (fn *FileNavigator) BuildTree(node *entities.FileNode) {
	entries, err := fn.repo.ReadDirectory(node.Path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		childNode := entities.NewFileNode(
			entry.Name(),
			filepath.Join(node.Path, entry.Name()),
			entry.IsDir(),
			node,
		)
		node.Children = append(node.Children, childNode)
	}
}

// GetVisibleNodes returns a list of nodes that are currently visible based on the expanded state
func (fn *FileNavigator) GetVisibleNodes(root *entities.FileNode) []*entities.FileNode {
	var nodes []*entities.FileNode
	var traverse func(node *entities.FileNode)

	traverse = func(node *entities.FileNode) {
		nodes = append(nodes, node)
		if node.IsDir && node.Expanded {
			for _, child := range node.Children {
				traverse(child)
			}
		}
	}

	traverse(root)
	return nodes
}

// ToggleExpand toggles the expanded state of a directory node
func (fn *FileNavigator) ToggleExpand(node *entities.FileNode) {
	if node.IsDir {
		node.Expanded = !node.Expanded
		if node.Expanded && len(node.Children) == 0 {
			fn.BuildTree(node)
		}
	}
}

// GetNodeLevel returns the depth level of a node in the tree
func (fn *FileNavigator) GetNodeLevel(node *entities.FileNode) int {
	level := 0
	current := node

	for current.Parent != nil {
		level++
		current = current.Parent
	}

	return level
}

// GetBreadcrumbs returns a list of nodes from root to the given node
func (fn *FileNavigator) GetBreadcrumbs(node *entities.FileNode) []*entities.FileNode {
	var path []*entities.FileNode
	current := node

	// Collect path from node to root
	for current != nil {
		path = append([]*entities.FileNode{current}, path...)
		current = current.Parent
	}

	return path
}

// MoveToNextDirectory finds the next directory in the visible nodes
func (fn *FileNavigator) MoveToNextDirectory(visibleNodes []*entities.FileNode, currentNode *entities.FileNode) *entities.FileNode {
	currentIndex := -1
	for i, node := range visibleNodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return currentNode
	}

	for i := currentIndex + 1; i < len(visibleNodes); i++ {
		if visibleNodes[i].IsDir {
			return visibleNodes[i]
		}
	}

	return currentNode
}

// MoveToPreviousDirectory finds the previous directory in the visible nodes
func (fn *FileNavigator) MoveToPreviousDirectory(visibleNodes []*entities.FileNode, currentNode *entities.FileNode) *entities.FileNode {
	currentIndex := -1
	for i, node := range visibleNodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return currentNode
	}

	for i := currentIndex - 1; i >= 0; i-- {
		if visibleNodes[i].IsDir {
			return visibleNodes[i]
		}
	}

	return currentNode
}
