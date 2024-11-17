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
    case tea.WindowSizeMsg:
        m.height = msg.Height
        return m, nil
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "w":
            m.copySelection()
            return m, tea.Quit
        case "K":
            m.moveBrotherCursorUp()
        case "J":
            m.moveBrotherCursorDown()
        case "up", "k":
            m.moveCursorUp()
        case "down", "j":
            m.moveCursorDown()
        case "enter", "a":
            m.toggleExpand()
            m.toggleSelect()
        }
    }
    return m, nil
}

// View renders the current state of the model as a string.
// It displays the file tree and provides user instructions at the bottom.
func (m model) View() string {
    visibleNodes := m.getVisibleNodes()
    var s strings.Builder

    // ヘルプメッセージの行数を考慮して表示可能な行数を計算
    displayLines := m.height - 4 // ヘルプメッセージが3行＋余裕1行

    // startIndexが範囲内に収まるように調整
    if m.startIndex > len(visibleNodes)-displayLines {
        m.startIndex = len(visibleNodes) - displayLines
    }
    if m.startIndex < 0 {
        m.startIndex = 0
    }

    // 表示するノードの範囲を決定
    endIndex := m.startIndex + displayLines
    if endIndex > len(visibleNodes) {
        endIndex = len(visibleNodes)
    }

    // ノードを表示
    for i := m.startIndex; i < endIndex; i++ {
        node := visibleNodes[i]
        level := m.getNodeLevel(node)
        s.WriteString(m.renderNodeLine(node, level))
    }

    // ヘルプメッセージを追加
    s.WriteString("\nHow to use")
    s.WriteString("\nPress 'w'/Ctrl+'c' to quit, 'Enter' to select a file or")
    s.WriteString(" expand/collapse a dir, up('k')/down('j') to move")
    return s.String()
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

func (m *model) renderNodeLine(node *fileNode, level int) string {
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
        } else {
            line += "📁 " + node.name + "\n"
        }
    } else {
        if node.selected {
            line += "[レ] "
        } else {
            line += "[ ] "
        }
        line += node.name + "\n"
    }
    return line
}

