# Public Dataset Registry

Public datasets are useful for research direction and model ideas, but they do
not replace personal calibration. Track them by purpose and limitation.

| Need | Candidate Source | Use | Limitation |
| --- | --- | --- | --- |
| Telomere relationships | [UK Biobank telomere field 22193](https://biobank.ndph.ox.ac.uk/ukb/field.cgi?id=22193) and linked studies | Literature context | Controlled cohort evidence, not a personal telomere model |
| Stem-cell differentiation | [NCBI GEO](https://www.ncbi.nlm.nih.gov/geo/) and [NCBI SRA](https://www.ncbi.nlm.nih.gov/sra) datasets | Research background | Dataset-specific license/access review; not direct intervention evidence |
| Single-cell states | [Human Cell Atlas Data Portal](https://data.humancellatlas.org/) | Cell-state vocabulary | Some data is controlled; not continuous whole-body monitoring |
| Wearable time series | Public wearable datasets in scientific repositories | Benchmark import/model patterns | Different devices and populations |
| Scaffold materials | Paper supplementary data, [Figshare](https://figshare.com/), [Zenodo](https://zenodo.org/) | Future tissue-engineering notes | Repository presence does not establish quality or clinical relevance |
| Clinical questionnaires | Published scales and validation papers | Survey design | Must avoid medical scoring claims |

## Registry Fields

When a dataset is added later, record:

- Name
- URL or citation
- License
- Data type
- Population
- Fields
- Intended project use
- Privacy concerns
- Known limitations

## Rule

Do not mix external cohort findings into personal predictions without labeling
them as population-level priors.

Portal availability does not grant redistribution rights. Record the exact
dataset, license, access tier, citation, population, consent restrictions, and
retrieval date before adding any derived fixture or result.

## Implemented Fixtures

The first public fixture is deterministic and synthetic:

```bash
make fixtures
go run ./cmd/flyto2 benchmark run -profile balanced -days 30
```

Committed fixture:

```text
examples/benchmark_balanced.csv
```

The benchmark report is a software regression guardrail only. It is not a
clinical accuracy claim and does not represent a population.
