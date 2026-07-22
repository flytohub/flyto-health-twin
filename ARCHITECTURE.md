# Architecture

## Runtime Shape

The repository has two build surfaces and no application server:

- `cmd/flyto2` is a Go CLI that parses commands, reads CSV files, invokes the
  internal domain package, and writes text, CSV, or JSON.
- `internal/twin` owns all data contracts, adapters, deterministic generation,
  prediction, evaluation, privacy filtering, registries, reports, equipment
  gates, workflow recipes, and simulation logic.
- `web` is a React/Vite static dashboard. It reads the generated
  `web/public/public-data.json`; it does not receive private records or call a
  health backend.
- `scripts` owns the Go declaration comment/reference gate. `web/scripts` owns
  TypeScript reference and built-site metadata/public-data validation.

## Data Flow

1. CSV input or the deterministic generator creates finite, unique,
   chronological `DailyRecord` values.
2. Privacy inspection rejects dangerous source headers and flags free-form
   notes before public processing.
3. `BaselineModel` predicts next-day aggregate response from prior records.
4. Evaluation compares predictions with actual daily aggregates and calculates
   mean absolute errors.
5. Public export clears private/equipment-gated values, projects an explicit
   JSON field allowlist, adds model and roadmap provenance, and becomes the
   static dashboard input.

## Extension Boundaries

`DeviceAdapter` and `Model` are internal interfaces used to compare candidate
implementations. Registry entries for unavailable devices describe planned
contracts; they do not imply working integrations. `CheckEquipmentGate()` is
the authority for whether a source is runnable in this prototype.

## Deployment

The Go CLI runs locally. The static `web/dist` output may be hosted on
Cloudflare Pages, Netlify, or Vercel using checked-in configuration.
`make web-data` regenerates the public artifact with a fixed timestamp before
the web build. Production providers inject canonical, robots, social, JSON-LD,
and sitemap URLs; previews remain `noindex`.

## Documentation And Verification

Go AST and TypeScript compiler-based generators require a comment for every
top-level declaration and produce source-linked references. `make verify` then
checks code, tests, all runnable CLI surfaces, privacy boundaries, generated
data, TypeScript, the production bundle, crawler policy, and forbidden public
record keys.

## Safety Boundary

The telomere command is an educational deterministic toy. Its output includes
an explicit non-clinical boundary and must not be presented as biological,
medical, or longevity evidence.
