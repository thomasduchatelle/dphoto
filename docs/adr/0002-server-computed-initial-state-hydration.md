# ADR-0002: Server-computed initial state, client hydration

- Status: accepted
- Date: 2026-01-31
- Scope: `web-nextjs` (NextJS App Router)

## Context

The catalog state management (`CatalogViewerState`, reducer, thunks) is lifted-and-shifted from `web/src/core/catalog/` and is validated by 90+ tests. It must keep working unchanged, but NextJS App Router wants data loaded server-side for fast first paint, while interactions after load must run client-side.

The same thunk declarations must therefore execute in two contexts: on the server (once, to compute an initial state) and on the client (repeatedly, dispatching into a live React reducer).

## Decision

- Server Components compute the initial state by executing the relevant thunk (e.g. `onPageRefresh`) through a server-side runner, `constructThunkFromDeclaration`, which accumulates the actions the thunk dispatches and reduces them against `initialCatalogState()` to produce the final state.
- The computed state is passed as a prop to a Client Component, which hydrates `useReducer(catalogReducer, initialState)` — **no `useEffect` for initial load**.
- The Client Component instantiates thunks via `useThunks` for subsequent interactions.
- Adapters are context-specific: a **server** adapter factory reads the access token from the session service; a **client** adapter factory is used inside `useThunks`. Server and client adapter code must not cross-import.
- Pure UI components receive state slices and handlers as props and hold no state of their own.

```
Server Component → computeInitialState (constructThunkFromDeclaration) → props
                                                                          ↓
Client Component → useReducer(reducer, initialState) + useThunks(...) → props
                                                                          ↓
Pure UI components (no state)
```

## Consequences

- One set of thunks/reducer serves SSR and client; no reimplementation.
- First paint reflects loaded data without a client round-trip.
- `constructThunkFromDeclaration` is a reusable primitive relied on by every page that loads server-side state.
- Care is required to keep server-only code (session service) out of the client bundle.
