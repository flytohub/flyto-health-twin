# Data Model

## Daily Record

The internal normalized data contract is one row per local calendar day. It can
contain private optional values and must pass through the public projection
before publication.

| Field | Type | Public by default | Notes |
| --- | --- | --- | --- |
| `date` | `YYYY-MM-DD` | yes | Daily bucket only |
| `sleep_score` | number | yes | Aggregate score, not raw timeline |
| `sleep_hours` | number | yes | Rounded daily value |
| `hrv` | number | yes | Daily aggregate |
| `resting_heart_rate` | number | yes | Daily aggregate |
| `steps` | integer | yes | Daily count |
| `exercise_minutes` | integer | yes | Daily aggregate |
| `stress_score` | number | yes | Self-report or device score |
| `fatigue_score` | number | yes | Self-report, 1 to 10 |
| `caffeine_servings` | number | yes | Rounded count |
| `water_liters` | number | yes | Rounded estimate |
| `weight_kg` | number | optional | Consider private by default |
| `blood_pressure_systolic` | integer | optional | Equipment-ready, private review first |
| `blood_pressure_diastolic` | integer | optional | Equipment-ready, private review first |
| `blood_glucose_mgdl` | number | optional | Equipment-ready, private review first |
| `body_temperature_c` | number | no | Equipment-gated illness/recovery context |
| `illness_score` | number | no | Private self-report marker used by the local model |
| `training_load` | number | no | Private device or manual estimate used by the local model |
| `source_id` | string | yes | Adapter/source provenance |
| `notes` | string | no | May contain identifying details |

## Prediction Result

| Field | Meaning |
| --- | --- |
| `target_date` | Date being predicted |
| `predicted_hrv` | Next-day HRV prediction |
| `predicted_resting_heart_rate` | Next-day RHR prediction |
| `predicted_fatigue_score` | Next-day fatigue prediction |
| `predicted_sleep_score` | Next-day sleep quality prediction |
| `recovery_state` | `recovered`, `normal`, or `strained` |
| `hints` | Human-readable likely error drivers |
| `missing_variables` | Data sources that may explain error |
| `model_id` | Stable model identifier |
| `model_version` | Model version used for the prediction |
| `feature_set` | Inputs the model considered |
| `input_start_date` / `input_end_date` | Input window |

## Evaluation Result

| Field | Meaning |
| --- | --- |
| `target_date` | Actual day compared |
| `prediction` | Prediction generated from previous data |
| `actual` | Actual daily record |
| `hrv_error` | Actual minus predicted |
| `rhr_error` | Actual minus predicted |
| `fatigue_error` | Actual minus predicted |
| `sleep_error` | Actual minus predicted |

## Future Extensions

- Device-specific import metadata
- Lab-value snapshots
- Illness and travel markers
- Training load
- Feature provenance
- Model version tracking

## Public Export

The public export is the dashboard-safe contract. It includes daily aggregates,
prediction history, error summary, model metadata, and missing-variable hints.
It does not include notes, weight, blood pressure, blood glucose, body
temperature, illness score, training load, raw GPS, exact sleep timeline, full
medical reports, device credentials, or raw time series.

Additional public-safe status fields:

| Field | Meaning |
| --- | --- |
| `roadmap_status.adapter_contract_count` | Number of adapter contracts tracked in the private/open-core registry |
| `roadmap_status.model_card_count` | Number of model cards tracked |
| `roadmap_status.dataset_candidate_count` | Number of dataset tracks recorded |
| `roadmap_status.workflow_recipe_count` | Number of Flyto2-native workflow recipes |
| `roadmap_status.equipment_gates` | Source ids and gate status only; no raw field mappings |
| `roadmap_status.simulation_boundary` | Safe simulation boundary statement |
| `benchmark` | Synthetic model regression summary for dashboard transparency |

Raw adapter contract fields stay in CLI registry output and docs, not in the
public dashboard export.
