# Client-side access-token refresh for long-lived SPA tabs

**Status**: ready
**Depends on**: `04-create-album-modal.md` — first write-story to land; subsequent write-stories (05, 06, 07, 08) should assume this is in place so they don't inherit the same 401-on-stale-tab assumption.
**Blocks**: `05`, `06`, `07`, `08` — ideally picked up before them, but not a hard gate; the modals degrade to "refresh the page and try again" without it.

## Description

The Next.js client talks to the REST API by relying on the httpOnly `dphoto-access-token` cookie being
auto-attached by the browser. The token is only refreshed server-side by the proxy
(`web-nextjs/proxy.ts`) on page navigation. Any client-driven activity on a stale tab breaks:

- **Writes (Create album, Share, Rename, …)**: the POST returns 401 and the thunk surfaces an error. The
  user has to reload to recover.
- **Lazy reads**: scrolling the media grid or stepping through the fullscreen viewer loads thumbnails and
  originals through `<img src>`. The browser attaches the cookie automatically, but if the access token
  has expired the images simply fail to load — the fetch wrapper never sees them.

A single reactive "refresh on 401 and retry" wrapper is not enough: `<img>` loads cannot be retried once
the browser has decided the response is broken. We need both a **proactive** hook that the UI calls just
before kicking off work that depends on a live cookie, and a **reactive** wrapper for API calls.

The refresh token cookie is httpOnly and scoped to the Next.js origin, so refresh must happen through a
Next.js route handler — the Go API has no visibility of these cookies. The route handler is a thin
wrapper around the existing `refreshSession` in `web-nextjs/libs/security/access-token-service.ts`.

## Context

- **Reuse `refreshSession`** (`web-nextjs/libs/security/access-token-service.ts`): it reads the refresh
  token from cookies, calls the OIDC provider, and returns rotated cookies. Do not duplicate this logic.
- **The access token cookie stays httpOnly**. `accessTokenSupplier` keeps returning `undefined` on the
  client; the browser attaches the cookie automatically. We are not moving the token into JS reach.
- **Expiry hint cookie**: `storeSession` in `web-nextjs/libs/security/backend-store.ts` writes a new
  non-httpOnly cookie `dphoto-access-token-expires-at` (unix seconds, string) alongside the three
  httpOnly cookies. This is the only piece of session information the client-side TypeScript reads
  directly. It carries no secret — just a timestamp — and lets `ensureFreshSession()` decide whether to
  trigger a refresh without a round trip. `clearFullSession` deletes it the same way as the other session
  cookies. Rotated uniformly by the proxy, `/auth/callback`, and `/auth/refresh`.
- **Per-cookie `httpOnly` override**: `appliesCookies` in `web-nextjs/libs/nextjs-cookies/` hard-codes
  `httpOnly: true` for every cookie it writes. Extend `CookieValue` with an optional `httpOnly` field so
  the expiry hint can be written as non-httpOnly while every other cookie keeps the secure default.
