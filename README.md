# 🌲 Partial Tree Copy

A CLI tool for selectively copying files from your project directory tree.

## Overview

**Partial Tree Copy** is a terminal-based utility that allows you to browse through your directory structure, select specific files, and copy their contents to your clipboard in a well-formatted manner. Perfect for code reviews, documentation, and sharing specific parts of your project.

## Features

- Terminal-based UI (TUI) for easy navigation and file selection
- **Browser-based GUI** for point-and-click file selection with content preview
- Tree-structured file browser
- Multi-file selection capabilities
- Formatted clipboard output with file headers
- Intuitive keyboard controls
- Efficient directory navigation

## Demo

```shell
❯ partial-tree-copy
Path: /src/update.go                              Selected Files (8):     
                                                                          
...                                                 1. .git/COMMIT_EDITMSG
      📁 logs                                       2. .git/FETCH_HEAD    
      📁 objects                                    3. .git/HEAD          
      [ ] packed-refs                               4. .git/ORIG_HEAD     
      📁 refs                                       5. .git/config        
    [ ] .gitignore                                  6. src/commands.go    
    [ ] LICENSE                                     7. src/main.go        
    [ ] README.md                                   8. src/models.go      
    📁 demo                                                               
    [ ] go.mod                                                            
    [ ] go.sum                                                            
    [ ] partial-tree-copy                                                 
    📂 src                                                                
      [✓] commands.go                                                     
      [✓] main.go                                                         
      [✓] models.go                                                       
    > [ ] update.go                                                       
                                                                          
How to use
Press 'w'/Ctrl+'c' to quit and copy, 'Space' to select file, 'Enter' to expand/collapse dir
Navigation: 'h'/'l' to switch panels, 'j'/'k' to move up/down, 'J'/'K' to jump between directories
```

[Copied Result](demo/realText.txt)

## Usage

### TUI Mode (default)

```bash
partial-tree-copy
```

Controls:
- `Enter` - Expand/collapse directory
- `j/k` - Move up/down
- `J/K` - Jump between directories
- `h/l` - Switch between panels
- `w` or `Ctrl+c` - Copy selected files and exit

### Web GUI Mode

```bash
partial-tree-copy --web
```

Opens a browser-based GUI where you can:
- Browse the file tree by clicking directories
- Preview file contents by clicking on files
- Select files with checkboxes
- Copy all selected files to clipboard with the "Copy to Clipboard" button

Use `--port` to specify a custom port (default: 8080):
```bash
partial-tree-copy --web --port 3000
```

## Installation

### Install with go install

```bash
go install github.com/makinzm/partial-tree-copy/cmd/partial-tree-copy@latest
```

The binary will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure this directory is in your `PATH`.

### Build from source

```bash
git clone git@github.com:makinzm/partial-tree-copy.git
cd partial-tree-copy
go build -o partial-tree-copy ./cmd/partial-tree-copy
```

### Global installation

```bash
# Copy to a location in your PATH
cp partial-tree-copy ~/.local/bin/
# Or
sudo cp partial-tree-copy /usr/local/bin/
```

## Use Cases

- Share relevant files during code reviews
- Extract code snippets for technical documentation
- Share problematic code sections with team members or AI
- Collect selected files from multiple repositories


## Contributing

Contributions are welcome! Feel free to submit bug reports, feature requests, or pull requests.

## Others

- Japanese Document : [ Go言語で「複数ファイルの内容をクリップボードへコピー」するCLIを作りました #Go - Qiita ]( https://qiita.com/making111/items/67220e315b93d50222d3 )
