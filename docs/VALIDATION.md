# Validation And Release Gates

`make verify` is the repository's complete local contract. It does not contact
device services or upload health data.

## Gate Sequence

| Target | Enforced contract |
| --- | --- |
| `fmt-check` | All Go under `cmd/`, `internal/`, and `scripts/` is formatted. |
| `lint` | `go vet`, complete Go comments/reference, and complete web JSDoc/reference. |
| `test` | Go model, privacy, registry, fixture, gate, and simulation regressions. |
| `build` | Compiles the `flyto2` CLI. |
| `privacy-check` | Both checked-in synthetic CSV fixtures pass source privacy checks. |
| `public-export-smoke` | A non-empty explicit public JSON export is generated. |
| `roadmap-smoke` | Registries, model report, benchmark, equipment gate, and toy simulation run. |
| `web-data` | Regenerates deterministic dashboard data with a fixed timestamp. |
| `web-check` | TypeScript and generated web reference are current. |
| `web-build` | Production Vite bundle and metadata/public-data validator pass. |

## Documentation Drift

`go run ./scripts/generate-reference.go` fails if a type, function, or method is
undocumented or if [go-symbols.md](generated/go-symbols.md) is stale. The
reference currently maps 131 Go declarations, including tests and the generator.

`npm --prefix web run docs:check` enforces JSDoc and the generated 26-declaration
[web reference](generated/web-symbols.md). Use the corresponding `-write` or
`docs:write` command after intentional source changes.

## Independent Audits

```bash
python3 ../.github/scripts/audit-documentation.py . --strict --json
flyto-index verify . --full-scan --strict --json
npm audit --prefix web --audit-level=high
```

CI runs pinned golangci-lint 2.11.4, installs the pinned Flyto2 Indexer release,
then runs `make lint`, the complete verify loop, the demo, and strict Indexer
verification on `main` and pull requests.
