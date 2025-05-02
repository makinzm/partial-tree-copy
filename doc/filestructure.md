```shell
partial-tree-copy/
├── cmd/
│   └── partial-tree-copy/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   └── file_node.go
│   │   └── repositories/
│   │       └── file_repository.go
│   ├── usecases/
│   │   ├── navigator/
│   │   │   └── file_navigator.go
│   │   ├── selector/
│   │   │   └── file_selector.go
│   │   └── copier/
│   │       └── file_copier.go
│   ├── adapters/
│   │   ├── repositories/
│   │   │   └── os_file_repository.go
│   │   ├── ui/
│   │   │   ├── tui/
│   │   │   │   ├── model.go
│   │   │   │   └── view.go
│   │   │   └── presenter.go
│   │   └── clipboard/
│   │       └── clipboard_service.go
│   └── app/
│       └── app.go
├── go.mod
├── go.sum
├── .gitignore
├── LICENSE
└── README.md
```
