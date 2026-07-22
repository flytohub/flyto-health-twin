# Contributing

Flyto2 welcomes small, reviewable contributions that keep the project
privacy-first and non-medical.

## Before You Change Code

1. Read `AGENTS.md`.
2. Read `PROJECT_CONTEXT.md`.
3. Check the relevant document under `docs/`.
4. Inspect existing tests for the package you are changing.
5. Use Flyto2 Indexer search and impact analysis on affected symbols.

## Verification

Run:

```bash
make verify
flyto-index verify . --full-scan --strict
```

For CLI behavior, also run:

```bash
make demo
```

## Scope Rules

Do not add medical, treatment, anti-aging, stem-cell, telomere, or clinical
accuracy claims.

Do not commit raw health exports, account tokens, device credentials, GPS
history, full medical reports, or identifying health notes.

Use synthetic data for examples and tests.

Run `go run ./scripts/generate-reference.go -write` or
`npm --prefix web run docs:write` after intentional declaration changes. Do not
hand-edit generated source references.
