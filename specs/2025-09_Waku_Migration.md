# Waku Migration Plan (2025-09)

This document records the decisions, known risks, and concrete task plan to migrate the web app from Create React App (CRA) to Waku (React + server-first). It is tailored to the current codebase structure and constraints.

Contents
- Context and scope
- Decisions (target architecture)
- Risks and assumptions
- Inventory of CRA → Waku hot spots (project-specific)
- Migration plan
  - Step 1: Anticipation (safe, forward-compatible refactors)
  - Step 2: Swap & Stabilize (cut over to Waku and ship)
  - Step 3: Completion & Cleanup (reach parity and remove legacy)
- Acceptance criteria (phase gates)
- Open items and additional files to collect
- Appendix: Notes and gotchas

## Context and scope

Current stack highlights (from the repository):
- React 18 (with React Router v6.2, BrowserRouter).
- Application providers: MUI v5 theme (Emotion), CssBaseline, LocalizationProvider (dayjs, fr), ApplicationContext.
- Authentication:
  - Tokens stored in localStorage (REFRESH_TOKEN_KEY).
  - Axios interceptor injects Bearer token (DPhotoApplication).
  - Token refresh scheduled with setTimeout; dispatch through security-reducer.
  - Relative API paths: /oauth/token, /oauth/logout.
- Global state: ApplicationContext with a per-app singleton DPhotoApplication (axios instance created at construction).
- Runtime config: Client fetch to /env-config.json inside ApplicationContextComponent useEffect.
- PWA: Manual registration of /serviceworker.js from public/index.html at window load.
- Dev proxy: package.json sets "proxy": "http://127.0.0.1:8080".
- Storybook: Webpack 5 builder + CRA preset.

Scope: Migrate to Waku with SSR streaming, keeping React Router and most app code, minimizing behavior change in Step 1 and Step 2. Deeper auth changes (cookies/HttpOnly) are optional in Step 3.

## Decisions (target architecture)

Runtime and rendering
- Start on Waku Node runtime (not Edge) to simplify axios and Node-compatible libs. Revisit Edge after stabilization.
- SSR with React 18 streaming. Use:
  - StaticRouter on server with the request URL.
  - BrowserRouter on client.
- Instantiate one DPhotoApplication + axios per request on the server to avoid cross-request leakage.

Auth (phased)
- Step 1–2: Retain localStorage-based refresh/access tokens and Authorization header injection via axios interceptor on the client. Avoid reading localStorage during SSR.
- Step 3 (optional): Move refresh (and ideally access) tokens to HttpOnly cookies; reduce or remove client-managed Authorization header.

Configuration
- Step 1: Keep /env-config.json.
- Step 2: Serve config from Waku server (embed in HTML as window.__ENV__ or via a config route), but keep CRA-compatible path temporarily if needed.
- Step 3: Remove /env-config.json reliance on the client; prefer server-injected initial config/state.

Styling
- Use Emotion SSR for MUI v5 to avoid FOUC and className mismatch. Extract critical CSS on server and inject into HTML head.

Routing and SPA behavior
- Preserve URLs and client-side navigation. Implement SSR for:
  - /albums
  - /albums/:owner/:album
  - /albums/:owner/:album/:encodedId/:filename

Static assets and PWA
- Serve /public assets via Waku’s static handler. Handle /serviceworker.js at the root path.
- Consider disabling service worker temporarily during Swap if it interferes with SSR navigation.

Build and dev
- Replace CRA dev server and proxy with Waku dev server and server routes for API proxy (if needed).
- Keep Storybook on webpack during migration; switch to builder-vite in Step 3.

## Risks and assumptions

- SSR hydration mismatches (dayjs locale, time zone, router state, dynamic content).
- Emotion SSR setup errors may cause FOUC or mismatched class names.
- Singleton axios/DPhotoApplication leaking tokens across requests if not re-scoped per request.
- NodeJS.Timeout typing issues across Node/browser when compiling with SSR.
- Service worker may cache SSR responses and interfere with navigation during canary rollout.
- Dev proxy removal may break relative API calls unless server routes are in place.

Assumptions:
- Backend remains reachable at the same base URLs (or via server proxy).
- We can add minimal Waku server code and deploy alongside current infra (blue/green).

## Inventory of CRA → Waku hot spots (project-specific)

Routing and providers
- BrowserRouter used globally in src/App.tsx; will need StaticRouter on server.
- Providers stack: DPhotoTheme, CssBaseline, LocalizationProvider, ApplicationContextComponent must be SSR-safe.

