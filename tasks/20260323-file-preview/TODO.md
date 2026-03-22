# File Preview Feature

## Goal
Add a file preview panel to the TUI that toggles with `p` key, replacing the selection list with file content preview.

## Steps
1. Write tests first (preview_test.go)
2. Add Model fields: PreviewMode, PreviewContent, PreviewScroll, FileRepo
3. Update NewModel to accept FileRepository
4. Wire FileRepository through presenter.go and app.go
5. Add `p` key handler and preview scroll in update.go
6. Add LoadPreviewContent helper in update_subcommands.go
7. Add buildPreviewView in view_subcommands.go
8. Update View() to use preview panel when active
9. Auto-update preview when navigating tree
