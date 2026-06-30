# Flyto2 Web

React Vite dashboard for the public Flyto2 demo.

The Flyto2 logo used by the dashboard lives at `public/brand/flyto-logo.png`.

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
