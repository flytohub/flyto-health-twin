# Research Plan

## Question

Can a privacy-preserving personal data loop predict next-day recovery signals
well enough to identify which variables matter for one person?

## Initial Hypotheses

1. Low sleep score and high stress predict lower next-day HRV.
2. High exercise minutes without enough sleep predict higher next-day fatigue.
3. Higher caffeine intake predicts lower sleep quality for some users.
4. Prediction errors expose missing variables such as illness, emotional stress,
   training load, alcohol, medication, or late meals.

## MVP Protocol

Collect daily aggregate data for 30 to 90 days:

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
- Weight
- Short notes

Each day:

1. Predict tomorrow's HRV, resting heart rate, fatigue, sleep quality, and
   recovery state.
2. Record tomorrow's actual values.
3. Compute error.
4. Record likely missing variables.
5. Publish only privacy-safe aggregates.

## Model Ladder

1. Baseline strain model: transparent rules and recent averages.
2. Weighted trend model: parameterized feature weights learned from one person.
3. Ensemble model: many plausible personal models, keep the closest ones.
4. Bayesian or particle filter model: update model probabilities as new actuals
   arrive.
5. Equipment-calibrated model: add device/lab features only after real access
   exists.

## Non-Goals

- Medical diagnosis
- Treatment recommendation
- Clinical decision support
- Anti-aging claims
- Stem-cell or telomere interventions
- Full-body simulation

## Future Research Gates

- Wearable-only gate: 30 to 90 days of complete daily data.
- Device gate: repeatable import from a named device with privacy review.
- Research gate: written collaborator scope and non-diagnostic protocol.
- Lab gate: approved data source, no raw identifiable files in the public repo.

## Success Criteria

Version 0.1 succeeds if it can:

- Run locally from synthetic data.
- Produce transparent baseline predictions.
- Compare predictions with actuals.
- Generate useful error explanations.
- Keep public data safe by default.
