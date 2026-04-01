package selector

import (
	"testing"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
)

// Why test FileSelector?
//
// FileSelector manages which files the user has picked. It is the single
// source of truth for selection state — if ToggleSelect has a bug (e.g.
// double-toggle doesn't deselect, or directories slip through), the user
// copies the wrong set of files. The sorting logic also matters: clipboard
// output order depends on it, so a broken sort means unpredictable output.

// The most basic contract: pressing Space on a file marks it as selected and
// adds it to the internal map. If this fails, no file can ever be copied.
func TestToggleSelect_SelectsFile(t *testing.T) {
	sel := NewFileSelector()
	node := entities.NewFileNode("main.go", "/src/main.go", false, nil)

	sel.ToggleSelect(node)

	if !node.Selected {
		t.Fatal("node should be selected after toggle")
	}
	if _, ok := sel.GetSelection()["/src/main.go"]; !ok {
		t.Fatal("node should be in selection map")
	}
}

// Double-toggle must deselect. Without this, users cannot undo an accidental
// selection — they'd have to restart the tool to remove a file from the list.
func TestToggleSelect_DeselectsFile(t *testing.T) {
	sel := NewFileSelector()
	node := entities.NewFileNode("main.go", "/src/main.go", false, nil)

	sel.ToggleSelect(node) // select
	sel.ToggleSelect(node) // deselect

	if node.Selected {
		t.Fatal("node should be deselected after second toggle")
	}
	if _, ok := sel.GetSelection()["/src/main.go"]; ok {
		t.Fatal("node should not be in selection map after deselect")
	}
}

// Directories must not be selectable — they can't be "copied" as file content.
// If a directory slips into the selection map, CopySelectionToClipboard would
// try to read a directory as a file and produce garbage or an error.
func TestToggleSelect_IgnoresDirectory(t *testing.T) {
	sel := NewFileSelector()
	dir := entities.NewFileNode("src", "/src", true, nil)

	sel.ToggleSelect(dir)

	if dir.Selected {
		t.Fatal("directory should not become selected")
	}
	if len(sel.GetSelection()) != 0 {
		t.Fatal("selection map should be empty after toggling a directory")
	}
}

// GetSelectedNodes must return files sorted by path so clipboard output is
// deterministic. Without sorting, the same selection could produce different
// clipboard text on each run (map iteration order is random in Go), making
// the tool's output unreliable for diffs or documentation.
func TestGetSelectedNodes_SortedByPath(t *testing.T) {
	sel := NewFileSelector()
	nodeC := entities.NewFileNode("c.go", "/src/c.go", false, nil)
	nodeA := entities.NewFileNode("a.go", "/src/a.go", false, nil)
	nodeB := entities.NewFileNode("b.go", "/src/b.go", false, nil)

	sel.ToggleSelect(nodeC)
	sel.ToggleSelect(nodeA)
	sel.ToggleSelect(nodeB)

	nodes := sel.GetSelectedNodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 selected nodes, got %d", len(nodes))
	}

	for i := 0; i < len(nodes)-1; i++ {
		if nodes[i].Path > nodes[i+1].Path {
			t.Fatalf("nodes not sorted: %s > %s", nodes[i].Path, nodes[i+1].Path)
		}
	}
}

// GetSelectedNodesInTreeOrder must respect the visible-nodes ordering, not
// alphabetical. The TUI sidebar shows selections in tree order — if this
// function sorts differently, the sidebar display won't match what the user
// sees in the main tree panel.
func TestGetSelectedNodesInTreeOrder(t *testing.T) {
	sel := NewFileSelector()

	// Simulate a visible-nodes list in tree order
	nodeA := entities.NewFileNode("a.go", "/a.go", false, nil)
	nodeB := entities.NewFileNode("b.go", "/b.go", false, nil)
	nodeC := entities.NewFileNode("c.go", "/c.go", false, nil)
	dir := entities.NewFileNode("dir", "/dir", true, nil)

	visible := []*entities.FileNode{dir, nodeC, nodeA, nodeB}

	// Select only A and C
	sel.ToggleSelect(nodeA)
	sel.ToggleSelect(nodeC)

	result := sel.GetSelectedNodesInTreeOrder(visible)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	// Tree order: C comes before A in the visible list
	if result[0] != nodeC || result[1] != nodeA {
		t.Fatal("nodes should follow tree (visible) order, not alphabetical")
	}
}

// A fresh selector must have zero selections. This catches initialization bugs
// where the map might be nil (causing panics) or pre-populated.
func TestGetSelection_EmptyByDefault(t *testing.T) {
	sel := NewFileSelector()
	if len(sel.GetSelection()) != 0 {
		t.Fatal("new selector should have empty selection")
	}
}
