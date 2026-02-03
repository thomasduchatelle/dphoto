# Story 1.5: Album Filtering

Status: ready-for-dev

## Story

As a user,  
I want to filter albums by owner,  
So that I can focus on my own albums or view all albums including shared ones.

## Acceptance Criteria

**Given** I am viewing the album list with multiple albums from different owners  
**When** I use the owner filter  
**Then** a filter control is displayed above the album grid with options:

- "All Albums" - shows all accessible albums (owned + shared)
- "My Albums" - shows only albums I own
- Specific owner names - shows albums owned by that specific owner (if albums are shared with me)

**And** the filter uses MUI ToggleButtonGroup component with brand blue (#185986) for selected state  
**And** the active filter state is visually indicated (FR41)  
**And** filtering is handled client-side using pure filter function (NO reducer, NO state management)  
**And** the filter updates immediately without API call  
**And** the filter state persists during navigation within the session using URL search params  
**And** "My Albums" filter shows only albums where I am the owner (FR1)  
**And** "All Albums" shows owned albums (FR1) and shared albums (FR2)  
**And** specific owner filter shows albums owned by that owner only (FR3)  
**And** the filter control is responsive and works on mobile, tablet, and desktop  
**And** keyboard navigation works for the filter control (tab, arrow keys, enter)  
**And** the filter state is preserved when returning to the page (URL-based state)  
**And** changing filter maintains scroll position in the album list  
**And** the component remains a pure UI component with props, NO useReducer from Story 1.2

## Tasks / Subtasks

- [ ] Create filter helper function (AC: client-side filtering, pure function)
    - [ ] Create `app/(authenticated)/_components/album-filters.ts`
    - [ ] Create `filterAlbumsByOwner(albums, filterType, currentUserId): Album[]`
    - [ ] FilterType: 'all' | 'mine' | string (ownerId)
    - [ ] 'all': return all albums
    - [ ] 'mine': return albums where album.ownerId === currentUserId
    - [ ] ownerId: return albums where album.ownerId === ownerId
    - [ ] Pure function, NO state management
    - [ ] Export filter types and function

- [ ] Create filter component (AC: MUI ToggleButtonGroup, responsive)
    - [ ] Create `app/(authenticated)/_components/AlbumFilterBar.tsx`
    - [ ] Accept props: currentFilter (string), onFilterChange (callback), ownerOptions (OwnerOption[])
    - [ ] OwnerOption interface: { ownerId: string, ownerName: string, albumCount: number }
    - [ ] Use MUI ToggleButtonGroup with value={currentFilter}
    - [ ] onChange handler: call onFilterChange(newValue)
    - [ ] Buttons: "All Albums", "My Albums", ...owner names
    - [ ] Selected state: brand blue #185986 background
    - [ ] Responsive layout:
        - Mobile (xs): Stack vertically
        - Desktop (sm+): Horizontal row
    - [ ] Use MUI sx prop for styling (NO inline styles)
    - [ ] Keyboard accessible (tab, arrow keys, enter)

- [ ] Update page to support filtering (AC: URL search params, persist filter)
    - [ ] Update `app/(authenticated)/page.tsx` (Server Component)
    - [ ] Read searchParams.filter (defaults to 'all')
    - [ ] Pass filter to client component as prop
    - [ ] Server Component remains simple: fetch + pass props

- [ ] Create client wrapper for filtering (AC: URL-based state)
    - [ ] Create `app/(authenticated)/_components/AlbumsPageClient.tsx`
    - [ ] Mark as 'use client'
    - [ ] Accept props: albums (Album[]), currentUserId (string), initialFilter (string)
    - [ ] Use useRouter and useSearchParams from next/navigation
    - [ ] Calculate ownerOptions from albums (extract unique owners with counts)
    - [ ] Filter albums using filterAlbumsByOwner(albums, initialFilter, currentUserId)
    - [ ] onFilterChange handler: router.push(`/?filter=${newValue}`)
    - [ ] Pass filtered albums to AlbumsPage component
    - [ ] Pass currentFilter and onFilterChange to AlbumFilterBar
    - [ ] NO useReducer, NO complex state management

- [ ] Update AlbumsPage layout (AC: filter above grid)
    - [ ] Update `app/(authenticated)/_components/AlbumsPage.tsx`
    - [ ] Accept filterBar prop (ReactNode)
    - [ ] Render filterBar above section title
    - [ ] Keep existing layout (background gradient, title, grid)
    - [ ] Maintain spacing with MUI sx prop

- [ ] Calculate owner options from albums (AC: specific owner names)
    - [ ] Create helper function `getOwnerOptions(albums, currentUserId): OwnerOption[]`
    - [ ] Extract unique owner IDs from albums (excluding current user)
    - [ ] Map to { ownerId, ownerName, albumCount }
    - [ ] Count albums per owner
    - [ ] Sort by album count descending
    - [ ] Return array for filter buttons

- [ ] Add TypeScript interfaces (AC: all)
    - [ ] FilterType: 'all' | 'mine' | string
    - [ ] OwnerOption: { ownerId: string, ownerName: string, albumCount: number }
    - [ ] AlbumFilterBarProps: { currentFilter: string, onFilterChange: (filter: string) => void, ownerOptions: OwnerOption[] }
    - [ ] AlbumsPageClientProps: { albums: Album[], currentUserId: string, initialFilter: string }
    - [ ] Update AlbumsPageProps to include filterBar?: ReactNode

- [ ] Write component tests (AC: all tests pass)
    - [ ] Create `app/(authenticated)/_components/__tests__/album-filters.test.ts`
        - [ ] Test filterAlbumsByOwner with 'all' - returns all albums
        - [ ] Test filterAlbumsByOwner with 'mine' - returns only owned albums
        - [ ] Test filterAlbumsByOwner with ownerId - returns albums by specific owner
        - [ ] Test with empty albums array
        - [ ] Test with no matching albums
    - [ ] Create `app/(authenticated)/_components/__tests__/AlbumFilterBar.test.tsx`
        - [ ] Test renders "All Albums" and "My Albums" buttons
        - [ ] Test renders owner name buttons
        - [ ] Test selected state styling
        - [ ] Test onChange callback when button clicked
        - [ ] Test keyboard navigation (tab, enter)
        - [ ] Test responsive layout (mobile/desktop)
    - [ ] Create `app/(authenticated)/_components/__tests__/AlbumsPageClient.test.tsx`
        - [ ] Test filters albums correctly on mount
        - [ ] Test calls router.push when filter changes
        - [ ] Test calculates owner options correctly
        - [ ] Test passes props to child components
    - [ ] Update `app/(authenticated)/_components/__tests__/AlbumsPage.test.tsx`
        - [ ] Test renders filter bar when provided
        - [ ] Test maintains existing layout

- [ ] Run tests and verify build (AC: all)
    - [ ] Run `npm run test` - verify all tests pass
    - [ ] Run `npm run build` - verify successful build
    - [ ] Run `npm run test:visual` - verify visual tests pass (if applicable)
    - [ ] Verify no TypeScript errors
    - [ ] Verify no missing imports
    - [ ] Verify filter buttons appear and work correctly
    - [ ] Verify URL updates when filter changes
    - [ ] Verify filter persists on page refresh
    - [ ] Verify scroll position maintained

## Dev Notes

### Story Context

**This is the FINAL story in Epic 1 - Album List Home Page.**

Story 1.5 adds client-side filtering to the album list, allowing users to filter by owner. This is the first story requiring Client Components, but we keep it
SIMPLE - no useReducer from Story 1.2, just URL-based filter state.

**Stories 1.1-1.4 provide the foundation:**

- Story 1.1: Material UI theme with brand color #185986
- Story 1.2: Catalog domain (unused in this story), fetch adapter, image loader
- Story 1.3: Basic album list with Server Component + pure UI (edge-to-edge photos, text overlay)
- Story 1.4: Album card enhancements (density indicators, sharing avatars)

**Story 1.5 builds on Story 1.3-1.4 by adding filtering WITHOUT complex state management.**

### Simplified Architecture (URL-Based State)

**Key Decision: NO useReducer for Simple Filtering**

From Arch's guidance: Story 1.2 migrated the full catalog state management with useReducer, actions, thunks. However, for Story 1.5, we keep it simple because:

1. Filtering is a simple operation (filter array client-side)
2. Filter state fits naturally in URL search params
3. No complex mutations or async operations
4. Don't need the reducer until Epic 3 (Album Management)

**Architecture Pattern:**

```typescript
// Server Component (page.tsx) - Read filter from URL
export default async function HomePage({ searchParams }: { searchParams: { filter?: string } }) {
    const albums = await fetchAlbumsAdapter()
    const currentUserId = await getCurrentUserId()  // From auth context
    const filter = searchParams.filter || 'all'
    
    return <AlbumsPageClient 
        albums={albums} 
        currentUserId={currentUserId} 
        initialFilter={filter} 
    />
}

// Client Component - Handle filter changes
'use client'
function AlbumsPageClient({ albums, currentUserId, initialFilter }: Props) {
    const router = useRouter()
    
    // Calculate owner options from albums
    const ownerOptions = getOwnerOptions(albums, currentUserId)
    
    // Filter albums client-side
    const filteredAlbums = filterAlbumsByOwner(albums, initialFilter, currentUserId)
    
    // Handle filter change - update URL
    const handleFilterChange = (newFilter: string) => {
        router.push(`/?filter=${newFilter}`)
    }
    
    return (
        <AlbumsPage 
            albums={filteredAlbums}
            filterBar={
                <AlbumFilterBar
                    currentFilter={initialFilter}
                    onFilterChange={handleFilterChange}
                    ownerOptions={ownerOptions}
                />
            }
        />
    )
}

// Pure UI Components (NO changes from Story 1.3-1.4)
function AlbumsPage({ albums, filterBar }: Props) {
    return (
        <Box>
            {filterBar}
            <Typography>Your Albums</Typography>
            <AlbumsGrid albums={albums} />
        </Box>
    )
}
```

**Key Points:**

- ✅ Server Component reads filter from URL search params
- ✅ Client Component handles filter state via URL
- ✅ Pure filter function (NO reducer, NO actions)
- ✅ URL-based state (persists across navigation)
- ✅ AlbumsPage, AlbumsGrid, AlbumCard remain pure UI (NO changes)
- ❌ NO useReducer in this story
- ❌ NO catalog actions/thunks
- ❌ NO client state management complexity

### Filter Function Design

**Pure Function - No State Management:**

```typescript
// album-filters.ts
export type FilterType = 'all' | 'mine' | string  // string is ownerId

export function filterAlbumsByOwner(
    albums: Album[], 
    filterType: FilterType, 
    currentUserId: string
): Album[] {
    if (filterType === 'all') {
        return albums  // All albums (owned + shared)
    }
    
    if (filterType === 'mine') {
        return albums.filter(album => album.ownerId === currentUserId)
    }
    
    // Specific owner ID
    return albums.filter(album => album.ownerId === filterType)
}

export interface OwnerOption {
    ownerId: string
    ownerName: string
    albumCount: number
}

export function getOwnerOptions(albums: Album[], currentUserId: string): OwnerOption[] {
    // Group albums by owner (excluding current user)
    const ownerMap = new Map<string, { name: string, count: number }>()
    
    albums.forEach(album => {
        if (album.ownerId !== currentUserId && album.owner) {
            const existing = ownerMap.get(album.ownerId) || { name: album.owner.name, count: 0 }
            ownerMap.set(album.ownerId, { ...existing, count: existing.count + 1 })
        }
    })
    
    // Convert to array and sort by album count descending
    return Array.from(ownerMap.entries())
        .map(([ownerId, data]) => ({
            ownerId,
            ownerName: data.name,
            albumCount: data.count,
        }))
        .sort((a, b) => b.albumCount - a.albumCount)
}
```

**Examples:**

```typescript
// User owns 10 albums, has 5 shared from Alice, 3 from Bob
const albums = [...30 albums...]

getOwnerOptions(albums, 'user-1')
// Returns: [
//   { ownerId: 'alice-id', ownerName: 'Alice', albumCount: 5 },
//   { ownerId: 'bob-id', ownerName: 'Bob', albumCount: 3 }
// ]

filterAlbumsByOwner(albums, 'all', 'user-1')  
// Returns: all 18 albums (10 owned + 8 shared)

filterAlbumsByOwner(albums, 'mine', 'user-1')  
// Returns: 10 owned albums

filterAlbumsByOwner(albums, 'alice-id', 'user-1')  
// Returns: 5 albums owned by Alice
```

### Filter Component Design

**MUI ToggleButtonGroup Pattern:**

```typescript
// AlbumFilterBar.tsx
import { ToggleButtonGroup, ToggleButton, Box } from '@mui/material'

interface AlbumFilterBarProps {
    currentFilter: string
    onFilterChange: (filter: string) => void
    ownerOptions: OwnerOption[]
}

export default function AlbumFilterBar({ 
    currentFilter, 
    onFilterChange, 
    ownerOptions 
}: AlbumFilterBarProps) {
    return (
        <Box sx={{ mb: 3 }}>
            <ToggleButtonGroup
                value={currentFilter}
                exclusive
                onChange={(e, value) => {
                    if (value !== null) {
                        onFilterChange(value)
                    }
                }}
                sx={{
                    flexWrap: 'wrap',
                    gap: 1,
                    '& .MuiToggleButton-root': {
                        border: '1px solid rgba(255, 255, 255, 0.12)',
                        color: 'rgba(255, 255, 255, 0.7)',
                        textTransform: 'none',
                        px: 2,
                        py: 1,
                        '&.Mui-selected': {
                            bgcolor: '#185986',  // Brand blue
                            color: '#ffffff',
                            '&:hover': {
                                bgcolor: '#1a6a9e',  // Slightly lighter on hover
                            },
                        },
                    },
                }}
            >
                <ToggleButton value="all">
                    All Albums
                </ToggleButton>
                
                <ToggleButton value="mine">
                    My Albums
                </ToggleButton>
                
                {ownerOptions.map(owner => (
                    <ToggleButton key={owner.ownerId} value={owner.ownerId}>
                        {owner.ownerName} ({owner.albumCount})
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        </Box>
    )
}
```

**Responsive Design:**

- Mobile (xs): ToggleButtonGroup wraps naturally with flexWrap
- Desktop (sm+): Buttons appear in a row
- Keyboard: Tab to navigate, Enter to select

**Visual Design:**

- Default: Dark border, light text
- Selected: Brand blue #185986 background, white text
- Hover: Slightly lighter blue (#1a6a9e)
- Shows album count per owner: "Alice (5)"

### URL-Based State Pattern

**Why URL Search Params?**

1. State persists across page refresh
2. Shareable URLs (e.g., `/?filter=mine`)
3. Back/forward navigation works naturally
4. Simple to implement (no complex state management)
5. Server Component can read initial state

**Implementation:**

```typescript
// Server Component - Read from URL
export default async function HomePage({ 
    searchParams 
}: { 
    searchParams: { filter?: string } 
}) {
    const albums = await fetchAlbumsAdapter()
    const currentUserId = await getCurrentUserId()
    const filter = searchParams.filter || 'all'
    
    return <AlbumsPageClient 
        albums={albums} 
        currentUserId={currentUserId} 
        initialFilter={filter} 
    />
}

// Client Component - Update URL on change
'use client'
function AlbumsPageClient({ albums, currentUserId, initialFilter }: Props) {
    const router = useRouter()
    
    const handleFilterChange = (newFilter: string) => {
        router.push(`/?filter=${newFilter}`)
        // Next.js will re-render with new searchParams
    }
    
    // ... rest of component
}
```

**Benefits:**

- Filter state in URL: `/?filter=mine`
- Page refresh maintains filter
- Back button works (browser history)
- No need for localStorage or cookies

### Current User ID

**How to Get Current User ID:**

From Story 1.3, we already have authentication context. The user ID is available from the existing auth backend (Google OAuth).

**Two Options:**

**Option 1: From API Response (Recommended)**

```typescript
// Server Component
export default async function HomePage({ searchParams }: Props) {
    const albums = await fetchAlbumsAdapter()
    
    // First album's owner ID if owned, otherwise check for current user endpoint
    // Or add currentUserId to API response
    const currentUserId = albums.find(a => a.ownerId)?.ownerId || 'unknown'
    
    // Better: API returns current user info
    const currentUser = await fetch('/api/v1/me').then(r => r.json())
    const currentUserId = currentUser.userId
    
    return <AlbumsPageClient albums={albums} currentUserId={currentUserId} initialFilter={filter} />
}
```

**Option 2: From Album Data (Simpler)**

```typescript
// Assume first owned album indicates current user
// Or derive from owner vs. sharedWith fields
function getCurrentUserId(albums: Album[]): string {
    // Find first album where user is owner (ownerId matches)
    // This works because albums include both owned and shared
    // Owned albums will have ownerId = current user
    const ownedAlbum = albums.find(a => a.owner)
    return ownedAlbum?.ownerId || 'unknown'
}
```

**For Story 1.5, use Option 2 (simpler, no new API call).**

### Component Structure

**File Layout:**

```
app/
├── (authenticated)/
│   ├── page.tsx                    # UPDATE: Read searchParams, pass to client
│   └── _components/
│       ├── AlbumsPageClient.tsx    # 🆕 Client wrapper for filtering
│       ├── AlbumFilterBar.tsx      # 🆕 Filter UI component
│       ├── album-filters.ts        # 🆕 Pure filter functions
│       ├── AlbumsPage.tsx          # UPDATE: Accept filterBar prop
│       ├── AlbumsGrid.tsx          # ✅ No changes (Story 1.3)
│       ├── AlbumCard.tsx           # ✅ No changes (Story 1.3-1.4)
│       ├── DensityIndicator.tsx    # ✅ No changes (Story 1.4)
│       ├── SharingAvatars.tsx      # ✅ No changes (Story 1.4)
│       └── __tests__/
│           ├── album-filters.test.ts         # 🆕
│           ├── AlbumFilterBar.test.tsx       # 🆕
│           ├── AlbumsPageClient.test.tsx     # 🆕
│           └── AlbumsPage.test.tsx           # UPDATE
```

**Changes Summary:**

- 🆕 NEW: AlbumsPageClient, AlbumFilterBar, album-filters
- UPDATE: page.tsx (read searchParams), AlbumsPage (accept filterBar)
- ✅ NO CHANGES: AlbumsGrid, AlbumCard, DensityIndicator, SharingAvatars

### Data Flow

```typescript
// 1. Server Component reads filter from URL
export default async function HomePage({ searchParams }: Props) {
    const albums = await fetchAlbumsAdapter()  // All albums (owned + shared)
    const filter = searchParams.filter || 'all'
    
    // Derive current user ID from albums
    const currentUserId = albums.find(a => a.owner)?.ownerId || 'unknown'
    
    return <AlbumsPageClient 
        albums={albums} 
        currentUserId={currentUserId} 
        initialFilter={filter} 
    />
}

// 2. Client Component handles filtering
'use client'
function AlbumsPageClient({ albums, currentUserId, initialFilter }: Props) {
    const router = useRouter()
    
    // Calculate owner options (Alice, Bob, etc.)
    const ownerOptions = getOwnerOptions(albums, currentUserId)
    
    // Filter albums client-side
    const filteredAlbums = filterAlbumsByOwner(albums, initialFilter, currentUserId)
    
    // Handle filter change - update URL
    const handleFilterChange = (newFilter: string) => {
        router.push(`/?filter=${newFilter}`)
    }
    
    return (
        <AlbumsPage 
            albums={filteredAlbums}
            filterBar={
                <AlbumFilterBar
                    currentFilter={initialFilter}
                    onFilterChange={handleFilterChange}
                    ownerOptions={ownerOptions}
                />
            }
        />
    )
}

// 3. AlbumsPage renders layout (updated to accept filterBar)
function AlbumsPage({ albums, filterBar }: Props) {
    return (
        <Box sx={{ 
            background: 'linear-gradient(135deg, #0a1520 0%, #12242e 50%, #0f1d28 100%)',
            minHeight: '100vh',
            p: 4,
        }}>
            {filterBar}  {/* NEW: Render filter above title */}
            
            <Typography variant="h6" sx={{ mb: 3 }}>
                Your Albums
            </Typography>
            
            <AlbumsGrid albums={albums} />
        </Box>
    )
}

// 4. AlbumsGrid and AlbumCard unchanged (Story 1.3-1.4)
```

### Material UI Components

**New Components for Story 1.5:**

- `ToggleButtonGroup` - from @mui/material
- `ToggleButton` - from @mui/material
- `Box` - already used in Story 1.3-1.4

**Already Used:**

- `Box`, `Typography`, `Paper`, `Button` (Story 1.3)
- `AvatarGroup`, `Avatar`, `Tooltip` (Story 1.4)

### Responsive Design

**Filter Bar Layout:**

- Mobile (xs <600px): Buttons wrap vertically with flexWrap
- Tablet (sm 600-960px): Buttons in horizontal row
- Desktop (md+ >960px): Buttons in horizontal row

**Spacing:**

- Filter bar margin bottom: 24px (mb: 3)
- Button gap: 8px (gap: 1)
- Button padding: 16px horizontal, 8px vertical (px: 2, py: 1)

**Album Grid (Unchanged from Story 1.3):**

- Mobile (xs): 1 column
- Tablet (sm): 2 columns
- Desktop (md): 3 columns
- Large (lg): 4 columns

### Testing Strategy

**Component Tests (Vitest + React Testing Library):**

```typescript
// album-filters.test.ts
describe('filterAlbumsByOwner', () => {
    const mockAlbums = [
        { albumId: '1', ownerId: 'user-1', name: 'My Album' },
        { albumId: '2', ownerId: 'alice', name: 'Alice Album' },
        { albumId: '3', ownerId: 'bob', name: 'Bob Album' },
        { albumId: '4', ownerId: 'user-1', name: 'My Album 2' },
    ]
    
    it('returns all albums when filter is "all"', () => {
        const result = filterAlbumsByOwner(mockAlbums, 'all', 'user-1')
        expect(result).toHaveLength(4)
    })
    
    it('returns only owned albums when filter is "mine"', () => {
        const result = filterAlbumsByOwner(mockAlbums, 'mine', 'user-1')
        expect(result).toHaveLength(2)
        expect(result[0].albumId).toBe('1')
        expect(result[1].albumId).toBe('4')
    })
    
    it('returns albums by specific owner', () => {
        const result = filterAlbumsByOwner(mockAlbums, 'alice', 'user-1')
        expect(result).toHaveLength(1)
        expect(result[0].albumId).toBe('2')
    })
    
    it('returns empty array when no albums match', () => {
        const result = filterAlbumsByOwner(mockAlbums, 'charlie', 'user-1')
        expect(result).toHaveLength(0)
    })
})

describe('getOwnerOptions', () => {
    it('returns unique owners excluding current user', () => {
        const mockAlbums = [
            { ownerId: 'user-1', owner: { name: 'Me' } },
            { ownerId: 'alice', owner: { name: 'Alice' } },
            { ownerId: 'alice', owner: { name: 'Alice' } },  // Duplicate
            { ownerId: 'bob', owner: { name: 'Bob' } },
        ]
        
        const result = getOwnerOptions(mockAlbums, 'user-1')
        
        expect(result).toHaveLength(2)
        expect(result[0]).toEqual({ ownerId: 'alice', ownerName: 'Alice', albumCount: 2 })
        expect(result[1]).toEqual({ ownerId: 'bob', ownerName: 'Bob', albumCount: 1 })
    })
    
    it('sorts by album count descending', () => {
        const mockAlbums = [
            { ownerId: 'alice', owner: { name: 'Alice' } },  // 1 album
            { ownerId: 'bob', owner: { name: 'Bob' } },      // 3 albums
            { ownerId: 'bob', owner: { name: 'Bob' } },
            { ownerId: 'bob', owner: { name: 'Bob' } },
        ]
        
        const result = getOwnerOptions(mockAlbums, 'user-1')
        
        expect(result[0].ownerId).toBe('bob')  // 3 albums first
        expect(result[1].ownerId).toBe('alice')  // 1 album second
    })
})

// AlbumFilterBar.test.tsx
describe('AlbumFilterBar', () => {
    const mockOwnerOptions = [
        { ownerId: 'alice', ownerName: 'Alice', albumCount: 5 },
        { ownerId: 'bob', ownerName: 'Bob', albumCount: 3 },
    ]
    
    it('renders "All Albums" and "My Albums" buttons', () => {
        render(
            <AlbumFilterBar 
                currentFilter="all" 
                onFilterChange={jest.fn()} 
                ownerOptions={[]} 
            />
        )
        
        expect(screen.getByText('All Albums')).toBeInTheDocument()
        expect(screen.getByText('My Albums')).toBeInTheDocument()
    })
    
    it('renders owner name buttons with album counts', () => {
        render(
            <AlbumFilterBar 
                currentFilter="all" 
                onFilterChange={jest.fn()} 
                ownerOptions={mockOwnerOptions} 
            />
        )
        
        expect(screen.getByText('Alice (5)')).toBeInTheDocument()
        expect(screen.getByText('Bob (3)')).toBeInTheDocument()
    })
    
    it('calls onFilterChange when button clicked', () => {
        const handleChange = jest.fn()
        render(
            <AlbumFilterBar 
                currentFilter="all" 
                onFilterChange={handleChange} 
                ownerOptions={mockOwnerOptions} 
            />
        )
        
        fireEvent.click(screen.getByText('My Albums'))
        expect(handleChange).toHaveBeenCalledWith('mine')
    })
    
    it('shows selected state styling for current filter', () => {
        const { container } = render(
            <AlbumFilterBar 
                currentFilter="mine" 
                onFilterChange={jest.fn()} 
                ownerOptions={[]} 
            />
        )
        
        const myAlbumsButton = screen.getByText('My Albums').closest('button')
        expect(myAlbumsButton).toHaveClass('Mui-selected')
    })
})

// AlbumsPageClient.test.tsx
describe('AlbumsPageClient', () => {
    const mockAlbums = [
        { albumId: '1', ownerId: 'user-1', name: 'My Album' },
        { albumId: '2', ownerId: 'alice', name: 'Alice Album' },
    ]
    
    it('filters albums based on initialFilter', () => {
        render(
            <AlbumsPageClient 
                albums={mockAlbums} 
                currentUserId="user-1" 
                initialFilter="mine" 
            />
        )
        
        // Should only show "My Album"
        expect(screen.getByText('My Album')).toBeInTheDocument()
        expect(screen.queryByText('Alice Album')).not.toBeInTheDocument()
    })
    
    it('calls router.push when filter changes', () => {
        const mockPush = jest.fn()
        jest.spyOn(require('next/navigation'), 'useRouter').mockReturnValue({ push: mockPush })
        
        render(
            <AlbumsPageClient 
                albums={mockAlbums} 
                currentUserId="user-1" 
                initialFilter="all" 
            />
        )
        
        fireEvent.click(screen.getByText('My Albums'))
        expect(mockPush).toHaveBeenCalledWith('/?filter=mine')
    })
})
```

**Run Tests:**

```bash
npm run test              # All unit tests
npm run build             # Verify build
npm run test:visual       # Visual tests (if applicable)
```

### Previous Story Patterns

**From Story 1.3 - Server Component + Pure UI:**

- ✅ Server Component fetches data
- ✅ Pure UI components receive props
- ✅ No client state in pure components

**From Story 1.4 - Component Enhancement:**

- ✅ Enhance existing components, don't recreate
- ✅ Maintain existing styling and behavior
- ✅ Add new features without breaking old ones

**Story 1.5 Follows Same Pattern:**

- ✅ Server Component reads URL, fetches data
- ✅ Client wrapper handles filter logic
- ✅ Pure UI components unchanged
- ✅ Simple filter function, no complex state

### Success Validation

**Before Raising PR:**

1. ✅ All tests pass: `npm run test`
2. ✅ Build succeeds: `npm run build`
3. ✅ Visual tests pass: `npm run test:visual` (if applicable)
4. ✅ No TypeScript errors
5. ✅ No missing imports
6. ✅ Filter buttons appear above album grid
7. ✅ "All Albums" shows all albums (owned + shared)
8. ✅ "My Albums" shows only owned albums
9. ✅ Owner buttons show albums by that owner
10. ✅ Selected filter has brand blue background
11. ✅ Filter state persists in URL (refresh maintains filter)
12. ✅ Changing filter updates URL
13. ✅ Keyboard navigation works (tab, enter)
14. ✅ Responsive layout works on mobile/tablet/desktop
15. ✅ All acceptance criteria met

### Critical Guardrails

❌ **DO NOT:**

1. Use useReducer from Story 1.2 (not needed for simple filtering)
2. Create complex state management (just URL params)
3. Change AlbumsGrid or AlbumCard components
4. Break existing hover effects or card styling
5. Add unnecessary client state (keep it simple)
6. Forget to read searchParams in Server Component
7. Forget to make AlbumsPageClient a Client Component
8. Use inline styles (MUI sx prop only)
9. Modify fetch adapter or domain logic
10. Reimplement filtering logic (use pure function)

✅ **DO:**

1. Use simple URL-based filter state
2. Create pure filter function (no side effects)
3. Use MUI ToggleButtonGroup for filter UI
4. Update page.tsx to read searchParams
5. Create AlbumsPageClient as Client Component
6. Pass filterBar as prop to AlbumsPage
7. Keep AlbumsGrid, AlbumCard unchanged
8. Use brand color #185986 for selected state
9. Handle filter change by updating URL
10. Test all filter options (all, mine, owner)
11. Maintain responsive design
12. Follow existing patterns from Story 1.3-1.4
13. Calculate owner options from albums
14. Show album count per owner
15. Support keyboard navigation

### Anti-Patterns

**From AGENTS.md - Simplicity Principle:**

❌ **Avoid:**

- Over-engineering with complex state management
- Using reducers when simple functions suffice
- Adding unnecessary abstractions
- Tight coupling between components

✅ **Follow:**

- Simple solutions for simple problems
- Pure functions for filtering
- URL-based state for persistence
- Loose coupling with props

**From nextjs.instructions.md - Testing Strategy:**

❌ **Avoid:**

- Testing implementation details
- Skipping edge cases (empty albums, no owner data)
- Not testing keyboard navigation

✅ **Follow:**

- Test behavior (filter works correctly)
- Test all filter options
- Test URL state persistence
- Test keyboard accessibility

### Commit Message Pattern

**From Git History:**

```
catalog/web - add album filtering by owner with URL-based state

- Create pure filter functions (filterAlbumsByOwner, getOwnerOptions)
- Implement AlbumFilterBar with MUI ToggleButtonGroup
- Create AlbumsPageClient to handle client-side filtering
- Update page.tsx to read filter from URL search params
- Update AlbumsPage to accept filterBar prop
- Filter options: All Albums, My Albums, specific owners
- Selected state uses brand blue #185986
- Filter state persists in URL across page refresh
- Keyboard navigation support (tab, enter)
- Responsive layout (mobile/desktop)
- All tests passing (X), build succeeds
```

**Pattern:** `catalog/web - <action description>`

### References

- **Story 1.1**: Material UI theme, Link wrapper, authentication layout
- **Story 1.2**: Catalog domain (unused), fetch adapter, image loader
- **Story 1.3**: Server Component + pure UI pattern, album list layout
- **Story 1.4**: Album card enhancements (density, avatars)
- **Epic 1, Story 1.5**: Album Filtering requirements (FR1, FR2, FR3, FR41)
- **Architecture.md**: Material UI patterns, responsive design, NextJS App Router
- **UX Design Specification**: Filter UI placement and behavior
- **nextjs.instructions.md**: Testing strategy, file structure, coding standards

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
