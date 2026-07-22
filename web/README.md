# Flyto2 Health Twin Web

React Vite dashboard for the public Flyto2 Health Twin research demo.

The Flyto2 logo used by the dashboard lives at `public/brand/flyto-logo.png`.

The dashboard reads `public/public-data.json`. That file contains redacted daily
aggregates, prediction/evaluation history, benchmark summary, and public-safe
roadmap status counts. It does not include raw adapter fields or private health
exports.

```bash
npm install
npm run dev
npm run build
npm run preview
```

The app reads `public/public-data.json`. From the repository root, regenerate it
with:

```bash
make web-data
```

The public JSON must remain redacted. Do not place raw wearable exports, account
tokens, GPS traces, full medical reports, or private health notes in `web/public`.

Production indexing requires an approved `PUBLIC_SITE_URL` or a detected
Cloudflare Pages, Netlify, or Vercel production URL. Preview builds remain
`noindex`. See [`docs/SEO.md`](../docs/SEO.md).
