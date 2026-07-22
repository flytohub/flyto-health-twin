# Feature Reference

## Daily Aggregate Import

`CSVAdapter.Import()` and `LoadCSV()` parse one privacy-scoped record per local
calendar day. `import csv` reports count, range, and source provenance without
persisting the input. Header identity, finite numbers, and unique dates are
validated before modeling. Required columns and public/private fields are
defined in [the data model](data-model.md).

## Synthetic Data

`SyntheticProfiles()` lists deterministic scenarios. `GenerateSyntheticRecords()`
creates at least three days from a named profile, date, and optional length;
`WriteDailyCSV()` emits a reusable fixture. Synthetic output is test data, not a
claim about a real person.

## Prediction And Evaluation

`BaselineModel` implements the internal `Model` interface. It emits next-day
HRV, resting heart rate, fatigue, sleep, recovery state, explanatory hints,
assumptions, and missing-variable labels. `EvaluateWithModel()` performs a
rolling comparison and `MeanAbsoluteErrors()` summarizes four targets.

## Benchmarks And Model Reports

Benchmark fixtures identify repeatable synthetic scenarios.
`BuildModelEvaluationReport()` attaches model, dataset, record-count, and error
provenance. The CLI can generate a fixed profile or report on supplied records.

## Public Export And Privacy Gate

`InspectCSVPrivacy()` detects blocked raw/private source fields.
`PrivacyWarnings()` flags free-form notes and `PublicRecord()` clears all
private and equipment-gated values. `PublicRecordJSON` is the final explicit
field allowlist.
Public export includes daily aggregates, evaluations, model metadata, error
summary, registry counts, equipment status, and safety text.

## Registries And Equipment Gates

Adapter, model, dataset, synthetic-profile, benchmark, and workflow registries
make implemented and planned capabilities discoverable. Registry presence does
not equal implementation. Equipment gates report runnable, blocked, and review
checks for each source.

## Workflow Recipes

`WorkflowRecipes()` describes Flyto2-native import, privacy, evaluation,
benchmark, and export steps. The recipes are contracts for orchestration; this
repository does not run a scheduler or hosted workflow service.

## Telomere Toy Simulation

`RunTelomereToy()` computes a bounded deterministic series from initial length,
division count, stress, and repair inputs. It rejects invalid ranges and labels
its result as non-clinical and non-predictive.

## Static Dashboard

The React/Vite dashboard renders only `public-data.json`. It presents records,
predictions, benchmark errors, roadmap readiness, and safety context with no
account, storage, mutation, or medical-decision workflow.

The UI includes explicit loading and error states, accessible SVG chart labels,
semantic tables, model provenance, public-data boundaries, and deployment-aware
search metadata. See [the web reference](WEB_UI.md).
