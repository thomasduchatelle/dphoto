# Story 1.6: Album Filtering

**Status**: ready-for-dev

---

## Story

As a **user**,  
I want **to filter albums by owner**,  
so that **I can focus on my own albums or view all albums including shared ones**.

## Acceptance Criteria

**Given** I am viewing the album list with multiple albums from different owners
**When** I use the owner filter
**Then** a filter control is displayed above the album grid
**And** the filter offers these options:

- "All Albums" - shows all accessible albums (owned + shared)
- "My Albums" - shows only albums I own
- One option per owner - shows albums owned by that specific owner

**And** the active filter is visually indicated
**And** when I select "All Albums", I see albums I own and albums shared with me
**And** when I select "My Albums", I see only albums where I am the owner
**And** when I select a specific owner, I see only albums owned by that person
**And** the filtered album list updates immediately when I change the filter
**And** the filter selection persists as I navigate within the session
**And** the filter is cleared when I log out
**And** my scroll position in the album list is maintained when I change filters
**And** the filter control works on mobile, tablet, and desktop
**And** I can navigate the filter using keyboard (tab, arrow keys, enter)

## Scope & Context

### Depends on

* **Story 1.1 (Project Foundation Setup)**: Material UI theme with brand color #185986, dark mode, ThemeProvider configured
* **Story 1.2 (State Management Migration)**: Complete catalog state management with actions, reducer, thunks for filtering operations
* **Story 1.3 (Basic Album List Display)**: Server-side state initialization pattern, CatalogProvider client component, basic album loading
* **Story 1.4 (UI Components & Layout)**: All UI components created (AlbumCard, AlbumGrid, AppLayout, feedback components)
* **Story 1.5 (Style the Home Page)**: Styled home page with albums displayed in grid, AppLayout integration, loading/error/empty states

### Expected outcomes

This story completes **Epic 1: Album List Home Page**, delivering the final interactive feature for album browsing. The DEV agent implementing Story 2.1 (Epic 2
start) will need:

* **Working filter state management pattern**: How filter selection modifies visible albums through reducer actions and selectors
* **Filter component structure**: Reusable pattern for dropdown/chip filters using MUI components
* **Session persistence pattern**: How filter state is maintained during navigation and cleared on logout
* **Accessibility pattern**: Keyboard navigation and ARIA labels for filter controls

This establishes the pattern for future filtering features (e.g., date range filters, tag filters).

### Out of scope

* DO NOT implement album creation functionality - this is Epic 3 (Album Management)
* DO NOT implement actual album page viewing - this is Epic 2 (Photo Viewing & Navigation)
* DO NOT implement album editing/deletion - this is Epic 3
* DO NOT implement sharing functionality - this is Epic 4
* DO NOT add random photo highlights - this is Epic 5
* DO NOT implement advanced filtering (date ranges, tags, search) - these are potential future enhancements

---

## Technical Context

### Architecture Reference

**File**: `/home/dush/dev/git/dphoto/specs/designs/architecture.md`

**Relevant architectural decisions:**

