# Client-side access-token refresh for long-lived SPA tabs

**Status**: ready
**Depends on**: `04-create-album-modal.md` — first write-story to land; subsequent write-stories (05, 06, 07, 08) should assume this is in place so they don't inherit the same 401-on-stale-tab assumption.
**Blocks**: `05`, `06`, `07`, `08` — ideally picked up before them, but not a hard gate; the modals degrade to "refresh the page and try again" without it.

## Description

The Next.js client talks to the REST API by relying on cookies auto-attached by the browser. The token
is only refreshed server-side by the proxy (`web-nextjs/proxy.ts`) on page navigation. Any client-driven
activity on a stale tab breaks:

- **Writes (Create album, Share, Rename, …)**: the POST returns 401 and the thunk surfaces an error. The
  user has to reload to recover.
- **Lazy image loads** (scrolling the media grid, stepping through the fullscreen viewer): the browser
  attaches the cookie automatically to `<img src>` requests, but if the access token has expired the
  images simply fail to load — no fetch wrapper ever sees them.

A single reactive "refresh on 401 and retry" wrapper is not enough: `<img>` loads cannot be retried once
the browser has decided the response is broken. We need both a **proactive** hook that the UI calls just
before kicking off work that depends on a live cookie, and a **reactive** wrapper for API calls.

## Context

Two orthogonal findings shape the design:

1. **The AWS Lambda authorizer in front of `/api/v1/*` is cached per Authorization-header value**
   (`identitySource: ['$request.header.Authorization']`, 30 min TTL). Routing API calls through the
   cookie alone would either break the cache (empty identity source → one Lambda invocation per request)
   or tie the cache key to the whole cookie header (churn from unrelated cookies). The client must send
   the access token in the `Authorization: Bearer <token>` header.

2. **The access token therefore has to be reachable from JS.** The simplest expression of that is to
   store it in a readable cookie (not httpOnly) — the browser handles persistence across tab reloads,
   the server rotates it the same way as before. An "endpoint to read the cookie" would add network
   overhead without changing the security profile: as soon as the token is reachable from JS, an XSS
   attacker can read it.

Security considerations driving the "readable access token" decision:
- **Refresh token stays httpOnly** (long-lived, 30-day cookie). An exfiltrated access token is bounded
  to its 1-hour TTL; an exfiltrated refresh token is a persistent compromise.
- **ID token stays httpOnly** (same 30-day lifetime as refresh).
- **XSS attack surface**: single-tenant app, no third-party JS, no user-generated HTML. The classic
  "bearer in memory" risk is near-zero here, and in exchange the authorizer cache keeps working.

## Design

- **`storeSession`** writes `dphoto-access-token` with `httpOnly: false` (readable by JS). Refresh token
  and user-info cookies stay httpOnly. `clearFullSession` clears the access-token cookie as non-httpOnly
  to match. `appliesCookies` honours a per-cookie `httpOnly` override (default `true`).
- **`POST /auth/refresh`** route handler at `web-nextjs/app/auth/refresh/route.ts`. Reuses
  `refreshSession` from `access-token-service.ts`. Returns `204` on success (with rotated cookies
  applied) or `401` on failure (with session-clearing cookies applied). Lives on the Next.js origin
  because the refresh token cookie is scoped there and httpOnly — the Go API has no visibility of it.
- **Client slice at `web-nextjs/libs/security/client/`** — a cohesive package with one integration test:
  - `readAccessTokenFromCookie()` — reads `document.cookie`, returns the raw token or `null`.
  - `ensureFreshSession({ thresholdSeconds })` — the proactive primitive. Reads the token from the
    cookie, decodes its JWT `exp` claim (no extra "hint cookie" — the JWT self-describes its expiry);
    if the token is comfortably in the future (default: 60 seconds), resolves immediately.
    Otherwise `POST`s `/auth/refresh` and resolves when it returns. Deduplicates: module-level inflight
    promise ensures N parallel calls trigger exactly one network refresh. Never throws.
  - `getAccessToken()` — the adapter-facing primitive. `await ensureFreshSession()`, then returns the
    current token from the cookie, or `undefined`.
  - `authenticatedFetch(input, init)` — reactive safety net. Sends the request as-is. On 401, force a
    refresh (`thresholdSeconds: Infinity`), rebuild the `Authorization: Bearer` header from the
    rotated cookie, retry **once**. Does not retry on any other status.
