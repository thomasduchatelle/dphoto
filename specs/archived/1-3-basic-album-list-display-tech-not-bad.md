# Story 1.3: Basic Album List Display

**Status**: ready-for-dev

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem
* What has been done to solve it
* Results and screenshots when possible

As a **user**,  
I want **to view my album list on the authenticated home page**,  
so that **I can see all albums I own and albums shared with me in chronological order**.

## Acceptance Criteria

**Given** I am an authenticated user with access to albums
**When** I navigate to the home page
**Then** the page is located at `app/(authenticated)/page.tsx` (Server Component)
**And** user profile information (name, picture) is displayed in the app header (FR35)
**And** the Server Component fetches albums from the existing REST API using fetch adapter

**And** each album card shows metadata (name, date range, media count, owner information) (FR4)
**And** albums are displayed in chronological order (newest first) (FR6)
**And** I can see albums I own (FR1) and albums shared with me (FR2)
**And** owner information is displayed on shared albums (FR29)
**And** a loading skeleton is shown while fetching albums (FR38)
**And** an empty state message is displayed when no albums exist (FR39)
**And** error messages are displayed if the API call fails with a "Try Again" button (FR40)
**And** the layout adapts responsively across mobile, tablet, and desktop (FR44)
**And** clicking an album card navigates to `/owners/[ownerId]/albums/[albumId]`

## Scope & Context

### Depends on

* **Story 1.1 - Project Foundation Setup** (`specs/stories/1-1-project-foundation-setup.md`): Material UI theme configured, dark mode enabled, brand color
  #185986 set, ThemeProvider wrapping application
* **Story 1.2 - State Management Migration** (`specs/stories/1-2-state-management-migration.md`): Catalog state management migrated to `domains/catalog/`,
  FetchCatalogAdapter created, all 230 tests passing, image loader configured

### Expected outcomes

The DEV agent implementing this story will deliver a functional home page that:

- Displays the album list using the migrated catalog state management from Story 1.2
- Uses Material UI components styled with the theme from Story 1.1
- Implements Server Component data fetching pattern with Client Component state hydration
- Shows albums in a responsive grid layout working on mobile, tablet, and desktop
- Handles loading, empty, and error states gracefully
- Provides navigation to individual album pages

This story provides the foundation for:

- **Story 1.4** will enhance album cards with random photo samples and density indicators
- **Story 1.5** will add filtering capabilities to this album list

### Out of scope

* DO NOT implement random photo samples on cards - that's Story 1.4
* DO NOT implement album density indicators - that's Story 1.4
* DO NOT implement filtering by owner - that's Story 1.5
* DO NOT implement album creation, editing, or deletion dialogs - those are Epic 3
* DO NOT implement the photo viewing pages - that's Epic 2
* DO NOT implement sharing management - that's Epic 4

---

<!-- Part to be completed by the Senior Dev -->

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

### Responsive Grid Layout

MUI Grid with responsive column configuration (FR44):

```tsx
<Grid container spacing={2}>
    {albums.map(album => (
        <Grid item xs={12} sm={6} md={4} lg={3} key={album.id}>
            <AlbumCard album={album} onClick={onClick}/>
        </Grid>
    ))}
</Grid>
```

Breakpoints (from Architecture.md):

- **xs (< 600px)**: 1 column (12/12)
- **sm (600-959px)**: 2 columns (6/12)
- **md (960-1279px)**: 3 columns (4/12)
- **lg (≥ 1280px)**: 4 columns (3/12)

### Navigation Handling

Clicking an album card navigates to album page:

- Route: `/owners/[ownerId]/albums/[albumId]`
- Implementation: MUI Card wrapped with NextJS Link OR Button with onClick calling `router.push()`
- Must use `@/components/Link` wrapper (from Story 1.1) for server component compatibility

### Album Card Display Requirements

Each album card displays (FR4):

- **Album name**: Typography variant="h6"
- **Date range**: "Jan 15, 2026 - Jan 20, 2026"
- **Media count**: "42 photos"
- **Owner info** (for shared albums): Avatar + name (FR29)

**NOT in this story** (out of scope):

- ❌ Random photo samples (Story 1.4)
- ❌ Density indicators (Story 1.4)
- ❌ Edit/delete buttons (Epic 3)

### Critical Success Factors

1. **Server-side thunk execution works**: `constructThunkFromDeclaration` must correctly execute thunks in server context, dispatch actions, reduce state, and
   return computed state
2. **State hydration seamless**: Client Component receives server state via props and hydrates useReducer without hydration errors
3. **Pure UI components**: AlbumCard and AlbumsGrid have NO internal state - only props
4. **Responsive layout works**: Grid adapts correctly across all breakpoints
5. **All tests pass**: Unit tests for components using Ladle stories

---

## Implementation Guidance

