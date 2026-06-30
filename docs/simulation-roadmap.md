# Simulation Roadmap

## Safe Purpose

The simulation track is for learning and research planning. It is not a wet-lab
protocol, medical intervention, or claim that the software models a real human
body.

## Topic Map

| Discussion Topic | Project Treatment |
| --- | --- |
| Chromosomes and telomeres | Biology explainer and optional concept simulator |
| Telomere shortening | Toy cell-division model only |
| Telomerase | Risk-aware concept model; include cancer-risk caveat |
| Stem-cell DNA repair | Research note: cells maintain their own genome, not whole-body repair |
| Cord blood / placenta cells | Research note: distinguish HSC, MSC, and actual approved uses |
| Local scaffold injection | Research note for tissue engineering, not implementation |
| Full-body digital human | Long-term research background; currently impossible at cell resolution |
| Many-model calibration | Implementable as ensemble/data-assimilation model over wearable data |

## Simulation Levels

1. Concept model: no real biological accuracy, explains mechanisms.
2. Literature-parameter model: uses public papers to test trends.
3. Personal wearable model: uses daily personal aggregates.
4. Equipment-calibrated model: adds device/lab features after access exists.
5. Research-grade model: requires collaborators, protocol, and approved data.

## Implementation Boundary

The current codebase should implement level 3 first. Levels 1 and 2 may be
added as separate educational packages later. Levels 4 and 5 wait for hardware
or collaborators.