- **Route handler** at `web-nextjs/app/auth/refresh/route.ts`: `POST` only, calls `refreshSession`,
  returns `204` on success (rotated cookies applied via `appliesCookies`) or `401` on failure
  (session-clearing cookies applied — same semantics as the proxy's failure path). Other methods → `405`.
- **`ensureFreshSession()`** — the proactive primitive, lives in
  `web-nextjs/libs/security/client/ensure-fresh-session.ts` ("use client"). Reads
  `dphoto-access-token-expires-at` from `document.cookie`; if the token is comfortably in the future,
  resolves immediately. Otherwise `POST /auth/refresh` and resolves when the server responds. The default
  freshness threshold is 60 seconds; call sites can pass `thresholdSeconds: Infinity` to force a refresh.
  Deduplicates: a module-level inflight promise ensures N parallel calls trigger exactly one network
  refresh. On failure resolves silently — the caller is best-effort; the next real request will see the
  401 and the wrapper (or the proxy, on next navigation) will handle it.
- **`authenticatedFetch()`** — the API wrapper, lives in `web-nextjs/libs/requests/authenticated-fetch.ts`.
  Pre-request: `await ensureFreshSession()` (no-op when fresh, one refresh otherwise). On `401`:
  `await ensureFreshSession({ thresholdSeconds: Infinity })` to force a refresh, then retry the original
  request once. Does not retry on any non-401 status, nor more than once. If the forced refresh also
  fails, the original `401` propagates — the next navigation will push the user through the proxy's
  auth flow.
- **`FetchCatalogAdapter`** uses `authenticatedFetch` instead of raw `fetch`. The server-side adapter
  factory is untouched (server components continue to rely on the proxy for refresh). The client adapter
  factory is wired to `authenticatedFetch`.
- `web-nextjs/proxy.ts:6` already excludes `auth` from the proxy matcher, so `/auth/refresh` bypasses the
  middleware. Verified in `proxy.test.ts`.

## Acceptance Criteria

- `storeSession` writes a `dphoto-access-token-expires-at` cookie with the absolute expiry (unix
  seconds, string), `httpOnly: false`, same `maxAge` as the access token. `clearFullSession` clears it.
  `appliesCookies` honours a per-cookie `httpOnly` override (default `true`).
- A `POST /auth/refresh` route handler exists at `web-nextjs/app/auth/refresh/route.ts` and:
  - On success: returns `204` with the rotated `dphoto-*` cookies (including the expiry hint) applied.
  - On failure: returns `401` with the session-clearing cookies applied.
  - Rejects any method other than `POST` with `405`.
- `ensureFreshSession(opts?)` lives under `web-nextjs/libs/security/` client bundle and:
  - Reads `dphoto-access-token-expires-at` from `document.cookie`.
  - Resolves immediately when the expiry is further than `thresholdSeconds` (default 60s) in the future.
  - Otherwise `POST`s `/auth/refresh` and resolves when it returns.
  - Deduplicates: N concurrent calls trigger exactly one refresh.
  - Never throws; failures resolve silently.
- `authenticatedFetch(input, init)` lives under `web-nextjs/libs/requests/` and:
  - Calls `ensureFreshSession()` before every request (proactive).
  - On a `401` response, calls `ensureFreshSession({ thresholdSeconds: Infinity })` to force a refresh
    and retries the original request once (reactive safety net).
  - Does not retry on any other status code and does not retry more than once per original call.
  - Shares the inflight promise with the proactive path so the "two parallel 401s" case still triggers a
    single refresh.
- `FetchCatalogAdapter` and the client adapter factory use `authenticatedFetch`. `accessTokenSupplier`
  stays returning `undefined`.
- Unit tests cover:
  - Route handler: `POST` success (204 + rotated cookies), `POST` failure (401 + cleared cookies), other
    methods → 405.
  - `ensureFreshSession`: no cookie → attempts refresh; fresh enough → no network call; expiring →
    refresh; two concurrent calls → one network call; `thresholdSeconds: Infinity` → always refresh.
  - `authenticatedFetch`: happy path (200), proactive refresh when expiry is close, reactive refresh on
    401, no retry on non-401, no double retry.
- `storeSession` tests updated to assert the hint cookie is written.
- `cd web-nextjs && npx vitest run` passes.
- Manual sanity check documented in the PR description: leave a tab open until the access-token expiry
  cookie is close to zero, scroll a long media list, verify thumbnails continue to load; click Create
  album, verify the request succeeds without a page reload.

## Out of scope

- Any change to the server-side (`newServerSideRestCatalogAdapter`) refresh path — server components
  continue to rely on the proxy.
- Pre-emptive refresh *inside* `authenticatedFetch` based on a stored expiry is already covered by
  `ensureFreshSession`; no second mechanism.
- Moving the access token out of the httpOnly cookie — explicitly rejected.
- Any UI change: the behaviour is invisible to the user when it works.
