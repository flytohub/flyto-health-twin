# Bugfix Workflow

Use this for defects and regressions.

1. Reproduce with synthetic or temporary data; never commit a real health input.
2. Trace the defect through import, chronology, model, export allowlist, payload,
   and dashboard as applicable.
3. Add the smallest regression that proves the privacy or behavior boundary.
4. Fix the owning layer without bypassing equipment gates.
5. Run complete verification and update operator-facing documents.
