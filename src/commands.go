package main

import (
    "os"
    "strings"
    "github.com/atotto/clipboard"
)

// moveCursorUp moves the cursor up in the tree view.
// If the cursor is at the first child of its parent, it moves the cursor to the parent node.
// Otherwise, it moves the cursor to the previous sibling node.
func (m *model) moveCursorUp() {
    parent := m.cursor.parent
    if parent == nil {
        return
    }
    // Find the index of the cursor in the parent's children.
    index := -1
    for i, child := range parent.children {
        if child == m.cursor {
            index = i
            break
        }
    }
    if index > 0 {
        m.cursor = parent.children[index-1]
    } else {
        m.cursor = parent
    }
}

// moveCursorDown moves the cursor down in the tree view.
// If the current node is a directory and expanded, it moves the cursor to its first child.
// Otherwise, it moves to the next sibling node if available.
func (m *model) moveCursorDown() {
    if m.cursor.isDir && m.cursor.expanded && len(m.cursor.children) > 0 {
        m.cursor = m.cursor.children[0]
    } else {
        parent := m.cursor.parent
        if parent == nil {
            return
        }
        index := -1
        for i, child := range parent.children {
            if child == m.cursor {
                index = i
                break
            }
        }
        if index < len(parent.children)-1 {
            m.cursor = parent.children[index+1]
        }
    }
}

// toggleExpand toggles the expanded state of the current node.
// If the node is a directory and expanded, it will collapse it. If collapsed, it will expand it.
// When expanding, it builds the tree if the children are not yet loaded.
func (m *model) toggleExpand() {
    if m.cursor.isDir {
        m.cursor.expanded = !m.cursor.expanded
        if m.cursor.expanded && len(m.cursor.children) == 0 {
            buildTree(m.cursor)
        }
    }
}

// toggleSelect toggles the selection state of the current node.
// It only selects non-directory files. Selected files are added to the selection map,
// and unselected files are removed from the map.
func (m *model) toggleSelect() {
    if !m.cursor.isDir {
        m.cursor.selected = !m.cursor.selected
        if m.cursor.selected {
            m.selection[m.cursor.path] = m.cursor
        } else {
            delete(m.selection, m.cursor.path)
        }
    }
}

// copySelection copies the contents of all selected files to the clipboard.
// It reads the content of each selected file, prepends the filename as a header, and concatenates them.
// The resulting string is written to the clipboard using the clipboard package.
func (m *model) copySelection() {
    var builder strings.Builder
    for _, node := range m.selection {
        content, err := os.ReadFile(node.path)
        if err != nil {
            continue
        }
        builder.WriteString("### " + node.name + "\n")
        builder.Write(content)
        builder.WriteString("\n\n")
    }
    clipboard.WriteAll(builder.String())
}
