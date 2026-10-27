# Cutover: delete `web/` and drop the `/nextjs` basepath

**Status**: done
**Depends on**: `01` through `09` — must only be picked up once every other issue in this feature is `done` and Thomas has validated `web-nextjs` as his primary interface.

## Description

Mechanical closure of the feature: retire the deprecated Waku-based `web/` project, drop the `/nextjs`
basepath from `web-nextjs`, and update the deployment so `web-nextjs` serves at `/`. Only pick this up when
issues 01–09 are all `done` and Thomas has validated `web-nextjs` as the primary interface.

## Context

- `web-nextjs/next.config.ts` currently sets `basePath: '/nextjs'`. The comment on that line warns that
  every `image src` using the basepath must be updated when it is removed — audit the codebase for any
  URL built on top of this basepath.
- `web/` is deployed as its own lambda (via `deployments/cdk/`) at the root path. The cutover switches that
  behaviour to serve `web-nextjs` at the root instead and removes the `web/` lambda.
- Google OAuth redirect URIs are configured out-of-band in the Google Cloud console. If any configured
  redirect URI contains `/nextjs`, it must be removed before deploy, otherwise sign-in breaks. Call this
  out explicitly in the PR description so Thomas can perform the console change — do not try to automate
  it from code.
- Documentation to update:
  - `AGENTS.md` — remove the `web/` section and the "deprecated" markers; `web-nextjs` becomes the only web
    project.
  - `README.md` — update to reflect the new structure.
  - `Makefile` — drop any `web/` targets.
  - `.github/workflows/` — drop any build/test/deploy steps that only concern `web/`.

## Acceptance Criteria

- The entire `web/` directory is deleted.
- No remaining reference to `web/` in `AGENTS.md`, `README.md`, `Makefile`, `.github/workflows/`, or
  `deployments/cdk/` (grep the repository).
- `basePath: '/nextjs'` is removed from `web-nextjs/next.config.ts`.
- Any URL-building code or image reference that assumed the `/nextjs` basepath is updated.
- `deployments/cdk/` is updated so `web-nextjs` serves at `/`; the `web/` lambda / CloudFront behaviour is
  removed; `npm run synth:test` passes.
- All test suites pass: `go test ./...` from the repo root, `cd api/lambdas && go test ./...`,
  `cd web-nextjs && npm run test`, `cd deployments/cdk && npm test && npm run synth:test`.
- The PR description lists any required human action (notably: removing `/nextjs` from configured Google
  OAuth redirect URIs) that must be performed before deploy.
- Once deployed, the home page is reachable at the root path and the sign-in flow still works.

## Out of scope

- Any functional change to `web-nextjs` — this story is purely retirement + path change.
- Changing the authentication provider (e.g. moving to Auth0 or Cognito) — a separate TODO item.
- Introducing a redirect from `/nextjs/*` to `/*` for old bookmarks — not worth it for a private family
  app; mention it in the PR for Thomas to decide.
