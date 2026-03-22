package web

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
)

// TreeNode represents a file/directory in the JSON tree response
type TreeNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []TreeNode `json:"children,omitempty"`
}

// Handler handles HTTP requests for the web UI
type Handler struct {
	rootDir string
	mux     *http.ServeMux
}

// NewHandler creates a new web UI handler rooted at the given directory
func NewHandler(rootDir string) *Handler {
	h := &Handler{
		rootDir: rootDir,
		mux:     http.NewServeMux(),
	}
	h.mux.HandleFunc("/api/tree", h.handleTree)
	h.mux.HandleFunc("/api/file", h.handleFile)
	h.mux.HandleFunc("/api/copy", h.handleCopy)
	h.mux.HandleFunc("/", h.handleIndex)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) handleTree(w http.ResponseWriter, r *http.Request) {
	tree := h.buildTree(h.rootDir, "")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}

func (h *Handler) buildTree(fullPath, relPath string) TreeNode {
	name := filepath.Base(fullPath)
	if relPath == "" {
		relPath = "."
	}

	node := TreeNode{
		Name: name,
		Path: relPath,
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return node
	}

	if !info.IsDir() {
		return node
	}

	node.IsDir = true
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return node
	}

	for _, entry := range entries {
		childRel := entry.Name()
		if relPath != "." {
			childRel = relPath + "/" + entry.Name()
		}
		child := h.buildTree(filepath.Join(fullPath, entry.Name()), childRel)
		node.Children = append(node.Children, child)
	}

	return node
}

func (h *Handler) handleFile(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		http.Error(w, "path parameter required", http.StatusBadRequest)
		return
	}

	// Prevent path traversal
	cleaned := filepath.Clean(relPath)
	if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.rootDir, cleaned)

	// Ensure the resolved path is within rootDir
	absRoot, _ := filepath.Abs(h.rootDir)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "failed to read file: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(content)
}

