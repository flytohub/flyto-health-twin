# Documentation Tooling

`generate-reference.go` parses Go AST for `cmd/`, `internal/`, and `scripts/`.
It requires comments on every type, function, and method and checks the generated
source-linked reference for drift.

Run `go run ./scripts/generate-reference.go -write` after intentional source
declaration changes.
