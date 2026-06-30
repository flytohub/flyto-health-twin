# Privacy Policy Draft

This repository is a software prototype. It should not contain private health
exports or personally identifying health data.

## Data Published by Default

Public examples may include:

- Synthetic demo data
- Daily aggregate metrics
- Rounded scores
- Prediction outputs
- Error summaries

## Data Not Published by Default

Do not publish:

- Raw GPS or route traces
- Exact home, work, gym, or clinic locations
- Exact sleep timelines
- Full heart-rate time series
- Full medical reports
- Raw Apple Health, Garmin, Fitbit, Oura, or similar exports
- Account credentials or API tokens
- Medication history
- Detailed diagnosis history
- Notes containing names, locations, or sensitive events

## Redaction Rules

Before making data public:

1. Aggregate to one row per day.
2. Remove location fields.
3. Remove raw timestamps unless they are needed for a daily bucket.
4. Round sensitive values where precision is unnecessary.
5. Review free-text notes manually.
6. Prefer synthetic demo data in the repository.

## CLI Enforcement

Use:

```bash
go run ./cmd/healthtwin privacy check -data examples/synthetic_daily.csv
go run ./cmd/healthtwin export public -data examples/synthetic_daily.csv -out public.json
```

The privacy checker blocks obvious raw/private headers such as GPS,
latitude/longitude, access tokens, full medical reports, diagnosis,
medication, exact sleep timelines, and raw heart-rate time series. Free-text
notes are always removed from public export.

## Medical Disclaimer

This project is not a medical device and does not provide medical diagnosis,
treatment, prevention, or clinical decision support.