Auth and side effects
- localStorage accesses in AuthenticateCase and LogoutCase.
- setTimeout scheduling in AuthenticateCase, security-reducer uses NodeJS.Timeout.
- DPhotoApplication holds axios interceptor; needs per-request lifetime and proper eject on revokeAccessToken.

Configuration
- axios.get("/env-config.json") in ApplicationContextComponent; SSR needs server-provided config.

i18n/date/locale
- dayjs.locale(fr) executed at module import; ensure server also sets locale before render to avoid mismatch.

PWA and static
- Manual service worker registration in public/index.html; may need to be deferred to client entry to avoid SSR contamination and easier toggling.

Development proxy and environment
- package.json "proxy" must be removed. Provide Waku proxy routes for /oauth/token and /oauth/logout or configure a reverse proxy.

Testing and Storybook
- CRA preset in Storybook; will need to move off in Step 3.

## Migration plan

### Step 1 — Anticipation (forward-compatible, no behavior change)

Routing and providers
- Extract AppProviders: Wrap DPhotoTheme, CssBaseline, LocalizationProvider, ApplicationContextComponent into a component that accepts a Router component/props. Keep BrowserRouter in the CRA entry; SSR will provide StaticRouter later.
- Add an ErrorBoundary around GeneralRouter (TODO already noted).

Auth and state safety
- Replace NodeJS.Timeout with ReturnType<typeof setTimeout> in:
  - web/src/core/security/security-reducer.ts
  - web/src/core/security/AuthenticateCase.ts
- Abstract token storage:
  - Define interface TokenStorage { get(): string|null; set(v: string): void; remove(): void }
  - Implement LocalStorageTokenStorage as default.
  - Update AuthenticateCase and LogoutCase to receive TokenStorage (inject via hook/factory). Default to LocalStorageTokenStorage in browser.
- Harden DPhotoApplication lifecycle:
  - Ensure renewRefreshToken sets interceptor once per instance.
  - Ensure revokeAccessToken ejects interceptor and clears internal state (and can be called multiple times).
  - Avoid any static/singleton DPhotoApplication; prepare for per-request instantiation.

Configuration
- Introduce ConfigService interface with two adapters:
  - BrowserConfigService → fetches /env-config.json (current behavior).
  - ServerConfigService (placeholder) → reads process.env and will inject initial config in Step 2.
- ApplicationContextComponent: read config via ConfigService injected from a context/provider with a default BrowserConfigService.

i18n and locale
- Ensure dayjs.locale('fr') is called in a centralized module and that SSR server can call it before rendering. No changes to functionality.

PWA registration
- Move service worker registration from public/index.html script tag to a client entry helper (e.g., src/serviceWorkerRegistration.ts). Keep behavior the same in CRA environment.
- Keep public/serviceworker.js path and logic unchanged for now.

Build and tooling
- No immediate change to Storybook or tests; log current assumptions and defer.

Deliverables Step 1
- New interfaces: TokenStorage, ConfigService.
- Updated types for timeouts.
- AppProviders extraction.
- Client-only SW registration helper.

### Step 2 — Swap & Stabilize (adopt Waku and ship)

Waku scaffolding
- Add Waku server entry:
  - Static asset serving for web/public.
  - SSR handler that renders the React app with:
    - StaticRouter using request URL.
    - A per-request DPhotoApplication instance.
    - Emotion SSR extraction for MUI v5 styles injected into head.
    - Initial config injection (window.__ENV__ or serialized initial state).
- Add proxy routes (temporary) on the Waku server:
  - POST /oauth/token → forwards to current backend.
  - POST /oauth/logout → forwards to current backend.

Client bootstrap
- Replace /env-config.json fetch with server-injected config on first load. Keep a compatibility route if needed.
- Ensure LocalStorageTokenStorage is only used on the client. Server uses a NoopTokenStorage.

Routing and hydration
- Use the AppProviders with BrowserRouter on the client, StaticRouter on server.
- Verify AuthenticatedRouter and GeneralRouter SSR safety:
  - RestoreAPIGatewayOriginalPath: consider handling ?path redirect on the server for full parity, or render a stable shell and immediate client redirect.

Build scripts
- package.json scripts:
  - dev: Start Waku dev server.
  - build: Build Waku server and client bundles.
  - start: Run the built Waku server.
  - Remove CRA-specific start/build/test during cutover (or keep temporarily, but CI should use Waku).
- Remove "proxy" field from package.json after proxy routes are set on Waku.

