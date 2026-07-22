# Agent Instructions

You are working on Flyto2, an open-source personal health digital
twin prototype with a Go core and React Vite public dashboard.

## Product Boundary

This project is a non-medical, non-diagnostic research prototype. It must not
claim to diagnose, treat, prevent disease, reverse aging, repair telomeres,
perform stem-cell therapy, or provide clinical accuracy.

The first useful loop is:

1. Import or simulate daily wearable and lifestyle data.
2. Predict next-day HRV, resting heart rate, fatigue, sleep quality, and
   recovery state.
3. Compare predictions with actual values.
4. Explain error and likely missing variables.
5. Improve the model and repeat.

Long-term biology topics from the originating discussion must be represented as
safe docs, toy simulations, or future research tracks unless real equipment,
collaborators, and approvals exist.

## Privacy Rules

Never commit private health exports, raw GPS, exact home/work locations, full
medical reports, account tokens, device credentials, or personally identifying
health history.

Public examples must use synthetic data or daily aggregates. Raw exports belong
under ignored local folders such as `data/` or `exports/`.

## Engineering Rules

- Pre-change exploration is required before edits: inspect the current package,
  tests, and relevant docs so changes follow the local shape. Use Flyto2
  Indexer search, context, and impact analysis on affected symbols.
- Post-change verification is required after edits: run `make verify` and record
  any command that could not run, then run
  `flyto-index verify . --full-scan --strict`.
- Read `PROJECT_CONTEXT.md` and the relevant `docs/` page before making
  product or model changes.
- Keep the open-source core independent from private Flyto2 services.
- Prefer small Go packages with clear data contracts.
- Keep the CLI useful and treat the web app as a public read-only dashboard.
- Use only standard-library dependencies in the Go core until a dependency has a
  clear reason. Keep web dependencies minimal and deployable as static Vite
  output.
- Add tests for prediction, aggregation, and privacy filters.
- Keep documentation current when the model scope or data contract changes.
- When working inside `/Users/chester/flytohub`, use `flyto-indexer` for
  secret scans, context checks, and architecture review where useful.

## Recommended Commands

```bash
make verify
go test ./...
go run ./cmd/flyto2 demo
go run ./cmd/flyto2 export public -data examples/synthetic_daily.csv -out -
go run ./cmd/flyto2 privacy check -data examples/synthetic_daily.csv
go run ./cmd/flyto2 benchmark run -profile balanced -days 30
go run ./cmd/flyto2 equipment gate
go run ./cmd/flyto2 workflow recipes
make web-dev
make web-build
flyto-index verify . --full-scan --strict
```

## Flyto2 Project Memory Contract

Every Flyto2 repository must keep this project-memory scaffold current:

- `AGENTS.md`: agent operating rules, repo-specific constraints, verification commands.
- `CLAUDE.md`: Claude-facing handoff rules when this repo is edited outside Codex.
- `PROJECT.md`: product purpose, owned surfaces, users, and non-goals.
- `ARCHITECTURE.md`: module boundaries, runtime shape, data flow, and integration points.
- `STATE.md`: current status, known risks, release/deploy state, and last verification.
- `ROADMAP.md`: near-term, later, and explicitly out-of-scope work.
- `tasks.md`: actionable checklist with owners/status when known.
- `DECISIONS.md`: durable architectural/product decisions with dates and rationale.
- `CHANGELOG.md`: user-visible or operator-visible changes.
- `docs/README.md`: index for durable docs in this repo.
- `workflows/*.md`: repeatable agent workflows for idea capture, planning, implementation, bugfix, refactor, investigation, and wrap-up.
- `handoffs/_registry.md`: index of handoffs; new handoffs use `YYYY-MM-DD-topic.md`.

When changing behavior, public copy, deployment, security posture, or frontend UX, update the relevant memory files in the same change. Do not leave stale brand, email, module count, route, or deployment information behind.

## Flyto2 Frontend Quality Gate

Any frontend, website, dashboard, extension webview, app screen, or generated UI in this repository must avoid these eight failures:

1. Ignoring accessibility: every interactive control needs keyboard access, visible focus, semantic HTML or ARIA, sufficient contrast, and useful alt/labels.
2. Missing responsive design: verify mobile, tablet, and desktop; no clipped text, overflow, hidden primary actions, or broken navigation.
3. Weak visual hierarchy: users must immediately see page purpose, primary action, status, and next step.
4. Template-looking UI: reuse Flyto2 design tokens and local components, but tailor layout and copy to the actual product surface.
5. Useless elements: remove decorative or placeholder UI that does not help the workflow, trust, navigation, or comprehension.
6. Unclear hierarchy: controls, cards, tables, panels, and modals must have clear grouping, spacing, headings, and state.
7. Unintuitive navigation: current location, back/forward paths, and cross-links to docs/blog/product pages must be obvious.
8. Hard-to-understand content: copy must be concrete, scannable, current, and consistent with Flyto2 terminology.

Frontend verification must include the relevant automated checks plus manual or screenshot review for responsive layout, accessibility states, navigation clarity, loading/empty/error states, and content readability. Public pages must preserve SEO basics: canonical URL, sitemap coverage, metadata, structured data when relevant, and no broken internal or external links.
