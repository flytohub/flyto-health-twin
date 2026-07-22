/** Validate built metadata, crawler policy, assets, and public health data. */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const webRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const output = path.join(webRoot, "dist");
const html = fs.readFileSync(path.join(output, "index.html"), "utf8");
const robots = fs.readFileSync(path.join(output, "robots.txt"), "utf8");

/** Stop the build with one explicit contract error. */
function require(condition, message) {
  if (!condition) {
    console.error(`build validation: ${message}`);
    process.exit(1);
  }
}

require(!html.includes("__SITE_URL__") && !html.includes("__ROBOTS__"), "SEO placeholders remain");
for (const marker of [
  "rel=\"canonical\"",
  "property=\"og:title\"",
  "property=\"og:url\"",
  "property=\"og:image\"",
  "name=\"twitter:card\"",
  "application/ld+json",
  "Flyto2 Health Twin",
]) {
  require(html.includes(marker), `index.html is missing ${marker}`);
}

const canonical = html.match(/rel="canonical" href="([^"]+)"/)?.[1];
const ogUrl = html.match(/property="og:url" content="([^"]+)"/)?.[1];
require(Boolean(canonical) && canonical === ogUrl, "canonical and og:url differ");
require(html.includes(`\"url\": \"${canonical}\"`), "structured-data URL differs from canonical");

const sitemapPath = path.join(output, "sitemap.xml");
if (robots.includes("Allow: /")) {
  require(html.includes("index,follow,max-image-preview:large"), "production build is not indexable");
  require(fs.existsSync(sitemapPath), "indexable build has no sitemap");
  const sitemap = fs.readFileSync(sitemapPath, "utf8");
  require(sitemap.includes(`<loc>${canonical}</loc>`), "sitemap URL differs from canonical");
  require(robots.includes(`Sitemap: ${canonical}sitemap.xml`), "robots sitemap differs from canonical");
} else {
  require(html.includes("noindex,nofollow"), "preview build must be noindex");
  require(!fs.existsSync(sitemapPath), "preview build must not publish a sitemap");
}

for (const match of html.matchAll(/(?:src|href)="\.\/([^"#?]+)"/g)) {
  require(fs.existsSync(path.join(output, match[1])), `missing built asset ${match[1]}`);
}

const publicData = JSON.parse(fs.readFileSync(path.join(output, "public-data.json"), "utf8"));
require(publicData.project === "flyto2", "public export project is not flyto2");
require(Array.isArray(publicData.records) && publicData.records.length >= 3, "public records are missing");
require(Array.isArray(publicData.evaluations) && publicData.evaluations.length > 0, "public evaluations are missing");
const forbiddenRecordFields = [
  "notes",
  "weight_kg",
  "blood_pressure_systolic",
  "blood_pressure_diastolic",
  "blood_glucose_mgdl",
  "body_temperature_c",
  "illness_score",
  "training_load",
];
const publicRecordObjects = [
  ...publicData.records,
  ...publicData.evaluations.map((evaluation) => evaluation.actual),
];
for (const record of publicRecordObjects) {
  for (const field of forbiddenRecordFields) {
    require(!Object.hasOwn(record, field), `public record leaked ${field}`);
  }
}

console.log(`build validation: PASS (${publicData.records.length} public records)`);
