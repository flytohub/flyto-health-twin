# Equipment Readiness

The code should be ready for devices before devices exist. That means stable
interfaces, provenance, and privacy rules, not fake device integrations.

## Adapter Contract

Each device adapter must declare:

- Source id
- Human-readable name
- Source kind
- Sync mode
- Privacy risk
- Capabilities
- Raw fields
- Aggregate fields
- Public export eligibility

## First Device Classes

| Device | Expected Fields | Public Default |
| --- | --- | --- |
| Smartwatch / ring | HRV, heart rate, sleep, steps, activity | Daily aggregates only |
| Blood pressure monitor | systolic, diastolic, pulse | Daily or session summary |
| Body scale | weight, body composition estimate | Optional, private by default |
| Sleep monitor | sleep score, duration, stages | Aggregate only |
| Blood glucose / CGM | glucose snapshots or curves | Private by default |
| Lab snapshot | lipid, glucose, inflammation, hormone markers | Private by default |

## Connection Rule

No adapter may write directly to public export. All sources must flow through:

```text
adapter import -> private/raw boundary -> daily aggregate -> privacy filter -> public export
```

## Hardware Gate

Before adding a real adapter:

1. Obtain sample export files or API docs.
2. Document private fields.
3. Add synthetic fixtures.
4. Add privacy checker coverage.
5. Add importer tests.
6. Prove `export public` omits raw/private data.

This gate is implemented in the CLI:

```bash
go run ./cmd/flyto2 equipment gate
go run ./cmd/flyto2 equipment gate -source wearable_daily_aggregate
```

`manual_csv` is ready for synthetic public demos. Real device contracts remain
`blocked_until_real_equipment_evidence` until the required sample exports,
privacy mapping, fixture coverage, importer tests, and public redaction proof
exist.