func (h *Handler) handleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Paths []string `json:"paths"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	absRoot, _ := filepath.Abs(h.rootDir)

	var builder strings.Builder
	for _, relPath := range req.Paths {
		cleaned := filepath.Clean(relPath)
		if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) {
			continue
		}
		fullPath := filepath.Join(h.rootDir, cleaned)
		absPath, _ := filepath.Abs(fullPath)
		if !strings.HasPrefix(absPath, absRoot) {
			continue
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		builder.WriteString("★★ The contents of " + relPath + " is below.\n")
		builder.Write(content)
		builder.WriteString("\n\n")
	}

	if err := clipboard.WriteAll(builder.String()); err != nil {
		http.Error(w, "failed to copy to clipboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

// StartServer starts the web UI server and opens the browser
func StartServer(rootDir string, port int) error {
	handler := NewHandler(rootDir)

	// Find available port if default is taken
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("failed to find available port: %w", err)
		}
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", actualPort)
	fmt.Printf("Partial Tree Copy Web UI: %s\n", url)

	// Open browser
	openBrowser(url)

	return http.Serve(listener, handler)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Partial Tree Copy</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #1a1b26; color: #c0caf5; height: 100vh; display: flex; flex-direction: column; }
  header { background: #24283b; padding: 12px 20px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #3b4261; }
  header h1 { font-size: 18px; color: #7aa2f7; }
  .header-right { display: flex; align-items: center; gap: 12px; }
  .selected-count { font-size: 14px; color: #9ece6a; }
  button { background: #7aa2f7; color: #1a1b26; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 600; }
  button:hover { background: #89b4fa; }
  button:disabled { background: #3b4261; color: #565f89; cursor: default; }
  .container { display: flex; flex: 1; overflow: hidden; }
  .tree-panel { width: 350px; min-width: 250px; overflow-y: auto; border-right: 1px solid #3b4261; padding: 8px 0; }
  .preview-panel { flex: 1; display: flex; flex-direction: column; }
  .preview-header { padding: 10px 16px; background: #24283b; border-bottom: 1px solid #3b4261; font-size: 13px; color: #565f89; }
  .preview-content { flex: 1; overflow: auto; padding: 0; }
  .preview-content pre { margin: 0; padding: 12px; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 13px; line-height: 1.6; white-space: pre; }
  .line-num { display: inline-block; width: 45px; text-align: right; padding-right: 12px; color: #3b4261; user-select: none; }
  .tree-item { display: flex; align-items: center; padding: 3px 8px; cursor: pointer; user-select: none; }
  .tree-item:hover { background: #24283b; }
  .tree-item.active { background: #283457; }
  .tree-toggle { width: 18px; text-align: center; font-size: 11px; color: #565f89; flex-shrink: 0; }
  .tree-icon { margin-right: 6px; font-size: 14px; flex-shrink: 0; }
  .tree-name { font-size: 13px; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .tree-check { width: 16px; height: 16px; margin-right: 6px; accent-color: #7aa2f7; flex-shrink: 0; }
  .no-preview { display: flex; align-items: center; justify-content: center; height: 100%; color: #565f89; font-size: 14px; }
  .toast { position: fixed; bottom: 20px; right: 20px; background: #9ece6a; color: #1a1b26; padding: 12px 20px; border-radius: 8px; font-weight: 600; opacity: 0; transition: opacity 0.3s; pointer-events: none; }
  .toast.show { opacity: 1; }
</style>
</head>
<body>
<header>
  <h1>Partial Tree Copy</h1>
  <div class="header-right">
    <span class="selected-count" id="selectedCount">0 files selected</span>
    <button id="copyBtn" disabled onclick="copySelected()">Copy to Clipboard</button>
  </div>
</header>
<div class="container">
  <div class="tree-panel" id="treePanel"></div>
  <div class="preview-panel">
    <div class="preview-header" id="previewHeader">Select a file to preview</div>
    <div class="preview-content" id="previewContent">
      <div class="no-preview">Click a file to view its contents</div>
    </div>
  </div>
</div>
<div class="toast" id="toast">Copied to clipboard!</div>

<script>
const state = { tree: null, selected: new Set(), activeFile: null };

async function init() {
  const res = await fetch('/api/tree');
  state.tree = await res.json();
  renderTree();
}

function renderTree() {
  const panel = document.getElementById('treePanel');
  panel.innerHTML = '';
  if (state.tree && state.tree.children) {
    state.tree.children.forEach(child => renderNode(child, panel, 0));
  }
}

function renderNode(node, parent, depth) {
  const item = document.createElement('div');
  item.className = 'tree-item' + (state.activeFile === node.path ? ' active' : '');
  item.style.paddingLeft = (8 + depth * 18) + 'px';

  if (node.isDir) {
    const toggle = document.createElement('span');
    toggle.className = 'tree-toggle';
    toggle.textContent = node._expanded ? '▼' : '▶';
    item.appendChild(toggle);

    const icon = document.createElement('span');
    icon.className = 'tree-icon';
    icon.textContent = node._expanded ? '📂' : '📁';
    item.appendChild(icon);

    const name = document.createElement('span');
    name.className = 'tree-name';
    name.textContent = node.name;
    item.appendChild(name);

    item.onclick = (e) => {
      e.stopPropagation();
      node._expanded = !node._expanded;
      renderTree();
    };
  } else {
    const toggle = document.createElement('span');
    toggle.className = 'tree-toggle';
    item.appendChild(toggle);

    const check = document.createElement('input');
    check.type = 'checkbox';
    check.className = 'tree-check';
    check.checked = state.selected.has(node.path);
    check.onclick = (e) => {
      e.stopPropagation();
      if (state.selected.has(node.path)) {
        state.selected.delete(node.path);
      } else {
        state.selected.add(node.path);
      }
      updateCount();
    };
    item.appendChild(check);

    const icon = document.createElement('span');
    icon.className = 'tree-icon';
    icon.textContent = getFileIcon(node.name);
    item.appendChild(icon);

    const name = document.createElement('span');
    name.className = 'tree-name';
    name.textContent = node.name;
    item.appendChild(name);

    item.onclick = (e) => {
      if (e.target.type === 'checkbox') return;
      previewFile(node.path);
    };
  }

  parent.appendChild(item);

  if (node.isDir && node._expanded && node.children) {
    node.children.forEach(child => renderNode(child, parent, depth + 1));
  }
}

function getFileIcon(name) {
  const ext = name.split('.').pop().toLowerCase();
  const icons = { go: '🔵', js: '🟡', ts: '🔷', py: '🐍', md: '📝', json: '📋', yaml: '⚙️', yml: '⚙️', html: '🌐', css: '🎨', sh: '🐚', mod: '📦', sum: '🔒' };
  return icons[ext] || '📄';
}

async function previewFile(path) {
  state.activeFile = path;
  renderTree();
  document.getElementById('previewHeader').textContent = path;
  try {
    const res = await fetch('/api/file?path=' + encodeURIComponent(path));
    if (!res.ok) throw new Error(await res.text());
    const text = await res.text();
    const lines = text.split('\n');
    const html = lines.map((line, i) =>
      '<span class="line-num">' + (i + 1) + '</span>' + escapeHtml(line)
    ).join('\n');
    document.getElementById('previewContent').innerHTML = '<pre>' + html + '</pre>';
  } catch (e) {
    document.getElementById('previewContent').innerHTML = '<div class="no-preview">Error: ' + escapeHtml(e.message) + '</div>';
  }
}

function escapeHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function updateCount() {
  const n = state.selected.size;
  document.getElementById('selectedCount').textContent = n + ' file' + (n !== 1 ? 's' : '') + ' selected';
  document.getElementById('copyBtn').disabled = n === 0;
}

async function copySelected() {
  const paths = Array.from(state.selected);
  try {
    const res = await fetch('/api/copy', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ paths })
    });
    if (!res.ok) throw new Error(await res.text());
    showToast();
  } catch (e) {
    alert('Copy failed: ' + e.message);
  }
}

function showToast() {
  const toast = document.getElementById('toast');
  toast.classList.add('show');
  setTimeout(() => toast.classList.remove('show'), 2000);
}

init();
</script>
</body>
</html>`
