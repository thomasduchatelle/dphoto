# Story 1.3: Basic Album List Display

**Status**: ready-for-dev

---

## Story

As a **user**,  
I want **to view my album list on the authenticated home page**,  
so that **I can see all albums I own and albums shared with me**.

## Acceptance Criteria

**Given** I am an authenticated user with access to albums
**When** I navigate to the home page at `/`
**Then** I see a list of albums displayed as clickable text links
**And** each album displays its name (FR4)
**And** albums are displayed in chronological order (newest first) (FR6)
**And** I can see albums I own (FR1) and albums shared with me (FR2)
**And** clicking an album navigates to `/albums/[ownerId]/[albumId]`
**And** when no albums exist, I see a message inviting me to create an album (FR38)
**And** if an error occurs loading albums, I see an error message with a "Try Again" button (FR39)

## Scope & Context

### Depends on

* **Story 1.1 - Project Foundation Setup** (`specs/stories/1-1-project-foundation-setup.md`): Material UI theme configured with brand color #185986, dark
  theme (#121212 background, #1e1e1e surface), breakpoint system established, ThemeProvider wrapping application
* **Story 1.2 - State Management Migration** (`specs/stories/1-2-state-management-migration.md`): Complete catalog state management migrated to
  `domains/catalog/`, FetchCatalogAdapter implementing server/client compatible API access, custom image loader configured, error boundaries created, 230 tests
  passing validating behavior preservation

### Expected outcomes

The DEV agent implementing Story 1.4 will need from this story:

* A working **server-side state initialization pattern** that loads albums using the migrated state management and passes initial state to client components
* A **CatalogProvider client component** that hydrates server state, instantiates thunks client-side, and provides state + handlers to child components
* A **HomePageContent pure component** that renders the album list accepting albums and handlers as props
* The **routing structure** for album navigation (`/albums/[ownerId]/[albumId]`) to be established for clickable album links

This story establishes the fundamental data flow pattern (Server → Client → Pure UI) that Story 1.4 will enhance with styled album cards, and Stories 1.5-1.6
will build upon with visual polish and filtering.

### Out of scope

* DO NOT implement styled album cards with photo thumbnails - that's Story 1.4
* DO NOT implement density indicators, sharing avatars, or owner information display - that's Story 1.4
* DO NOT implement album filtering by owner - that's Story 1.6
* DO NOT implement random photo highlights - that's Epic 5
* DO NOT implement responsive grid layouts with multiple columns - that's Story 1.4
* DO NOT implement album management dialogs (create, edit, delete) - that's Epic 3
* DO NOT implement the actual album page at `/albums/[ownerId]/[albumId]` - that's Epic 2

---

## Technical Design

## Technical Design

### Overview

This story implements the authenticated home page displaying the album list using the **Server Component + Client Component pattern** defined in Architecture.md
Decision #2. The Server Component fetches initial data and executes thunks server-side, then passes the computed state to a Client Component that hydrates the
useReducer for subsequent user interactions.

### Architecture Pattern: Server-Side Thunk Execution

**Critical requirement from Story 1.2**: Story 1.2 introduced the need for **server-side thunk execution** to compute initial state before passing to client
components. This story implements that pattern for the first time.

From Story 1.2 MUST DO section:

```tsx
// Server Component executes thunk to compute initial state
export async function Page() {
    const loadedCatalogState = await constructThunkFromDeclaration(
        catalogThunks.onPageRefresh,
        initialCatalogState(),
        {adapterFactory: newServerAdapterFactory()}
    )
    return <CatalogProvider initialState={loadedCatalogState}>
        {/* Client component receives loaded state */}
    </CatalogProvider>
}
```

### Component Architecture Strategy

Following Architecture.md Decision #3 (Colocation Principle):

- **Server Component** (`app/(authenticated)/page.tsx`): Fetches data, executes thunks server-side, passes state as props
- **Client Component** (`_components/AlbumListClient.tsx`): Hydrates state with useReducer, provides handlers for future interactions
- **Pure UI Components** (`_components/AlbumCard.tsx`, `_components/AlbumsGrid.tsx`): Receive data and handlers as props, NO internal state

### Server-Side Thunk Execution Implementation

**New infrastructure required** (from Story 1.2 MUST DO):

1. **`libs/dthunks/server/constructThunkFromDeclaration.ts`**:
    - Accepts: `ThunkDeclaration`, `initialState`, `factoryArgs`
    - Creates a dispatch function that collects actions
    - Instantiates the thunk with server adapters
    - Reduces collected actions against initial state
    - Returns final computed state

2. **`domains/catalog/adapters/server-adapter-factory.ts`**:
    - Creates `FetchCatalogAdapter` with server-side `AccessTokenHolder`
    - Uses `getAccessTokenHolder()` from `libs/security/session-service.ts`
    - Enables thunks to call API from server context

### Data Flow Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ SERVER COMPONENT: app/(authenticated)/page.tsx             │
│                                                             │
│  1. Get initial state                                       │
│     const state = initialCatalogState(authenticatedUser)   │
│                                                             │
│  2. Execute onPageRefresh thunk (server-side)              │
│     const loadedState = await constructThunkFromDeclaration(│
│         catalogThunks.onPageRefresh,                       │
│         state,                                              │
│         {adapterFactory: newServerAdapterFactory()}        │
│     )                                                       │
│     // Internally: calls API, dispatches actions, reduces  │
│                                                             │
│  3. Pass computed state to Client Component                 │
│     <AlbumListClient initialState={loadedState} />        │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ CLIENT COMPONENT: _components/AlbumListClient.tsx          │
│ 'use client'                                                │
│                                                             │
│  1. Hydrate reducer with server state                       │
│     const [state, dispatch] = useReducer(                   │
│         catalogReducer,                                     │
│         props.initialState  // ← from server               │
│     )                                                       │
│                                                             │
│  2. Prepare handlers for future interactions                │
│     const handlers = useThunks(                             │
│         catalogThunks,                                      │
│         {adapterFactory: newClientAdapterFactory(), dispatch},│
│         state                                               │
│     )                                                       │
│                                                             │
│  3. Render pure UI components                               │
│     <AlbumsGrid albums={state.albums}                      │
│                 onAlbumClick={handlers.navigateToAlbum} /> │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ PURE UI: _components/AlbumsGrid.tsx                         │
│                                                             │
│  Props: albums[], onAlbumClick                              │
│  - Maps albums to AlbumCard components                      │
│  - Responsive MUI Grid layout                               │
│  - NO state, NO business logic                              │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ PURE UI: _components/AlbumCard.tsx                          │
│                                                             │
│  Props: album object, onClick handler                       │
│  - MUI Card displaying metadata                             │
│  - Brand color styling                                      │
│  - NO state, NO business logic                              │
└─────────────────────────────────────────────────────────────┘
```

### State Management Integration

From Story 1.2, we have:

- **CatalogViewerState**: Holds albums, selected album, dialogs, errors
- **catalogReducer**: Reduces actions to produce new state
- **catalogThunks**: Business logic coordinating API calls and actions
- **FetchCatalogAdapter**: Fetches data from REST API

This story uses:

- **catalogThunks.onPageRefresh**: Loads initial albums list
- **Actions**: `catalogRefreshStarted`, `catalogRefreshed`, `catalogRefreshFailed`
- **Selector**: `selectVisibleAlbums(state)` (returns albums filtered/sorted)

### Loading, Error, and Empty States

**Loading State**: Display skeleton cards while Server Component fetches data

- Implementation: NextJS automatic loading.tsx handling OR
- Manual: Check if `state.albums === undefined` → show `<LoadingSkeleton />`

**Error State**: Display error message with retry button

- Condition: `state.error !== null`
- Component: `<ErrorMessage error={state.error} onRetry={handlers.retry} />`
- Must provide "Try Again" button (FR40)

**Empty State**: Display helpful message when no albums exist

- Condition: `state.albums.length === 0`
- Component: `<EmptyState message="No albums found" />`
- Should guide user on next steps (FR39)

### Critical Success Factors

1. **Server-side thunk execution MUST work** - This is the foundational pattern described in Story 1.2's "MUST DO PART OF STORY 1.3" section. Without it, the
   entire architecture fails.

2. **State hydration MUST preserve behavior** - The 230 migrated tests validate state management works. Server-loaded state must hydrate client-side reducer
   without breaking any existing behavior.

3. **Simple UI ONLY** - Resist temptation to add styled cards, thumbnails, or visual complexity. Use Typography + Link components only. Story 1.4 adds the
   visual polish.

4. **Navigation MUST work** - Clicking album links must navigate to `/albums/[ownerId]/[albumId]` even though that page only shows a placeholder.

5. **Error handling MUST be complete** - Empty state (no albums), loading state, and error state with retry all functional.

### Implementation Guidance

This technical guidance has been validated by the lead developer. Following it significantly increases the chance of getting your PR accepted. Any deviation
required to complete the story must be reported.

#### Coding standards

You must follow the coding standard instructions from these files:

* `.github/instructions/nextjs.instructions.md`

#### Tasks to complete

Implementing this story will require implementing the following tasks, but is not limited to it:

* [ ] **Create server-side thunk execution library at `libs/dthunks/server/`**
    * CRITICAL: This implements the pattern described in Story 1.2 "MUST DO PART OF STORY 1.3 IMPLEMENTATION"
    * Create `libs/dthunks/server/constructThunkFromDeclaration.ts`
    * Implement function signature:
      `constructThunkFromDeclaration<TState, TArgs, TFactories>(declaration: ThunkDeclaration<TState, TArgs, TFactories>, initialState: TState, factories: TFactories): (...args: TArgs) => Promise<TState>`
    * The function must:
        * Accept a thunk declaration (from `catalogThunks.onPageRefresh`)
        * Accept an initial state (from `initialCatalogState()`)
        * Accept factory args containing adapter instances
        * Return an executable thunk that accumulates dispatched actions
        * Reduce all dispatched actions against the initial state
        * Return the final computed state
    * Reference existing client-side implementation: `web-nextjs/libs/dthunks/react/useThunks.ts`
    * Must work in Server Component context (no React hooks)
    * Export from `libs/dthunks/server/index.ts`

* [ ] **Create server adapter factory at `domains/catalog/adapters/server-adapter-factory.ts`**
    * Export function `newServerAdapterFactory(): CatalogAdapter`
    * Use `getAccessTokenHolder()` from `@/libs/security/session-service` to retrieve access token
    * Instantiate `FetchCatalogAdapter` with the access token holder
    * This factory is used ONLY in server components for initial state loading
    * Must NOT be used in client components (client uses `catalog-factories.ts`)

* [ ] **Create Server Component for home page at `app/(authenticated)/page.tsx`**
    * Mark as `export default async function Page()`
    * Import `constructThunkFromDeclaration` from `@/libs/dthunks/server`
    * Import `catalogThunks.onPageRefresh` from `@/domains/catalog/thunks`
    * Import `initialCatalogState` from `@/domains/catalog/language/initial-catalog-state`
    * Import `newServerAdapterFactory` from `@/domains/catalog/adapters/server-adapter-factory`
    * Execute thunk:
      `const catalogState = await constructThunkFromDeclaration(catalogThunks.onPageRefresh, initialCatalogState(), {adapter: newServerAdapterFactory()})(undefined)`
    * Pass loaded state to client component: `<CatalogProvider initialState={catalogState}><HomePageContent /></CatalogProvider>`
    * Handle errors by letting them bubble to error boundary (`app/(authenticated)/error.tsx`)

* [ ] **Create CatalogProvider client component at `app/(authenticated)/_components/CatalogProvider/`**
    * Mark as `'use client'`
    * Accept prop: `initialState: CatalogViewerState`
    * Accept prop: `children: ReactNode`
    * Use `useReducer(catalogReducer, initialState)` to hydrate state with server-loaded data
    * Instantiate client-side thunks using pattern from `web/src/components/catalog-react/CatalogViewerProvider.tsx`:
        * `const factoryArgs = useMemo(() => ({adapter: getAdapter(), dispatch}), [dispatch])`
        * `const handlers = useThunks(catalogThunks, factoryArgs, state)`
    * Import `getAdapter` from `@/domains/catalog/catalog-factories` (existing client-side factory)
    * Provide state + handlers to children via Context or props (start with props - simpler)
    * Render `children` function: `{children(state, handlers)}`
    * NO useEffect for initial load - server already loaded the data

* [ ] **Create HomePageContent pure UI component at `app/(authenticated)/_components/HomePageContent/`**
    * Accept props extracted from CatalogViewerState:
        * `albums: Album[]` (from `state.albums`)
        * `isLoading: boolean` (from `state.loading`)
        * `error: CatalogError | null` (from `state.error`)
        * `onRetry: () => void` (handler to retry loading)
    * Render using Material UI components ONLY:
        * `Typography` for headings and text
        * `Link` from `@/components/Link` for clickable album names
        * `Box` for simple layout containers
        * `CircularProgress` for loading indicator
        * `Button` for retry action
    * Display states:
        * **Loading**: `<CircularProgress />` with "Loading albums..." text
        * **Error**: Display error message + `<Button onClick={onRetry}>Try Again</Button>`
        * **Empty**: "No albums found. Create your first album to get started." (message text per FR38)
        * **Success**: List of albums as simple text links
    * Album list format (SIMPLE TEXT ONLY):
      ```tsx
      {albums.map(album => (
        <Box key={album.albumId} sx={{marginBottom: 1}}>
          <Link href={`/albums/${album.ownerId}/${album.albumId}`} prefetch={false}>
            <Typography>{album.name}</Typography>
          </Link>
        </Box>
      ))}
      ```
    * NO styled cards, NO thumbnails, NO density indicators (Story 1.4)
    * Sort albums by date (newest first) - use selector if needed, or ensure onPageRefresh returns sorted albums

* [ ] **Create album page placeholder route at `app/(authenticated)/albums/[ownerId]/[albumId]/page.tsx`**
    * Mark as `export default async function AlbumPage({params})`
    * Accept Next.js params: `{ownerId: string, albumId: string}`
    * Display simple placeholder page:
      ```tsx
      <Box sx={{padding: 3}}>
        <Typography variant="h4">Album: {params.albumId}</Typography>
        <Typography>Owner: {params.ownerId}</Typography>
        <Typography sx={{marginTop: 2}}>
          Album viewing will be implemented in Epic 2.
        </Typography>
        <Link href="/" prefetch={false}>
          <Button>Back to Albums</Button>
        </Link>
      </Box>
      ```
    * This enables navigation testing from home page
    * Actual album implementation is Epic 2

* [ ] **Update error boundary at `app/(authenticated)/error.tsx`** (if needed)
    * Ensure it displays user-friendly error message for catalog loading failures
    * Provide "Try Again" button that calls `reset()`
    * Provide "Return to Home" link
    * Already created in Story 1.2, but verify it meets requirements

* [ ] **Create unit tests for server thunk construction**
    * Test file: `libs/dthunks/server/constructThunkFromDeclaration.test.ts`
    * Test that actions are accumulated and reduced correctly
    * Test that final state reflects all dispatched actions
    * Test error propagation
    * Use vitest framework

* [ ] **Create visual test for HomePageContent using Ladle**
    * Test file: `app/(authenticated)/_components/HomePageContent/HomePageContent.stories.tsx`
    * Stories to create:
        * `Default` - with 3-5 sample albums
        * `Loading` - loading state with spinner
        * `Error` - error state with retry button
        * `Empty` - no albums, empty state message
    * Follow Ladle patterns from `nextjs.instructions.md`
    * NO need to test navigation (Next.js Link tested by framework)

* [ ] **Verify all existing tests still pass**
    * Run `npm run test` - all 230+ tests must pass
    * NO test logic changes allowed - migration preserved behavior
    * Fix only import issues if any arise

#### Target files structure

You will be expected to make changes on the following files:

```
web-nextjs/
├── libs/
│   └── dthunks/
│       └── server/
│           ├── constructThunkFromDeclaration.ts    # NEW: Server-side thunk execution
│           ├── constructThunkFromDeclaration.test.ts # NEW: Unit tests
│           └── index.ts                            # NEW: Exports
│
├── domains/
│   └── catalog/
│       └── adapters/
│           └── server-adapter-factory.ts           # NEW: Server-side adapter factory
│
├── app/
│   └── (authenticated)/
│       ├── page.tsx                                # NEW: Server Component - loads initial state
│       ├── error.tsx                               # UPDATE: Verify meets requirements
│       │
│       ├── _components/
│       │   ├── CatalogProvider/
│       │   │   ├── index.tsx                       # NEW: Client component - hydrates state
│       │   │   └── CatalogProvider.stories.tsx     # NEW: Visual tests
│       │   │
│       │   └── HomePageContent/
│       │       ├── index.tsx                       # NEW: Pure UI component
│       │       └── HomePageContent.stories.tsx     # NEW: Visual tests
│       │
│       └── albums/
│           └── [ownerId]/
│               └── [albumId]/
│                   └── page.tsx                    # NEW: Placeholder album page
```

#### Important Implementation Notes

**Server-Side Thunk Execution Pattern:**

This is the CRITICAL architecture piece. Study the existing client-side implementation in `web-nextjs/libs/dthunks/react/useThunks.ts` to understand the
pattern:

1. Thunk declarations contain factory functions that create executable thunks
2. Factory args provide adapters and dispatch function
3. Thunks dispatch actions through the provided dispatch function
4. On client: dispatch updates React state immediately
5. On server: dispatch accumulates actions, then reduces them all at once to compute final state

**DO NOT mix server and client contexts:**

- Server Component (`page.tsx`) calls `constructThunkFromDeclaration` with `newServerAdapterFactory()`
- Client Component (`CatalogProvider`) calls `useThunks` with `getAdapter()` from client factories
- Adapters are instantiated differently (server: direct session service; client: from factories)

**Keep UI Simple:**

- Resist adding styled album cards - that's Story 1.4
- Typography + Link components only
- Basic Box containers for layout
- Material UI default spacing and typography
- Focus on FUNCTIONALITY, not visual polish

**State Hydration:**

- Server passes computed state to client as prop
- Client initializes useReducer with that state (second argument)
- NO useEffect for initial load - state already loaded
- Handlers instantiated client-side to enable future interactions

**Error Handling:**

- Let server errors bubble to error boundary (`app/(authenticated)/error.tsx`)
- Client errors handled by CatalogProvider state (error in CatalogViewerState)
- Retry handler re-executes onPageRefresh thunk (already available from thunks)

**Testing Strategy:**

- Unit test server thunk construction (actions reduced correctly)
- Visual test UI component states (loading, error, empty, success)
- NO integration tests needed - existing 230 tests validate state management
- Run `npm run test` frequently to catch regressions

**What NOT to Do:**

- ❌ DO NOT implement styled album cards (Story 1.4)
- ❌ DO NOT add photo thumbnails (Story 1.4)
- ❌ DO NOT implement density indicators (Story 1.4)
- ❌ DO NOT implement filtering (Story 1.6)
- ❌ DO NOT implement random photo highlights (Epic 5)
- ❌ DO NOT implement actual album page (Epic 2)
- ❌ DO NOT bypass error boundaries with try-catch in page.tsx
- ❌ DO NOT modify existing state management logic (230 tests validate it)
- ❌ DO NOT use client-side data fetching (useEffect) - server loads initial state

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem ?
* What has been done to solve it ?k
* Results and screenshots when possible
