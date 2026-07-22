# DECISIONS.md

## 2026-07-16 - Adopt Flyto2 Workspace Memory Standard

Decision: `flyto-health-twin` follows the Flyto2 project memory scaffold and frontend quality gate.

Rationale: All 27 Flyto2 repositories need consistent handoff context, durable decisions, and UI quality constraints.

Consequences:

- Root memory files must stay current.
- UI changes must avoid the eight forbidden frontend failures in `AGENTS.md`.
- Handoffs must be registered in `handoffs/_registry.md`.

## 2026-07-22 - Public Export Is An Explicit Allowlist

Decision: public JSON uses `PublicRecordJSON`; the broader `DailyRecord` helper
also clears all private and equipment-gated values.

Rationale: zeroing only notes and weight left health values available to
accidental direct serialization.

Consequences:

- Blood pressure, glucose, temperature, illness, and training load are private
  by default.
- Build validation rejects those keys in public record and actual objects.
- Missing-variable labels may name a signal but cannot publish its value.

## 2026-07-22 - Declaration Documentation Is Generated And Enforced

Decision: Go AST and TypeScript compiler tooling require source comments and
generate complete source-linked references.

Rationale: manually maintained API lists did not cover internal helpers, tests,
CLI handlers, UI components, or formatting functions.

## 2026-07-22 - Preview Builds Are Not Indexable

Decision: canonical, robots, sitemap, social, and structured URLs become
indexable only when an approved production deployment URL is available.

Rationale: the repository currently has no verified live dashboard URL, and a
404 or preview canonical would damage search quality.
