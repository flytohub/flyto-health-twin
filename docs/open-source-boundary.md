# Open-Source Boundary

## Open Core

The public repository should contain:

- Go data contracts
- CSV and JSON importers
- Synthetic datasets
- Privacy filters
- Daily aggregation
- Baseline prediction models
- Error analysis
- CLI
- React Vite public dashboard
- Local-only examples
- Documentation

## Private Layer

The following should stay outside the public repository unless explicitly
redacted or reimplemented as safe mocks:

- Personal raw exports
- Device credentials
- Research partner lists
- Sponsor or vendor agreements
- Cloud sync credentials
- Hosted dashboards with real user data
- Advanced proprietary model weights
- Internal Flyto2 entitlement or billing code

## Flytohub Relationship

This project may live inside the `/Users/chester/flytohub` workspace as a
separate git repository. It should not depend on private Flyto2 services at
runtime.

Flyto2 resources may be used around the project:

- `flyto-indexer` for code review, secret scan, and architecture checks.
- `flyto-core` for local workflow/module execution and verification demos.
- `flyto-cloud` for optional hosted scheduling, collaboration, and public
  orchestration after privacy-safe exports exist.
- Existing Flyto2 docs and gates as engineering patterns.

The open-source project must remain usable by someone who only clones this
repository.
