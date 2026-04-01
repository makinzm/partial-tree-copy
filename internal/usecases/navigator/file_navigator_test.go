package navigator

import (
	"fmt"
	"testing"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// Why test FileNavigator?
//
// FileNavigator is the core tree-traversal engine. It decides which nodes the
// user sees (GetVisibleNodes), how deep each node is (GetNodeLevel), how to
// jump between directories, and how breadcrumbs are built. A bug here means
// the user sees the wrong files, selects the wrong items, or gets lost in the
// tree. None of this logic is covered by the existing web-server tests.

// --- mock repository ---

type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string { return m.name }
func (m mockDirEntry) IsDir() bool  { return m.isDir }

type mockFileRepo struct {
	currentDir string
	dirs       map[string][]repositories.DirEntry
}

func (m *mockFileRepo) GetCurrentDirectory() (string, error) { return m.currentDir, nil }
func (m *mockFileRepo) ReadDirectory(path string) ([]repositories.DirEntry, error) {
	entries, ok := m.dirs[path]
	if !ok {
		return nil, fmt.Errorf("directory not found: %s", path)
	}
	return entries, nil
}
func (m *mockFileRepo) ReadFile(string) ([]byte, error)              { return nil, nil }
func (m *mockFileRepo) GetRelativePath(string, string) (string, error) { return "", nil }
func (m *mockFileRepo) WriteToClipboard(string) error                  { return nil }

// --- helpers ---

// buildTestTree creates:
//
//	root/
//	  dirA/
//	    file1.go
//	    file2.go
//	  dirB/
//	  file3.go
func buildTestTree() (*entities.FileNode, *FileNavigator) {
	repo := &mockFileRepo{
		currentDir: "/root",
		dirs: map[string][]repositories.DirEntry{
			"/root": {
				mockDirEntry{"dirA", true},
				mockDirEntry{"dirB", true},
				mockDirEntry{"file3.go", false},
			},
			"/root/dirA": {
				mockDirEntry{"file1.go", false},
				mockDirEntry{"file2.go", false},
			},
			"/root/dirB": {},
		},
	}

	nav := NewFileNavigator(repo)
	root, _ := nav.BuildRootNode()
	return root, nav
}

// --- tests ---

// Verifies that BuildRootNode reads the filesystem and creates the correct
// direct children. If this breaks, the entire tree is empty or malformed on
// startup and nothing else works.
func TestBuildRootNode_CreatesCorrectChildren(t *testing.T) {
	root, _ := buildTestTree()

	if len(root.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(root.Children))
	}

	// dirA should have 2 children (built lazily, but BuildTree is called recursively from BuildRootNode only for root)
	dirA := root.Children[0]
	if dirA.Name != "dirA" || !dirA.IsDir {
		t.Fatalf("expected dirA directory, got %s (isDir=%v)", dirA.Name, dirA.IsDir)
	}
}

// A collapsed root must hide all descendants. Without this guarantee the TUI
// would render a full tree even when the user hasn't expanded anything yet.
func TestGetVisibleNodes_CollapsedRoot(t *testing.T) {
	root, nav := buildTestTree()
	// Root is not expanded by default
	root.Expanded = false

	visible := nav.GetVisibleNodes(root)
	if len(visible) != 1 {
		t.Fatalf("collapsed root should show 1 node, got %d", len(visible))
	}
}

// Expanding root should reveal only its direct children, not grandchildren.
// A bug here (e.g. recursive expansion) would flood the screen with every file.
func TestGetVisibleNodes_ExpandedRoot(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	visible := nav.GetVisibleNodes(root)
	// root + dirA + dirB + file3.go = 4 (dirA/dirB not expanded)
	if len(visible) != 4 {
		t.Fatalf("expected 4 visible nodes, got %d", len(visible))
	}
}

// Expanding a subdirectory should add its children to the visible list while
// keeping everything else intact. Tests the recursive traversal path that
// differs from the single-level expansion above.
func TestGetVisibleNodes_NestedExpand(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	dirA := root.Children[0]
	nav.ToggleExpand(dirA) // expand dirA

	visible := nav.GetVisibleNodes(root)
	// root + dirA + file1 + file2 + dirB + file3 = 6
	if len(visible) != 6 {
		t.Fatalf("expected 6 visible nodes, got %d", len(visible))
	}
}

