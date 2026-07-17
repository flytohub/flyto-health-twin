# Flyto2 Health Twin

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Website](https://img.shields.io/badge/website-flyto2.com-8B5CF6)](https://flyto2.com)
[![GitHub](https://img.shields.io/badge/GitHub-flytohub%2Fflyto--health--twin-181717?logo=github)](https://github.com/flytohub/flyto-health-twin)

Open-source personal health digital twin research prototype for privacy-first
wearable data experiments. It imports or simulates daily lifestyle data,
predicts next-day response metrics, compares predictions with actual outcomes,
and reports transparent error analysis.

Use it to study digital twin data models, synthetic health datasets, Go CLI
pipelines, React Vite dashboards, registry-based model adapters, privacy export
rules, and benchmarkable prediction workflows.

This is not a medical product. It does not diagnose, treat, prevent disease,
claim clinical accuracy, reverse aging, repair telomeres, or recommend medical
interventions.

Official links: [flyto2.com](https://flyto2.com) ·
[Repository](https://github.com/flytohub/flyto-health-twin) ·
[Security](mailto:security@flyto2.com)

## Current Scope

- Daily aggregate health records
- Synthetic demo data
- Baseline next-day prediction
- Error analysis
- Privacy-first public data rules
- Go CLI
- React Vite public dashboard
- Synthetic data generator and benchmark fixture
- Adapter, model, dataset, and workflow registries
- Real-equipment integration gate
- Safety-scoped biology toy simulation

## Installation

Install Go 1.22 or newer and Node 24 or newer, then clone the repository.

```bash
git clone https://github.com/flytohub/flyto-health-twin.git
cd flyto-health-twin
go test ./...
npm install --prefix web
```

## Quick Start

```bash
make verify
go run ./cmd/flyto2 demo
make web-dev
```

Example output:

```text
date        hrv_pred hrv_actual hrv_err rhr_pred rhr_actual fatigue_pred fatigue_actual recovery
2026-06-03  48.7     47.0       -1.7    62.6     64.0       4.2          5.0            strained
```

## Repository Layout

```text
cmd/flyto2/        CLI entrypoint
internal/twin/         Data model, CSV import, prediction, evaluation
examples/              Synthetic demo data only
docs/                  Research, privacy, and data model docs
web/                   React Vite public dashboard
```

## CLI

```bash
go run ./cmd/flyto2 demo
go run ./cmd/flyto2 import csv -data examples/synthetic_daily.csv
go run ./cmd/flyto2 predict -data examples/synthetic_daily.csv
go run ./cmd/flyto2 evaluate -data examples/synthetic_daily.csv -limit 5
go run ./cmd/flyto2 export public -data examples/synthetic_daily.csv -out -
go run ./cmd/flyto2 privacy check -data examples/synthetic_daily.csv
go run ./cmd/flyto2 generate synthetic -profile balanced -days 30 -out examples/benchmark_balanced.csv
go run ./cmd/flyto2 registry adapters
go run ./cmd/flyto2 registry models
go run ./cmd/flyto2 registry datasets
go run ./cmd/flyto2 registry workflows
go run ./cmd/flyto2 report model -data examples/synthetic_daily.csv
go run ./cmd/flyto2 benchmark run -profile balanced -days 30
go run ./cmd/flyto2 equipment gate
go run ./cmd/flyto2 simulate telomere -divisions 24
```

## Web Dashboard

The public dashboard is a React Vite app. It reads only
`web/public/public-data.json`, which is generated from the privacy-safe public
export.

The public logo asset is stored at `web/public/brand/flyto-logo.png`.

```bash
make web-install
make web-data
make web-dev
```

For production:

```bash
make web-build
```

Deploy `web/dist` to Cloudflare Pages, Netlify, Vercel, or GitHub Pages. See
`docs/deployment.md`.

## Public Data Rule

Public dashboards should show only daily aggregate or synthetic data. Do not
publish raw GPS, exact sleep timelines, full medical reports, account tokens,
device credentials, or identifiable health history.

## Development

```bash
make verify
go test ./...
go run ./cmd/flyto2 -data examples/synthetic_daily.csv -limit 7
npm --prefix web run build
```

## Usage

The CLI expects one CSV row per day. Start with the synthetic sample:

```bash
go run ./cmd/flyto2 evaluate -data examples/synthetic_daily.csv
```

Use `-limit` to print fewer rows:

```bash
go run ./cmd/flyto2 evaluate -data examples/synthetic_daily.csv -limit 3
```

The current baseline model uses recent daily aggregates and simple strain
signals. It is intentionally transparent so errors can be inspected before any
more complex model is added.

For public dashboards, use `export public`; it strips private notes and omits
raw/private fields by design.

## Roadmap Closure Artifacts

The repo includes runnable artifacts for the non-device parts of the long-term
plan:

- `examples/benchmark_balanced.csv` is a deterministic synthetic benchmark
  fixture.
- `registry adapters` records implemented and future adapter contracts.
- `registry models` records implemented, planned, and equipment-gated model
  cards.
- `benchmark run` produces a model regression report with non-clinical pass
  thresholds.
- `equipment gate` blocks real-device adapters until sample exports, privacy
  mapping, fixtures, importer tests, and redaction proof exist.
- `workflow recipes` describes Flyto2-native automation using `flyto-core` and
  optional `flyto-cloud` orchestration.
- `simulate telomere` is an educational toy simulation only, separate from
  personal wearable predictions.

## Contributing

See `CONTRIBUTING.md`. Contributions must preserve the non-medical scope and
the privacy rules in `docs/privacy-policy.md`.

## License

Apache-2.0.
