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

