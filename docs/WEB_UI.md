# Web Dashboard Reference

The React/Vite dashboard is a static, read-only view over the generated
`public-data.json` allowlist. It has no account, write, upload, device-sync, or
medical-decision workflow.

## Runtime States

| State | Owner | Visible result |
| --- | --- | --- |
| Loading | `App()` | Branded status screen while the public JSON request is pending. |
| Error | `App()` | Branded error screen with the HTTP or parsing failure. |
| Ready | `Dashboard()` | Metrics, charts, provenance, privacy boundary, evaluations, and roadmap counts. |

The fetch uses `import.meta.env.BASE_URL`, so relative Cloudflare, Netlify,
Vercel, and subpath builds load the same generated payload and logo.

## Dashboard Sections

- Summary metrics show latest recovery state and HRV, resting-heart-rate, and
  sleep mean absolute error.
- Three SVG charts compare actual and predicted HRV, sleep, and fatigue.
  Empty series render an explicit no-evaluation message rather than invalid SVG
  coordinates.
- Model trace exposes model ID/version, input window, feature names, assumptions,
  and export generation time.
- Public data boundary states exactly what is shown and omitted.
- Recent evaluations show signed model errors and heuristic drivers.
- Roadmap items render registry counts, benchmark status, gate counts, and the
  simulation safety boundary supplied by the Go export.

## Components And Helpers

`StatusScreen`, `MetricCard`, `PanelHeader`, `LineChart`, `TagList`, `TextList`,
`BoundaryItem`, `RoadmapItem`, and `StatePill` are presentation components.
`formatState`, `formatNumber`, `formatSigned`, and `formatDateTime` normalize
labels and values. The complete source-linked list of 26 TypeScript declarations
is generated in [web-symbols.md](generated/web-symbols.md).

## Accessibility And Responsive Behavior

Charts expose `role="img"` and an actual-versus-predicted accessible label.
Tables retain semantic headings and use a horizontal wrapper at narrow widths.
Loading and error states preserve one page heading. The dashboard must be
verified at mobile, tablet, and desktop sizes after visual changes, with no
clipped text, hidden primary content, image failures, or browser-console errors.

## Data Boundary

The build validator inspects `records[]` and every evaluation's `actual` object.
It rejects notes, weight, blood pressure, blood glucose, body temperature,
illness score, and training load keys. Missing-variable labels may name those
signals without exposing values.
