# ADR-0005: Browser reaches the API directly with its own access token

- Status: accepted
- Date: 18 Oct 2026
- Scope: `web-nextjs/libs/security/client`, `web-nextjs/app/auth/refresh`, `api/lambdas/authorizer`, `deployments/cdk/lib/access/lambda-authoriser-construct.ts`

## Context

The Next.js SPA calls the Go REST API behind an API Gateway v2 HTTP + a Lambda authorizer. Two caches are
on the critical path and both need stable identity keys to work:

- **Browser image cache**: the media grid and the full-screen viewer load thumbnails and originals through
  `<img src>`. If a token is encoded into the URL (as a query parameter or path segment), every token
  rotation invalidates the browser cache and every page reload re-downloads the same bytes.
- **AWS Lambda authorizer cache**: `HttpLambdaAuthorizer` keys its 30-minute cache on `identitySource`.
  With `$request.header.Authorization` as the identity source, the authorizer Lambda is invoked once per
  token value — a scrolling grid firing dozens of API calls per minute pays the Lambda cost once per
  session, not per call. With an empty identity source (as the media authorizer uses today), the Lambda
  is invoked on every request.

The browser cannot attach a custom `Authorization` header to an `<img src>` request — only cookies ride
automatically. Conversely, routing all API calls through the cookie would either break the authorizer
cache (empty identity source → one Lambda invocation per request) or key the cache on the full
`Cookie` header (churn from any unrelated cookie breaks the hit rate).

Prior iterations of this stack tried to reconcile the two with a single `cookie-or-header` authorizer and
a cookie-only client. The cache either misbehaved or had to be disabled. This ADR picks the one-shape-per-
endpoint resolution.

## Decision

- **Each endpoint accepts exactly one token transport.**
  - `GET /api/v1/owners/{owner}/medias/*` (image bytes): **cookie only**, served by the uncached
    `mediaAuthorizer`. Image URLs stay stable, browser cache works.
  - Every other `/api/v1/*` endpoint: **`Authorization: Bearer <token>` header only**, served by the
    cached `authorizer` with `identitySource: ['$request.header.Authorization']` (default in the current
    CDK). The 30-min cache is keyed per token value.
- **The access token lives in a cookie that is readable by both server and client.** `storeSession`
  writes `dphoto-access-token` with `httpOnly: false`. Both axes work from this single source of truth:
  the browser auto-attaches it to `<img>` requests, and the SPA reads it from `document.cookie` to
  build the `Authorization` header. No endpoint, no extra round trip, no second copy to keep in sync.
- **The refresh and ID tokens stay httpOnly.** They are long-lived (30 days). An exfiltrated refresh
  token is a persistent compromise; an exfiltrated access token is bounded by its 1-hour TTL.
- **A Next.js route handler, `POST /auth/refresh`, is the single mechanism to rotate the cookies from
  the browser.** It reuses `refreshSession` from `libs/security/access-token-service.ts` — the same code
  path the server-side proxy uses on navigation. On success it returns `204` with rotated cookies; on
  failure, `401` with session-clearing cookies.
- **Two primitives on the client, in `web-nextjs/libs/security/client`:**
  - `ensureFreshSession()` reads the current access-token cookie, decodes its JWT `exp` claim, and if
    the token is comfortably in the future resolves immediately. Otherwise it `POST`s `/auth/refresh`
    and awaits the server's reply. The server is trusted — the function does **not** re-validate the
    rotated cookie. If the HTTP call returns non-OK, `ensureFreshSession` throws
    `SessionExpiredException`. Concurrent callers share a single in-flight promise so parallel scroll
    events trigger exactly one refresh.
  - `getAccessToken()` awaits `ensureFreshSession()`, reads the cookie, and returns the token value
    (throws `SessionExpiredException` if the cookie is absent — in practice unreachable, because
    `ensureFreshSession` would have thrown first).
- **`SessionExpiredException` is the single failure signal.** `FetchCatalogAdapter` on the client throws
  it on any `401` response. The UI catches it at the thunk layer and routes to a dedicated
  `/auth/session-expired` page that offers a "Sign In Again" link back to `/`; the proxy then kicks off
  the normal OIDC flow.

### How the authentication flow works

1. **First login.** The user hits any protected page. The proxy sees no session cookies and starts the
   OIDC authorization code flow. The `/auth/callback` route exchanges the code for tokens and
   `storeSession` writes all three cookies (access non-httpOnly, refresh + id httpOnly).
2. **Revive from refresh-token-only.** The access-token cookie expires naturally (1 h `maxAge`) while
   the tab is open or on return. The next `ensureFreshSession()` call sees no access-token cookie,
   `POST /auth/refresh` with the still-present refresh-token cookie, server rotates all three cookies,
   the SPA proceeds.
3. **Keep a valid access token during activity.** UI call sites (lazy image loaders, viewer navigation)
   call `ensureFreshSession()` before kicking off work that depends on the cookie — cheap when the token
   is fresh, debounced when several events fire in the same frame. API call sites use `getAccessToken()`
   to produce the `Authorization` header — same in-flight dedup. The server-side proxy keeps refreshing
   on page navigation; the two paths are independent and never fight for the same cookie.

## Consequences

- The access token is reachable from JS, so an XSS vulnerability can exfiltrate it. The window is bounded
  by the 1-hour TTL and the refresh token stays protected. For this single-tenant app with no
  third-party JS and no user-generated HTML rendered unsafely, the trade-off buys an authorizer cache
  that is otherwise impossible to make work.
- Two authorizers remain in the CDK (`authorizer`, `mediaAuthorizer`) because the two transports require
  different cache profiles. The authorizer Lambda keeps extracting the token from either the
  `Authorization` header or the `dphoto-access-token` cookie — both paths are live, one per endpoint.
- UI call sites must opt into calling `ensureFreshSession()` where lazy loads happen (`<img>` rows,
  viewer next/prev). This is a one-liner but it is the responsibility of the view component, not the
  slice.
- Any thunk that calls the client adapter can throw `SessionExpiredException`. The UI's thunk wiring
  catches it once (`useSessionAwareThunks`) and navigates to `/auth/session-expired`, so individual
  thunks do not need to handle the case.
- If we ever change the browser-image-cache or authorizer-cache guarantees (e.g. move to server-side
  rendered images), the one-transport-per-endpoint rule stops being load-bearing and the cookie can go
  back to httpOnly. The ADR is the record to revisit at that point.
