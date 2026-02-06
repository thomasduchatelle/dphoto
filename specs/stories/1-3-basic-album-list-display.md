# Story 1.3: Basic Album List Display

**Status**: ready-for-dev

---

## Story

As a **user**,  
I want **to view my album list on the authenticated home page**,  
so that **I can see all albums I own and albums shared with me in chronological order**.

## Acceptance Criteria

**Given** I am an authenticated user with access to albums  
**When** I navigate to the home page  
**Then** the page is located at `app/(authenticated)/page.tsx` (Server Component)  
**And** authentication is verified via existing Google OAuth backend (FR33, FR37)  
**And** user session is maintained across navigation (FR34)  
**And** user profile information (name, picture) is displayed in the app header (FR35)  
**And** the system identifies if I'm an owner (can create albums) (FR36)  
**And** the Server Component fetches albums from the existing REST API using fetch adapter  
**And** initial album data is passed as props to Client Component wrapper  
**And** Client Component (`_components/AlbumListClient.tsx`) initializes catalog state using useReducer with migrated reducer  
**And** album cards are displayed in a responsive grid:

- Mobile (xs): 1 column
- Tablet (sm): 2 columns
- Desktop (md): 3 columns
- Large Desktop (lg): 4 columns

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

* **Story 1.1 (Project Foundation Setup)**: Material UI theme configuration, dark mode with brand color #185986, MUI breakpoints
* **Story 1.2 (State Management Migration)**:
    - Migrated catalog state management from `web/src/core/catalog/`
    - Fetch adapter replacing axios for API calls
    - Custom image loader configuration
    - Error boundaries (`app/error.tsx` and `app/(authenticated)/error.tsx`)
    - Not-found pages

### Expected outcomes

After this story is complete:

1. **Working home page** at `app/(authenticated)/page.tsx` that displays album list
2. **Pure UI components** in `_components/` folder:
    - `AlbumListClient.tsx` - Client wrapper with state initialization
    - `AlbumsGrid.tsx` - Responsive grid layout for album cards
    - `AlbumCard.tsx` - Basic album card showing metadata (no photo samples yet)
    - `LoadingSkeleton.tsx` - Loading state for albums
    - `EmptyState.tsx` - Empty state when no albums exist
3. **Responsive grid layout** working across all breakpoints (xs/sm/md/lg)
4. **Navigation** from album card to album detail page route
5. **State management** pattern established: Server Component fetches → Client Component initializes state → Pure UI renders
6. **Error handling** for failed API calls with retry mechanism

These components will be enhanced in Story 1.4 with random photo samples and density indicators.

### Out of scope

* DO NOT implement random photo samples on album cards (Story 1.4)
* DO NOT implement density color-coding/temperature indicators (Story 1.4)
* DO NOT implement user avatars for sharing status (Story 1.4)
* DO NOT implement album filtering by owner (Story 1.5)
* DO NOT implement Create Album functionality (Story 3.1)
* DO NOT implement album edit/delete actions (Stories 3.2-3.4)
* DO NOT implement sharing functionality (Stories 4.1-4.2)
* DO NOT implement random photo highlights section on home page (Story 5.1)

Focus: Basic album list display with metadata only. Visual enhancements come in Story 1.4.

---

## Implementation Guidance

This technical guidance has been validated by the lead developer, following it significantly increases the chance of getting your PR accepted. Any infringement
required to complete the story must be reported.

### Coding standards

You must follow the coding standard instructions from these files:

* @.github/instructions/nextjs.instructions.md

### Tasks to complete

Implementing this story will require implementing the following tasks, but is not limited to it:

* [ ] Create Server Component page at `app/(authenticated)/page.tsx`
    * Fetches albums using `FetchCatalogAdapter` from `domains/catalog/adapters/fetch-adapter.ts`
    * Computes initial `CatalogViewerState` with albums loaded
    * Passes initial state to Client Component as props
    * Uses `cache: 'no-store'` for dynamic data fetching
    * Handles authentication verification through existing middleware

* [ ] Create Client Component wrapper `app/(authenticated)/_components/AlbumListClient/index.tsx`
    * Receives `initialState` as prop from Server Component
    * Initializes catalog state with `useReducer(catalogReducer, initialState)` from `domains/catalog/actions.ts`
    * Instantiates thunks client-side for user interactions
    * Passes state and handlers as props to pure UI components
    * Must NOT contain UI rendering logic (only state management wrapper)

* [ ] Create pure UI component `app/(authenticated)/_components/AlbumsGrid/index.tsx`
    * Receives `albums: Album[]` from state as prop
    * Receives handlers for user interactions as props
    * Displays responsive grid using Material UI Box with sx prop
    * Grid columns: xs: 1, sm: 2, md: 3, lg: 4
    * Must be pure presentational component with NO internal state management
    * Renders AlbumCard components in grid

