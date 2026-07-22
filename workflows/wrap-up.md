# Wrap-Up Workflow

Use this before finishing a task.

1. Run `make verify`, web dependency audit, documentation strict audit, and
   Flyto2 Indexer strict verification.
2. Scan all files for retired branding and non-`@flyto2.com` public addresses.
3. For UI changes, verify mobile and desktop bounds, images, loading/error
   states, and browser console output.
4. Confirm preview noindex and production-shaped canonical/robots/sitemap gates.
5. Update state, tasks, changelog, and handoff; fetch, commit, push `main`, and
   verify the remote SHA.
