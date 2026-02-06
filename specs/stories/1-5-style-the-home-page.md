# Story 1.5: Style the Home Page

**Status**: ready-for-dev

---

## Story

As a **user**,  
I want to **see a beautifully designed home page with my albums**,  
so that **I can enjoy using the application and easily browse my collection**.

## Acceptance Criteria

**Given** the UI components are created (Story 1.4) and albums are loaded (Story 1.3)
**When** I view the home page
**Then** the page uses the AppLayout component with header showing my profile
**And** albums are displayed using AlbumCard components in an AlbumGrid
**And** each album card shows:

- Album name as prominent heading
- Date range in readable format
- Media count
- Owner information (avatar, name) when album is shared with me
- Density color indicator
- Sharing avatars when I've shared the album with others

**And** discrete loading feedback is shown when page loads:

- PageLoadingIndicator when initial page displays
- NavigationLoadingIndicator when navigating to album after clicking card

**And** when no albums exist, EmptyState invites me to create an album
**And** if loading fails, ErrorDisplay shows the error with "Try Again" button
**And** the page is responsive and adapts across mobile, tablet, and desktop
**And** the overall design uses dark theme with brand color #185986
**And** visual regression tests capture the styled home page in different states:

- Home page with albums
- Home page empty state
- Home page error state
- Home page loading state
- Mobile, tablet, desktop viewports

## Scope & Context

### Depends on

