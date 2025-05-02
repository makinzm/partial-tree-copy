package main

import (
	"path/filepath"
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
			m.moveToPreviousDirectory() // Changed from moveBrotherCursorUp
		case "J":
			m.moveToNextDirectory() // Changed from moveBrotherCursorDown
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
	// 表示すべき行数
	maxLines := m.maxVisibleRows - 4 // ヘルプテキスト用に4行確保

	// 全ノードの表示情報を取得
	visibleNodes := m.getVisibleNodes()

	// カーソル位置を探す
	cursorIdx := -1
	for i, node := range visibleNodes {
		if node == m.cursor {
			cursorIdx = i
			break
		}
	}

	// 表示する範囲を決定（カーソルを中心に表示）
	startIdx := 0
	endIdx := len(visibleNodes)

	if cursorIdx >= 0 && len(visibleNodes) > maxLines {
		// カーソル行を中心に表示
		halfHeight := maxLines / 2

		if cursorIdx > halfHeight {
			startIdx = cursorIdx - halfHeight
		}

		if startIdx+maxLines > len(visibleNodes) {
			startIdx = len(visibleNodes) - maxLines
		}

		if startIdx < 0 {
			startIdx = 0
		}

		endIdx = startIdx + maxLines
		if endIdx > len(visibleNodes) {
			endIdx = len(visibleNodes)
		}
	}

	// 表示する行を構築
	var s strings.Builder

	// 現在のカーソル位置のパス情報を追加
	s.WriteString("CurrentPosition: " + m.getBreadcrumbs() + "\n")

	// 前が省略されていることを示す
	if startIdx > 0 {
		s.WriteString("...\n")
	}

	// ノードを表示
	for i := startIdx; i < endIdx; i++ {
		node := visibleNodes[i]
		level := getNodeLevel(node)
		s.WriteString(renderSingleNode(node, level, node == m.cursor))
	}

	// 後ろが省略されていることを示す
	if endIdx < len(visibleNodes) {
		s.WriteString("...\n")
	}

	// ヘルプテキストを追加
	s.WriteString("\nHow to use\n")
	s.WriteString("Press 'w'/Ctrl+'c' to quit, 'Enter' to expand/collapse dir or select file\n")
	s.WriteString("Navigation: up('k')/down('j') to move, 'J'/'K' to jump between directories\n")

	return s.String()
}

// 単一ノードの表示を行う関数
func renderSingleNode(node *fileNode, level int, isCursor bool) string {
	prefix := strings.Repeat("  ", level)
	line := prefix

	// カーソル位置表示の改善
	if isCursor {
		line += "> "
	} else {
		line += "   "
	}

	if node.isDir {
		if node.expanded {
			line += "📂 " + filepath.Base(node.path) + "\n"
		} else {
			line += "📁 " + filepath.Base(node.path) + "\n"
		}
	} else {
		if node.selected {
			line += "[レ] " + filepath.Base(node.path) + "\n"
		} else {
			line += "[ ] " + filepath.Base(node.path) + "\n"
		}
	}

	return line
}

// ノードの階層レベルを取得
func getNodeLevel(node *fileNode) int {
	level := 0
	current := node

	for current.parent != nil {
		level++
		current = current.parent
	}

	return level
}

// getBreadcrumbs returns a breadcrumb navigation string for the current cursor position.
func (m *model) getBreadcrumbs() string {
	if m.cursor == nil {
		return ""
	}

	var path []*fileNode
	node := m.cursor

	// カーソルから根までのパスを収集
	for node != nil {
		path = append([]*fileNode{node}, path...)
		node = node.parent
	}

	// パンくずナビゲーションを構築
	var breadcrumbs strings.Builder
	breadcrumbs.WriteString("current position: ")

	for i, node := range path {
		// ルートノードの場合は特別な処理
		if i == 0 {
			breadcrumbs.WriteString("/")
			continue
		}

		// 最後のノード（カーソル位置）は強調表示
		if i == len(path)-1 {
			breadcrumbs.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Render(node.name))
		} else {
			breadcrumbs.WriteString(node.name)
		}

		// 最後のノード以外には区切り文字を追加
		if i < len(path)-1 {
			breadcrumbs.WriteString(" / ")
		}
	}

	return breadcrumbs.String()
}
