import { rm, writeFile } from "node:fs/promises";
import path from "node:path";
import { defineConfig, type Plugin } from "vite";
import react from "@vitejs/plugin-react";

const repositoryUrl = "https://github.com/flytohub/flyto-health-twin";

/** Normalize an absolute deployment URL and retain a trailing slash. */
function normalizeSiteUrl(value?: string) {
  if (!value) return null;
  const url = new URL(value.startsWith("http") ? value : `https://${value}`);
  if (url.protocol !== "https:") {
    throw new Error("PUBLIC_SITE_URL must use HTTPS");
  }
  return url.href.endsWith("/") ? url.href : `${url.href}/`;
}

/** Resolve only a production deployment URL, excluding preview branches. */
function productionSiteUrl() {
  if (process.env.PUBLIC_SITE_URL) return normalizeSiteUrl(process.env.PUBLIC_SITE_URL);
  if (process.env.CF_PAGES_BRANCH === "main") return normalizeSiteUrl(process.env.CF_PAGES_URL);
  if (process.env.CONTEXT === "production") return normalizeSiteUrl(process.env.URL);
  if (process.env.VERCEL_ENV === "production") {
    return normalizeSiteUrl(process.env.VERCEL_PROJECT_PRODUCTION_URL);
  }
  return null;
}

/** Inject deployment-correct metadata and write matching crawler assets. */
function seoPlugin(): Plugin {
  const siteUrl = productionSiteUrl();
  const canonical = siteUrl ?? repositoryUrl;
  const indexable = siteUrl !== null;
  const socialImage = siteUrl
    ? new URL("brand/flyto-logo.png", siteUrl).href
    : "https://flyto2.com/assets/img/og-image.png";
  return {
    name: "flyto2-health-seo",
    transformIndexHtml(html) {
      return html
        .replaceAll("__SITE_URL__", canonical)
        .replaceAll("__ROBOTS__", indexable ? "index,follow,max-image-preview:large" : "noindex,nofollow")
        .replaceAll("__SOCIAL_IMAGE__", socialImage);
    },
    async closeBundle() {
      const output = path.resolve("dist");
      if (siteUrl) {
        await writeFile(
          path.join(output, "robots.txt"),
          `User-agent: *\nAllow: /\n\nSitemap: ${siteUrl}sitemap.xml\n`,
        );
        await writeFile(
          path.join(output, "sitemap.xml"),
          `<?xml version="1.0" encoding="UTF-8"?>\n` +
            `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
            `  <url><loc>${siteUrl}</loc></url>\n` +
            `</urlset>\n`,
        );
      } else {
        await writeFile(path.join(output, "robots.txt"), "User-agent: *\nDisallow: /\n");
        await rm(path.join(output, "sitemap.xml"), { force: true });
      }
    },
  };
}

export default defineConfig({
  base: "./",
  plugins: [react(), seoPlugin()],
  build: {
    outDir: "dist",
    sourcemap: true,
  },
});
