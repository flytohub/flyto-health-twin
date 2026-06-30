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
  tests, and relevant docs so changes follow the local shape.
- Post-change verification is required after edits: run `make verify` and record
  any command that could not run.
- Read `PROJECT_CONTEXT.md` and the relevant `docs/` page before making
  product or model changes.
- Keep the open-source core independent from private Flyto services.
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
```
