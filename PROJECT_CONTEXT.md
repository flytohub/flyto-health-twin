# Project Context

## Summary

Flyto Health Twin is an open-source personal health digital twin prototype.
The project uses wearable data, lifestyle logs, and optional periodic health
metrics to build a low-cost personal physiological response model.

The goal is not to simulate a full human body. The practical first version is a
daily prediction and error-analysis loop for personal recovery signals.

The long arc is to preserve every research direction from the originating
conversation while implementing them in safe layers:

1. Biology concepts become educational simulations and research notes.
2. Low-cost personal signals become the first working model.
3. Equipment and lab data become future adapters when real hardware or
   approved research access exists.

## Positioning

Use this framing:

> Open personal health digital twin experiment using privacy-preserving daily
> aggregates, transparent prediction, and error analysis.

Avoid this framing:

> Longevity system, anti-aging product, stem-cell experiment, telomere repair,
> medical diagnosis, treatment recommendation, or clinical decision support.

## Inputs

Initial input fields:

- Date
- Sleep score
- Sleep hours
- HRV
- Resting heart rate
- Steps
- Exercise minutes
- Stress score
- Fatigue score
- Caffeine servings
- Water intake
- Body weight
- Notes

Optional future inputs:

- Blood pressure
- Blood glucose
- Body temperature
- Illness marker
- Training load
- Lab values

## First Predictions

- Tomorrow HRV
- Tomorrow resting heart rate
- Tomorrow fatigue score
- Tomorrow sleep quality
- Recovery state

## Core Loop

1. Collect daily aggregate data.
2. Generate a next-day prediction.
3. Record actual next-day values.
4. Calculate prediction error.
5. Explain likely missing variables.
6. Publish a privacy-safe dashboard.
7. Improve the model.

## Architecture Principle

The system should always preserve this chain:

```text
raw/private source
  -> privacy filter
  -> daily aggregate
  -> model input window
  -> prediction with model version
  -> actual outcome
  -> error analysis
  -> missing variable report
  -> public export
```

No future device, workflow, or dashboard may bypass the privacy filter.

## Open-Source Boundary

Open-source:

- Data schemas
- Synthetic data
- Importers for CSV or JSON
- Daily aggregation
- Privacy filters
- Baseline prediction engine
- Error analysis
- CLI and local dashboard
- Documentation

Private or optional commercial layer:

- Raw personal health data
- Device account credentials
- Cloud sync
- Team research dashboard
- Device vendor integrations
- Advanced hosted AI analysis
- Partner or sponsor data

## First Milestone

Version 0.1 should prove a transparent local loop:

```text
synthetic_daily.csv
  -> load records
  -> predict next day
  -> compare with actuals
  -> print prediction errors and missing-variable hints
```

## Roadmap Documents

- `docs/master-plan.md` — full phased plan from the shared discussion.
- `docs/simulation-roadmap.md` — telomere, stem-cell, scaffold, and digital
  human topics as safe simulation/research tracks.
- `docs/equipment-readiness.md` — future device and lab adapter boundary.
- `docs/public-dataset-registry.md` — public data sources to use later.
- `docs/workflow-automation.md` — n8n/Flyto workflow automation boundary.
- `docs/collaboration-plan.md` — how to approach researchers and device teams.
