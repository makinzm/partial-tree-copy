```
★★ The contents of src/main.go is below.
package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    rootPath, err := os.Getwd()
    if err != nil {
        fmt.Println("ディレクトリを取得できません:", err)
        os.Exit(1)
    }

    rootNode := &fileNode{
        name:  rootPath,
        path:  rootPath,
        isDir: true,
    }
    buildTree(rootNode)

    m := model{
        root:      rootNode,
        cursor:    rootNode,
        selection: make(map[string]*fileNode),
    }

    p := tea.NewProgram(m)
    if err := p.Start(); err != nil {
        fmt.Println("エラー:", err)
        os.Exit(1)
    }
}



★★ The contents of src/update.go is below.
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
    s := m.renderNode(m.root, 0)
    s += "\nHow to use"
    s += "\nPress 'w'/Ctrl+'c' to quit, 'Enter' to select a file or expand/collapse a dir, up('k')/down('j') to move"
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


★★ The contents of src/commands.go is below.
package main

import (
    "os"
    "strings"
    "path/filepath"

    "github.com/atotto/clipboard"
)

// moveBrotherCursorUp moves the cursor to the previous sibling node.
func (m *model) moveBrotherCursorUp() {
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

// moveBrotherCursorDown moves the cursor to the next sibling node.
// If the cursor is at the last sibling, it moves to the parent's next sibling node.
func (m *model) moveBrotherCursorDown() {
    // 親ノードを取得
    parent := m.cursor.parent
    if parent == nil {
        return
    }

    // 現在のカーソルのインデックスを取得
    index := -1
    for i, child := range parent.children {
        if child == m.cursor {
            index = i
            break
        }
    }

    // 次の兄弟ノードに移動
    if index < len(parent.children)-1 {
        m.cursor = parent.children[index+1]
        return
    }

    // 次の親ノードの兄弟ノードに移動する
    for parent != nil {
        grandParent := parent.parent
        if grandParent == nil {
            return
        }

        index = -1
        for i, child := range grandParent.children {
            if child == parent {
                index = i
                break
            }
        }

        // 親ノードの次の兄弟ノードがある場合に移動
        if index < len(grandParent.children)-1 {
            m.cursor = grandParent.children[index+1]
            return
        }

        // さらに上の親ノードに移動して探索を続ける
        parent = grandParent
    }
}

// moveCursorDown moves the cursor down in the tree view.
// If the cursor is at the last child of its parent, it moves the cursor to the parent's next sibling node.
// Otherwise, it moves the cursor to the next sibling node if available.
func (m *model) moveCursorUp() {
    // 表示されているノード一覧を取得
    visibleNodes := m.getVisibleNodes()

    // 現在のカーソルのインデックスを取得
    index := -1
    for i, node := range visibleNodes {
        if node == m.cursor {
            index = i
            break
        }
    }

    // インデックスが0より大きい場合、上に移動
    if index > 0 {
        m.cursor = visibleNodes[index-1]
    }
}

// moveCursorDown moves the cursor down in the tree view.
// If the current node is a directory and expanded, it moves the cursor to its first child.
// Otherwise, it moves to the next sibling node if available.
func (m *model) moveCursorDown() {
    // 展開されているディレクトリの場合、最初の子ノードに移動
    if m.cursor.isDir && m.cursor.expanded && len(m.cursor.children) > 0 {
        m.cursor = m.cursor.children[0]
        return
    }

    // 親ノードを取得
    parent := m.cursor.parent
    if parent == nil {
        return
    }

    // 現在のカーソルのインデックスを取得
    index := -1
    for i, child := range parent.children {
        if child == m.cursor {
            index = i
            break
        }
    }

    // 次の兄弟ノードに移動
    if index < len(parent.children)-1 {
        m.cursor = parent.children[index+1]
        return
    }

    // 次の親ノードに移動する処理を追加
    for parent != nil {
        grandParent := parent.parent
        if grandParent == nil {
            return
        }

        index = -1
        for i, child := range grandParent.children {
            if child == parent {
                index = i
                break
            }
        }

        if index < len(grandParent.children)-1 {
            m.cursor = grandParent.children[index+1]
            return
        }

        parent = grandParent
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
    currentDir, err := os.Getwd()
    if err != nil {
        return
    }

    var builder strings.Builder
    for _, node := range m.selection {
        // 現在のディレクトリからの相対パスを取得
        relativePath, err := filepath.Rel(currentDir, node.path)
        if err != nil {
            continue
        }

        content, err := os.ReadFile(node.path)
        if err != nil {
            continue
        }

        // 相対パスを含めてコピー内容に追加
        builder.WriteString("★★ The contents of " + relativePath + " is below.\n")
        builder.Write(content)
        builder.WriteString("\n\n")
    }

    clipboard.WriteAll(builder.String())
}



```
