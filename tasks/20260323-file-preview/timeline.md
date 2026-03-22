# Timeline: File Preview Feature

## 2026-03-23

### Step 1: Write tests first
- Creating preview_test.go with tests for toggle, content loading, and scroll
- Tests will fail initially since the feature doesn't exist yet

- Tests written, compilation fails as expected:
  - unknown field FileRepo, PreviewMode, PreviewContent, PreviewScroll in Model
  - m.LoadPreviewContent undefined
  - m.buildPreviewView undefined

### Step 2: Implement feature
- Adding fields to Model struct
- Updating NewModel signature
- Wiring FileRepository through presenter and app
- Adding preview key handler and content loading
- Adding preview view rendering
