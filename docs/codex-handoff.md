# Codex Handoff

Use this prompt when starting a fresh Codex session:

```text
You are building Flyto Health Twin, an open-source personal health digital twin
prototype in Go. Read AGENTS.md and PROJECT_CONTEXT.md first.

Build an MVP that imports or simulates wearable/lifestyle daily aggregates,
predicts next-day HRV, resting heart rate, fatigue, and sleep quality, compares
predictions with actual values, and explains likely missing variables.

Do not make medical, treatment, anti-aging, stem-cell, telomere, or clinical
accuracy claims. Keep all public data privacy-safe and use synthetic demo data
by default.
```

## First Implementation Tasks

1. Keep the CLI stable.
2. Add JSON import support beside CSV.
3. Add a simple local dashboard that consumes `export public`.
4. Add the weighted trend model behind the `Model` interface.
5. Add model-version metadata to every persisted prediction output.
6. Keep future device adapters behind the privacy filter.
