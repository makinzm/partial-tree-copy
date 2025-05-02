package entities

// FileNode represents a node in the file tree structure.
// It contains information about the file or directory, including its name, path, whether it is a directory,
// its expanded state, its child nodes, selection state, and a reference to its parent node.
type FileNode struct {
	Name     string      // Name of the file or directory
	Path     string      // Full path of the file or directory
	IsDir    bool        // Indicates if the node is a directory
	Expanded bool        // Indicates if the directory is expanded in the tree view
	Children []*FileNode // List of child nodes (files/directories within this directory)
	Selected bool        // Indicates if the node is selected
	Parent   *FileNode   // Reference to the parent node
}

// NewFileNode creates a new FileNode with the given properties
func NewFileNode(name, path string, isDir bool, parent *FileNode) *FileNode {
	return &FileNode{
		Name:     name,
		Path:     path,
		IsDir:    isDir,
		Expanded: false,
		Children: []*FileNode{},
		Selected: false,
		Parent:   parent,
	}
}
