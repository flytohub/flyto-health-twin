# Dashboard Source

`App.tsx` owns the read-only public dashboard and all loading, error, metrics,
chart, provenance, privacy, evaluation, and roadmap views. `main.tsx` mounts it.
`styles.css` owns responsive presentation.

Every top-level TypeScript declaration requires JSDoc and appears in
[`docs/generated/web-symbols.md`](../../docs/generated/web-symbols.md).
