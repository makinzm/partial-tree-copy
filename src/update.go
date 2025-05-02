package main

import (
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model's state accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "w":
			m.copySelection()
			return m, tea.Quit
		case "L", "l":
			// 右側パネルにフォーカスを移動（選択があれば）
			if len(m.selection) > 0 {
				m.focusRight = true
			}
		case "H", "h":
			// 左側パネルにフォーカスを移動
			m.focusRight = false
		case "up", "k":
			if m.focusRight {
				// 右側パネルの場合は選択リストをスクロールアップ
				if m.rightScroll > 0 {
					m.rightScroll--
				}
			} else {
				// 左側パネルの場合は通常のカーソル移動
				m.moveCursorUp()
			}
		case "down", "j":
			if m.focusRight {
				// 右側パネルの場合は選択リストをスクロールダウン
				if m.rightScroll < len(m.getAllSelectedNodes())-1 {
					m.rightScroll++
				}
			} else {
				// 左側パネルの場合は通常のカーソル移動
				m.moveCursorDown()
			}
		case "K":
			if !m.focusRight {
				m.moveToPreviousDirectory()
			}
		case "J":
			if !m.focusRight {
				m.moveToNextDirectory()
			}
		case "enter":
			if !m.focusRight {
				// 左側パネルの場合のみ
				if m.cursor.isDir {
					m.toggleExpand()
				} else {
					m.toggleSelect()
				}
			}
		case "space":
			if !m.focusRight {
				// 左側パネルの場合のみ
				m.toggleSelect()
			}
		}
	}
	return m, nil
}

// View renders the current state of the model as a string.
func (m model) View() string {
	// 表示すべき行数
	maxLines := m.maxVisibleRows - 4 // ヘルプテキスト用に4行確保

	// 左側（ツリービュー）を構築
	leftView := m.buildTreeView(maxLines)

	// 右側（選択ファイルリスト）を構築 - 全ての選択ファイルを表示
	rightView := m.buildSelectionView(maxLines)

	// ツリービューの最大幅を計算（固定幅）
	treeViewWidth := 50 // 適切な幅に調整

	// 左側のビューを固定幅に調整
	paddedLeftView := leftView

	// フォーカス状態に応じてスタイルを変更（全体ではなくカーソルのみ）
	leftStyle := lipgloss.NewStyle().Width(treeViewWidth)
	rightStyle := lipgloss.NewStyle()

	// フォーカス表示は個別のカーソルのみに適用する（全体に色を付けない）
	if m.focusRight {
		// 右側がフォーカスされている場合（ただし色は付けない）
		// カーソルのみ色付けは buildSelectionView() 内で行う
	} else {
		// 左側がフォーカスされている場合
		// カーソルのみ色付けは renderSingleNode() 内で行う
	}

	// 両方を固定位置で結合
	combinedView := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(paddedLeftView),
		rightStyle.Render(rightView))

	// ヘルプテキストを追加
	helpText := "\nHow to use\n" +
		"Press 'w'/Ctrl+'c' to quit and copy, 'Space' to select file, 'Enter' to expand/collapse dir\n" +
		"Navigation: 'h'/'l' to switch panels, 'j'/'k' to move up/down, 'J'/'K' to jump between directories"

	return combinedView + helpText
}

// ツリービューを構築する関数
func (m *model) buildTreeView(maxLines int) string {
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
	s.WriteString("Path: " + m.getBreadcrumbs() + "\n\n")

	// 前が省略されていることを示す
	if startIdx > 0 {
		s.WriteString("...\n")
	}

	// ノードを表示
	for i := startIdx; i < endIdx; i++ {
		node := visibleNodes[i]
		level := getNodeLevel(node)
		s.WriteString(renderSingleNode(node, level, node == m.cursor, !m.focusRight))
	}

	// 後ろが省略されていることを示す
	if endIdx < len(visibleNodes) {
		s.WriteString("...\n")
	}

	return s.String()
}

