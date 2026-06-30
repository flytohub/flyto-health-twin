# Master Plan

This plan turns the originating discussion into an implementation ladder. Every
topic is preserved, but risky biology topics are kept as safe research context
until proper equipment, collaborators, and approvals exist.

## Phase 0: Ground Rules

- Non-medical, non-diagnostic, non-treatment prototype.
- No anti-aging, stem-cell, telomere-repair, or clinical accuracy claims.
- No home wet-lab instructions.
- Program and model logic can be open-source; private health data does not go
  into the repository.

## Phase 1: Working Low-Cost Personal Model

Build the daily loop first:

```text
daily aggregate
  -> next-day prediction
  -> actual outcome
  -> error analysis
  -> missing-variable hints
```

Inputs:

- Sleep score and hours
- HRV
- Resting heart rate
- Steps
- Exercise minutes
- Stress and fatigue score
- Caffeine and water
- Weight, blood pressure, glucose, illness, and training load when available

Outputs:

- Tomorrow HRV
- Tomorrow resting heart rate
- Tomorrow fatigue
- Tomorrow sleep score
- Recovery state

## Phase 2: Public Experiment Log

Create dashboard-safe exports:

- Daily aggregate table
- Prediction history
- Prediction vs actual error
- Missing-variable report
- Model version and feature set
- Privacy-safe collaboration summary

Do not publish raw health exports, GPS, exact sleep timelines, device tokens,
full medical reports, diagnosis history, or identifying notes.

## Phase 3: Model Calibration

Add models without changing the data contract:

- Baseline strain model
- Weighted trend model
- Ensemble model
- Bayesian calibration
- Particle filter / data assimilation

The project measures quality by prediction error and reproducibility, not by
claiming a human-similarity percentage.

## Phase 4: Equipment Readiness

When equipment exists, connect through adapters:

- Smartwatch or smart ring
- Blood pressure monitor
- Body scale
- Sleep monitor
- Blood glucose meter or CGM
- Periodic lab snapshots

Each adapter must define source id, capabilities, sync mode, privacy risk, raw
fields, aggregate fields, and public-export eligibility.

## Phase 5: Research Collaboration

Use the grounded pitch:

> Privacy-preserving personal physiological response model using wearable,
> lifestyle, and optional health metrics to predict recovery signals and expose
> missing biomarkers through error analysis.

Target collaborators:

- Biomedical engineering
- Medical informatics
- Preventive medicine
- Family medicine
- Rehabilitation or sports medicine
- Wearable device vendors
- Student research teams

## Phase 6: Biology Simulation Track

Keep biology topics separate from the wearable model:

- Chromosome and telomere education
- Cell division and telomere shortening concept simulation
- Stem-cell differentiation research notes
- Injectable scaffold / hydrogel research notes
- Digital human limitations

These modules can teach concepts and track public research, but they must not
claim real intervention, clinical safety, or human-body accuracy.

## Phase 7: Automation

Use n8n, Flyto workflows, or GitHub Actions only around privacy-safe outputs:

- Daily import check
- Prediction run
- Public export generation
- Error report
- Collaboration packet generation
- Device adapter smoke test

Automation must not move raw private data into public artifacts.
