# Deployment

The public dashboard is a React Vite app under `web/`. It is safe to deploy as
a static site because it only reads `web/public/public-data.json`, which is a
redacted public export.

## Local Development

```bash
make web-install
make web-dev
```

Open the local Vite URL printed by the command.

## Production Build

```bash
make web-data
npm --prefix web run build
```

The deployable output is:

```text
web/dist
```

Set `PUBLIC_SITE_URL` to the approved absolute HTTPS dashboard URL when the
hosting provider cannot expose its production URL automatically. Without a
verified production URL the output is intentionally `noindex` and has no
sitemap. See [the SEO contract](SEO.md).

## Free Hosting Options

### Cloudflare Pages

Use these settings:

```text
Framework preset: Vite
Root directory: /
Build command: npm ci --prefix web && make web-data && npm --prefix web run build
Build output directory: web/dist
```

Cloudflare's `CF_PAGES_URL` is used only when `CF_PAGES_BRANCH=main`; preview
branches remain excluded from search indexing.

### Netlify

`netlify.toml` is included. Connect the GitHub repository and keep the default
settings from the file:

```text
Build command: npm ci --prefix web && make web-data && npm --prefix web run build
Publish directory: web/dist
```

### Vercel

`vercel.json` is included:

```text
Install command: npm ci --prefix web
Build command: make web-data && npm --prefix web run build
Output directory: web/dist
```

### GitHub Pages

GitHub Pages is also viable for this static Vite app. Set `PUBLIC_SITE_URL` to
the final Pages URL in the deployment job, run:

```bash
npm ci --prefix web
make web-data
npm --prefix web run build
```

Then publish `web/dist`.

## Privacy Gate

Never point the public dashboard at raw exports. The only public data path is:

```text
raw/private source -> Go privacy filter -> public JSON export -> React Vite dashboard
```

Use this before publishing updated data:

```bash
make verify
```
