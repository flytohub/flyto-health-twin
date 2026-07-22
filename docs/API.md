# Go Package Reference

`internal/twin` is an internal package: code outside this Go module cannot
import it. These exported contracts are stable within the CLI and test suite.
The [generated symbol reference](generated/go-symbols.md) maps every type,
function, method, test, and documentation generator to its source comment.

## Core Data And Models

- `DailyRecord` is internal and may contain private optional values;
  `Prediction`, `Evaluation`, and `ErrorSummary` are described in
  [the data model](data-model.md).
- `Model` requires `ID()`, `Version()`, and `Predict([]DailyRecord)`.
- `BaselineModel` is the shipped deterministic implementation.
- `PredictNext(records)` delegates to the baseline model.
- `Evaluate(records)` and `EvaluateWithModel(model, records)` perform rolling
  next-day evaluation.
- `MeanAbsoluteErrors(evaluations)` returns HRV, resting-heart-rate, fatigue,
  and sleep MAE in that order.

## Input And Synthetic Fixtures

- `DeviceCapability`, `DeviceSource`, and `DeviceAdapter` define adapter
  metadata and import behavior.
- `CSVAdapter.Source()` describes the implemented daily CSV source;
  `CSVAdapter.Import(path)` loads it.
- `LoadCSV(path)` normalizes header case, rejects empty or duplicate headers,
  requires `date`, rejects non-finite numbers and duplicate dates, and returns
  sorted daily aggregate records.
- `SyntheticProfile` and `SyntheticProfiles()` describe deterministic fixture
  scenarios.
- `GenerateSyntheticRecords(profileID, start, days)` creates records;
  `WriteDailyCSV(writer, records)` serializes them.

## Privacy And Public Export

- `PrivacyIssue` describes a blocked or review finding.
- `InspectCSVPrivacy(path)` checks source headers before parsing.
- `PrivacyWarnings(record)` checks values that need review.
- `PublicRecord(record)` clears notes, weight, blood pressure, blood glucose,
  body temperature, illness score, and training load. JSON publication still
  uses `PublicRecordJSON` so those fields are absent instead of zero-valued.
- `PublicExport`, `PublicRoadmapStatus`, `PublicEquipmentGateStatus`,
  `PublicRecordJSON`, `PublicPrediction`, and `PublicEvaluation` define JSON.
- `BuildPublicExport()` uses the current UTC time;
  `BuildPublicExportAt()` accepts a reproducible timestamp.
- `WritePublicExport()` and `WritePublicExportAt()` encode those forms to an
  `io.Writer`.

## Reports, Registries, And Gates

- `BenchmarkFixture` and `BenchmarkFixtures()` enumerate regression scenarios.
- `ModelEvaluationReport`, `BuildBenchmarkReport()`,
  `BuildModelEvaluationReport()`, and
  `BuildModelEvaluationReportFromEvaluations()` produce model evidence.
- `AdapterContract`, `ModelCard`, and `DatasetCandidate` are public-safe
  capability metadata returned by `AdapterContracts()`, `ModelRegistry()`, and
  `DatasetRegistry()`; `FindAdapterContract(id)` performs exact lookup.
- `GateCheck`, `EquipmentGateReport`, `CheckEquipmentGate(id)`, and
  `CheckAllEquipmentGates()` distinguish implemented sources from gated plans.

## Workflows And Simulation

- `WorkflowStep`, `WorkflowRecipe`, and `WorkflowRecipes()` publish recipe
  metadata without executing it.
- `TelomereToyParams`, `TelomereToyPoint`, `TelomereToyResult`, and
  `RunTelomereToy(params)` implement the bounded educational simulation.

Errors are returned for invalid inputs and surfaced as fatal CLI exits. None of
these interfaces provide clinical validation or medical advice.

`BaselineModel.Predict()` and `EvaluateWithModel()` require records in strictly
increasing date order. CSV loading sorts input and rejects duplicate dates;
direct package callers must provide the same invariant.
