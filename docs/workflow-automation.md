# Flyto2-Native Workflow Automation

Automation should help run the research loop. It must not bypass privacy
boundaries.

## Suitable Automation

- Daily import check
- Data completeness check
- Prediction run
- Prediction vs actual evaluation
- Public JSON export
- Weekly error report
- Collaboration packet generation
- Device adapter smoke test

## Flyto2 Runtime Direction

Flyto2 should use the existing Flyto2 stack as its native automation path:

- `flyto-core` for local workflow execution, module checks, browser/YAML
  verification, and reproducible smoke tests.
- `flyto-cloud` for hosted scheduling, workspace permissions, collaboration,
  auditability, and public dashboard orchestration when a hosted layer is
  needed.
- GitHub Actions for open-source CI checks only.

```text
daily aggregate file
  -> privacy check
  -> prediction
  -> evaluation
  -> public export
  -> report
```

Third-party automation tools can be optional bridges later, but they are not the
default architecture and must not be required to run the open-source project.
Agents may generate workflow drafts, but high-risk actions need human review.

## Implemented Recipes

Workflow recipes are exposed as JSON so `flyto-core` or `flyto-cloud` can
consume them later without requiring another automation platform:

```bash
go run ./cmd/flyto2 workflow recipes
go run ./cmd/flyto2 registry workflows
```

Current recipes:

- Daily public export refresh
- Benchmark regression report
- Equipment adapter gate review

## Hard Boundary

Automation must not:

- Publish raw private data
- Store device credentials in repo files
- Move full health reports to public artifacts
- Treat model output as medical advice
- Claim intervention, treatment, or anti-aging effects