// Children are loaded lazily on first expand (not at tree construction time).
// This is critical for large repos — without lazy loading, BuildRootNode would
// recursively read the entire filesystem, making startup unusably slow.
func TestToggleExpand_LoadsChildrenLazily(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	dirA := root.Children[0]
	// Before toggle, dirA has no children (not yet built for subdirectories)
	// Actually BuildTree is called once for root, so dirA's children are not populated.
	// Wait — BuildRootNode calls BuildTree(rootNode) which only builds direct children of root.
	// dirA children are NOT populated yet.
	if len(dirA.Children) != 0 {
		t.Fatalf("dirA should have 0 children before expand, got %d", len(dirA.Children))
	}

	nav.ToggleExpand(dirA)

	if !dirA.Expanded {
		t.Fatal("dirA should be expanded after toggle")
	}
	if len(dirA.Children) != 2 {
		t.Fatalf("dirA should have 2 children after expand, got %d", len(dirA.Children))
	}
}

// ToggleExpand on a file must be a no-op. If files could be "expanded", the
// TUI would attempt to read a file as a directory, causing errors or panics.
func TestToggleExpand_FileNodeIsIgnored(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	fileNode := root.Children[2] // file3.go
	nav.ToggleExpand(fileNode)

	if fileNode.Expanded {
		t.Fatal("file node should not become expanded")
	}
}

// GetNodeLevel drives indentation in the tree view. Wrong levels make the
// hierarchy visually misleading — child files would appear at the same indent
// as their parent directory.
func TestGetNodeLevel(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true
	nav.ToggleExpand(root.Children[0]) // expand dirA

	dirA := root.Children[0]
	file1 := dirA.Children[0]

	if nav.GetNodeLevel(root) != 0 {
		t.Fatalf("root level should be 0, got %d", nav.GetNodeLevel(root))
	}
	if nav.GetNodeLevel(dirA) != 1 {
		t.Fatalf("dirA level should be 1, got %d", nav.GetNodeLevel(dirA))
	}
	if nav.GetNodeLevel(file1) != 2 {
		t.Fatalf("file1 level should be 2, got %d", nav.GetNodeLevel(file1))
	}
}

// Breadcrumbs show the user where they are in the tree (root > dirA > file1).
// Incorrect ordering or missing segments would disorient the user when deep
// in a nested directory structure.
func TestGetBreadcrumbs(t *testing.T) {
	root, nav := buildTestTree()
	nav.ToggleExpand(root.Children[0])

	file1 := root.Children[0].Children[0]
	crumbs := nav.GetBreadcrumbs(file1)

	if len(crumbs) != 3 {
		t.Fatalf("expected 3 breadcrumbs (root > dirA > file1), got %d", len(crumbs))
	}
	if crumbs[0] != root || crumbs[1] != root.Children[0] || crumbs[2] != file1 {
		t.Fatal("breadcrumb order is wrong")
	}
}

// The J key jumps to the next directory. This must skip over files and stop
// at the end rather than wrapping or crashing. A bug means the user presses J
// and lands on a file or gets stuck in an infinite loop.
func TestMoveToNextDirectory(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	visible := nav.GetVisibleNodes(root)
	// visible: root, dirA, dirB, file3.go

	next := nav.MoveToNextDirectory(visible, root)
	if next.Name != "dirA" {
		t.Fatalf("next dir from root should be dirA, got %s", next.Name)
	}

	next = nav.MoveToNextDirectory(visible, root.Children[0]) // from dirA
	if next.Name != "dirB" {
		t.Fatalf("next dir from dirA should be dirB, got %s", next.Name)
	}

	// No directory after dirB — should stay
	next = nav.MoveToNextDirectory(visible, root.Children[1])
	if next.Name != "dirB" {
		t.Fatalf("should stay at dirB when no next dir, got %s", next.Name)
	}
}

// The K key jumps to the previous directory. Same risks as MoveToNext — must
// skip files and stay put at the top boundary instead of going out of bounds.
func TestMoveToPreviousDirectory(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	visible := nav.GetVisibleNodes(root)

	prev := nav.MoveToPreviousDirectory(visible, root.Children[1]) // from dirB
	if prev.Name != "dirA" {
		t.Fatalf("prev dir from dirB should be dirA, got %s", prev.Name)
	}

	prev = nav.MoveToPreviousDirectory(visible, root) // from root — no previous
	if prev != root {
		t.Fatal("should stay at root when no previous dir")
	}
}

// Edge case: the current node isn't in the visible list (e.g. it was collapsed
// out of view). The function must return the node as-is rather than panicking
// on a -1 index.
func TestMoveToNextDirectory_NodeNotInList(t *testing.T) {
	root, nav := buildTestTree()
	root.Expanded = true

	visible := nav.GetVisibleNodes(root)
	orphan := entities.NewFileNode("orphan", "/orphan", true, nil)

	result := nav.MoveToNextDirectory(visible, orphan)
	if result != orphan {
		t.Fatal("should return the same node when not found in visible list")
	}
}
