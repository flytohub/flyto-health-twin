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

## Free Hosting Options

### Cloudflare Pages

Recommended first deploy target while GitHub Actions is blocked.

Use these settings:

```text
Framework preset: Vite
Root directory: /
Build command: npm ci --prefix web && make web-data && npm --prefix web run build
Build output directory: web/dist
```

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

GitHub Pages is also viable for this static Vite app, but this repository's
current GitHub Actions runs are blocked by the account billing state. After that
is fixed, add a Pages workflow that runs:

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