// 選択ファイルのリストを構築する関数
func (m *model) buildSelectionView(maxLines int) string {
	var s strings.Builder

	// タイトルを追加（2桁以上の数にも対応）
	s.WriteString("Selected Files (" + strconv.Itoa(len(m.selection)) + "):\n\n")

	// 選択がない場合
	if len(m.selection) == 0 {
		s.WriteString("No files selected\n")
		return s.String()
	}

	// 選択されたノードを取得（ツリー構造に関係なく、全ての選択を表示）
	selectedNodes := m.getAllSelectedNodes()

	// 表示開始位置（スクロール位置）
	startIdx := m.rightScroll
	if len(selectedNodes) > 0 && startIdx >= len(selectedNodes) {
		startIdx = len(selectedNodes) - 1
	}

	// 最大表示数を決定
	visibleCount := maxLines - 2 // タイトル行と空行の分を引く
	endIdx := startIdx + visibleCount
	if endIdx > len(selectedNodes) {
		endIdx = len(selectedNodes)
	}

	// 前がスクロールされていることを示す
	if startIdx > 0 {
		s.WriteString("...\n")
	}

	// 選択ファイルを表示
	for i := 0; i < len(selectedNodes); i++ {
		// 表示範囲外はスキップ
		if i < startIdx || i >= endIdx {
			continue
		}

		node := selectedNodes[i]

		// 相対パスを作成
		relPath := getRelativePath(node.path, m.root.path)

		// 表示行を追加（相対パスのみ）
		numStr := strconv.Itoa(i + 1)
		padding := " "
		if i+1 >= 10 { // 2桁になる場合はパディングを調整
			padding = ""
		}

		line := numStr + "." + padding + relPath

		// 現在のスクロール位置を強調表示（右側パネルがフォーカスされている場合）
		if m.focusRight && i == m.rightScroll {
			// カーソル位置だけをピンク色で強調表示
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("> " + line)
		} else {
			line = "  " + line
		}

		s.WriteString(line + "\n")
	}

	// 後ろがスクロールされていることを示す
	if endIdx < len(selectedNodes) {
		s.WriteString("...\n")
	}

	return s.String()
}

// 全ての選択されたノードを取得する関数（ツリー上の表示状態に関わらず）
// 全ての選択されたノードを取得する関数（安定した順序で取得）
func (m *model) getAllSelectedNodes() []*fileNode {
	var selectedNodes []*fileNode

	// 選択マップから全てのノードを取得
	for _, node := range m.selection {
		selectedNodes = append(selectedNodes, node)
	}

	// パスでソート（安定した順序を提供）
	sortNodesByPath(selectedNodes)

	return selectedNodes
}

// ノードをパスでソートする関数
func sortNodesByPath(nodes []*fileNode) {
	// ノードをパスでソート
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].path > nodes[j].path {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}

// ツリー上で選択されたノードを取得する関数（表示順）
func (m *model) getSelectedNodesInTreeOrder() []*fileNode {
	var selectedNodes []*fileNode

	// 表示されているすべてのノードを取得
	visibleNodes := m.getVisibleNodes()

	// 選択されているノードだけを抽出
	for _, node := range visibleNodes {
		if node.selected {
			selectedNodes = append(selectedNodes, node)
		}
	}

	return selectedNodes
}

// 相対パスを取得する関数
func getRelativePath(path, rootPath string) string {
	// パスがルートパスで始まる場合、相対パスを返す
	if strings.HasPrefix(path, rootPath) {
		relPath, err := filepath.Rel(rootPath, path)
		if err == nil {
			return relPath
		}
	}

	// デフォルトでファイル名を返す
	return filepath.Base(path)
}

// 単一ノードの表示を行う関数
func renderSingleNode(node *fileNode, level int, isCursor bool, isFocused bool) string {
	prefix := strings.Repeat("  ", level)
	line := prefix

	// カーソル位置表示（フォーカスがある場合のみ色付け）
	if isCursor {
		if isFocused {
			// フォーカスがある側のカーソルを色付け
			line += lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("> ")
		} else {
			line += "> "
		}
	} else {
		line += "  "
	}

	if node.isDir {
		if node.expanded {
			line += "📂 " + filepath.Base(node.path) + "\n"
		} else {
			line += "📁 " + filepath.Base(node.path) + "\n"
		}
	} else {
		if node.selected {
			line += "[✓] " + filepath.Base(node.path) + "\n"
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
			breadcrumbs.WriteString("/")
		}
	}

	return breadcrumbs.String()
}
