package copier

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makinzm/partial-tree-copy/internal/domain/entities"
	"github.com/makinzm/partial-tree-copy/internal/domain/repositories"
)

// Why test FileCopier?
//
// FileCopier produces the final clipboard text that the user pastes into
// code reviews, AI prompts, or documentation. The output format is a contract:
// it must have the "★★" header, the correct relative path, and the file
// content. If the format drifts (e.g. missing newline, wrong path), every
// downstream consumer breaks silently. These tests lock in that contract.

// --- mock repository ---

type mockFileRepo struct {
	currentDir    string
	files         map[string][]byte
	clipboardText string
}

func (m *mockFileRepo) GetCurrentDirectory() (string, error) { return m.currentDir, nil }
func (m *mockFileRepo) ReadDirectory(string) ([]repositories.DirEntry, error) {
	return nil, nil
}
func (m *mockFileRepo) ReadFile(path string) ([]byte, error) {
	content, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}
func (m *mockFileRepo) GetRelativePath(target, base string) (string, error) {
	// Simple mock: strip base prefix
	rel := strings.TrimPrefix(target, base+"/")
	return rel, nil
}
func (m *mockFileRepo) WriteToClipboard(content string) error {
	m.clipboardText = content
	return nil
}

// The golden-path test: one file selected, copied to clipboard. Verifies the
// exact output format including the "★★" header, relative path, content, and
// trailing newlines. This format is a user-facing contract — AI tools and
// reviewers parse it, so even a missing newline is a breaking change.
func TestCopySelectionToClipboard_SingleFile(t *testing.T) {
	repo := &mockFileRepo{
		currentDir: "/project",
		files: map[string][]byte{
			"/project/main.go": []byte("package main"),
		},
	}
	cp := NewFileCopier(repo)

	node := entities.NewFileNode("main.go", "/project/main.go", false, nil)
	selection := map[string]*entities.FileNode{
		node.Path: node,
	}

	if err := cp.CopySelectionToClipboard(selection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "★★ The contents of main.go is below.\npackage main\n\n"
	if repo.clipboardText != expected {
		t.Fatalf("clipboard mismatch.\nwant: %q\ngot:  %q", expected, repo.clipboardText)
	}
}

// An empty selection must not error and must write empty string to clipboard.
// Without this, pressing "copy" with nothing selected could panic on nil map
// iteration or overwrite the user's existing clipboard with garbage.
func TestCopySelectionToClipboard_EmptySelection(t *testing.T) {
	repo := &mockFileRepo{currentDir: "/project"}
	cp := NewFileCopier(repo)

	selection := map[string]*entities.FileNode{}
	if err := cp.CopySelectionToClipboard(selection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.clipboardText != "" {
		t.Fatalf("expected empty clipboard for empty selection, got %q", repo.clipboardText)
	}
}

// A file might be deleted between selection and copy (e.g. by another process).
// The copier must skip it gracefully rather than returning an error or crashing,
// so the user still gets the rest of their selection.
func TestCopySelectionToClipboard_SkipsMissingFile(t *testing.T) {
	repo := &mockFileRepo{
		currentDir: "/project",
		files:      map[string][]byte{}, // no files
	}
	cp := NewFileCopier(repo)

	node := entities.NewFileNode("gone.go", "/project/gone.go", false, nil)
	selection := map[string]*entities.FileNode{node.Path: node}

	if err := cp.CopySelectionToClipboard(selection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Missing file should be silently skipped
	if repo.clipboardText != "" {
		t.Fatalf("expected empty clipboard when file is missing, got %q", repo.clipboardText)
	}
}

// Verifies that nested files use their relative path (src/lib.go) in the
// header, not the absolute path (/project/src/lib.go). Absolute paths leak
// the user's directory structure and break portability when sharing snippets.
func TestCopySelectionToClipboard_FormatContainsStarHeader(t *testing.T) {
	repo := &mockFileRepo{
		currentDir: "/project",
		files: map[string][]byte{
			"/project/src/lib.go": []byte("package lib"),
		},
	}
	cp := NewFileCopier(repo)

	node := entities.NewFileNode("lib.go", "/project/src/lib.go", false, nil)
	selection := map[string]*entities.FileNode{node.Path: node}

	if err := cp.CopySelectionToClipboard(selection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(repo.clipboardText, "★★ The contents of src/lib.go is below.") {
		t.Fatalf("clipboard should contain star header with relative path, got %q", repo.clipboardText)
	}
}
