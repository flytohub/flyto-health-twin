# CLI Reference

Run commands from the repository root with Go 1.22 or newer:

```bash
go run ./cmd/flyto2 help
```

Commands read local files unless `-out` names a file. They make no network
requests.

| Command | Flags and defaults | Output or effect |
| --- | --- | --- |
| `demo` | `-limit 0` | Evaluates `examples/synthetic_daily.csv`; zero prints all rows. |
| `evaluate` | `-data examples/synthetic_daily.csv`, `-limit 0` | Text predictions plus mean absolute errors. This is also the default command when the first argument is a flag. |
| `predict` | `-data examples/synthetic_daily.csv` | JSON next-day prediction. |
| `import csv` | `-data examples/synthetic_daily.csv` | Prints record count, date range, and adapter ID; does not store records. |
| `privacy check` | `-data examples/synthetic_daily.csv` | Prints pass or JSON issues; exits non-zero when issues exist. |
| `generate synthetic` | `-profile balanced`, `-days 0`, `-start 2026-06-01`, `-out -` | Writes deterministic CSV to stdout or a local file. Zero days uses the profile default. |
| `export public` | `-data examples/synthetic_daily.csv`, `-out -`, `-generated-at ''` | Writes privacy-filtered JSON. An RFC 3339 timestamp makes output reproducible. |
| `registry` | optional `adapters`, `models`, `datasets`, `synthetic`, `benchmarks`, or `workflows` | Writes one registry, or all registries when omitted, as JSON. |
| `report model` | `-data examples/synthetic_daily.csv`, `-dataset synthetic_daily_v0`, `-out -` | Writes an evaluation report for the baseline model. |
| `benchmark run` | `-profile balanced`, `-days 30`, `-start 2026-06-01`, `-out -` | Generates records and writes a deterministic benchmark report. |
| `equipment gate` | `-source ''` | Writes all gate reports or one named adapter report. |
| `workflow recipes` | none | Writes the Flyto2-native recipe registry. |
| `simulate telomere` | `-initial-kb 10`, `-divisions 24`, `-stress 0.35`, `-repair 0.1`, `-out -` | Writes the educational toy series as JSON. |

`-out -` means stdout. A named `-out` path is created or truncated. Invalid
commands, flags, dates, profiles, ranges, or files terminate with a non-zero
exit status.

## Verification Commands

| Command | Scope |
| --- | --- |
| `make test` | Go unit and roadmap contract tests. |
| `make privacy-check` | Fixture privacy gates. |
| `make roadmap-smoke` | Registry, report, benchmark, gate, and simulation CLI smoke tests. |
| `make web-check` | TypeScript validation. |
| `make web-build` | Production Vite build. |
| `make verify` | Formatting, vet, tests, builds, privacy, smoke, public-data, and web checks. |
