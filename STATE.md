# STATE.md

Last reviewed: 2026-07-22

## Current State

The deterministic CLI, internal domain package, public export, registries,
equipment gates, simulation boundary, and static dashboard are implemented as
an experimental research prototype. The source-linked reference covers 131 Go
and 26 TypeScript declarations. Real device adapters remain gated.

CSV loading rejects empty/duplicate headers, missing dates, non-finite numbers,
and duplicate daily records. Modeling requires strict date order. Public export
uses an explicit record allowlist and omits private/equipment-gated values.

The static build generates production canonical, social, structured, robots,
and sitemap metadata only when a verified provider URL is available. No live
dashboard URL is currently declared by this repository; preview and CI builds
remain `noindex`.

## Known Risks

- Optional numeric zero currently represents missing sensor data in model hints;
  a future schema version should distinguish measured zero from unavailable.
- JSON import and weighted-trend model entries remain planned, not implemented.
- A production dashboard URL still needs hosting assignment and live SEO/performance
  verification before indexing.

## Verification

- `make verify`: passed on 2026-07-22, including 17 Go tests, CLI/privacy
  smokes, generated references, TypeScript, Vite, and public-data checks.
- `golangci-lint v2.11.4 run ./...`: 0 issues.
- `npm audit --prefix web --audit-level=high`: 0 vulnerabilities.
- Documentation strict audit: 5 source areas and 9 feature surfaces, no errors.
- Flyto2 Indexer strict audit: 17/17 checks, documentation/README/inline 100.
- Browser verification: 390x844, 768x1024, and 1440x900; no page overflow,
  missing images, or console warnings/errors. The table scroll remains scoped to
  its wrapper on tablet/mobile.
