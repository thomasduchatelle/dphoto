# Client-side access-token refresh for long-lived SPA tabs

**Status**: ready
**Depends on**: `04-create-album-modal.md` — first write-story to land; subsequent write-stories (05, 06, 07, 08) should assume this is in place so they don't inherit the same 401-on-stale-tab assumption.
**Blocks**: `05`, `06`, `07`, `08` — ideally picked up before them, but not a hard gate; the modals degrade to "refresh the page and try again" without it.

## Description

The Next.js client adapters talk to the REST API by relying on the httpOnly `dphoto-access-token` cookie
being auto-attached by the browser (`web-nextjs/domains/catalog/adapters/client-adapter-factory.ts` returns
`undefined` from `accessTokenSupplier`). The access token is refreshed server-side by the proxy
(`web-nextjs/proxy.ts`) only on page navigation, when `aboutToExpire` becomes true.

Consequence: a user keeps a page open past the access-token TTL, clicks a write action (Create album,
Share, Rename, Edit dates, Delete), the POST returns 401, the thunk surfaces an error, and the user has to
reload the page to recover. The data cookies are still valid (refresh token has a much longer TTL) — the
UX break is purely a missing client-side refresh trigger.

This story closes that gap once, in a shared place, so every current and future client-side mutation
inherits a working refresh path.

## Context

- Server-side refresh already exists and must be reused: `refreshSession` in
  `web-nextjs/libs/security/access-token-service.ts` reads the refresh token from cookies, calls the
  OIDC provider, and returns the new cookies to set.
- The proxy at `web-nextjs/proxy.ts:29` wires `refreshSession` on navigation. It is the model to follow —
  **do not duplicate the refresh logic**; expose it through a route handler.
- `appliesCookies` from `web-nextjs/libs/nextjs-cookies` is how refreshed cookies are attached to a
  response. The new route handler must use it so the httpOnly cookies get rotated the same way.
- The client adapter must stay cookie-based: `accessTokenSupplier` keeps returning `undefined`. We are
  **not** moving the access token into JS-readable storage — that would weaken the httpOnly guarantee.
- Shape to implement:
  1. A Next.js Route Handler at `web-nextjs/app/auth/refresh/route.ts` accepting `POST`, calling
     `refreshSession`, and returning `204` on success (with rotated cookies applied) or `401` on failure
     (no body needed; cookies cleared on failure as `refreshSession` already does).
  2. A small shared `fetch` wrapper in `web-nextjs/libs/requests/` that, on any `401` from an authenticated
     call, POSTs `/auth/refresh` once, then retries the original request. If the refresh also returns
     `401`, let the original `401` propagate (the next navigation will push the user through the proxy's
     auth flow). The wrapper must not retry on non-401 errors and must not retry more than once per call.
  3. The catalog `FetchCatalogAdapter` (and any other client adapter) uses this wrapper instead of raw
     `fetch` so every mutation inherits the behaviour.
- Keep the route handler path aligned with the existing `/auth/*` convention (see `web-nextjs/app/auth/`)
  and make sure the proxy's `skipProxyForPageMatching` regex continues to let `/auth/refresh` through
  without triggering the proxy — it already excludes `auth`, verify this.
- Concurrency: if two requests race and both see a 401, both will try to refresh. A single in-flight
  refresh promise shared across callers (module-level `let inflight: Promise<Response> | null`) is enough
  — do not pull in a dedicated library for this.

## Acceptance Criteria

- A `POST /auth/refresh` route handler exists at `web-nextjs/app/auth/refresh/route.ts` and:
  - On success: returns `204` with the rotated `dphoto-*` cookies applied via `appliesCookies`.
  - On failure: returns `401` with the session-clearing cookies applied (same semantics as the proxy's
    failure path).
  - Rejects any method other than `POST` with `405`.
- A shared client fetch wrapper lives under `web-nextjs/libs/requests/` and:
  - On a `401` response, calls `POST /auth/refresh` and retries the original request once.
  - Does not retry on any other status code.
  - Does not retry more than once per original call.
  - Deduplicates concurrent refreshes: N parallel calls that all 401 trigger exactly one refresh.
  - If the refresh itself returns `401`, the original `401` is returned to the caller unchanged.
- `FetchCatalogAdapter` (and any other client-side adapter in `web-nextjs/domains/*/adapters/`) uses the
  wrapper. `accessTokenSupplier` stays as it is on the client (returns `undefined`); the server-side
  adapter is untouched.
- Unit tests cover: happy path (200 passes through), 401 → refresh 204 → retry 200, 401 → refresh 401 →
  original 401 returned, two concurrent 401s trigger one refresh, non-401 errors do not trigger refresh.
- `cd web-nextjs && npm run test` passes.
- Manual sanity check documented in the PR description: leave a tab open until the access token expires,
  click Create album, verify the request succeeds without a page reload.

## Out of scope

- Any change to the server-side (`newServerSideRestCatalogAdapter`) refresh path — server components
  continue to rely on the proxy.
- Pre-emptive refresh (refreshing before a request based on a stored expiry) — a reactive retry-on-401
  is simpler and only pays the cost when it matters.
- Moving the access token out of the httpOnly cookie — explicitly rejected: keeps the token out of JS
  reach.
- Any UI change: the behaviour is invisible to the user when it works, and when the refresh fails the
  existing error path (thunk → error state) is good enough.
