package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create test file structure
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test Project"), 0644)
	os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	os.WriteFile(filepath.Join(dir, "src", "util.go"), []byte("package main\n\nfunc hello() string { return \"hi\" }"), 0644)

	return dir
}

func TestTreeEndpoint(t *testing.T) {
	dir := setupTestDir(t)
	handler := NewHandler(dir)

	req := httptest.NewRequest("GET", "/api/tree", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tree TreeNode
	if err := json.Unmarshal(w.Body.Bytes(), &tree); err != nil {
		t.Fatalf("failed to parse tree JSON: %v", err)
	}

	if !tree.IsDir {
		t.Error("root should be a directory")
	}

	if len(tree.Children) == 0 {
		t.Error("root should have children")
	}

	// Check that we have README.md and src/
	names := make(map[string]bool)
	for _, child := range tree.Children {
		names[child.Name] = true
	}
	if !names["README.md"] {
		t.Error("should contain README.md")
	}
	if !names["src"] {
		t.Error("should contain src/")
	}
}

func TestFileEndpoint(t *testing.T) {
	dir := setupTestDir(t)
	handler := NewHandler(dir)

	// Read a valid file
	req := httptest.NewRequest("GET", "/api/file?path=README.md", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "# Test Project") {
		t.Error("expected file content")
	}

	// Try to escape the root directory
	req = httptest.NewRequest("GET", "/api/file?path=../../etc/passwd", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("path traversal should be rejected, got %d", w.Code)
	}
}

func TestCopyEndpoint(t *testing.T) {
	dir := setupTestDir(t)
	handler := NewHandler(dir)

	// Note: clipboard won't work in test env, but we can test the format
	body := `{"paths": ["README.md", "src/main.go"]}`
	req := httptest.NewRequest("POST", "/api/copy", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// May fail due to clipboard not available in CI, but should at least parse request
	// We accept both 200 (clipboard worked) and 500 (clipboard not available)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", w.Code)
	}
}

func TestIndexPage(t *testing.T) {
	dir := setupTestDir(t)
	handler := NewHandler(dir)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Partial Tree Copy") {
		t.Error("index page should contain app title")
	}
}
