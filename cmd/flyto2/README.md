# Flyto2 Health Twin CLI

`main.go` parses local flags and delegates all domain behavior to
`internal/twin`. It may read or create local files but makes no network requests.
Fatal input and encoding errors return a non-zero process status.

Run `go run ./cmd/flyto2 help` from the repository root.
