# Story 1.3: Basic Album List Display

Status: ready-for-dev

## Story

As a user,  
I want to view my album list on the authenticated home page,  
So that I can see all albums I own and albums shared with me in chronological order.

## Acceptance Criteria

**Given** I am an authenticated user  
**When** I navigate to the home page `/`  
**Then** I see a list of albums displayed in a responsive grid:

- Mobile (xs <600px): 1 column
- Tablet (sm 600-960px): 2 columns
- Desktop (md 960-1280px): 3 columns
- Large (lg >1280px): 4 columns

**And** each album card shows:

- Album name (Georgia serif font, 22px, weight 300)
- Date range formatted as "MMM DD – MMM DD, YYYY" (Courier New monospace, 13px, uppercase)
- Media count as "X photos" (cyan #6ab9de color, 13px)
- 4 random photo thumbnails in 2x2 grid (edge-to-edge, 8px gap)
- Text overlay with dark gradient at bottom: `rgba(10, 21, 32, 0.98)`

**And** album cards have hover interaction:

- Card lifts with shadow: `0 12px 40px rgba(24, 89, 134, 0.4)`
- Photos darken: `brightness(0.85)`
- Text overlay becomes blue gradient: `rgba(24, 89, 134, 0.95)`
- Album count changes from cyan (#6ab9de) to white

**And** albums are sorted chronologically by start date (newest first)

**And** clicking an album navigates to `/owners/{ownerId}/albums/{folderName}`

**And** the page shows appropriate states:

- **Loading:** Skeleton grid with placeholder cards
- **Empty:** Message "No albums yet" with helpful text
- **Error:** Error message with "Try Again" button that refreshes the page

**And** the page uses Server Component to fetch initial data:

- Fetch albums using adapter from Story 1.2
- Build minimal state structure (albums list only)
- Pass specific props to pure UI components (NOT full CatalogViewerState)
- No client state management (no useReducer, no handlers)

**And** the page background uses dark blue gradient: `linear-gradient(135deg, #0a1520 0%, #12242e 50%, #0f1d28 100%)`

**And** images use Next.js Image component with:

- Custom loader from Story 1.2 (width ≤360 → 360, width >360 → 2400)
- `fill` prop for responsive sizing
- `sizes` prop for responsive hints
- `style={{ objectFit: 'cover' }}` for proper aspect ratio

**And** all styling uses Material UI `sx` prop (NO inline styles, NO Tailwind)

**And** navigation uses `@/components/Link` wrapper (from Story 1.1) to avoid Next.js v16 Client Component restriction

## Tasks / Subtasks

- [ ] Create Server Component page (AC: Server Component, fetch data)
    - [ ] Create `app/(authenticated)/page.tsx` with async function
    - [ ] Import fetch adapter from `domains/catalog/adapters/fetch-adapter.ts`
    - [ ] Fetch albums using `fetchAlbumsAdapter()`
    - [ ] Handle try-catch for error state
    - [ ] Handle empty albums array for empty state
    - [ ] Add `export const revalidate = 0` to disable caching (dynamic data)
    - [ ] Pass albums array to pure UI components
    - [ ] Return appropriate component based on state (loading/error/empty/success)

- [ ] Create pure UI components in `_components/` (AC: pure UI, colocation)
    - [ ] Create `app/(authenticated)/_components/AlbumsPage.tsx`
        - [ ] Accept `albums` array as prop
        - [ ] Render page container with background gradient
        - [ ] Render "Your Albums" section title
        - [ ] Render `<AlbumsGrid>` component
    - [ ] Create `app/(authenticated)/_components/AlbumsGrid.tsx`
        - [ ] Accept `albums` array as prop
        - [ ] Use MUI Box with responsive grid layout
        - [ ] gridTemplateColumns: xs: 1, sm: 2, md: 3, lg: 4
        - [ ] gap: 4 (32px)
        - [ ] Map albums to `<AlbumCard>` components
    - [ ] Create `app/(authenticated)/_components/AlbumCard.tsx`
        - [ ] Accept specific album props (NOT full state): albumId, name, start, end, totalCount, ownerId, folderName, mediaIds
        - [ ] Wrap in Link component from `@/components/Link`
        - [ ] href: `/owners/${ownerId}/albums/${folderName}`
        - [ ] Render 2x2 photo grid with 4 random photos
        - [ ] Use Next.js Image with fill, sizes, custom loader
        - [ ] Render text overlay with album metadata
        - [ ] Apply hover effects via `sx` prop
        - [ ] Format date range: "MMM DD – MMM DD, YYYY"
        - [ ] Format media count: "X photos"
    - [ ] Create `app/(authenticated)/_components/LoadingSkeleton.tsx`
        - [ ] Use MUI Skeleton component
        - [ ] Render grid of skeleton cards matching album card layout
        - [ ] Show 8 placeholder cards
    - [ ] Create `app/(authenticated)/_components/EmptyState.tsx`
        - [ ] Use MUI Paper/Typography
        - [ ] Center message: "No albums yet"
        - [ ] Add helpful subtext about uploading photos via CLI
    - [ ] Create `app/(authenticated)/_components/ErrorDisplay.tsx`
        - [ ] Use MUI Paper/Typography/Button
        - [ ] Display error message
        - [ ] Add "Try Again" button with brand blue (#185986)
        - [ ] Button onClick triggers page refresh: `window.location.reload()`

- [ ] Implement album card styling (AC: design from HTML mockup)
    - [ ] Background gradient on page: `linear-gradient(135deg, #0a1520 0%, #12242e 50%, #0f1d28 100%)`
    - [ ] Section title styling: 13px, uppercase, letter-spacing 0.15em, underline accent
    - [ ] Album card container: cursor pointer, transition 0.4s
    - [ ] Photo grid: 2x2 layout, 8px gap, padding 0 (edge-to-edge)
    - [ ] Each photo: aspect-ratio 1, gradient placeholder background
    - [ ] Text overlay: absolute position at bottom, dark gradient
    - [ ] Album name: Georgia serif, 22px, weight 300
    - [ ] Date range: Courier New monospace, 13px, uppercase
    - [ ] Media count: cyan #6ab9de, 13px, weight 400
    - [ ] Hover transform: translateY(-6px)
    - [ ] Hover shadow: `0 12px 40px rgba(24, 89, 134, 0.4)`
    - [ ] Hover photos: brightness 0.85
    - [ ] Hover overlay: `rgba(24, 89, 134, 0.95)` gradient
    - [ ] Hover count: white color

- [ ] Implement image loading (AC: Next.js Image, custom loader)
    - [ ] Use Image from `next/image`
    - [ ] src prop: mediaId (passed to custom loader)
    - [ ] fill prop: true
    - [ ] sizes prop: "(max-width: 768px) 50vw, 170px"
    - [ ] alt prop: "" (decorative images)
    - [ ] style prop: `{{ objectFit: 'cover' }}`
    - [ ] Parent Box: position relative, aspectRatio 1
    - [ ] Custom loader already configured in Story 1.2

- [ ] Handle data fetching edge cases (AC: loading, error, empty)
    - [ ] Loading state: Show LoadingSkeleton while fetching
    - [ ] Empty state: Show EmptyState when albums.length === 0
    - [ ] Error state: Show ErrorDisplay with try-catch error
    - [ ] Success state: Show AlbumsPage with albums data

- [ ] Implement date formatting (AC: date range format)
    - [ ] Create helper function `formatDateRange(start, end)`
    - [ ] Parse ISO dates from album.start and album.end
    - [ ] Format as "MMM DD – MMM DD, YYYY" (e.g., "Jul 15 – Jul 22, 2026")
    - [ ] Handle same-year ranges (don't repeat year)
    - [ ] Use JavaScript Intl.DateTimeFormat or date-fns

- [ ] Add component prop interfaces (AC: specific props, NOT full state)
    - [ ] AlbumsPageProps: `{ albums: Album[] }`
    - [ ] AlbumsGridProps: `{ albums: Album[] }`
    - [ ] AlbumCardProps: `{ albumId, name, start, end, totalCount, ownerId, folderName, mediaIds: string[] }`
    - [ ] LoadingSkeletonProps: none
    - [ ] EmptyStateProps: none
    - [ ] ErrorDisplayProps: `{ error: Error | unknown }`

- [ ] Write component tests (AC: all tests pass)
    - [ ] Create `app/(authenticated)/_components/__tests__/AlbumCard.test.tsx`
        - [ ] Test renders album metadata correctly
        - [ ] Test renders 4 photo thumbnails
        - [ ] Test Link href is correct
        - [ ] Test date formatting
        - [ ] Test media count formatting
    - [ ] Create `app/(authenticated)/_components/__tests__/AlbumsGrid.test.tsx`
        - [ ] Test renders correct number of album cards
        - [ ] Test responsive grid layout
        - [ ] Test passes props to AlbumCard
    - [ ] Create `app/(authenticated)/_components/__tests__/AlbumsPage.test.tsx`
        - [ ] Test renders section title
        - [ ] Test renders AlbumsGrid
    - [ ] Create `app/(authenticated)/_components/__tests__/LoadingSkeleton.test.tsx`
        - [ ] Test renders skeleton cards
    - [ ] Create `app/(authenticated)/_components/__tests__/EmptyState.test.tsx`
        - [ ] Test renders empty message
    - [ ] Create `app/(authenticated)/_components/__tests__/ErrorDisplay.test.tsx`
        - [ ] Test renders error message
        - [ ] Test "Try Again" button triggers refresh

- [ ] Run tests and verify build (AC: all)
    - [ ] Run `npm run test` - verify all tests pass
    - [ ] Run `npm run build` - verify successful build
    - [ ] Run `npm run test:visual` - verify visual tests pass (if applicable)
    - [ ] Verify no TypeScript errors
    - [ ] Verify no missing imports
    - [ ] Verify no console errors in dev mode

## Dev Notes

### Story Context

**This is the first feature story building on top of the foundation stories (1.1, 1.2).**

Story 1.1 created Material UI theme, Link wrapper, and authentication layout.  
Story 1.2 migrated catalog state management, fetch adapter, image loader, and error boundaries.

Story 1.3 implements the album list page using these foundations WITHOUT client-side state (no useReducer yet).

### Simplified Architecture (No Client State Yet)

**From Arch's Clarification:**

> "No state is required on the client yet, so no need to have client components. Developer will only need to have the server component that fetches the album
> data and pass the `CatalogViewerState` structure as parameter (using thunks and reducer as documented in the architecture), and then use pure UI components."

**Architecture Pattern:**

```typescript
// Server Component (page.tsx) - Fetches data
export default async function HomePage() {
    const albums = await fetchAlbumsAdapter()
    return <AlbumsPage albums = {albums}
    />
}

// Pure UI Component (NO hooks, NO state)
function AlbumsPage({albums}: { albums: Album[] }) {
    return <AlbumsGrid albums = {albums}
    />
}
```

**Key Points:**

- ✅ Server Component fetches data using adapter from Story 1.2
- ✅ Pass specific props to pure UI (albums array, NOT full CatalogViewerState)
- ❌ NO Client Components needed for Story 1.3
- ❌ NO useReducer, NO hooks, NO handlers
- ❌ NO passing full state to all components

### Component Props Design

**From Arch's Clarification:**

> "Pure components must only define the properties they are using: we don't want to have the full state passed around everywhere."

**Component Props Pattern:**

```typescript
// ❌ DON'T: Pass entire state to every component
function AlbumCard({state}: { state: CatalogViewerState }) { ...
}

// ✅ DO: Pass only properties needed
interface AlbumCardProps {
    albumId: string
    name: string
    start: string
    end: string
    totalCount: number
    ownerId: string
    folderName: string
    mediaIds: string[]  // First 4 media IDs for thumbnails
}

function AlbumCard(props: AlbumCardProps) {
    // Component only knows about its specific data
}
```

**Benefits:**

- Clear component boundaries
- Easy to test in isolation
- No unnecessary re-renders
- Better code maintainability

### Design Direction (Edge-to-Edge Photos)

**From `specs/designs/ux-design-direction-final.html`:**

**Key Design Decisions:**

1. **Photos ARE the Card (No Card Background):**
    - 4 photos cover entire card edge-to-edge
    - NO separate card background color
    - Only 8px gap between photos
    - Text overlay positioned over bottom photos

2. **Background Gradient:**
   ```css
   background: linear-gradient(135deg, #0a1520 0%, #12242e 50%, #0f1d28 100%)
   ```

3. **Album Card Grid:**
   ```css
   display: grid;
   grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
   gap: 32px;
   ```

4. **Photo Grid Within Card:**
   ```css
   display: grid;
   grid-template-columns: repeat(2, 1fr);
   gap: 8px;
   padding: 0;  /* Edge-to-edge */
   ```

5. **Text Overlay:**
    - Position: absolute bottom
    - Background: `rgba(10, 21, 32, 0.98)` (dark)
    - Hover: `rgba(24, 89, 134, 0.95)` (brand blue)
    - Padding: 32px 24px 24px

6. **Typography:**
    - Album name: Georgia serif, 22px, weight 300
    - Date range: Courier New monospace, 13px, uppercase
    - Media count: 13px, weight 400, cyan #6ab9de → white on hover

7. **Hover Effects:**
    - Transform: translateY(-6px)
    - Shadow: `0 12px 40px rgba(24, 89, 134, 0.4)`
    - Photos: brightness(0.85)
    - Overlay: blue gradient
    - Count: white color

### Data Flow

```typescript
// 1. Server Component fetches albums
export default async function HomePage() {
    try {
        // Use adapter from Story 1.2
        const albums: Album[] = await fetchAlbumsAdapter()

        // Handle states
        if (albums.length === 0) {
            return <EmptyState / >
        }

        return <AlbumsPage albums = {albums}
        />

    } catch (error) {
        return <ErrorDisplay error = {error}
        />
    }
}

// Add this to disable caching (dynamic data)
export const revalidate = 0

// 2. AlbumsPage renders layout
function AlbumsPage({albums}: { albums: Album[] }) {
    return (
        <Box sx = {
    {
        background: 'linear-gradient(...)'
    }
}>
    <Typography variant = "h6" > Your
    Albums < /Typography>
    < AlbumsGrid
    albums = {albums}
    />
    < /Box>
)
}

// 3. AlbumsGrid renders cards
function AlbumsGrid({albums}: { albums: Album[] }) {
    return (
        <Box sx = {
    {
        display: 'grid', gridTemplateColumns
    :
        {...
        }
    }
}>
    {
        albums.map(album => (
            <AlbumCard
                key = {album.albumId}
        albumId = {album.albumId}
        name = {album.name}
        start = {album.start}
        end = {album.end}
        totalCount = {album.totalCount}
        ownerId = {album.ownerId}
        folderName = {album.folderName}
        mediaIds = {album.mediaIds?.slice(0, 4) || []}
        />
    ))
    }
    </Box>
)
}

// 4. AlbumCard renders single card
function AlbumCard(props: AlbumCardProps) {
    const {ownerId, folderName, name, start, end, totalCount, mediaIds} = props

    return (
        <Link href = {`/owners/${ownerId}/albums/${folderName}`
}>
    <Box sx = {
    {
        position: 'relative', cursor
    :
        'pointer', transition
    :
        '0.4s'
    }
}>
    {/* 2x2 Photo Grid */
    }
    <Box sx = {
    {
        display: 'grid', gridTemplateColumns
    :
        'repeat(2, 1fr)', gap
    :
        1
    }
}>
    {
        mediaIds.map(mediaId => (
            <Box key = {mediaId}
        sx = {
        {
            position: 'relative', aspectRatio
        :
            '1'
        }
    }>
        <Image src = {mediaId}
        fill
        sizes = "170px"
        alt = ""
        style = {
        {
            objectFit: 'cover'
        }
    }
        />
        < /Box>
    ))
    }
    </Box>

    {/* Text Overlay */
    }
    <Box sx = {
    {
        position: 'absolute', bottom
    :
        0, left
    :
        0, right
    :
        0,
    ...
    }
}>
    <Typography variant = "h6"
    sx = {
    {
        fontFamily: 'Georgia, serif', fontSize
    :
        22, fontWeight
    :
        300
    }
}>
    {
        name
    }
    </Typography>
    < Box
    sx = {
    {
        display: 'flex', justifyContent
    :
        'space-between'
    }
}>
    <Typography sx = {
    {
        fontFamily: 'Courier New, monospace', fontSize
    :
        13, textTransform
    :
        'uppercase'
    }
}>
    {
        formatDateRange(start, end)
    }
    </Typography>
    < Typography
    sx = {
    {
        fontSize: 13, color
    :
        '#6ab9de'
    }
}>
    {
        totalCount
    }
    photos
    < /Typography>
    < /Box>
    < /Box>
    < /Box>
    < /Link>
)
}
```

### Material UI Responsive Grid

**MUI sx Prop Pattern:**

```typescript
<Box sx = {
{
    display: 'grid',
        gridTemplateColumns
:
    {
        xs: 'repeat(1, 1fr)',    // Mobile <600px: 1 column
            sm
    :
        'repeat(2, 1fr)',    // Tablet 600-960px: 2 columns
            md
    :
        'repeat(3, 1fr)',    // Desktop 960-1280px: 3 columns
            lg
    :
        'repeat(4, 1fr)',    // Large >1280px: 4 columns
    }
,
    gap: 4,  // 32px (theme spacing 4 = 32px)
}
}>
```

**Breakpoints (from MUI theme):**

- xs: 0px
- sm: 600px
- md: 960px
- lg: 1280px

**Spacing System:**

- gap: 4 = 32px (8 * 4)
- gap: 1 = 8px (photos within card)

### Next.js Image Component

**Image Loader Configuration (Already Done in Story 1.2):**

```typescript
// libs/image-loader.ts
export default function imageLoader({src, width}: ImageLoaderProps): string {
    const targetWidth = width <= 360 ? 360 : 2400
    return `/api/v1/media/${src}/image?width=${targetWidth}`
}

// next.config.ts
const nextConfig = {
    images: {
        loader: 'custom',
        loaderFile: './libs/image-loader.ts',
    },
}
```

**Usage in AlbumCard:**

```typescript
import Image from 'next/image'

<Box sx = {
{
    position: 'relative', aspectRatio
:
    '1'
}
}>
<Image
    src = {mediaId}                    // mediaId passed to custom loader
alt = ""                            // Empty alt for decorative images
fill                              // Fill parent container
sizes = "(max-width: 768px) 50vw, 170px"  // Responsive sizing hints
style = {
{
    objectFit: 'cover'
}
}   // Maintain aspect ratio
/>
< /Box>
```

**Critical Requirements:**

- ✅ Parent MUST have `position: relative` or `position: absolute`
- ✅ Parent MUST have defined dimensions (aspectRatio: '1' for square)
- ✅ Use `fill` prop (NOT width/height props)
- ✅ Provide `sizes` for responsive hints
- ✅ Use `style` for objectFit (NOT sx prop)

### Link Wrapper (Next.js v16 Restriction)

**From Story 1.1 - Link Component Wrapper:**

```typescript
// components/Link.tsx (Already created)
'use client'
import Link, {LinkProps} from 'next/link'

export default Link
```

**Usage in AlbumCard:**

```typescript
// ❌ DON'T: Direct Next.js Link import
import Link from 'next/link'

<Button component = {Link}
href = "/about" > Click < /Button>

// ✅ DO: Use wrapper from Story 1.1
import Link from '@/components/Link'

<Link href = {`/owners/${ownerId}/albums/${folderName}`
}>
<Box>Album
Card
Content < /Box>
< /Link>
```

**Why:** Next.js v16 throws error "Functions cannot be passed directly to Client Components" when passing Link to MUI component prop. The wrapper solves this.

### Date Formatting

**Format: "MMM DD – MMM DD, YYYY"**

```typescript
function formatDateRange(start: string, end: string): string {
    const startDate = new Date(start)
    const endDate = new Date(end)

    const formatOptions: Intl.DateTimeFormatOptions = {
        month: 'short',
        day: 'numeric',
    }

    const startFormatted = startDate.toLocaleDateString('en-US', formatOptions)
    const endFormatted = endDate.toLocaleDateString('en-US', formatOptions)
    const year = endDate.getFullYear()

    return `${startFormatted} – ${endFormatted}, ${year}`
}

// Examples:
// formatDateRange('2026-07-15', '2026-07-22') → "Jul 15 – Jul 22, 2026"
// formatDateRange('2025-12-23', '2025-12-26') → "Dec 23 – Dec 26, 2025"
```

**Note:** Use en-dash (–) NOT hyphen (-) in date range.

### Album Type Structure

**From Story 1.2 - Catalog Domain:**

```typescript
// domains/catalog/language/Album.ts
interface Album {
    albumId: string          // Unique ID
    name: string             // Album name
    start: string            // ISO date "2026-07-15"
    end: string              // ISO date "2026-07-22"
    folderName: string       // URL-safe folder name
    ownerId: string          // Owner user ID
    totalCount: number       // Total media count
    mediaIds?: string[]      // Array of media IDs (first 4 for thumbnails)
    temperature?: number     // Photos per day (Story 1.4+)
    relativeTemperature?: number  // Normalized temperature (Story 1.4+)
    owner?: OwnerDetails     // Owner info (Story 1.4+)
    sharedWith?: UserDetails[]    // Shared users (Story 1.4+)
}
```

**For Story 1.3, We Use:**

- albumId, name, start, end, folderName, ownerId, totalCount, mediaIds
- temperature, owner, sharedWith are for future stories

### Fetch Adapter Usage

**From Story 1.2 - Already Implemented:**

```typescript
// domains/catalog/adapters/fetch-adapter.ts
export async function fetchAlbumsAdapter(): Promise<Album[]> {
    const response = await fetch('/api/v1/albums', {
        credentials: 'include',  // Session cookies for auth
    })

    if (!response.ok) {
        throw new CatalogError(`Failed to fetch albums: ${response.status}`)
    }

    const data = await response.json()
    return data.albums  // Already sorted by start date descending
}
```

**Usage in Server Component:**

```typescript
import {fetchAlbumsAdapter} from '@/domains/catalog/adapters/fetch-adapter'

export default async function HomePage() {
    try {
        const albums = await fetchAlbumsAdapter()
        return <AlbumsPage albums = {albums}
        />
    } catch (error) {
        return <ErrorDisplay error = {error}
        />
    }
}

// CRITICAL: Add this to disable Next.js caching
export const revalidate = 0
```

### Error Boundaries

**Already Created in Story 1.2:**

- ✅ `app/error.tsx` - Root catch-all
- ✅ `app/(authenticated)/error.tsx` - Authenticated routes
- ✅ `app/not-found.tsx` - 404 page

**For Story 1.3 - Component-Level Errors:**

Use try-catch in Server Component:

```typescript
export default async function HomePage() {
    try {
        const albums = await fetchAlbumsAdapter()
        // ... success path
    } catch (error) {
        return <ErrorDisplay error = {error}
        />
    }
}
```

**ErrorDisplay Component:**

```typescript
function ErrorDisplay({error}: { error: Error | unknown }) {
    const message = error instanceof Error ? error.message : 'Unknown error'

    return (
        <Paper sx = {
    {
        padding: 4, textAlign
    :
        'center'
    }
}>
    <Typography variant = "h6"
    color = "error" >
        Failed
    to
    load
    albums
    < /Typography>
    < Typography
    variant = "body2"
    color = "text.secondary"
    sx = {
    {
        mt: 2
    }
}>
    {
        message
    }
    </Typography>
    < Button
    variant = "contained"
    sx = {
    {
        mt: 3, bgcolor
    :
        '#185986'
    }
}
    onClick = {()
=>
    window.location.reload()
}
>
    Try
    Again
    < /Button>
    < /Paper>
)
}
```

### Testing Strategy

**Component Tests (Vitest + React Testing Library):**

```typescript
// AlbumCard.test.tsx
import {render, screen} from '@testing-library/react'
import AlbumCard from '../AlbumCard'

describe('AlbumCard', () => {
    const mockProps = {
        albumId: 'album-1',
        name: 'Beach Trip',
        start: '2026-07-15',
        end: '2026-07-22',
        totalCount: 47,
        ownerId: 'owner-1',
        folderName: 'beach-trip',
        mediaIds: ['media-1', 'media-2', 'media-3', 'media-4'],
    }

    it('renders album name', () => {
        render(<AlbumCard {...mockProps}
        />)
        expect(screen.getByText('Beach Trip')).toBeInTheDocument()
    })

    it('renders date range correctly', () => {
        render(<AlbumCard {...mockProps}
        />)
        expect(screen.getByText(/Jul 15 – Jul 22, 2026/i)).toBeInTheDocument()
    })

    it('renders media count', () => {
        render(<AlbumCard {...mockProps}
        />)
        expect(screen.getByText('47 photos')).toBeInTheDocument()
    })

    it('renders 4 photo thumbnails', () => {
        render(<AlbumCard {...mockProps}
        />)
        const images = screen.getAllByRole('img')
        expect(images).toHaveLength(4)
    })

    it('links to correct album page', () => {
        render(<AlbumCard {...mockProps}
        />)
        const link = screen.getByRole('link')
        expect(link).toHaveAttribute('href', '/owners/owner-1/albums/beach-trip')
    })
})
```

**Run Tests:**

```bash
npm run test              # All unit tests
npm run build             # Verify build
npm run test:visual       # Visual tests (if applicable)
```

### File Structure

```
app/
├── (authenticated)/
│   ├── page.tsx                    # 🆕 Server Component - home page
│   └── _components/                # 🆕 Colocated components
│       ├── AlbumsPage.tsx          # Layout wrapper
│       ├── AlbumsGrid.tsx          # Responsive grid
│       ├── AlbumCard.tsx           # Individual card
│       ├── LoadingSkeleton.tsx     # Loading state
│       ├── EmptyState.tsx          # Empty state
│       ├── ErrorDisplay.tsx        # Error state
│       └── __tests__/              # Component tests
│           ├── AlbumCard.test.tsx
│           ├── AlbumsGrid.test.tsx
│           ├── AlbumsPage.test.tsx
│           ├── LoadingSkeleton.test.tsx
│           ├── EmptyState.test.tsx
│           └── ErrorDisplay.test.tsx
│
├── error.tsx                       # ✅ Already exists (Story 1.2)
├── not-found.tsx                   # ✅ Already exists (Story 1.2)
└── layout.tsx                      # ✅ Already exists (Story 1.1)

components/
├── Link.tsx                        # ✅ Already exists (Story 1.1)
└── theme/                          # ✅ Already exists (Story 1.1)
    └── theme.ts

domains/
└── catalog/                        # ✅ Already migrated (Story 1.2)
    ├── language/
    │   └── Album.ts
    └── adapters/
        └── fetch-adapter.ts

libs/
├── image-loader.ts                 # ✅ Already exists (Story 1.2)
└── daction/                        # ✅ Already migrated (Story 1.2)
```

### Previous Story Context

**Story 1.1 - Project Foundation:**

- ✅ Material UI 6.x installed with dark theme
- ✅ Brand color #185986 configured
- ✅ Link wrapper created for Next.js v16 compatibility
- ✅ ThemeProvider set up in root layout
- ✅ AppRouterCacheProvider configured
- ✅ All tests passing (39)

**Story 1.2 - State Management Migration:**

- ✅ Catalog domain migrated (221+ tests)
- ✅ Fetch adapter created (fetchAlbumsAdapter)
- ✅ Image loader configured (360/2400 width mapping)
- ✅ Error boundaries created
- ✅ daction library migrated
- ✅ All tests passing, build succeeds

**Story 1.3 - Building On Top:**

- 🆕 First feature story (album list display)
- 🆕 Server Component only (no client state)
- 🆕 Pure UI components with specific props
- 🆕 Edge-to-edge photo design
- 🆕 Responsive grid layout

### Success Validation

**Before Raising PR:**

1. ✅ All tests pass: `npm run test`
2. ✅ Build succeeds: `npm run build`
3. ✅ Visual tests pass: `npm run test:visual` (if applicable)
4. ✅ No TypeScript errors
5. ✅ No missing imports
6. ✅ Page loads and displays albums correctly
7. ✅ Responsive grid works on mobile/tablet/desktop
8. ✅ Hover effects work smoothly
9. ✅ Clicking album navigates correctly
10. ✅ Loading/empty/error states display correctly
11. ✅ Images load progressively using custom loader
12. ✅ Date formatting is correct
13. ✅ All acceptance criteria met

### Critical Guardrails

❌ **DO NOT:**

1. Create Client Components for Story 1.3 (no 'use client')
2. Use useReducer or hooks (Server Component only)
3. Pass full CatalogViewerState to components
4. Use inline styles or Tailwind classes
5. Import Next.js Link directly (use @/components/Link)
6. Forget `export const revalidate = 0` (caching)
7. Use width/height props with Image fill
8. Forget position: relative on Image parent
9. Skip try-catch error handling
10. Reimplement state management

✅ **DO:**

1. Use async Server Component in page.tsx
2. Fetch data using adapter from Story 1.2
3. Pass specific props to pure UI components
4. Use MUI sx prop for all styling
5. Use @/components/Link wrapper for navigation
6. Add `export const revalidate = 0` to page
7. Use Image fill with sizes prop
8. Provide position: relative on parent
9. Handle loading, error, empty states
10. Follow edge-to-edge photo design from HTML mockup
11. Use Georgia serif for album names
12. Use Courier monospace for dates
13. Format dates as "MMM DD – MMM DD, YYYY"
14. Apply hover effects via sx prop
15. Test all components thoroughly

### Anti-Patterns

**From AGENTS.md - Testing Strategy:**

❌ **Avoid:**

- Excessive comments that paraphrase code
- Tight coupling between tests and implementation
- Testing implementation details vs. behavior
- Skipping edge cases (empty, error, loading)

✅ **Follow:**

- Test behavior, not implementation
- Clear test descriptions
- Mock external dependencies (fetch adapter)
- Test all component states
- Maintain low coupling with code structure

### Commit Message Pattern

**From Git History:**

```
catalog/web - implement basic album list display with responsive grid

- Create Server Component to fetch albums
- Build pure UI components (AlbumsPage, AlbumsGrid, AlbumCard)
- Implement edge-to-edge photo design with text overlay
- Add responsive grid layout (1/2/3/4 columns)
- Handle loading, error, empty states
- Use Next.js Image with custom loader
- Format dates as "MMM DD – MMM DD, YYYY"
- Apply hover effects (lift, shadow, blue overlay)
- All tests passing (X), build succeeds
```

**Pattern:** `catalog/web - <action description>`

## Dev Agent Record

### Agent Model Used

_To be filled by dev agent_

### Debug Log References

_To be filled by dev agent_

### Completion Notes List

_To be filled by dev agent_

### File List

_To be filled by dev agent_

## Change Log

_To be filled by dev agent_
