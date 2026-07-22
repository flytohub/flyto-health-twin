# Search And Deployment Metadata

The dashboard uses deployment-aware SEO to avoid advertising a URL that does
not exist. The repository does not currently declare a verified live dashboard
URL; its GitHub homepage points to the general `https://flyto2.com` site.

## Production URL Resolution

The Vite build uses the first applicable production value:

1. `PUBLIC_SITE_URL`
2. Cloudflare `CF_PAGES_URL` when `CF_PAGES_BRANCH=main`
3. Netlify `URL` when `CONTEXT=production`
4. Vercel `VERCEL_PROJECT_PRODUCTION_URL` when `VERCEL_ENV=production`

Values must be absolute HTTPS URLs. Preview and local builds deliberately emit
`noindex,nofollow`, disallow crawling, and omit a sitemap. A production build
emits indexable robots metadata, matching canonical/Open Graph/JSON-LD URLs,
`robots.txt`, and a one-page `sitemap.xml`.

## Page Metadata

The page includes a descriptive title and meta description, one static fallback
`h1`, Open Graph and Twitter cards, an absolute social image, and
`SoftwareApplication` JSON-LD with repository and Apache-2.0 license links. The
React application replaces the fallback content while preserving the same
product heading.

The sitemap intentionally contains one URL because the build has one public
HTML route. Do not create duplicate or fictional sitemap entries.

## Verification

`npm --prefix web run build` validates placeholder replacement, canonical URL
agreement, crawler policy, sitemap behavior, local assets, required public data,
and private-field exclusion. Test a production-shaped build with:

```bash
PUBLIC_SITE_URL=https://your-approved-flyto2-host/ npm --prefix web run build
```

After assigning a real domain, update the repository homepage and verify the
live response, canonical, robots, sitemap, structured data, Core Web Vitals, and
mobile layout. A successful local build cannot prove DNS or hosting state.
