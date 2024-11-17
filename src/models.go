package main

import (
    "os"
    "path/filepath"
)

// fileNode represents a node in the file tree structure.
// It contains information about the file or directory, including its name, path, whether it is a directory,
// its expanded state, its child nodes, selection state, and a reference to its parent node.
type fileNode struct {
    name     string       // Name of the file or directory
    path     string       // Full path of the file or directory
    isDir    bool         // Indicates if the node is a directory
    expanded bool         // Indicates if the directory is expanded in the tree view
    children []*fileNode  // List of child nodes (files/directories within this directory)
    selected bool         // Indicates if the node is selected
    parent   *fileNode    // Reference to the parent node
}

// model represents the state of the file tree viewer.
// It holds the root node of the tree, the current cursor position, and a map of selected nodes.
type model struct {
    root      *fileNode              // Root node of the file tree
    cursor    *fileNode              // Current position of the cursor in the tree
    selection map[string]*fileNode   // Map of selected nodes (key is the file path)
}

// buildTree populates the child nodes of the given fileNode.
// It reads the contents of the directory specified by the node's path and creates child nodes for each entry.
// If an error occurs (e.g., permission denied), it simply returns without adding children.
func buildTree(node *fileNode) {
    entries, err := os.ReadDir(node.path)
    if err != nil {
        return
    }
    for _, entry := range entries {
        childNode := &fileNode{
            name:   entry.Name(),
            path:   filepath.Join(node.path, entry.Name()),
            isDir:  entry.IsDir(),
            parent: node,
        }
        node.children = append(node.children, childNode)
    }
}

// getVisibleNodes returns a list of nodes that are currently visible based on the expanded state.
func (m *model) getVisibleNodes() []*fileNode {
    var nodes []*fileNode
    var traverse func(node *fileNode)

    traverse = func(node *fileNode) {
        nodes = append(nodes, node)
        if node.isDir && node.expanded {
            for _, child := range node.children {
                traverse(child)
            }
        }
    }

    traverse(m.root)
    return nodes
}
