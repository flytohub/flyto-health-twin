# Implementation Workflow

Use this while making changes.

1. Keep domain behavior in `internal/twin`, CLI parsing in `cmd/flyto2`, and
   read-only presentation in `web`.
2. Add tests for malformed input, chronology, model behavior, privacy, or gates
   before changing the corresponding contract.
3. Preserve `PublicRecordJSON` as an explicit field allowlist.
4. Regenerate Go and web symbol references after declaration changes.
5. Update API, CLI, data, feature, validation, state, and changelog documents.
6. Run `make verify` and Flyto2 Indexer strict verification.