- **`FetchCatalogAdapter`** keeps its `accessTokenSupplier` contract; `newClientAdapterFactory` injects
  `getAccessToken` as the supplier and `authenticatedFetch` as the fetch impl. The adapter already
  builds the `Authorization: Bearer <token>` header from the supplier's return value. The server-side
  factory is untouched (server components read cookies through `loadSession` and the proxy refreshes).
- `web-nextjs/proxy.ts:6` already excludes `auth` from the proxy matcher, so `/auth/refresh` bypasses
  the middleware. Verified in `proxy.test.ts`.

## UI integration (follow-up stories)

- Mutation entry points: no change needed — they go through `FetchCatalogAdapter`, which uses
  `authenticatedFetch`.
- Lazy image loaders (media grid scroll, fullscreen next/prev): call `ensureFreshSession()` before
  mounting a new batch of thumbnails. The cookie is then fresh and the browser's auto-attached cookie
  on `<img src>` requests succeeds.

## Acceptance Criteria

- `storeSession` writes the access-token cookie as non-httpOnly; refresh-token and user-info stay
  httpOnly. `clearFullSession` clears the access-token cookie as non-httpOnly. `appliesCookies` honours
  a per-cookie `httpOnly` override (default `true`).
- A `POST /auth/refresh` route handler exists at `web-nextjs/app/auth/refresh/route.ts` and:
  - On success: returns `204` with the rotated cookies applied.
  - On failure: returns `401` with the session-clearing cookies applied.
- The client slice `web-nextjs/libs/security/client/` exports `ensureFreshSession`, `getAccessToken`,
  `authenticatedFetch`, `readAccessTokenFromCookie` from a single `index.ts`. HTTP calls to
  `/auth/refresh` are internal to the slice.
- `ensureFreshSession(opts?)`:
  - Decides freshness from the JWT's `exp` claim read off the readable access-token cookie.
  - Resolves immediately when fresh beyond `thresholdSeconds` (default 60s).
  - Otherwise POSTs `/auth/refresh` and resolves when it returns.
  - Deduplicates: N concurrent calls → one network call.
  - Never throws; failures resolve silently.
- `getAccessToken()`:
  - Awaits `ensureFreshSession()` then returns the current token or `undefined`.
- `authenticatedFetch(input, init)`:
  - Passes the request through unchanged.
  - On 401, forces a refresh, rebuilds the `Authorization` header from the rotated cookie, retries once.
  - Does not retry on non-401 or more than once.
- `FetchCatalogAdapter` and the client adapter factory use `getAccessToken` and `authenticatedFetch`.
- One integration-style test under `web-nextjs/libs/security/client/` exercises the whole slice: cookie
  reads, proactive refresh at threshold, dedup across concurrent calls, 401 retry with rotated header,
  no retry on 500, silent failure.
- `access-token-service.test.ts` asserts the service-level clearing path (including the non-httpOnly
  access-token flag).
- `cd web-nextjs && npx vitest run` passes.
- Manual sanity check in the PR description: open a tab, let the access-token expire, click Create
  album, verify the request succeeds without a page reload; scroll a long media list past the expiry
  and verify thumbnails continue to load.

## Out of scope

- Any change to the server-side (`newServerSideRestCatalogAdapter`) refresh path — server components
  continue to rely on the proxy.
- Switching media URLs to a query-param token. Image loads continue to use the auto-attached cookie
  (the `mediaAuthorizer` is intentionally uncached and extracts the token from the cookie header).
- Any UI change: wiring `ensureFreshSession()` into scroll/viewer call sites belongs to the
  stories that touch those UIs.
