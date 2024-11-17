package main

import (
  "strings"

  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

// Init initializes the model. It returns a command (tea.Cmd), but here it simply returns nil.
// This function is called when the program starts.
func (m model) Init() tea.Cmd {
    return nil
}

// Update handles messages and updates the model's state accordingly.
// It processes keyboard inputs (tea.KeyMsg) and performs actions like copying, moving the cursor,
// toggling expansion, and selecting files.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "w":
            m.copySelection()
            return m, tea.Quit
        case "up":
            m.moveCursorUp()
        case "down":
            m.moveCursorDown()
        case "enter":
            m.toggleExpand()
        case "a":
            m.toggleSelect()
        }
    }
    return m, nil
}

// View renders the current state of the model as a string.
// It displays the file tree and provides user instructions at the bottom.
func (m model) View() string {
    s := m.renderNode(m.root, 0)
    s += "\n'w'で終了、'a'で選択、'Enter'で展開/縮小、上下キーで移動"
    return s
}

// renderNode recursively renders the file tree starting from the given node.
// It takes the current node and its depth level in the tree as arguments.
// It returns a string representation of the node and its children with appropriate indentation and icons.
func (m *model) renderNode(node *fileNode, level int) string {
    prefix := strings.Repeat("  ", level)
    line := prefix

    if node == m.cursor {
        line += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("> ")
    } else {
        line += "  "
    }

    if node.isDir {
        if node.expanded {
            line += "📂 " + node.name + "\n"
            for _, child := range node.children {
                line += m.renderNode(child, level+1)
            }
        } else {
            line += "📁 " + node.name + "\n"
        }
    } else {
        if node.selected {
            line += "[レ] "
        } else {
            line += "[ ]"
        }
        line += node.name + "\n"
    }
    return line
}