- **State Management (Decision #2)**: Filter state lives in `CatalogViewerState`, actions dispatch filter changes, selectors compute visible albums
- **Component Architecture (Decision #3)**: Filter component colocated in `app/(authenticated)/_components/` (page-specific)
- **MUI Integration (Decision #1)**: Use MUI Select, Chip, or ToggleButtonGroup for filter UI

**State management pattern from Story 1.3:**

```tsx
// Server Component initializes state
const catalogState = await constructThunkFromDeclaration(...)

// Client Component hydrates and provides handlers
const [state, dispatch] = useReducer(catalogReducer, initialState)
const handlers = useThunks(catalogThunks, {adapter, dispatch}, state)

    // Pure UI receives filtered albums
    < HomePageContent
albums = {selectVisibleAlbums(state)}
...
/>
```

### Product Requirements

**File**: `/home/dush/dev/git/dphoto/specs/designs/prd.md`

**Functional Requirements addressed:**

- **FR3**: Users can filter albums by owner (my albums, all albums, specific owner)
- **FR10**: System indicates selected albums and active filters

**Non-Functional Requirements:**

- **NFR5**: Focus management must be clear and logical during navigation and in dialogs
- **NFR6**: Keyboard shortcuts must not conflict with browser defaults
- **NFR22**: Focus management must be clear and logical during navigation and in dialogs

### UX Design Specification

**File**: `/home/dush/dev/git/dphoto/specs/designs/ux-design-specification.md`

**Key UX requirements:**

- Filter control displayed above album grid
- Active filter visually indicated (brand blue #185986)
- Filter options: "All Albums", "My Albums", + one per owner
- Immediate update when filter changes
- Keyboard navigation: Tab to filter, arrow keys to select, Enter to confirm
- Responsive: works on mobile, tablet, desktop
- Touch-friendly on mobile (minimum 44px touch targets)

### Epic Context

**File**: `/home/dush/dev/git/dphoto/specs/designs/epics.md`

**Epic 1: Album List Home Page** - FINAL STORY

This story completes Epic 1, delivering the full album browsing experience with filtering capabilities.

**Story progression:**

1. Story 1.1 ✅ - Material UI foundation
2. Story 1.2 ✅ - State management migration
3. Story 1.3 ✅ - Basic album loading
4. Story 1.4 ✅ - UI components created
5. Story 1.5 ✅ - Styled home page
6. **Story 1.6 (THIS STORY)** - Album filtering (COMPLETES EPIC 1)

After this story, Epic 1 is complete and Epic 2 (Photo Viewing & Navigation) can begin.

### Existing State Management (from Story 1.2)

**Location**: `web-nextjs/domains/catalog/`

The migrated state management from `web/src/core/catalog/` includes:

**CatalogViewerState** (`language/CatalogViewerState.ts`):

```typescript
interface CatalogViewerState {
    albums: Album[]
    selectedAlbumKey?: AlbumKey
    filter: {
        ownerFilter: 'all' | 'mine' | string // 'all', 'mine', or specific ownerId
    }
    loading: boolean
    error: CatalogError | null
    // ... other fields
}
```

**Actions** (`actions.ts`):

- `filterByOwner(ownerFilter: string)` - dispatches filter change
- `catalogRefreshStarted` - loading state
- `catalogRefreshed(albums: Album[])` - albums loaded
- `catalogRefreshFailed(error: CatalogError)` - error state

**Selectors** (`language/selectors.ts` or inline):

- `selectVisibleAlbums(state: CatalogViewerState): Album[]` - computes filtered albums based on `state.filter.ownerFilter`
- `selectUniqueOwners(state: CatalogViewerState): Owner[]` - extracts unique owners from all accessible albums for filter dropdown

**Thunks** (`thunks.ts`):

- May include `applyOwnerFilter` thunk that dispatches `filterByOwner` action
- Or filtering may be purely client-side (dispatch action, selector computes)

### Project Structure (Current State)

After Stories 1.1-1.5:

```
web-nextjs/
├── app/
│   └── (authenticated)/
│       ├── page.tsx                            # Server Component - loads initial state
│       └── _components/
│           ├── CatalogClient/
│           │   └── index.tsx                   # Client - hydrates state, provides handlers
│           └── HomePageContent/
│               ├── index.tsx                   # Pure UI - renders albums
│               └── HomePageContent.stories.tsx # Ladle stories
│
├── components/
│   ├── layout/
│   │   ├── AppLayout/                          # From Story 1.4
│   │   └── AppHeader/                          # From Story 1.4
│   ├── albums/
│   │   ├── AlbumCard/                          # From Story 1.4
│   │   └── AlbumGrid/                          # From Story 1.4
│   ├── feedback/                               # From Story 1.4
│   └── user/                                   # From Story 1.4
│
├── domains/
│   └── catalog/                                # From Story 1.2
│       ├── language/
│       │   ├── CatalogViewerState.ts
│       │   └── selectors.ts (if exists)
│       ├── actions.ts                          # Reducer + actions
│       └── thunks.ts                           # Coordinated operations
│
└── libs/
    └── dthunks/
        └── server/                             # From Story 1.3
```

---

## Technical Design

### Overview

This story implements **album filtering by owner** using the existing filter state management from Story 1.2's migration. The `CatalogViewerState` already
contains `albumFilterOptions`, `albumFilter`, and the derived `albums` array. This story creates the UI component to display filter options and dispatches
actions to change the active filter.

**Key Insight**: The filtering logic and state management already exist (migrated in Story 1.2). This story is **purely about adding the UI control** that
allows users to interact with the existing filter system.

### Architecture Pattern

Following **Architecture.md Decision #2** (State Management):

```
Server Component (page.tsx)
   ↓ (loads albums + filter options)
   ↓ (passes initialState with albumFilterOptions populated)
Client Component (CatalogClient)
   ├─ useReducer(catalogReducer, initialState)
   ├─ useThunks(catalogThunks, {adapter, dispatch}, state)
   └─ Pure UI → HomePageContent
        ├─ AlbumFilterControl (NEW)
        │    ↓ (displays filter options)
        │    ↓ (calls onFilterChange handler)
        └─ AlbumGrid with filtered albums
             ↓ (state.albums already filtered by reducer)
```

**Critical Understanding**: The `state.albums` array is **already filtered** by the reducer based on `state.albumFilter`. The UI component simply displays the
current filter and allows users to select a different one.

### Existing Filter State Structure

From `/web-nextjs/domains/catalog/language/catalog-state.ts` (migrated in Story 1.2):

```typescript
interface CatalogViewerState {
    allAlbums: Album[]              // All accessible albums (unfiltered)
    albumFilterOptions: AlbumFilterEntry[]  // Available filter options
    albumFilter: AlbumFilterEntry    // Currently active filter
    albums: Album[]                  // Filtered albums (derived from allAlbums + albumFilter)
    // ... other fields
}

interface AlbumFilterEntry {
    criterion: AlbumFilterCriterion  // The actual filter logic
    avatars: string[]                // User avatars to display
    name: string                     // Display name ("All Albums", "My Albums", "John Doe")
}

interface AlbumFilterCriterion {
    owners: Owner[]      // Empty array = all albums
    selfOwned?: boolean  // True = only albums owned by current user
}
```

**Filter Options Expected** (from Story Context line 147-162):

1. "All Albums" - `{owners: [], selfOwned: false}` → shows all accessible albums
2. "My Albums" - `{selfOwned: true}` → shows only albums where current user is owner
3. One per unique owner - `{owners: ['owner-id']}` → shows albums owned by specific owner

### Existing Actions and Reducer

From Story 1.2 migration, the reducer already handles filter changes:

**Action Pattern** (likely exists as `action-albumFilterChanged.ts` or similar):

```typescript
const albumFilterChanged = (newFilter: AlbumFilterEntry): Action<CatalogViewerState> => ({
    type: 'albumFilterChanged',
    payload: {filter: newFilter},
    reduce: (state) => ({
        ...state,
        albumFilter: newFilter,
        albums: state.allAlbums.filter(albumMatchCriterion(newFilter.criterion))
    })
})
```

**Selector for Filtered Albums** (may exist as `selector-visibleAlbums.ts`):

```typescript
const selectVisibleAlbums = (state: CatalogViewerState): Album[] => {
    return state.albums // Already filtered by reducer
}
```

### Component Design: AlbumFilterControl

**Purpose**: Display filter options and allow user to select active filter

**Location**: `app/(authenticated)/_components/AlbumFilterControl/index.tsx`

**Why page-specific**: Only used on home page for now. Move to `components/shared/` if Epic 2 reuses it.

**Props Interface**:

```typescript
interface AlbumFilterControlProps {
    filterOptions: AlbumFilterEntry[]  // Available filters
    activeFilter: AlbumFilterEntry     // Currently selected filter
    onFilterChange: (filter: AlbumFilterEntry) => void  // Handler to change filter
}
```

**MUI Component Choice**: **ToggleButtonGroup** or **Select**

**Recommendation**: ToggleButtonGroup for better visual indication of active filter

- Buttons show clearly which filter is active
- Better for small number of options (typically 3-5)
- Touch-friendly on mobile (large tap targets)
- Horizontal on desktop, can stack on mobile

**Alternative**: Select (dropdown) if many filter options (10+)

**Visual Design Requirements** (from UX Design Specification):

- Active filter visually indicated with brand blue (#185986) background or border
- Keyboard navigation: Tab to focus, arrow keys to select, Enter to confirm
- Responsive: works on mobile, tablet, desktop
- Minimum 44px touch targets on mobile
- Display above album grid (not floating, not in header)

### UI Component Structure

**AlbumFilterControl Layout**:

```
┌─────────────────────────────────────────────────┐
│  Filter by Owner:                               │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────┐│
│  │ All Albums ✓ │ │  My Albums   │ │ John Doe ││
│  └──────────────┘ └──────────────┘ └──────────┘│
└─────────────────────────────────────────────────┘
        (ToggleButtonGroup with 3 options)
```

**Responsive Behavior**:

- Desktop (md+): Horizontal button group, all options visible
- Tablet (sm): Horizontal, may wrap if many options
- Mobile (xs): Vertical stack or horizontal scroll

**Accessibility Requirements** (from PRD NFR5, NFR6, NFR22):

- ARIA role="group" with aria-label="Album filter by owner"
- Each button has aria-pressed indicating active state
- Keyboard navigation: Tab, arrow keys, Enter
- Focus visible with brand blue outline
- No conflicts with browser shortcuts

### Integration with HomePageContent

**Update HomePageContent component** (from Story 1.5):

Add filter control above AlbumGrid:

```tsx
// app/(authenticated)/_components/HomePageContent/index.tsx
import {AlbumFilterControl} from '../AlbumFilterControl'

function HomePageContent({
                             albums,
                             filterOptions,    // NEW
                             activeFilter,     // NEW
                             onFilterChange,   // NEW
                             loading,
                             error,
                             onRetry,
                             onAlbumClick
                         }: HomePageContentProps) {
    if (loading) return <PageLoadingIndicator/>
    if (error) return <ErrorDisplay error={error} onRetry={onRetry}/>
    if (albums.length === 0) return <EmptyState/>

    return (
        <Box>
            <AlbumFilterControl
                filterOptions={filterOptions}
                activeFilter={activeFilter}
                onFilterChange={onFilterChange}
            />
            <Box sx={{marginTop: 4}}>  {/* Spacing between filter and grid */}
                <AlbumGrid>
                    {albums.map(album => (
                        <AlbumCard key={...} album={album} onClick={onAlbumClick}/>
                    ))}
                </AlbumGrid>
            </Box>
        </Box>
    )
}
```

**Update CatalogClient to pass filter props**:

```tsx
// app/(authenticated)/_components/CatalogClient/index.tsx
'use client'

function CatalogClient({initialState, user}) {
    const [state, dispatch] = useReducer(catalogReducer, initialState)
    const handlers = useThunks(catalogThunks, {adapter, dispatch}, state)

    const handleFilterChange = (newFilter: AlbumFilterEntry) => {
        dispatch(albumFilterChanged(newFilter))
    }

    return (
        <AppLayout user={user}>
            <HomePageContent
                albums={state.albums}
                filterOptions={state.albumFilterOptions}  // NEW
                activeFilter={state.albumFilter}          // NEW
                onFilterChange={handleFilterChange}       // NEW
                loading={state.loading}
                error={state.error}
                onRetry={handlers.onPageRefresh}
                onAlbumClick={handleAlbumClick}
            />
        </AppLayout>
    )
}
```

### Session Persistence Strategy

**Requirement** (from Acceptance Criteria line 28-29): "Filter selection persists as I navigate within the session" and "Filter is cleared when I log out"

**Implementation**: Filter state persists in React state (useReducer)

- User selects filter → dispatches action → state.albumFilter updates
- Navigate to album page → state remains in memory (React keeps it)
- Browser back button → state restored (React state persistence)
- Log out → unmount app → state cleared
- Page refresh → server initializes with default filter ("All Albums")

**NO localStorage or sessionStorage needed** - React state handles session persistence naturally within the Single Page Application context.

**Default Filter on Initial Load**: "All Albums" (show all accessible albums)

### Scroll Position Maintenance

**Requirement** (from Acceptance Criteria line 30): "My scroll position in the album list is maintained when I change filters"

**Implementation**: Browser maintains scroll position automatically when DOM elements remain in place

- AlbumGrid container element stays in same position
- Filter change updates albums array → re-renders cards
- No full page navigation → scroll position preserved
- NO additional code needed (browser handles it)

**Edge Case**: If filtered results are fewer and viewport is taller, scroll position may be at bottom with whitespace. This is acceptable behavior.

### Data Flow

```
User clicks filter button
    ↓
AlbumFilterControl calls props.onFilterChange(newFilter)
    ↓
CatalogClient's handleFilterChange dispatches action
    ↓
dispatch(albumFilterChanged(newFilter))
    ↓
Reducer updates state:
    - state.albumFilter = newFilter
    - state.albums = allAlbums.filter(matchesCriterion)
    ↓
HomePageContent re-renders with new filtered albums
    ↓
AlbumGrid displays updated album cards
    ↓
Scroll position maintained (no DOM unmount)
```

### Filter Option Generation

**Server-Side** (in onPageRefresh thunk or during state initialization):

```typescript
// Generate filter options from loaded albums
const generateFilterOptions = (
    allAlbums: Album[],
    currentUserId: string
): AlbumFilterEntry[] => {
    const options: AlbumFilterEntry[] = []

    // 1. "All Albums" option
    options.push({
        name: "All Albums",
        criterion: {owners: [], selfOwned: false},
        avatars: []
    })

    // 2. "My Albums" option
    options.push({
        name: "My Albums",
        criterion: {selfOwned: true},
        avatars: [currentUserPicture]
    })

    // 3. One option per unique owner (excluding current user)
    const uniqueOwners = new Map<string, OwnerDetails>()
    allAlbums.forEach(album => {
        if (album.ownedBy && !uniqueOwners.has(album.albumId.owner)) {
            uniqueOwners.set(album.albumId.owner, album.ownedBy)
        }
    })

    uniqueOwners.forEach((ownerDetails, ownerId) => {
        options.push({
            name: ownerDetails.name,
            criterion: {owners: [ownerId]},
            avatars: ownerDetails.users.map(u => u.picture).filter(Boolean)
        })
    })

    return options
}
```

**Note**: This logic may already exist in the migrated code from Story 1.2. Verify before implementing.

### Visual Design Specifications

**AlbumFilterControl Styling** (following Story 1.4 patterns):

**Typography**:

- Label "Filter by Owner": Typography variant body1, text secondary color
- Button text: Typography variant button, text primary color when inactive, white when active

**Colors** (from theme.ts):

- Inactive button: background transparent, border rgba(255,255,255,0.2), text primary
- Active button: background brand blue (#185986), border brand blue, text white
- Hover (inactive): background rgba(255,255,255,0.05)
- Focus: outline 2px solid brand blue

**Spacing**:

- Label to buttons: 8px (theme.spacing(1))
- Between buttons: 0 (ToggleButtonGroup handles it)
- Filter control to AlbumGrid: 32px (theme.spacing(4))

**Layout**:

```typescript
sx = {
{
    display: 'flex',
        flexDirection
:
    {
        xs: 'column', sm
    :
        'row'
    }
,  // Stack on mobile
    alignItems: {
        xs: 'stretch', sm
    :
        'center'
    }
,
    gap: 1,  // 8px
        marginBottom
:
    4  // 32px spacing to grid below
}
}
```

**ToggleButtonGroup Styling**:

```typescript
sx = {
{
    display: 'flex',
        flexWrap
:
    'wrap',  // Allow wrapping on small screens
        '& .MuiToggleButton-root'
:
    {
        minWidth: '120px',
            minHeight
    :
        '44px',  // Touch target requirement
            textTransform
    :
        'none',  // Don't uppercase button text
            color
    :
        'text.primary',
            borderColor
    :
        'rgba(255, 255, 255, 0.2)',
            '&.Mui-selected'
    :
        {
            backgroundColor: 'primary.main',  // Brand blue #185986
                color
        :
            '#ffffff',
                '&:hover'
        :
            {
                backgroundColor: 'primary.dark'
            }
        }
    }
}
}
```

### Critical Success Factors

1. **Filter State Already Works** - Story 1.2 migrated the filtering logic. This story only adds the UI control. Verify existing actions/reducer handle
   `albumFilterChanged`.

2. **UI Component is Pure** - AlbumFilterControl receives all data via props, calls callback on selection. NO direct state manipulation, NO business logic.

3. **Scroll Position Maintains Naturally** - Don't overthink this. Browser preserves scroll when elements stay in DOM. No special code needed.

4. **Session Persistence Works via React State** - Filter selection lives in useReducer state. Persists during navigation (SPA), clears on logout (unmount). NO
   localStorage.

5. **Responsive and Accessible** - ToggleButtonGroup works on all breakpoints, keyboard navigable, 44px touch targets, clear focus indicators.

6. **Default Filter is "All Albums"** - Initial state on page load shows all accessible albums (owned + shared).

---

## Implementation Guidance

This technical guidance has been validated by the lead developer, following it significantly increases the chance of getting your PR accepted. Any infringement
required to complete the story must be reported.

### Coding standards

You must follow the coding standard instructions from these files:

* `.github/instructions/nextjs.instructions.md`

### Tasks to complete

Implementing this story will require implementing the following tasks, but is not limited to it:

#### Phase 1: Verify Existing Filter State Management

* [ ] **Verify filter state structure in catalog state**
    * Location: `web-nextjs/domains/catalog/language/catalog-state.ts`
    * Confirm interface `CatalogViewerState` has:
        * `allAlbums: Album[]`
        * `albumFilterOptions: AlbumFilterEntry[]`
        * `albumFilter: AlbumFilterEntry`
        * `albums: Album[]` (derived/filtered array)
    * Confirm `AlbumFilterEntry` and `AlbumFilterCriterion` interfaces exist
    * Confirm `albumMatchCriterion(criterion)` filter function exists

* [ ] **Locate or create filter change action**
    * Search for existing action: `action-albumFilterChanged.ts` or similar
    * If doesn't exist, create: `domains/catalog/actions/action-albumFilterChanged.ts`
    * Action signature:
      ```typescript
      export const albumFilterChanged = (
          filter: AlbumFilterEntry
      ): Action<CatalogViewerState, { filter: AlbumFilterEntry }> => ({
          type: 'albumFilterChanged',
          payload: { filter },
          reduce: (state) => ({
              ...state,
              albumFilter: filter,
              albums: state.allAlbums.filter(albumMatchCriterion(filter.criterion))
          })
      })
      ```
    * Export from `domains/catalog/actions/index.ts`

* [ ] **Verify initial state includes filter options**
    * Location: `domains/catalog/language/initial-catalog-state.ts`
    * Confirm function `initialCatalogState()` returns state with:
        * `albumFilterOptions` initialized (at minimum with "All Albums" option)
        * `albumFilter` set to default ("All Albums")
        * `albums` initialized as empty array (will be populated by onPageRefresh thunk)
    * Default "All Albums" filter:
      ```typescript
      const defaultFilter: AlbumFilterEntry = {
          name: "All Albums",
          criterion: { owners: [], selfOwned: false },
          avatars: []
      }
      ```

* [ ] **Verify onPageRefresh thunk populates filter options**
    * Location: `domains/catalog/thunks/` (likely `thunk-onPageRefresh.ts` or similar)
    * After loading albums, thunk should dispatch action to populate `albumFilterOptions`
    * Options should include:
        1. "All Albums" (show all)
        2. "My Albums" (selfOwned: true)
        3. One per unique owner from loaded albums
    * If logic missing, add action `albumFilterOptionsGenerated(options: AlbumFilterEntry[])`
    * Generate options using helper function (see Technical Design "Filter Option Generation")

#### Phase 2: Create AlbumFilterControl Component

* [ ] **Create AlbumFilterControl component**
    * Location: `app/(authenticated)/_components/AlbumFilterControl/index.tsx`
    * Pure component accepting props:
      ```typescript
      interface AlbumFilterControlProps {
          filterOptions: AlbumFilterEntry[]
          activeFilter: AlbumFilterEntry
          onFilterChange: (filter: AlbumFilterEntry) => void
      }
      
      export const AlbumFilterControl = (props: AlbumFilterControlProps) => { ... }
      ```
    * Import MUI components:
      ```typescript
      import { Box, ToggleButtonGroup, ToggleButton, Typography } from '@mui/material'
      ```
    * Structure:
        * Label: "Filter by Owner" (Typography variant body1, text secondary)
        * ToggleButtonGroup with exclusive selection
        * Map filterOptions to ToggleButton elements
        * Selected button: `value === activeFilter.name`
        * onChange handler calls `props.onFilterChange(selectedFilter)`
    * Styling (MUI sx prop):
        * Flex layout: column on xs, row on sm+
        * Button min-width: 120px, min-height: 44px (touch target)
        * Selected button: background brand blue (#185986), text white
        * Inactive button: transparent background, border rgba(255,255,255,0.2)
        * Focus: brand blue outline
    * Accessibility:
        * ToggleButtonGroup has aria-label="Album filter by owner"
        * Each button aria-pressed indicates selected state
        * Keyboard: Tab, Arrow keys, Enter all work (MUI handles)

* [ ] **Create Ladle stories for AlbumFilterControl**
    * Location: `app/(authenticated)/_components/AlbumFilterControl/AlbumFilterControl.stories.tsx`
    * Follow Ladle patterns from `nextjs.instructions.md`
    * Story: Default with 3 options
      ```typescript
      import { action } from '@ladle/react'
      
      const SAMPLE_OPTIONS: AlbumFilterEntry[] = [
          { name: "All Albums", criterion: { owners: [], selfOwned: false }, avatars: [] },
          { name: "My Albums", criterion: { selfOwned: true }, avatars: ["user-pic.jpg"] },
          { name: "John Doe", criterion: { owners: ["john"], selfOwned: false }, avatars: ["john-pic.jpg"] }
      ]
      
      export const Default = () => (
          <AlbumFilterControl
              filterOptions={SAMPLE_OPTIONS}
              activeFilter={SAMPLE_OPTIONS[0]}
              onFilterChange={action('onFilterChange')}
          />
      )
      ```
    * Story: My Albums selected
      ```typescript
      export const MyAlbumsSelected = () => (
          <AlbumFilterControl
              filterOptions={SAMPLE_OPTIONS}
              activeFilter={SAMPLE_OPTIONS[1]}
              onFilterChange={action('onFilterChange')}
          />
      )
      ```
    * Story: Many options (5+ owners)
      ```typescript
      const MANY_OPTIONS = [
          SAMPLE_OPTIONS[0],
          SAMPLE_OPTIONS[1],
          { name: "Alice", criterion: { owners: ["alice"] }, avatars: [] },
          { name: "Bob", criterion: { owners: ["bob"] }, avatars: [] },
          { name: "Charlie", criterion: { owners: ["charlie"] }, avatars: [] },
      ]
      
      export const ManyOptions = () => (
          <AlbumFilterControl
              filterOptions={MANY_OPTIONS}
              activeFilter={MANY_OPTIONS[0]}
              onFilterChange={action('onFilterChange')}
          />
      )
      ```
    * Story: Mobile viewport (use Ladle width control)
    * Story: Desktop viewport

#### Phase 3: Update HomePageContent Component

* [ ] **Update HomePageContent to include filter control**
    * Location: `app/(authenticated)/_components/HomePageContent/index.tsx`
    * Add props to interface:
      ```typescript
      interface HomePageContentProps {
          albums: Album[]
          filterOptions: AlbumFilterEntry[]     // NEW
          activeFilter: AlbumFilterEntry        // NEW
          onFilterChange: (filter: AlbumFilterEntry) => void  // NEW
          loading: boolean
          error: Error | null
          currentUserId: string
          onRetry: () => void
          onAlbumClick: (albumId: string, ownerId: string) => void
      }
      ```
    * Import AlbumFilterControl:
      ```typescript
      import { AlbumFilterControl } from '../AlbumFilterControl'
      ```
    * Update render to include filter above grid:
      ```tsx
      if (loading) return <PageLoadingIndicator />
      if (error) return <ErrorDisplay error={error} onRetry={onRetry} />
      if (albums.length === 0 && activeFilter.criterion.owners.length === 0) {
          // Empty state only if "All Albums" selected and no albums exist
          return <EmptyState ... />
      }
      
      return (
          <Box>
              <AlbumFilterControl
                  filterOptions={filterOptions}
                  activeFilter={activeFilter}
                  onFilterChange={onFilterChange}
              />
              {albums.length === 0 ? (
                  <Typography sx={{ marginTop: 4, textAlign: 'center', color: 'text.secondary' }}>
                      No albums match this filter.
                  </Typography>
              ) : (
                  <Box sx={{ marginTop: 4 }}>
                      <AlbumGrid>
                          {albums.map(album => (
                              <AlbumCard key={...} album={album} onClick={onAlbumClick} />
                          ))}
                      </AlbumGrid>
                  </Box>
              )}
          </Box>
      )
      ```

* [ ] **Update HomePageContent Ladle stories to include filter**
    * Location: `app/(authenticated)/_components/HomePageContent/HomePageContent.stories.tsx`
    * Add filter props to all existing stories:
        * filterOptions: SAMPLE_OPTIONS (same as AlbumFilterControl stories)
        * activeFilter: SAMPLE_OPTIONS[0] (or vary by story)
        * onFilterChange: action('onFilterChange')
    * Add new story: Filtered with no matches
      ```typescript
      export const FilteredNoMatches = () => (
          <HomePageContent
              albums={[]}
              filterOptions={SAMPLE_OPTIONS}
              activeFilter={SAMPLE_OPTIONS[1]}  // "My Albums" selected but none exist
              onFilterChange={action('onFilterChange')}
              loading={false}
              error={null}
              currentUserId="user-123"
              onRetry={action('onRetry')}
              onAlbumClick={action('onAlbumClick')}
          />
      )
      ```

#### Phase 4: Update CatalogClient Component

* [ ] **Update CatalogClient to pass filter props to HomePageContent**
    * Location: `app/(authenticated)/_components/CatalogClient/index.tsx`
    * Import action if needed:
      ```typescript
      import { albumFilterChanged } from '@/domains/catalog/actions/action-albumFilterChanged'
      ```
    * Create filter change handler:
      ```typescript
      const handleFilterChange = useCallback((newFilter: AlbumFilterEntry) => {
          dispatch(albumFilterChanged(newFilter))
      }, [dispatch])
      ```
    * Pass new props to HomePageContent:
      ```tsx
      <HomePageContent
          albums={state.albums}
          filterOptions={state.albumFilterOptions}   // NEW
          activeFilter={state.albumFilter}           // NEW
          onFilterChange={handleFilterChange}        // NEW
          loading={state.albumsLoaded}
          error={state.error}
          currentUserId={...}
          onRetry={handlers.onPageRefresh}
          onAlbumClick={handleAlbumClick}
      />
      ```

#### Phase 5: Testing & Verification

* [ ] **Run Ladle to verify visual states**
    * Execute: `npm run ladle`
    * Open http://localhost:61000
    * Navigate to AlbumFilterControl stories
    * Verify all filter options display correctly
    * Verify active state shows brand blue background
    * Verify hover states work
    * Test keyboard navigation (Tab, Arrow keys, Enter)
    * Test responsive behavior using Ladle viewport controls
    * Take screenshots for documentation

* [ ] **Run existing unit tests**
    * Execute: `npm run test`
    * Verify all 230+ tests from Story 1.2 still pass
    * If filter change action was created, add unit test:
      ```typescript
      // action-albumFilterChanged.test.ts
      describe('albumFilterChanged', () => {
          it('should update active filter and recompute albums', () => {
              const allAlbums: Album[] = [
                  // Sample albums (owned and shared)
              ]
              const initialState: CatalogViewerState = {
                  ...initialCatalogState(),
                  allAlbums,
                  albums: allAlbums
              }
              
              const myAlbumsFilter: AlbumFilterEntry = {
                  name: "My Albums",
                  criterion: { selfOwned: true },
                  avatars: []
              }
              
              const action = albumFilterChanged(myAlbumsFilter)
              const newState = catalogReducer(initialState, action)
              
              expect(newState.albumFilter).toEqual(myAlbumsFilter)
              expect(newState.albums).toEqual(
                  allAlbums.filter(a => albumIsOwnedByCurrentUser(a))
              )
          })
      })
      ```

* [ ] **Manual testing in browser**
    * Run dev server: `npm run dev`
    * Navigate to home page: http://localhost:3000/
    * Verify filter control displays above album grid
    * Click "All Albums" → see all albums (owned + shared)
    * Click "My Albums" → see only owned albums
    * Click owner name → see only albums from that owner
    * Verify active filter has brand blue background
    * Verify scroll position maintained when changing filters
    * Navigate to album page (click card) → back button → filter selection preserved
    * Test responsive breakpoints (resize browser):
        * Mobile (400px): Vertical stack or horizontal scroll
        * Tablet (700px): Horizontal layout
        * Desktop (1200px): Horizontal layout
    * Test keyboard navigation:
        * Tab to filter control
        * Arrow keys to navigate options
        * Enter to select
        * Verify focus indicator (brand blue outline)
    * Test empty filter results:
        * If no albums match filter, see "No albums match this filter" message
        * If no albums at all (on "All Albums"), see empty state invitation

* [ ] **Verify visual design compliance**
    * Active filter: brand blue (#185986) background, white text
    * Inactive filters: transparent background, border rgba(255,255,255,0.2)
    * Focus indicator: 2px solid brand blue outline
    * Touch targets: minimum 44px height on mobile
    * Spacing: 32px between filter and grid (theme.spacing(4))
    * Typography: consistent with design system

#### Phase 6: Edge Cases and Error Handling

* [ ] **Handle edge case: No filter options available**
    * If `filterOptions` is empty array, hide AlbumFilterControl
    * Conditional render in HomePageContent:
      ```tsx
      {filterOptions.length > 0 && (
          <AlbumFilterControl ... />
      )}
      ```

* [ ] **Handle edge case: Filtered albums empty**
    * Display message: "No albums match this filter."
    * Do NOT show EmptyState (that's for truly empty catalog)
    * User can change filter to see other albums

* [ ] **Handle edge case: Single filter option**
    * If only "All Albums" exists (no shared albums, current user not owner):
        * Still display filter control (shows "All Albums" selected)
        * OR hide filter control if only 1 option (design decision)
    * Recommendation: Always show filter even with 1 option (consistency)

* [ ] **Verify filter persists during navigation**
    * User selects "My Albums"
    * Clicks album card → navigates to album page
    * Presses back button
    * Filter should still be "My Albums" (React state preserved)
    * Filtered albums still displayed

### Target files structure

You will be expected to make changes on the following files:

```
web-nextjs/
├── domains/
│   └── catalog/
│       ├── actions/
│       │   ├── action-albumFilterChanged.ts     # NEW or VERIFY exists
│       │   └── index.ts                         # UPDATE: export new action
│       │
│       └── language/
│           ├── catalog-state.ts                 # VERIFY: filter types exist
│           └── initial-catalog-state.ts         # VERIFY: default filter set
│
├── app/
│   └── (authenticated)/
│       └── _components/
│           ├── AlbumFilterControl/
│           │   ├── index.tsx                    # NEW: Filter UI component
│           │   └── AlbumFilterControl.stories.tsx  # NEW: Ladle stories
│           │
│           ├── HomePageContent/
│           │   ├── index.tsx                    # UPDATE: Add filter control above grid
│           │   └── HomePageContent.stories.tsx  # UPDATE: Add filter props to stories
│           │
│           └── CatalogClient/
│               └── index.tsx                    # UPDATE: Pass filter props to HomePageContent
```

### Important Implementation Notes

**Filter State Already Exists**:

- Story 1.2 migrated `CatalogViewerState` with filter fields
- `albumFilterOptions`, `albumFilter`, `albums` already in state
- Reducer likely already handles filter changes
- Verify before creating new actions/logic

**Component is Pure**:

- AlbumFilterControl receives data, calls callback
- NO useState, NO useEffect, NO business logic
- ALL filtering logic in reducer and selectors
- UI only displays and reports user selection

**Scroll Position Maintenance**:

- Browser automatically preserves scroll when DOM stays mounted
- Filter change re-renders AlbumGrid in place
- NO special code needed (ScrollRestoration, window.scrollTo, etc.)
- DO NOT overthink this requirement

**Session Persistence**:

- Filter state lives in React useReducer
- Persists during SPA navigation (state in memory)
- Clears on logout (component unmount)
- Page refresh resets to default "All Albums"
- NO localStorage, NO sessionStorage needed

**Default Filter**:

- Initial page load shows "All Albums" (all accessible)
- Most inclusive default - user sees everything they can access
- Explicitly set in `initialCatalogState()`

**Filter Option Generation**:

- Server-side during onPageRefresh thunk
- Extract unique owners from loaded albums
- Generate AlbumFilterEntry for each owner
- Always include "All Albums" and "My Albums"
- May need to dispatch action to populate `state.albumFilterOptions`

**Empty Filter Results**:

- Distinguished from empty catalog
- Empty catalog (no albums at all): Show EmptyState with "Create Album" invitation
- Empty filter (no albums match): Show message "No albums match this filter"
- User can change filter to see other albums

**Accessibility**:

- MUI ToggleButtonGroup handles most accessibility (ARIA, keyboard)
- Ensure aria-label on group: "Album filter by owner"
- Verify keyboard nav works: Tab, Arrow keys, Enter
- Focus indicator brand blue (#185986)
- Touch targets 44px minimum on mobile

**Responsive Design**:

- ToggleButtonGroup wraps on narrow screens
- Buttons min-width 120px, min-height 44px
- Flex layout: column on xs, row on sm+
- Test all breakpoints in Ladle and browser

**TypeScript Type Safety**:

- Use types from `domains/catalog/language/catalog-state.ts`
- AlbumFilterEntry, AlbumFilterCriterion already defined
- NO `any` types
- Export prop interfaces for components

**What NOT to Do**:

- ❌ DO NOT implement localStorage persistence (React state is sufficient)
- ❌ DO NOT add scroll position tracking code (browser handles it)
- ❌ DO NOT create new filter logic (use existing from Story 1.2 migration)
- ❌ DO NOT add business logic to UI components
- ❌ DO NOT implement album creation (Epic 3)
- ❌ DO NOT implement advanced filtering (date ranges, tags, search - out of scope)
- ❌ DO NOT modify existing tests from Story 1.2 (should still pass)
- ❌ DO NOT add animations beyond MUI defaults
- ❌ DO NOT put filter control in AppHeader (should be above album grid on page)

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem?
* What has been done to solve it?
* Results and screenshots when possible