This technical guidance has been validated by the lead developer, following it significantly increases the chance of getting your PR accepted. Any infringement
required to complete the story must be reported.

### Coding standards

You must follow the coding standard instructions from these files:

* `@.github/instructions/nextjs.instructions.md`

### Tasks to complete

Implementing this story will require implementing the following tasks, but is not limited to it:

* [ ] **Create server-side thunk execution utility** (`libs/dthunks/server/constructThunkFromDeclaration.ts`)
    * Accept `ThunkDeclaration<State, Args>`, `initialState: State`, `factoryArgs: FactoryArgs`
    * Create dispatch function that collects actions in array
    * Instantiate thunk using declaration's factory with provided factoryArgs
    * Execute thunk with dispatch and args
    * Reduce collected actions against initial state using generic reducer
    * Return final computed state
    * Must handle async thunks (return Promise<State>)
    * Reference client implementation: `libs/dthunks/react/useThunks.ts`

* [ ] **Create server adapter factory** (`domains/catalog/adapters/server-adapter-factory.ts`)
    * Export function `newServerAdapterFactory(): MasterCatalogAdapter`
    * Instantiate `FetchCatalogAdapter` with `getAccessTokenHolder()` from `@/libs/security/session-service`
    * Return adapter instance for use in server-side thunk execution
    * This enables thunks to call REST API from Server Component context

* [ ] **Create client adapter factory** (`domains/catalog/adapters/client-adapter-factory.ts`)
    * Export function `newClientAdapterFactory(): MasterCatalogAdapter`
    * Instantiate `FetchCatalogAdapter` with `getAccessTokenHolder()` from `@/libs/security/session-service`
    * Return adapter instance for use in client-side thunk execution
    * NOTE: May be identical to server factory, but separated for future divergence

* [ ] **Implement Server Component home page** (`app/(authenticated)/page.tsx`)
    * NO 'use client' directive - this is a Server Component
    * Get authenticated user from session
    * Create initial state with `initialCatalogState(authenticatedUser)`
    * Execute `constructThunkFromDeclaration(catalogThunks.onPageRefresh, initialState, {adapterFactory: newServerAdapterFactory()})`
    * Pass computed state to `<AlbumListClient initialState={loadedState} />`
    * Handle errors with try-catch, throw to trigger error boundary if needed

* [ ] **Implement Client Component wrapper** (`app/(authenticated)/_components/AlbumListClient.tsx`)
    * MUST have 'use client' directive at top
    * Props: `{ initialState: CatalogViewerState }`
    * Use `const [state, dispatch] = useReducer(catalogReducer, props.initialState)`
    * Instantiate handlers with `useThunks(catalogThunks, {adapterFactory: newClientAdapterFactory(), dispatch}, state)`
    * Render loading state: `if (!state.albums) return <LoadingSkeleton />`
    * Render error state: `if (state.error) return <ErrorMessage error={state.error} onRetry={handlers.onPageRefresh} />`
    * Render empty state: `if (state.albums.length === 0) return <EmptyState />`
    * Render main content: `<AlbumsGrid albums={state.albums} onAlbumClick={(album) => router.push(`/owners/${album.ownerId}/albums/${album.id}`)} />`

* [ ] **Create AlbumsGrid component** (`app/(authenticated)/_components/AlbumsGrid/index.tsx`)
    * Pure UI component (NO 'use client' needed if no hooks)
    * Props: `{ albums: Album[], onAlbumClick: (album: Album) => void }`
    * Use MUI Grid with responsive columns: xs={12}, sm={6}, md={4}, lg={3}
    * Map albums to `<AlbumCard album={album} onClick={() => onAlbumClick(album)} />`
    * Add spacing prop for gap between cards