* **Story 1.1 (Project Foundation Setup)**: Material UI with dark theme (#185986), ThemeProvider configured
* **Story 1.2 (State Management Migration)**: Catalog state management migrated, fetch adapter created, image loader configured
* **Story 1.3 (Basic Album List Display)**: Albums loaded from API, basic text list display on home page `/`
* **Story 1.4 (UI Components & Layout)**: All UI components created and tested:
    - AppLayout & AppHeader (application structure with user profile)
    - AlbumCard & AlbumGrid (album display with metadata)
    - PageLoadingIndicator, NavigationLoadingIndicator, ErrorDisplay, EmptyState (feedback components)
    - UserAvatar, SharedByIndicator (user components)

### Expected outcomes

This story completes Epic 1's visual implementation. The home page demonstrates the design system in action and establishes the visual standards for subsequent
stories. The next story (1.6 Album Filtering) will add interactive filtering controls to this styled foundation.

**Key deliverables for next stories:**

- Fully styled home page implementing the dark theme and brand color system
- Complete catalog state integration pattern showing how server state initialization and client component hydration work together
- Visual regression tests establishing the baseline for design consistency
- Responsive layout patterns demonstrating mobile/tablet/desktop breakpoints
- Loading and error state patterns to be reused in album page and other views

### Out of scope

* DO NOT implement album filtering - this is covered by Story 1.6
* DO NOT implement album creation dialogs - covered by Epic 3
* DO NOT implement album page navigation - covered by Epic 2
* DO NOT add animations beyond Material UI defaults - polish is not in this story's scope
* DO NOT implement random photo highlights - covered by Epic 5

---

## Technical Context

### Architecture Reference

**File**: `/home/dush/dev/git/dphoto/specs/designs/architecture.md`

**Key architectural decisions:**

- **State Management (Decision 2)**: Lift-and-shift catalog state from `web/src/core/catalog/` - Initialize state server-side, pass to pure UI components as
  props
- **Component Architecture (Decision 3)**: Colocation principle - components in `_components/` subfolder
- **Routing Structure (Decision 4)**: Home page at `/`, albums at `/owners/[ownerId]/[albumId]`
- **Layout Architecture (Decision 6)**: Material UI `sx` prop with theme breakpoints for responsive layouts

### Product Requirements

**File**: `/home/dush/dev/git/dphoto/specs/designs/prd.md`

**Functional Requirements addressed:**

- FR1: Album owners can view all albums they own
- FR2: All users can view albums shared with them
- FR4: Users can view album metadata (name, date range, media count, owner information)
- FR5: Users can see which users an album is shared with
- FR6: Users can view albums in chronological order
- FR7: Users can see visual indicators of album activity density (temperature)
- FR36: Users can view their profile information (name, picture)
- FR37: System displays discrete loading feedback when a page is loaded in the background
- FR38: System invites user to the next step when no albums exist
- FR39: System shows explicit errors when something goes wrong
- FR40: System provides responsive layouts for mobile, tablet, and desktop devices

**Non-Functional Requirements:**

- NFR8: Layouts adapt across mobile (<600px), tablet (600-960px), desktop (>960px) breakpoints

### UX Design Specification

**File**: `/home/dush/dev/git/dphoto/specs/designs/ux-design-specification.md`

**Key UX requirements:**

- Dark-first theme with brand color #185986
- Photos as primary visual focus against dark background
- Album cards with random photo samples (3-4 thumbnails)
- Density color-coding: warmer colors (high density >10 photos/day), neutral (3-10), cooler (<3)
- Timeline/chronological navigation (newest first)
- Discrete loading feedback (no big spinners or skeleton screens)
- Responsive grid: Mobile (1 col), Tablet (2 cols), Desktop (3 cols), Large Desktop (4 cols)

### Epic Context

**File**: `/home/dush/dev/git/dphoto/specs/designs/epics.md`

**Epic 1: Album List Home Page**
Users can view their album list on the home page with album cards showing metadata and density indicators, and filter albums by owner with a fully functional
filter.

**Story progression:**

1. Story 1.1 ✅ - Material UI foundation with dark theme
2. Story 1.2 ✅ - State management migration
3. Story 1.3 ✅ - Basic album loading (text list)
4. Story 1.4 ✅ - UI components created and tested
5. **Story 1.5 (THIS STORY)** - Apply components to home page with full styling
6. Story 1.6 (NEXT) - Add album filtering functionality

### Project Structure Context

**Current NextJS structure:**

```
web-nextjs/
├── app/
│   ├── (authenticated)/
│   │   ├── layout.tsx              # Should wrap with CatalogContext
│   │   ├── page.tsx                # Home page - compute server state
│   │   └── _components/            # Page-specific components
│   │       └── CatalogClient.tsx   # Client component with state hydration
├── components/
│   ├── layout/
│   │   ├── AppLayout.tsx           # From Story 1.4
│   │   └── AppHeader.tsx           # From Story 1.4
│   ├── albums/
│   │   ├── AlbumCard.tsx           # From Story 1.4
│   │   └── AlbumGrid.tsx           # From Story 1.4
│   ├── feedback/
│   │   ├── PageLoadingIndicator.tsx    # From Story 1.4
│   │   ├── NavigationLoadingIndicator.tsx  # From Story 1.4
│   │   ├── ErrorDisplay.tsx        # From Story 1.4
│   │   └── EmptyState.tsx          # From Story 1.4
│   ├── user/
│   │   ├── UserAvatar.tsx          # From Story 1.4
│   │   └── SharedByIndicator.tsx   # From Story 1.4
│   └── theme/
│       └── theme.ts                # From Story 1.1
├── domains/
│   └── catalog/                    # From Story 1.2
│       ├── language/               # State types
│       ├── actions.ts              # Reducer
│       └── adapters/
│           └── fetch-adapter.ts    # API adapter
└── libs/
    └── image-loader.ts             # From Story 1.2
```

### State Management Pattern

**From Architecture Decision 2:**

```tsx
// app/(authenticated)/page.tsx (Server Component)
export default async function HomePage() {
    // 1. Compute initial catalog state on server
    const catalogState = await computeInitialCatalogState()

    // 2. Pass to client component
    return <CatalogClient initialState={catalogState}/>
}

// app/(authenticated)/_components/CatalogClient.tsx (Client Component)
'use client'

function CatalogClient({initialState}) {
    // 3. Hydrate state with reducer
    const [state, dispatch] = useReducer(catalogReducer, initialState)

    // 4. Instantiate thunks with handlers
    const handlers = useThunks(catalogThunks, {adapter, dispatch}, state)

    // 5. Render pure UI components with props
    return (
        <AppLayout user={state.user}>
            <AlbumsPageContent
                albums={state.visibleAlbums}
                loading={state.loading}
                error={state.error}
                onRetry={handlers.loadAlbums}
            />
        </AppLayout>
    )
}

// app/(authenticated)/_components/AlbumsPageContent.tsx (Pure UI)
function AlbumsPageContent({albums, loading, error, onRetry}) {
    if (loading) return <PageLoadingIndicator/>
    if (error) return <ErrorDisplay error={error} onRetry={onRetry}/>
    if (albums.length === 0) return <EmptyState/>

    return <AlbumGrid albums={albums}/>
}
```

### Coding Standards

**File**: `/home/dush/dev/git/dphoto/.github/instructions/nextjs.instructions.md`

**Key standards:**

- Always use TypeScript with strict mode
- Material UI `sx` prop for all styling (no inline styles, no CSS modules)
- Component colocation: page-specific in `_components/`, shared in `components/`
- Visual regression tests using Playwright
- Server Components by default, Client Components only when needed (state, events, browser APIs)

---

## Technical Design

### Overview

This story **integrates the visual components from Story 1.4 into the home page**, creating the complete user-facing experience. We update the existing server
component (`page.tsx`) from Story 1.3 to use the styled components instead of basic text. This completes the assembly of server state initialization (Story
1.3), visual components (Story 1.4), and page composition (Story 1.5).

The core pattern remains unchanged from Story 1.3: Server Component → Client Component with state hydration → Pure UI components. We're enhancing the UI layer
with the styled components created in Story 1.4.

### Architecture Pattern

Following **Architecture.md Decision #2** (Server-Side State Initialization):

```
Server Component (page.tsx)
   ↓ (execute thunk, compute state)
   ↓ (pass as initialState prop)
Client Component (CatalogClient.tsx) 
   ├─ useReducer(catalogReducer, initialState)
   ├─ useThunks(catalogThunks, {adapter, dispatch}, state)
   └─ Pure UI → HomePageContent.tsx
        ├─ Loading state → PageLoadingIndicator
        ├─ Error state → ErrorDisplay
        ├─ Empty state → EmptyState  
        └─ Success state → AlbumGrid with AlbumCard[]
```

**Critical from Story 1.3**: The server-side thunk execution infrastructure (`constructThunkFromDeclaration` and `newServerAdapterFactory`) is already
implemented. This story focuses on **composing the UI components** rather than creating new infrastructure.

### Component Composition Strategy

**Layer 1 - Server Component (`app/(authenticated)/page.tsx`)**

- Already exists from Story 1.3
- Executes `catalogThunks.onPageRefresh` server-side
- Computes initial `CatalogViewerState`
- Passes to `<CatalogClient initialState={catalogState} />`
- **Update needed**: May need to ensure user data is included in state for AppHeader

**Layer 2 - Client Component (`app/(authenticated)/_components/CatalogClient.tsx`)**

- Already exists from Story 1.3 (possibly named `CatalogProvider` or `AlbumListClient`)
- Hydrates state with `useReducer(catalogReducer, initialState)`
- Instantiates handlers with `useThunks(catalogThunks, factoryArgs, state)`
- **Update needed**: Replace basic rendering with AppLayout wrapper
- Pass selected state slices and handlers to `<HomePageContent />`

**Layer 3 - Pure UI Component (`app/(authenticated)/_components/HomePageContent.tsx`)**

- New component to create
- Receives: `albums`, `loading`, `error`, `onRetry`, `onAlbumClick`
- Conditional rendering based on state:
    - Loading → `<PageLoadingIndicator />`
    - Error → `<ErrorDisplay error={error} onRetry={onRetry} />`
    - Empty → `<EmptyState />` with album creation invitation
    - Success → `<AlbumGrid>{albums.map(album => <AlbumCard />)}</AlbumGrid>`
- NO state management, NO business logic

**Layer 4 - AppLayout Integration**

- Wrap entire page with `<AppLayout user={user}>` in CatalogClient
- AppLayout displays AppHeader with user profile
- Content area receives `<HomePageContent />` as children

### State Management Integration

**From Story 1.3, we have:**

- `CatalogViewerState` with albums, selected album, loading, error states
- `catalogReducer` reducing actions
- `catalogThunks` with `onPageRefresh` thunk
- `constructThunkFromDeclaration` for server-side execution
- `newServerAdapterFactory` creating server-compatible adapter

**This story uses:**

- **State slices**: Extract from `CatalogViewerState`:
    - `state.albums` → list of albums to display
    - `state.loading` → boolean for loading state
    - `state.error` → error object or null
    - `state.authenticatedUser` → user data for AppHeader
- **Selector (optional)**: `selectVisibleAlbums(state)` if albums need filtering/sorting
- **Handlers**:
    - `handlers.onPageRefresh` for retry on error
    - `handlers.navigateToAlbum` or direct Next.js navigation for album clicks

**User Data Flow:**

- Server component retrieves `authenticatedUser` from session (existing from Story 1.1)
- Include in initial catalog state or pass separately to CatalogClient
- CatalogClient passes to `<AppLayout user={...}>`
- AppHeader displays user avatar and name

### AlbumCard Data Mapping

**From CatalogViewerState album to AlbumCard props:**

```typescript
// CatalogViewerState album structure (from Story 1.2 migration)
interface Album {
    albumId: string;
    ownerId: string;
    name: string;
    startDate: string; // ISO date
    endDate: string;   // ISO date
    mediaCount: number;
    folderName?: string;
    sharedWith?: SharedUser[];
}

// AlbumCard props (from Story 1.4)
interface AlbumCardProps {
    album: {
        albumId: string;
        ownerId: string;
        name: string;
        startDate: string;
        endDate: string;
        mediaCount: number;
    };
    owner?: {       // Present when album shared with current user
        name: string;
        email: string;
        picture?: string;
    };
    sharedWith?: Array<{  // Present when current user owns album
        name: string;
        email: string;
        picture?: string;
    }>;
    onClick: (albumId: string, ownerId: string) => void;
}
```

**Mapping Logic:**

1. Check if `album.ownerId === currentUser.id`:
    - If true: user owns album → pass `sharedWith` array
    - If false: album shared with user → pass `owner` information
2. Format dates for display using AlbumCard's internal formatting
3. Calculate density using AlbumCard's internal calculation
4. onClick handler navigates to `/owners/${ownerId}/${albumId}`

### Responsive Layout Implementation

**AlbumGrid Breakpoints (from Story 1.4):**

```typescript
sx = {
{
    display: 'grid',
        gridTemplateColumns
:
    {
        xs: '1fr',                    // Mobile: 1 column
            sm
    :
        'repeat(2, 1fr)',         // Tablet: 2 columns
            md
    :
        'repeat(3, 1fr)',         // Desktop: 3 columns
            lg
    :
        'repeat(4, 1fr)',         // Large: 4 columns
    }
,
    gap: 4, // 32px
        maxWidth
:
    '1920px',
        margin
:
    '0 auto',
}
}
```

**AppLayout Padding (from Story 1.4):**

```typescript
sx = {
{
    padding: {
        xs: 2,  // 16px on mobile
            sm
    :
        3,  // 24px on tablet
            md
    :
        4,  // 32px on desktop
    }
}
}
```

**AppHeader Responsive Behavior:**

- Mobile (xs): Logo + UserAvatar only (hide user name)
- Tablet/Desktop (sm+): Logo + UserAvatar + user name

### Visual Design Implementation

**Typography (from UX Design Specification):**

- Album names: 22px, serif (Georgia fallback), weight 300
- Metadata: 13px, monospace (Courier New fallback), uppercase, letter-spacing 0.1em
- Implemented within AlbumCard component (Story 1.4)

**Color System (from UX Design & Story 1.1 theme):**

- Background: #121212
- Surface (cards): #1e1e1e
- Primary (brand blue): #185986
- Text primary: #ffffff
- Text secondary: rgba(255, 255, 255, 0.7)

**Density Color-Coding (from AlbumCard component):**

- High density (>10 photos/day): #ff6b6b (warm red)
- Medium density (3-10 photos/day): #ffd43b (neutral yellow)
- Low density (<3 photos/day): #51cf66 (cool green)

**Spacing (from theme.spacing()):**

- Base unit: 8px
- Card padding: 24px (theme.spacing(3))
- Grid gap: 32px (theme.spacing(4))
- Section spacing: 32-48px

**Dark Theme Application:**

- Handled by ThemeProvider from Story 1.1
- All MUI components inherit dark theme palette
- Photos pop dramatically against dark background

### Loading & Error State Patterns

**Loading States:**

1. **Initial Page Load**:
    - Server component computes state (may take time for API call)
    - Display `<PageLoadingIndicator />` if `state.loading === true`
    - Thin progress bar at top, discrete and non-blocking

2. **Navigation Loading**:
    - When user clicks AlbumCard to navigate to album page
    - Next.js shows `<NavigationLoadingIndicator />` during transition
    - Small spinner, inline or overlay

**Error States:**

- `<ErrorDisplay />` receives error object with message, code, details
- Displays user-friendly message with technical details collapsible
- "Try Again" button calls `handlers.onPageRefresh` to retry
- Error remains displayed until retry succeeds or user navigates away

**Empty State:**

- When `albums.length === 0` and no error
- `<EmptyState />` with PhotoLibrary icon
- Title: "No albums found"
- Message: "Create your first album to get started."
- Action button: "Create Album" → opens creation dialog (Story 1.6 or Epic 3)

### Critical Success Factors

1. **Component Integration Must Be Clean**
    - Server → Client → Pure UI data flow is clear and unidirectional
    - No prop drilling beyond one level
    - State transformations happen in selectors, not in render functions
    - Type safety maintained throughout the chain

2. **Responsive Layout Must Work Seamlessly**
    - Mobile, tablet, desktop layouts all feel native
    - Touch targets minimum 44px on mobile
    - Text remains readable at all breakpoints
    - No horizontal scrolling at any viewport width

3. **Loading States Must Be Discrete**
    - PageLoadingIndicator is thin progress bar, not full-page blocker
    - Initial load feels fast (server-side state initialization)
    - No jarring skeleton screens or layout shifts
    - Progressive enhancement, not blocking loading

4. **Error Handling Must Be Complete**
    - All error scenarios have clear recovery paths
    - Network errors, permission errors, empty states all handled
    - User is never stuck without forward action
    - Technical details available but not overwhelming

5. **Visual Consistency Must Be Maintained**
    - Brand blue (#185986) used consistently
    - Dark theme applied throughout
    - Spacing follows theme tokens
    - Typography scales correctly across breakpoints

---

## Implementation Guidance

This technical guidance has been validated by the lead developer, following it significantly increases the chance of getting your PR accepted. Any infringement
required to complete the story must be reported.

### Coding standards

You must follow the coding standard instructions from these files:

* `.github/instructions/nextjs.instructions.md`

### Tasks to complete

Implementing this story will require implementing the following tasks, but is not limited to it:

#### Phase 1: Update Server Component

* [ ] **Update `app/(authenticated)/page.tsx` to include user data in catalog state**
    * Retrieve authenticated user from session (existing pattern from Story 1.1)
    * Pass user data to CatalogClient either:
        - As separate prop: `<CatalogClient initialState={catalogState} user={authenticatedUser} />`
        - Or include in catalog state if `CatalogViewerState` has user field
    * Ensure server-side thunk execution (from Story 1.3) includes loading albums
    * Verify state includes: albums list, loading flag, error object, user data

#### Phase 2: Update Client Component

* [ ] **Update `app/(authenticated)/_components/CatalogClient.tsx` (or create if doesn't exist)**
    * Mark as `'use client'`
    * Accept props: `initialState: CatalogViewerState`, optionally `user` if passed separately
    * Hydrate state: `const [state, dispatch] = useReducer(catalogReducer, initialState)`
    * Instantiate handlers: `const handlers = useThunks(catalogThunks, {adapter: getAdapter(), dispatch}, state)`
    * Import AppLayout: `import { AppLayout } from '@/components/layout/AppLayout'`
    * Wrap entire render with AppLayout:
      ```tsx
      return (
        <AppLayout user={user}>
          <HomePageContent 
            albums={state.albums}
            loading={state.loading}
            error={state.error}
            onRetry={handlers.onPageRefresh}
            onAlbumClick={handleAlbumClick}
          />
        </AppLayout>
      )
      ```
    * Define `handleAlbumClick` using Next.js navigation:
      ```tsx
      const router = useRouter()
      const handleAlbumClick = (albumId: string, ownerId: string) => {
        router.push(`/owners/${ownerId}/${albumId}`)
      }
      ```

#### Phase 3: Create Pure UI Component

* [ ] **Create `app/(authenticated)/_components/HomePageContent/index.tsx`**
    * Pure component accepting props:
      ```tsx
      interface HomePageContentProps {
        albums: Album[];
        loading: boolean;
        error: Error | null;
        currentUserId: string; // For determining ownership
        onRetry: () => void;
        onAlbumClick: (albumId: string, ownerId: string) => void;
      }
      ```
    * Import all UI components from Story 1.4:
      ```tsx
      import { AlbumGrid } from '@/components/albums/AlbumGrid'
      import { AlbumCard } from '@/components/albums/AlbumCard'
      import { PageLoadingIndicator } from '@/components/feedback/PageLoadingIndicator'
      import { ErrorDisplay } from '@/components/feedback/ErrorDisplay'
      import { EmptyState } from '@/components/feedback/EmptyState'
      ```
    * Implement conditional rendering:
      ```tsx
      if (loading) return <PageLoadingIndicator message="Loading albums..." />
      if (error) return <ErrorDisplay error={error} onRetry={onRetry} />
      if (albums.length === 0) return <EmptyState 
        title="No albums found"
        message="Create your first album to get started."
        icon={<PhotoLibraryIcon />}
      />
      ```
    * Render albums in grid:
      ```tsx
      return (
        <AlbumGrid>
          {albums.map(album => {
            const isOwner = album.ownerId === currentUserId
            return (
              <AlbumCard
                key={album.albumId}
                album={album}
                owner={isOwner ? undefined : getOwnerInfo(album)}
                sharedWith={isOwner ? album.sharedWith : undefined}
                onClick={() => onAlbumClick(album.albumId, album.ownerId)}
              />
            )
          })}
        </AlbumGrid>
      )
      ```
    * Helper function for owner info:
      ```tsx
      const getOwnerInfo = (album: Album) => {
        // Extract owner info from album metadata
        // May need to fetch from user registry or include in album object
        return {
          name: album.ownerName || 'Unknown',
          email: album.ownerEmail || '',
          picture: album.ownerPicture
        }
      }
      ```

* [ ] **Create Ladle stories for HomePageContent: `HomePageContent.stories.tsx`**
    * Story: Default with 3-5 albums
      ```tsx
      export const Default = () => (
        <HomePageContent
          albums={SAMPLE_ALBUMS}
          loading={false}
          error={null}
          currentUserId="user-123"
          onRetry={action('onRetry')}
          onAlbumClick={action('onAlbumClick')}
        />
      )
      ```
    * Story: Loading state
      ```tsx
      export const Loading = () => (
        <HomePageContent
          albums={[]}
          loading={true}
          error={null}
          currentUserId="user-123"
          onRetry={action('onRetry')}
          onAlbumClick={action('onAlbumClick')}
        />
      )
      ```
    * Story: Error state
      ```tsx
      export const Error = () => (
        <HomePageContent
          albums={[]}
          loading={false}
          error={new Error('Failed to load albums')}
          currentUserId="user-123"
          onRetry={action('onRetry')}
          onAlbumClick={action('onAlbumClick')}
        />
      )
      ```
    * Story: Empty state
      ```tsx
      export const Empty = () => (
        <HomePageContent
          albums={[]}
          loading={false}
          error={null}
          currentUserId="user-123"
          onRetry={action('onRetry')}
          onAlbumClick={action('onAlbumClick')}
        />
      )
      ```
    * Story: Mobile viewport (set via Ladle width control)
    * Story: Tablet viewport
    * Story: Desktop viewport

#### Phase 4: Data Mapping Helpers

* [ ] **Create helper functions for album data transformation (if needed)**
    * Location: `app/(authenticated)/_components/HomePageContent/helpers.ts` or inline
    * Function to determine ownership:
      ```tsx
      const isAlbumOwner = (album: Album, currentUserId: string): boolean => {
        return album.ownerId === currentUserId
      }
      ```
    * Function to extract owner information from shared album:
      ```tsx
      const getOwnerFromAlbum = (album: Album): OwnerInfo | undefined => {
        if (!album.ownerName) return undefined
        return {
          name: album.ownerName,
          email: album.ownerEmail || '',
          picture: album.ownerPicture
        }
      }
      ```
    * Function to format album for AlbumCard (if needed):
      ```tsx
      const mapAlbumToCardProps = (
        album: Album, 
        currentUserId: string
      ): AlbumCardProps => {
        const isOwner = isAlbumOwner(album, currentUserId)
        return {
          album: {
            albumId: album.albumId,
            ownerId: album.ownerId,
            name: album.name,
            startDate: album.startDate,
            endDate: album.endDate,
            mediaCount: album.mediaCount,
          },
          owner: isOwner ? undefined : getOwnerFromAlbum(album),
          sharedWith: isOwner ? album.sharedWith : undefined,
          onClick: (albumId, ownerId) => { /* navigation */ }
        }
      }
      ```

#### Phase 5: Album Page Placeholder (if not from Story 1.3)

* [ ] **Verify album page route exists: `app/(authenticated)/owners/[ownerId]/[albumId]/page.tsx`**
    * If doesn't exist, create placeholder:
      ```tsx
      export default async function AlbumPage({ 
        params 
      }: { 
        params: { ownerId: string; albumId: string } 
      }) {
        return (
          <Box sx={{ padding: 3 }}>
            <Typography variant="h4">Album: {params.albumId}</Typography>
            <Typography>Owner: {params.ownerId}</Typography>
            <Typography sx={{ marginTop: 2 }}>
              Album viewing will be implemented in Epic 2.
            </Typography>
            <Button
              component={Link}
              href="/"
              prefetch={false}
            >
              Back to Albums
            </Button>
          </Box>
        )
      }
      ```
    * This enables navigation testing from home page
    * Actual album implementation is Epic 2

#### Phase 6: Verify AppLayout Integration

* [ ] **Ensure AppLayout is used correctly**
    * Verify AppLayout wraps HomePageContent in CatalogClient
    * Verify user prop is passed with correct shape: `{name, email, picture}`
    * Test AppHeader displays user avatar and name correctly
    * Test responsive behavior: name hides on mobile (xs), shows on tablet/desktop (sm+)
    * Verify logo links to "/"
    * Verify fixed header doesn't obscure content (content has top padding)

#### Phase 7: Testing & Verification

* [ ] **Run Ladle to verify all visual states**
    * Execute: `npm run ladle`
    * Open http://localhost:61000
    * Navigate to HomePageContent stories
    * Verify all states render correctly:
        - Default with albums in grid layout
        - Loading with PageLoadingIndicator
        - Error with ErrorDisplay and retry button
        - Empty with EmptyState invitation
    * Test responsive viewports using Ladle controls
    * Take screenshots for documentation

* [ ] **Run existing unit tests**
    * Execute: `npm run test`
    * Verify all 230+ tests from Story 1.2 still pass
    * NO modifications to existing tests should be needed
    * Fix any import issues if they arise

* [ ] **Manual testing in browser**
    * Run dev server: `npm run dev`
    * Navigate to home page: http://localhost:3000/
    * Verify albums display in grid layout
    * Verify responsive breakpoints work (resize browser)
    * Click album card → navigates to album placeholder page
    * Verify AppHeader displays user info correctly
    * Test loading state (may need to simulate slow network)
    * Test error state (may need to simulate API failure)
    * Test empty state (temporarily modify to return empty albums)

* [ ] **Verify visual design compliance**
    * Brand blue (#185986) used on:
        - Primary buttons
        - Links
        - Focus indicators
        - Selected states
    * Dark theme throughout (#121212 background, #1e1e1e surfaces)
    * Typography matches specifications (serif album names, monospace metadata)
    * Density colors display correctly (red/yellow/green based on calculation)
    * Spacing follows theme tokens
    * No layout shifts or horizontal scrolling

### Target files structure

You will be expected to make changes on the following files:

```
web-nextjs/
├── app/
│   └── (authenticated)/
│       ├── page.tsx                                    # UPDATE: Include user in state/props
│       │
│       ├── _components/
│       │   ├── CatalogClient/
│       │   │   └── index.tsx                           # UPDATE: Wrap with AppLayout, pass to HomePageContent
│       │   │
│       │   └── HomePageContent/
│       │       ├── index.tsx                           # NEW: Pure UI component with conditional rendering
│       │       ├── HomePageContent.stories.tsx         # NEW: Ladle stories (default, loading, error, empty)
│       │       └── helpers.ts                          # NEW (optional): Data mapping helpers
│       │
│       └── owners/
│           └── [ownerId]/
│               └── [albumId]/
│                   └── page.tsx                        # VERIFY: Placeholder exists (from Story 1.3)
```

### Important Implementation Notes

**Reuse Story 1.3 Infrastructure:**

- Server-side thunk execution (`constructThunkFromDeclaration`) is already implemented
- Server adapter factory (`newServerAdapterFactory`) exists
- CatalogClient component exists (may be named differently)
- DO NOT reimplement these - UPDATE existing components

**Compose, Don't Create:**

- Story 1.4 created all visual components (AlbumCard, AlbumGrid, AppLayout, etc.)
- This story COMPOSES them into the page
- NO new styled components should be created
- Focus on integration and data flow

**State Extraction:**

- Use selectors if available: `selectVisibleAlbums(state)`
- Or extract directly: `state.albums`, `state.loading`, `state.error`
- Keep transformations minimal - defer to selectors
- Only map data where absolutely necessary for component props

**User Data Flow:**

- User data from session (Story 1.1 authentication)
- May already be in catalog state from Story 1.3
- If not, pass as separate prop to CatalogClient
- CatalogClient passes to AppLayout
- AppHeader renders user profile

**Album Ownership Logic:**

- Current user ID needed to determine ownership
- `album.ownerId === currentUserId` → user owns album
- If owns: pass `sharedWith` array to AlbumCard
- If doesn't own: pass `owner` info to AlbumCard
- May require additional user data in album object

**Navigation Pattern:**

- Use Next.js App Router navigation: `useRouter()` hook
- Import: `import { useRouter } from 'next/navigation'`
- Navigate on album click: `router.push(`/owners/${ownerId}/${albumId}`)`
- Client component only (use in CatalogClient or HomePageContent with callback)

**Responsive Testing:**

- Test all breakpoints: xs (<600px), sm (600-960px), md (960-1280px), lg (>1280px)
- Mobile: 1 column, name hidden in AppHeader
- Tablet: 2 columns, name visible in AppHeader
- Desktop: 3 columns
- Large: 4 columns
- Use browser DevTools responsive mode or Ladle viewport controls

**Error Handling:**

- All error states must have clear recovery actions
- "Try Again" button calls `handlers.onPageRefresh`
- Error Display shows user-friendly message + technical details
- Empty state provides clear invitation to create album
- Loading states are discrete and non-blocking

**TypeScript Type Safety:**

- Define explicit interfaces for all props
- NO `any` types
- Use types from domains/catalog for state and albums
- Export prop interfaces for reuse in stories

**Performance Considerations:**

- Server-side state initialization minimizes client-side loading
- Progressive image loading in AlbumCard (if implemented in Story 1.4)
- Avoid unnecessary re-renders (React.memo if needed)
- Lazy load dialogs and heavy components (future stories)

**What NOT to Do:**

- ❌ DO NOT create new styled components (use Story 1.4 components)
- ❌ DO NOT reimplement server-side thunk execution (exists from Story 1.3)
- ❌ DO NOT add filtering logic (Story 1.6)
- ❌ DO NOT implement album creation dialogs (Epic 3)
- ❌ DO NOT implement actual album page (Epic 2)
- ❌ DO NOT add random photo highlights (Epic 5)
- ❌ DO NOT add animations beyond MUI defaults
- ❌ DO NOT modify domain logic or state management (only UI composition)
- ❌ DO NOT put business logic in UI components
- ❌ DO NOT fetch data in client components (server initialized state)
- ❌ DO NOT modify existing tests from Story 1.2 (should still pass)

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem?
* What has been done to solve it?
* Results and screenshots when possible
