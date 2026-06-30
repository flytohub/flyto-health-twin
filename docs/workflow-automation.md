# Workflow Automation

Automation should help run the research loop. It must not bypass privacy
boundaries.

## Suitable Automation

- Daily import check
- Data completeness check
- Prediction run
- Prediction vs actual evaluation
- Public JSON export
- Weekly error report
- Collaboration packet generation
- Device adapter smoke test

## n8n / Agent Direction

n8n, Flyto workflows, or other agents can act as workflow runtimes:

```text
daily aggregate file
  -> privacy check
  -> prediction
  -> evaluation
  -> public export
  -> report
```

Agents may generate workflow drafts, but high-risk actions need human review.

## Hard Boundary

Automation must not:

- Publish raw private data
- Store device credentials in repo files
- Move full health reports to public artifacts
- Treat model output as medical advice
- Claim intervention, treatment, or anti-aging effects
