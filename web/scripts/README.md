# Web Build Gates

`generate-reference.mjs` checks JSDoc and the generated TypeScript inventory.
`validate-build.mjs` checks deployment metadata, crawler policy, built assets,
public payload shape, and forbidden health record keys.

Both run through `npm run check` or `npm run build`.
