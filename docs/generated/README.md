# Generated References

`go-symbols.md` is the complete source-linked inventory of Go types, functions,
and methods in `cmd/`, `internal/`, and `scripts/`. Generate it with:

```bash
go run ./scripts/generate-reference.go -write
```

`make lint` checks that every declaration has a source comment and that this
reference matches current code. Do not edit generated references by hand.