Stabilization
- Fix any hydration mismatches:
  - Check components that depend on useLocation/useSearchParams.
  - Ensure dayjs locale matches on server and client.
  - Confirm no browser-only APIs run during SSR render.
- Evaluate service worker during SSR rollout:
  - Disable registration if it causes navigation cache issues; re-enable after stabilization.

Deliverables Step 2
- Working Waku SSR server with streaming, static assets, and Emotion SSR.
- The app renders and navigates for /albums and nested album/media routes.
- Authentication functions on client (login, refresh, logout).
- CI updated to build and start Waku app.

### Step 3 — Completion & Cleanup (parity and de-CRA)

Auth hardening (optional but recommended)
- Switch refresh (and possibly access) tokens to secure HttpOnly cookies set by backend:
  - Update AuthenticationAPIAdapter to rely on same-origin cookie auth.
  - Simplify axios interceptor (might be unnecessary if access token is server-managed).
- Remove LocalStorageTokenStorage usage once cookies are adopted.

Configuration finalization
- Remove client fetch to /env-config.json.
- Store initial config in server render and hydrate into context.

Storybook and testing
- Migrate Storybook to 7+ with builder-vite; remove @storybook/preset-create-react-app.
- Keep Jest or move to Vitest. Ensure ESM/TS transforms and JSDOM settings are correct. Update setupTests.ts as needed.

Cleanup
- Remove CRA-specific dependencies (react-scripts, CRA eslint preset).
- Remove package.json "proxy".
- Confirm browserslist/polyfills strategy with the new bundler/runtime.
- Review and remove any dead code related to former CRA setup.

Deliverables Step 3
- Cookie-based auth (if opted).
- Finalized config approach.
- Updated Storybook and test tooling.
- Removal of CRA artifacts.

## Acceptance criteria (phase gates)

Step 1 (Anticipation) done when:
- No behavioral changes; app still runs under CRA.
- Code compiles with ReturnType<typeof setTimeout>.
- TokenStorage and ConfigService abstractions exist and are wired.
- AppProviders extracted and used by CRA entry.
- Service worker registration moved to client helper but still active.

Step 2 (Swap & Stabilize) done when:
- App builds and runs with Waku SSR.
- Initial page requests SSR-render the HTML for /albums and nested routes without hydration errors.
- Login, token refresh, and logout work on client.
- Static assets (favicon, manifest, service worker, images) are served correctly.
- CI uses Waku build and start; CRA proxy removed.

Step 3 (Completion & Cleanup) done when:
- CRA-specific dependencies and configs are removed.
- Config served and hydrated from server only.
- Storybook works post-migration; tests green.
- Optional: Cookie-based auth deployed and localStorage token storage removed.

## Open items and additional files to collect

To refine tasks and avoid surprises:
- public/serviceworker.js
- public/manifest.json
- GoogleLoginButton component
- CatalogAPIAdapter implementation (to check base URLs and assumptions)
- Any .env* files (variable names only; redact secrets)
- Root package.json (full) and lockfile
- tsconfig.json (root) and any path alias usage
- Any craco/config-overrides (confirm absence)

## Appendix: Notes and gotchas

- Emotion SSR (MUI v5): Use @emotion/server’s extractCriticalToChunks and constructStyleTagsFromChunks to inject styles. Ensure identical Emotion cache key on server and client to avoid className mismatch.
- NodeJS.Timeout vs browser: Use ReturnType<typeof setTimeout> everywhere shared between server and client.
- Avoid reading localStorage during SSR: AuthenticateCase.restoreSession should only run on client effects; server must not attempt to access window/localStorage.
- Service worker: During migration, SW can cache HTML and JS in ways that conflict with SSR changes. Be ready to disable temporarily.
- Dev proxy removal: Ensure Waku proxy routes (or reverse proxy) handle /oauth/token and /oauth/logout before removing CRA proxy to avoid 404s.
- dayjs locale: Call dayjs.locale('fr') on both server (early in request handler) and client (as today) to maintain consistent formatting.

Ownership
- Lead dev: approves architectural decisions, reviews SSR setup.
- Frontend dev(s): implement Anticipation refactors, SSR wiring, Emotion SSR.
- Platform/DevOps: update CI/CD, deploy blue/green Waku alongside CRA for canary.
- QA: verifies SSR routes, auth flows, and regression on catalog/media pages.

Rollout
- Blue/green: Deploy Waku server under a parallel host or path. Route a small percentage (or internal-only) to Waku initially. Ensure session stickiness or auth domain parity when testing.

Status
- Document initialized 2025-09. Update this spec as decisions change; link PRs to checklist items.
