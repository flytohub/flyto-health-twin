# Refactor Workflow

Use this for structure changes that should preserve behavior.

1. State the CLI, model, JSON, privacy, and UI behavior that must remain stable.
2. Use Flyto2 Indexer impact results to identify all callers and references.
3. Preserve deterministic fixtures, model versioning, strict chronology, and
   public allowlist behavior.
4. Regenerate references and compare benchmark/public export output.
5. Update architecture or decisions when ownership changes.
