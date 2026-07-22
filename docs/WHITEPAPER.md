# Flyto2 Health Twin Whitepaper

## Abstract

Flyto2 Health Twin is an open-source research prototype for privacy-first,
next-day modeling over daily aggregate lifestyle data. It favors transparent
baselines, deterministic synthetic fixtures, explicit error analysis, and a
public export allowlist over clinical claims or opaque recommendations.

## Research Scope

The system imports or generates daily aggregate records, predicts next-day
HRV, resting heart rate, fatigue, sleep, and recovery state, and compares those
predictions with observed values. Registries expose model, adapter, dataset,
workflow, benchmark, and equipment readiness without presenting planned items
as implemented.

This project is not a medical device and does not diagnose, treat, prevent
disease, estimate clinical risk, or recommend interventions. The telomere model
is a bounded deterministic toy simulation, not a biological prediction.

## Architecture

The Go CLI drives internal/twin, where typed records pass chronology,
finiteness, and privacy checks before modeling and evaluation. Public export
converts records through an explicit allowlist. The static React dashboard reads
only generated public-data.json; it has no account, upload, mutation, or
backend path.

The [Go API](API.md), [CLI reference](CLI.md), and
[web UI reference](WEB_UI.md) explain the supported interfaces. Generated
references map all maintained Go and TypeScript declarations to source.

## Privacy Model

Private health inputs belong in ignored local paths. The privacy inspector
rejects blocked raw fields, warnings identify free-form notes, and public
conversion clears private or equipment-gated values. A public build is only as
safe as this allowlist and its tests; source datasets must never be copied into
the web tree.

## Modeling And Evaluation

The baseline is intentionally inspectable. Rolling evaluation and mean
absolute errors expose performance rather than hiding uncertainty. Synthetic
benchmarks make regressions reproducible but do not establish population or
clinical validity.

## Verification And Evolution

Go tests, deterministic fixtures, privacy checks, generated docs, web build,
SEO deployment rules, and strict Indexer checks form the gate. New adapters or
models require provenance, privacy classification, failure behavior, benchmark
evidence, and non-medical safety language before registry status can advance.

