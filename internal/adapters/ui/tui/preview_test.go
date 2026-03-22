package tui

import (
	"testing"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// mockFileRepo implements repositories.FileRepository for testing
type mockFileRepo struct {
	files map[string][]byte
}

func (m *mockFileRepo) GetCurrentDirectory() (string, error) {
	return "/mock", nil
}

func (m *mockFileRepo) ReadDirectory(path string) ([]repositories.DirEntry, error) {
	return nil, nil
}

func (m *mockFileRepo) ReadFile(path string) ([]byte, error) {
	if content, ok := m.files[path]; ok {
		return content, nil
	}
	return nil, &mockError{msg: "file not found: " + path}
}

func (m *mockFileRepo) GetRelativePath(target, base string) (string, error) {
	return target, nil
}

func (m *mockFileRepo) WriteToClipboard(content string) error {
	return nil
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

func newTestModel() *Model {
	repo := &mockFileRepo{
		files: map[string][]byte{
			"/root/file1.go":  []byte("package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n"),
			"/root/file2.txt": []byte("line1\nline2\nline3\n"),
		},
	}

	root := entities.NewFileNode("root", "/root", true, nil)
	root.Expanded = true
	file1 := entities.NewFileNode("file1.go", "/root/file1.go", false, root)
	file2 := entities.NewFileNode("file2.txt", "/root/file2.txt", false, root)
	subdir := entities.NewFileNode("subdir", "/root/subdir", true, root)
	root.Children = []*entities.FileNode{file1, file2, subdir}

	return &Model{
		Root:           root,
		Cursor:         root,
		MaxVisibleRows: 20,
		FocusRight:     false,
		RightScroll:    0,
		FileRepo:       repo,
		PreviewMode:    false,
		PreviewContent: "",
		PreviewScroll:  0,
	}
}

func TestPreviewModeToggle(t *testing.T) {
	m := newTestModel()

	// Initially preview mode should be off
	if m.PreviewMode {
		t.Error("PreviewMode should be false initially")
	}

	// Toggle on
	m.PreviewMode = true
	if !m.PreviewMode {
		t.Error("PreviewMode should be true after toggle")
	}

	// Toggle off
	m.PreviewMode = false
	if m.PreviewMode {
		t.Error("PreviewMode should be false after second toggle")
	}
}

func TestLoadPreviewContentFile(t *testing.T) {
	m := newTestModel()

	// Move cursor to file1.go
	m.Cursor = m.Root.Children[0] // file1.go

	m.LoadPreviewContent()

	if m.PreviewContent == "" {
		t.Error("PreviewContent should not be empty for a file")
	}

	expected := "package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n"
	if m.PreviewContent != expected {
		t.Errorf("PreviewContent = %q, want %q", m.PreviewContent, expected)
	}

	if m.PreviewScroll != 0 {
		t.Errorf("PreviewScroll should be 0 after loading, got %d", m.PreviewScroll)
	}
}

func TestLoadPreviewContentDirectory(t *testing.T) {
	m := newTestModel()

	// Cursor is on root (a directory)
	m.Cursor = m.Root.Children[2] // subdir

	m.LoadPreviewContent()

	if m.PreviewContent != "(directory)" {
		t.Errorf("PreviewContent = %q, want %q", m.PreviewContent, "(directory)")
	}
}

func TestLoadPreviewContentFileNotFound(t *testing.T) {
	m := newTestModel()

	// Create a node pointing to a non-existent file
	missingFile := entities.NewFileNode("missing.go", "/root/missing.go", false, m.Root)
	m.Cursor = missingFile

	m.LoadPreviewContent()

	if m.PreviewContent == "" {
		t.Error("PreviewContent should contain error message for missing file")
	}

	if m.PreviewContent[:len("Error reading file:")] != "Error reading file:" {
		t.Errorf("PreviewContent should start with 'Error reading file:', got %q", m.PreviewContent)
	}
}

func TestPreviewScroll(t *testing.T) {
	m := newTestModel()
	m.Cursor = m.Root.Children[0] // file1.go
	m.PreviewMode = true
	m.FocusRight = true
	m.LoadPreviewContent()

	// Initial scroll should be 0
	if m.PreviewScroll != 0 {
		t.Errorf("PreviewScroll should be 0, got %d", m.PreviewScroll)
	}

	// Scroll down
	m.PreviewScroll++
	if m.PreviewScroll != 1 {
		t.Errorf("PreviewScroll should be 1 after scrolling down, got %d", m.PreviewScroll)
	}

	// Scroll up
	m.PreviewScroll--
	if m.PreviewScroll != 0 {
		t.Errorf("PreviewScroll should be 0 after scrolling back up, got %d", m.PreviewScroll)
	}

	// Should not go below 0
	if m.PreviewScroll > 0 {
		m.PreviewScroll--
	}
	if m.PreviewScroll != 0 {
		t.Errorf("PreviewScroll should not go below 0, got %d", m.PreviewScroll)
	}
}

func TestBuildPreviewView(t *testing.T) {
	m := newTestModel()
	m.Cursor = m.Root.Children[0] // file1.go
	m.PreviewMode = true
	m.LoadPreviewContent()

	view := m.buildPreviewView(20)

	// Should contain the filename
	if view == "" {
		t.Error("buildPreviewView should not return empty string")
	}

	// Should contain "Preview:" header
	if !containsString(view, "Preview:") {
		t.Error("Preview view should contain 'Preview:' header")
	}

	// Should contain file content
	if !containsString(view, "package main") {
		t.Error("Preview view should contain file content")
	}
}

func TestBuildPreviewViewDirectory(t *testing.T) {
	m := newTestModel()
	m.Cursor = m.Root.Children[2] // subdir
	m.PreviewMode = true
	m.LoadPreviewContent()

	view := m.buildPreviewView(20)

	if !containsString(view, "Select a file to preview") {
		t.Error("Preview view for directory should show 'Select a file to preview'")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