* [ ] **Create AlbumCard component** (`app/(authenticated)/_components/AlbumCard/index.tsx`)
    * Pure UI component
    * Props: `{ album: Album, onClick: () => void }`
    * Use MUI Card, CardContent, Typography components
    * Display: album name (h6), date range (body2), media count (body2)
    * If album is shared (album.ownerId !== currentUserId), display owner info: Avatar + name
    * Apply brand color (#185986) for interactive elements (hover state)
    * Use Card's onClick prop for navigation
    * Style with sx prop for responsive behavior

* [ ] **Create LoadingSkeleton component** (`app/(authenticated)/_components/LoadingSkeleton/index.tsx`)
    * Use MUI Skeleton component
    * Display grid of skeleton cards matching AlbumCard layout
    * Show 4-6 skeleton items for visual consistency
    * Use same Grid layout as AlbumsGrid

* [ ] **Create ErrorMessage component** (`app/(authenticated)/_components/ErrorMessage/index.tsx`)
    * Props: `{ error: CatalogError, onRetry: () => void }`
    * Display error message with MUI Alert severity="error"
    * Include "Try Again" button (MUI Button variant="contained" color="primary")
    * Center content in page
    * Provide helpful message guiding user

* [ ] **Create EmptyState component** (`app/(authenticated)/_components/EmptyState/index.tsx`)
    * Props: `{ message?: string }`
    * Display message with MUI Typography variant="h6" and icon
    * Center content in page
    * Provide helpful guidance (e.g., "Upload photos via CLI to create albums")

* [ ] **Create Ladle stories for each component**
    * `AlbumCard.stories.tsx`: Default, Owned, Shared, Loading states
    * `AlbumsGrid.stories.tsx`: With albums, Empty, Few items (1-3), Many items
    * `LoadingSkeleton.stories.tsx`: Default
    * `ErrorMessage.stories.tsx`: Network error, Permission error
    * `EmptyState.stories.tsx`: Default, Custom message
    * Follow patterns from `.github/instructions/nextjs.instructions.md`

* [ ] **Write unit tests for selectors if needed**
    * Test `selectVisibleAlbums` if it performs filtering/sorting logic
    * Use vitest following existing test patterns in `domains/catalog/`

* [ ] **Update exports in catalog domain** (`domains/catalog/index.ts`)
    * Export `onPageRefresh` thunk declaration for server-side use
    * Export relevant selectors used by components
    * Ensure clean public API

### Target files structure

You will be expected to make changes on the following files:

```
web-nextjs/
├── app/(authenticated)/
│   ├── page.tsx                                    # NEW: Server Component - fetches data, executes thunk
│   └── _components/
│       ├── AlbumListClient.tsx                     # NEW: Client Component - hydrates state
│       ├── AlbumsGrid/
│       │   ├── index.tsx                           # NEW: Pure UI - responsive grid
│       │   └── AlbumsGrid.stories.tsx              # NEW: Ladle visual tests
│       ├── AlbumCard/
│       │   ├── index.tsx                           # NEW: Pure UI - album display
│       │   └── AlbumCard.stories.tsx               # NEW: Ladle visual tests
│       ├── LoadingSkeleton/
│       │   ├── index.tsx                           # NEW: Loading state UI
│       │   └── LoadingSkeleton.stories.tsx         # NEW: Ladle visual tests
│       ├── ErrorMessage/
│       │   ├── index.tsx                           # NEW: Error state UI
│       │   └── ErrorMessage.stories.tsx            # NEW: Ladle visual tests
│       └── EmptyState/
│           ├── index.tsx                           # NEW: Empty state UI
│           └── EmptyState.stories.tsx              # NEW: Ladle visual tests
│
├── libs/dthunks/server/
│   ├── constructThunkFromDeclaration.ts            # NEW: Server-side thunk execution
│   ├── constructThunkFromDeclaration.test.ts       # NEW: Unit tests
│   └── index.ts                                    # NEW: Exports
│
├── domains/catalog/
│   ├── adapters/
│   │   ├── server-adapter-factory.ts               # NEW: Server adapter factory
│   │   └── client-adapter-factory.ts               # NEW: Client adapter factory
│   └── index.ts                                    # UPDATED: Export thunks, selectors
```

### Important Implementation Notes

**Server-Side Thunk Execution Pattern:**

This is a NEW pattern being introduced in this story. Pay careful attention to:

- `constructThunkFromDeclaration` must correctly instantiate and execute thunks
- Actions must be collected and reduced sequentially
- Final state must include all data fetched by the thunk
- Reference the client-side `useThunks` implementation for patterns

**State Hydration:**

- Server state passed as prop → Client Component useReducer initial state
- NO hydration mismatches allowed (server/client render must match)
- State must be serializable (no functions, class instances)

**Pure Component Architecture:**

- AlbumCard and AlbumsGrid receive ONLY props (albums, onClick)
- NO useState, NO useEffect, NO business logic
- All data transformations happen in selectors
- All handlers come from parent via props

**MUI Integration:**

- Use Material UI components exclusively (Card, Grid, Typography, Button, Skeleton, Alert)
- Apply theme via sx prop for responsive behavior
- Use brand color (#185986) for interactive elements
- Follow MUI Grid responsive patterns

**Testing Strategy:**

- Pure UI components: Ladle stories showing all states
- Thunk execution utility: Vitest unit tests
- State management: Already tested in Story 1.2 (90+ tests passing)
- NO integration tests required for this story

**What NOT to Do:**

- ❌ DO NOT implement random photo samples (Story 1.4)
- ❌ DO NOT implement density indicators (Story 1.4)
- ❌ DO NOT implement filtering (Story 1.5)
- ❌ DO NOT implement album management dialogs (Epic 3)
- ❌ DO NOT put state management logic in UI components
- ❌ DO NOT create custom responsive breakpoints (use MUI defaults)
- ❌ DO NOT bypass Server Component pattern (must fetch on server first)

---