* [ ] Create pure UI component `app/(authenticated)/_components/AlbumCard/index.tsx`
    * Receives `album: Album` and `onClick: (albumId: AlbumId) => void` as props
    * Displays album metadata: name, date range (formatted), media count, owner information (if shared)
    * Uses Material UI Card component with brand color (#185986) for accents
    * Handles click to navigate to `/owners/[ownerId]/albums/[albumId]` route
    * Shows owner information when `album.ownedBy` is defined (shared albums)
    * Must be pure presentational component with NO internal state
    * Create Ladle stories testing: Default, OwnedAlbum, SharedAlbum states

* [ ] Create pure UI component `app/(authenticated)/_components/LoadingSkeleton/index.tsx`
    * Uses Material UI Skeleton component
    * Displays grid of skeleton cards matching AlbumsGrid layout
    * Receives `count?: number` prop (default: 8)
    * Must be pure presentational component

* [ ] Create pure UI component `app/(authenticated)/_components/EmptyState/index.tsx`
    * Uses Material UI Typography and Box components
    * Displays centered message: "No albums yet"
    * Optionally displays action button for owners to create album (if supported)
    * Receives `isOwner: boolean` as prop
    * Must be pure presentational component
    * Create Ladle stories testing: ForOwner, ForViewer states

* [ ] Update or create error boundary `app/(authenticated)/error.tsx`
    * Uses NextJS error boundary pattern with "use client" directive
    * Receives `error: Error` and `reset: () => void` props
    * Displays error message with "Try Again" button calling reset()
    * Uses Material UI components for styling
    * Provides clear user feedback on failure

* [ ] Add selector `domains/catalog/navigation/selector-albumList.ts`
    * Exports `selectAlbumList` function accepting `CatalogViewerState`
    * Returns filtered albums based on `state.albums` (already filtered by criterion)
    * Returns chronological order (newest first) - already sorted by adapter
    * Must be tested with actions: given initial state with albums -> dispatches filter action -> selector returns filtered albums

* [ ] Add selector `domains/catalog/navigation/selector-loadingState.ts`
    * Exports `selectIsLoadingAlbums` function accepting `CatalogViewerState`
    * Returns `!state.albumsLoaded` for loading skeleton display
    * Must be tested with actions: given initial state loading -> dispatch albumsLoaded -> selector returns false

### Target files structure

You will be expected to make changes on the following files:

```
web-nextjs/
├── app/
│   └── (authenticated)/
│       ├── page.tsx                              # NEW: Server Component - fetches albums, passes to client
│       ├── error.tsx                             # UPDATE: Error boundary for authenticated routes
│       └── _components/
│           ├── AlbumListClient/
│           │   └── index.tsx                     # NEW: Client wrapper with useReducer
│           ├── AlbumsGrid/
│           │   ├── index.tsx                     # NEW: Pure UI - responsive grid
│           │   └── AlbumsGrid.stories.tsx        # NEW: Ladle stories
│           ├── AlbumCard/
│           │   ├── index.tsx                     # NEW: Pure UI - album metadata display
│           │   └── AlbumCard.stories.tsx         # NEW: Ladle stories (Default, Owned, Shared)
│           ├── LoadingSkeleton/
│           │   ├── index.tsx                     # NEW: Pure UI - loading state
│           │   └── LoadingSkeleton.stories.tsx   # NEW: Ladle stories
│           └── EmptyState/
│               ├── index.tsx                     # NEW: Pure UI - empty state
│               └── EmptyState.stories.tsx        # NEW: Ladle stories (Owner, Viewer)
│
└── domains/
    └── catalog/
        └── navigation/
            ├── selector-albumList.ts             # NEW: selector with test
            └── selector-loadingState.ts          # NEW: selector with test
```

### Critical implementation notes

**State Management Pattern:**

- Server Component ONLY fetches data and computes initial state
- Client Component ONLY manages state with useReducer and instantiates thunks
- Pure UI components ONLY receive props and render (NO state management)
- Follow the pattern established in Story 1.2 state migration

**Testing Requirements:**

- All selectors MUST be tested with their corresponding actions (TDD principle)
- All UI components MUST have Ladle stories covering relevant states
- Use fake implementations for adapters, NOT mocks
- Tests must be robust to refactoring (test behavior, not implementation)

**Material UI Usage:**

- Use MUI Box with sx prop for responsive grids
- Responsive object syntax: `{ xs: 1, sm: 2, md: 3, lg: 4 }`
- Use brand color #185986 from theme for primary accents
- Use MUI Card, Typography, Skeleton components as appropriate

**Data Flow:**

- Albums already filtered by `albumFilter` criterion in state
- Albums already sorted chronologically (newest first) by FetchCatalogAdapter
- Owner information (`album.ownedBy`) is undefined for owned albums, defined for shared albums
- Navigation to album detail uses route: `/owners/[ownerId]/albums/[albumId]`

**DO NOT:**

- Add business logic or state management inside UI components
- Create new state management patterns (use migrated catalog state)
- Add comments in code (use chat for communication)
- Use `any` types (always explicit types)
- Implement features not in acceptance criteria (NO photo samples, NO density indicators yet)

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem
* What has been done to solve it
* Results and screenshots when possible
