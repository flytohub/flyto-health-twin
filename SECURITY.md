# Security Policy

Privately report vulnerabilities to `security@flyto2.com`. Include the affected
commit, reproduction steps, expected impact, and whether any health-related
data was exposed. Do not include real credentials or personal health records in
the report.

Do not open a public issue for a vulnerability or suspected private-data leak.
The maintained version is the latest commit on `main`; older snapshots do not
receive separate security fixes.

## Scope

Security-sensitive surfaces include CSV path and content handling, public-data
redaction, generated web artifacts, dependency supply chain, and static-host
headers. Planned device adapters and registry entries are not live integrations.

## Data Incidents

If real health data, credentials, locations, or identifying notes enter Git
history or a public artifact, stop publication, revoke affected credentials,
preserve only sanitized diagnostic evidence, and report privately. Deleting a
working-tree file does not remove it from existing commits, caches, or deployed
artifacts.

## Response Targets

- Acknowledge a complete report within 2 business days.
- Confirm severity and remediation ownership within 7 days when reproducible.
- Coordinate disclosure only after affected commits and deployments are fixed.

These are response targets, not a service-level agreement.

## Supply Chain

Go runtime code uses the standard library. Frontend versions are pinned in
`web/package.json` and `web/package-lock.json`; CI runs dependency, lint, build,
privacy, and strict repository checks. Never place deployment tokens or device
credentials in workflow files, examples, generated references, or public data.

This project is a research prototype, not a medical device or clinical system.
Do not use its outputs for diagnosis or treatment decisions.
