# Story 1.4: Album Card Enhancements

Status: ready-for-dev

## Story

As a user,  
I want to see preview photos and activity indicators on each album card,  
So that I can quickly understand what's in each album before clicking.

## Acceptance Criteria

**Given** the basic album list is displayed (from Story 1.3)  
**When** I view an album card  
**Then** each album card displays 3-4 random photo thumbnails in a preview grid  
**And** random photos are already included in the albums data from the API (mediaIds field)  
**And** thumbnails use the Next.js Image component with the custom loader  
**And** thumbnails request width=360 (appropriate for card preview size, maps to backend MiniatureCachedWidth)  
**And** thumbnail images use blur placeholders via Next.js Image placeholder prop  
**And** density color-coding is applied based on photos-per-day calculation:

- High density (>10 photos/day): Warmer color accent (#e57373 red tint) or bolder display
- Medium density (3-10 photos/day): Neutral color (brand blue #185986)
- Low density (<3 photos/day): Cooler color (#4a9ece light cyan) or lighter display

**And** visual density indicator appears on the card as a subtle accent bar or glow (FR7)  
**And** sharing status is displayed with user avatars when album is shared (FR43, FR5)  
**And** user avatars display profile pictures from owner and sharedWith data  
**And** avatar tooltip shows user name on hover  
**And** the card maintains responsive layout on mobile, tablet, and desktop  
**And** images load progressively without blocking card rendering  
**And** card appearance uses MUI sx prop for styling (NO inline styles)  
**And** brand color #185986 is used for primary interactive elements  
**And** the album card from Story 1.3 is enhanced with these features (NOT recreated)

## Tasks / Subtasks

- [ ] Enhance AlbumCard component with photo thumbnails (AC: thumbnails, blur placeholder)
    - [ ] Update `app/(authenticated)/_components/AlbumCard.tsx` to display 4 random photos
    - [ ] Photos are already in album.mediaIds (first 4), no additional fetch needed
    - [ ] Render 2x2 grid of photos edge-to-edge with 8px gap
    - [ ] Use Next.js Image component with fill, sizes props
    - [ ] Set sizes="(max-width: 768px) 50vw, 170px" for responsive hints
    - [ ] Custom loader will map width ≤360 → 360 (already configured in Story 1.2)
    - [ ] Use placeholder="empty" (Next.js will handle blur via loader)
    - [ ] Parent Box: position relative, aspectRatio 1
    - [ ] style={{ objectFit: 'cover' }} for proper aspect ratio

- [ ] Add density color-coding calculation (AC: density color-coding)
    - [ ] Create helper function `calculateDensity(totalCount, start, end): number`
    - [ ] Formula: `totalCount / numberOfDays(start, end)`
    - [ ] numberOfDays: `Math.ceil((new Date(end) - new Date(start)) / (1000 * 60 * 60 * 24)) + 1`
    - [ ] Create helper function `getDensityColor(density: number): string`
    - [ ] High density (>10): return '#e57373' (red tint)
    - [ ] Medium density (3-10): return '#185986' (brand blue)
    - [ ] Low density (<3): return '#4a9ece' (light cyan)

- [ ] Add visual density indicator to card (AC: visual density indicator)
    - [ ] Create `app/(authenticated)/_components/DensityIndicator.tsx`
    - [ ] Accept props: density (number), color (string)
    - [ ] Render subtle accent bar at top or left edge of card
    - [ ] Use MUI Box with height 3px or width 3px
    - [ ] Background color from density color
    - [ ] Optional: Add subtle box-shadow glow with same color
    - [ ] Integrate into AlbumCard component

- [ ] Add sharing status avatars (AC: sharing status, user avatars)
    - [ ] Create `app/(authenticated)/_components/SharingAvatars.tsx`
    - [ ] Accept props: owner (OwnerDetails), sharedWith (UserDetails[])
    - [ ] Render MUI AvatarGroup component
    - [ ] Display owner avatar first (if owner data exists)
    - [ ] Display shared user avatars (max 3 visible, +N for overflow)
    - [ ] Each avatar: MUI Avatar with src from profile picture URL
    - [ ] Wrap each avatar in MUI Tooltip with user name
    - [ ] Position AvatarGroup in top-right corner of card (absolute positioning)
    - [ ] Background: semi-transparent dark overlay for visibility
    - [ ] Integrate into AlbumCard component
    - [ ] Handle missing avatar data gracefully (show initials fallback)

- [ ] Update AlbumCard component integration (AC: all)
    - [ ] Import DensityIndicator and SharingAvatars components
    - [ ] Calculate density from totalCount, start, end
    - [ ] Get density color from getDensityColor helper
    - [ ] Pass density and color to DensityIndicator
    - [ ] Pass owner and sharedWith to SharingAvatars
    - [ ] Ensure existing hover effects still work
    - [ ] Maintain edge-to-edge photo design from Story 1.3
    - [ ] Keep all existing styling (Georgia font, Courier date, cyan count)

- [ ] Update TypeScript interfaces (AC: owner, sharedWith)
    - [ ] Update `domains/catalog/language/Album.ts` if needed
    - [ ] Ensure Album interface includes:
        - owner?: OwnerDetails
        - sharedWith?: UserDetails[]
    - [ ] Create `domains/catalog/language/OwnerDetails.ts` if not exists:
        - userId: string
        - name: string
        - profilePicture?: string
    - [ ] Create `domains/catalog/language/UserDetails.ts` if not exists:
        - userId: string
        - name: string
        - email: string
        - profilePicture?: string
    - [ ] Update AlbumCardProps interface to include owner and sharedWith

- [ ] Write component tests (AC: all tests pass)
    - [ ] Update `app/(authenticated)/_components/__tests__/AlbumCard.test.tsx`
        - [ ] Test renders 4 photo thumbnails
        - [ ] Test density indicator appears with correct color
        - [ ] Test sharing avatars display when album is shared
        - [ ] Test owner avatar displays
        - [ ] Test avatar tooltips show user names
        - [ ] Test no sharing avatars when album not shared
    - [ ] Create `app/(authenticated)/_components/__tests__/DensityIndicator.test.tsx`
        - [ ] Test renders with high density color
        - [ ] Test renders with medium density color
        - [ ] Test renders with low density color
    - [ ] Create `app/(authenticated)/_components/__tests__/SharingAvatars.test.tsx`
        - [ ] Test renders owner avatar
        - [ ] Test renders shared user avatars
        - [ ] Test displays max 3 avatars with overflow
        - [ ] Test tooltips show correct names
        - [ ] Test handles missing profile pictures (initials fallback)

- [ ] Run tests and verify build (AC: all)
    - [ ] Run `npm run test` - verify all tests pass
    - [ ] Run `npm run build` - verify successful build
    - [ ] Run `npm run test:visual` - verify visual tests pass (if applicable)
    - [ ] Verify no TypeScript errors
    - [ ] Verify no missing imports
    - [ ] Verify density indicators display correctly
    - [ ] Verify sharing avatars display correctly

## Dev Notes

### Story Context

**This story enhances the album cards created in Story 1.3 with:**

1. **Photo thumbnails** - 4 random photos in 2x2 grid (already in Story 1.3 design)
2. **Density indicators** - color-coded visual feedback based on photos-per-day
3. **Sharing avatars** - display owner and shared users with profile pictures

**Stories 1.1, 1.2, 1.3 provide the foundation:**

- Story 1.1: Material UI theme with brand color #185986
- Story 1.2: Catalog domain with Album types, fetch adapter, image loader
- Story 1.3: Basic album list with card layout (edge-to-edge photos, text overlay, hover effects)

**Story 1.4 builds on Story 1.3 by ENHANCING the AlbumCard, not recreating it.**

### Architecture Pattern (Server Component + Pure UI)

**From Story 1.3 - Same Pattern:**

```typescript
// Server Component (page.tsx) - NO CHANGES NEEDED
export default async function HomePage() {
    const albums = await fetchAlbumsAdapter()
    return <AlbumsPage albums={albums} />
}

// Pure UI Component - ENHANCE AlbumCard
function AlbumCard(props: AlbumCardProps) {
    const { mediaIds, totalCount, start, end, owner, sharedWith } = props
    
    // Calculate density
    const density = calculateDensity(totalCount, start, end)
    const densityColor = getDensityColor(density)
    
    return (
        <Link href={`/owners/${ownerId}/albums/${folderName}`}>
            <Box sx={{ position: 'relative', ... }}>
                <DensityIndicator density={density} color={densityColor} />
                
                {/* 2x2 Photo Grid (already implemented in Story 1.3) */}
                <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1 }}>
                    {mediaIds.slice(0, 4).map(mediaId => (
                        <Box key={mediaId} sx={{ position: 'relative', aspectRatio: '1' }}>
                            <Image 
                                src={mediaId}
                                fill
                                sizes="(max-width: 768px) 50vw, 170px"
                                alt=""
                                style={{ objectFit: 'cover' }}
                            />
                        </Box>
                    ))}
                </Box>
                
                {/* Text Overlay (already implemented) */}
                <Box sx={{ position: 'absolute', bottom: 0, ... }}>
                    <Typography>{name}</Typography>
                    <Typography>{formatDateRange(start, end)}</Typography>
                    <Typography>{totalCount} photos</Typography>
                </Box>
                
                {/* NEW: Sharing Avatars */}
                <SharingAvatars owner={owner} sharedWith={sharedWith} />
            </Box>
        </Link>
    )
}
```

**Key Points:**

- ✅ Story 1.3 already implemented photo thumbnails (2x2 grid, edge-to-edge)
- ✅ Story 1.4 adds density indicator and sharing avatars
- ✅ NO changes to Server Component (page.tsx)
- ✅ NO client state needed (still pure UI)
- ❌ NO new API calls (data already in albums)

### Component Props Design

**Update AlbumCardProps (from Story 1.3):**

```typescript
interface AlbumCardProps {
    albumId: string
    name: string
    start: string          // ISO date
    end: string            // ISO date
    totalCount: number
    ownerId: string
    folderName: string
    mediaIds: string[]     // ✅ Already used in Story 1.3
    owner?: OwnerDetails   // 🆕 Story 1.4
    sharedWith?: UserDetails[]  // 🆕 Story 1.4
}

interface OwnerDetails {
    userId: string
    name: string
    profilePicture?: string
}

interface UserDetails {
    userId: string
    name: string
    email: string
    profilePicture?: string
}
```

**Pass from parent (AlbumsGrid):**

```typescript
function AlbumsGrid({ albums }: { albums: Album[] }) {
    return (
        <Box sx={{ display: 'grid', ... }}>
            {albums.map(album => (
                <AlbumCard
                    key={album.albumId}
                    albumId={album.albumId}
                    name={album.name}
                    start={album.start}
                    end={album.end}
                    totalCount={album.totalCount}
                    ownerId={album.ownerId}
                    folderName={album.folderName}
                    mediaIds={album.mediaIds?.slice(0, 4) || []}
                    owner={album.owner}           // 🆕 Story 1.4
                    sharedWith={album.sharedWith} // 🆕 Story 1.4
                />
            ))}
        </Box>
    )
}
```

### Density Calculation

**Formula:**

```typescript
function calculateDensity(totalCount: number, start: string, end: string): number {
    const numberOfDays = Math.ceil(
        (new Date(end).getTime() - new Date(start).getTime()) / (1000 * 60 * 60 * 24)
    ) + 1  // Include both start and end days
    
    return totalCount / numberOfDays
}

function getDensityColor(density: number): string {
    if (density > 10) return '#e57373'  // High: red tint (warm)
    if (density >= 3) return '#185986'   // Medium: brand blue
    return '#4a9ece'                     // Low: light cyan (cool)
}
```

**Examples:**

- 100 photos / 5 days = 20 photos/day → High density → #e57373
- 50 photos / 10 days = 5 photos/day → Medium density → #185986
- 10 photos / 7 days = 1.4 photos/day → Low density → #4a9ece

### Density Indicator Component

**Visual Design:**

- Subtle accent bar at top edge of card
- 3px height, full width
- Background color from density calculation
- Optional: subtle box-shadow glow for emphasis

```typescript
interface DensityIndicatorProps {
    density: number
    color: string
}

function DensityIndicator({ density, color }: DensityIndicatorProps) {
    return (
        <Box
            sx={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                height: '3px',
                bgcolor: color,
                boxShadow: `0 0 8px ${color}`,  // Subtle glow
                zIndex: 1,
            }}
        />
    )
}
```

**Alternative Design (Left Edge Bar):**

```typescript
function DensityIndicator({ density, color }: DensityIndicatorProps) {
    return (
        <Box
            sx={{
                position: 'absolute',
                top: 0,
                left: 0,
                bottom: 0,
                width: '3px',
                bgcolor: color,
                boxShadow: `0 0 8px ${color}`,
                zIndex: 1,
            }}
        />
    )
}
```

Choose the design that works best visually. Top edge is recommended for subtlety.

### Sharing Avatars Component

**Visual Design:**

- Position: absolute top-right corner of card
- MUI AvatarGroup with max 3 visible avatars
- Owner avatar first (primary), then shared users
- Tooltip on hover showing user name
- Semi-transparent background for visibility over photos

```typescript
interface SharingAvatarsProps {
    owner?: OwnerDetails
    sharedWith?: UserDetails[]
}

function SharingAvatars({ owner, sharedWith }: SharingAvatarsProps) {
    const hasSharing = sharedWith && sharedWith.length > 0
    
    if (!hasSharing && !owner) {
        return null  // No sharing info to display
    }
    
    return (
        <Box
            sx={{
                position: 'absolute',
                top: 8,
                right: 8,
                zIndex: 2,
                bgcolor: 'rgba(10, 21, 32, 0.8)',  // Dark semi-transparent
                borderRadius: 1,
                padding: 0.5,
            }}
        >
            <AvatarGroup max={4} sx={{ flexDirection: 'row' }}>
                {owner && (
                    <Tooltip title={owner.name} placement="top">
                        <Avatar
                            src={owner.profilePicture}
                            alt={owner.name}
                            sx={{ width: 32, height: 32, border: '2px solid #185986' }}
                        >
                            {getInitials(owner.name)}
                        </Avatar>
                    </Tooltip>
                )}
                
                {sharedWith?.map(user => (
                    <Tooltip key={user.userId} title={user.name} placement="top">
                        <Avatar
                            src={user.profilePicture}
                            alt={user.name}
                            sx={{ width: 32, height: 32 }}
                        >
                            {getInitials(user.name)}
                        </Avatar>
                    </Tooltip>
                ))}
            </AvatarGroup>
        </Box>
    )
}

function getInitials(name: string): string {
    return name
        .split(' ')
        .map(n => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
}
```

**MUI Components:**

- `AvatarGroup` - groups avatars with overlap
- `Avatar` - profile picture or initials fallback
- `Tooltip` - shows user name on hover

### Data Flow (No Changes to Server Component)

**Story 1.3 Already Fetches Albums:**

```typescript
// page.tsx (NO CHANGES)
export default async function HomePage() {
    const albums = await fetchAlbumsAdapter()  // ✅ Already includes owner and sharedWith
    return <AlbumsPage albums={albums} />
}
```

**fetch-adapter.ts Already Returns Full Album Data:**

The fetch adapter from Story 1.2 should already return:

```typescript
interface Album {
    albumId: string
    name: string
    start: string
    end: string
    folderName: string
    ownerId: string
    totalCount: number
    mediaIds?: string[]       // ✅ Already used
    owner?: OwnerDetails      // ✅ Should exist in API response
    sharedWith?: UserDetails[] // ✅ Should exist in API response
}
```

**If owner/sharedWith are missing from fetch-adapter.ts, update it in this story.**

### Image Loading (Already Configured in Story 1.2)

**Custom Image Loader:**

```typescript
// libs/image-loader.ts (Already exists from Story 1.2)
export default function imageLoader({ src, width }: ImageLoaderProps): string {
    const targetWidth = width <= 360 ? 360 : 2400
    return `/api/v1/media/${src}/image?width=${targetWidth}`
}
```

**For Story 1.4:**

- Album card photos are small (~170px per photo in 2x2 grid)
- Next.js will request width around 170-340px
- Custom loader will map to 360 (MiniatureCachedWidth)
- Perfect for card preview thumbnails

**Usage (Already in Story 1.3):**

```typescript
<Image
    src = {mediaId}
fill
sizes = "(max-width: 768px) 50vw, 170px"
alt = ""
style = {
{
    objectFit: 'cover'
}
}
/>
```

**No changes needed to image loading.**

### Material UI Components Used

**New Components for Story 1.4:**

- `AvatarGroup` - from @mui/material
- `Avatar` - from @mui/material
- `Tooltip` - from @mui/material
- `Box` - already used in Story 1.3

**Already Used from Story 1.3:**

- `Box` - layout container
- `Typography` - text elements
- `Paper`, `Button` - error states

### Responsive Design

**Album Card Grid (Story 1.3 - Unchanged):**

- Mobile (xs <600px): 1 column
- Tablet (sm 600-960px): 2 columns
- Desktop (md 960-1280px): 3 columns
- Large (lg >1280px): 4 columns

**Photo Grid Within Card (Story 1.3 - Unchanged):**

- Always 2x2 grid (4 photos)
- 8px gap between photos
- Edge-to-edge design

**Sharing Avatars (Story 1.4):**

- Size: 32x32px on all devices
- Position: top-right, 8px padding
- Max 4 avatars visible (+N overflow)

**Density Indicator (Story 1.4):**

- Top edge: 3px height, full width
- OR Left edge: 3px width, full height
- Same on all devices

### Testing Strategy

**Component Tests (Vitest + React Testing Library):**

```typescript
// AlbumCard.test.tsx (UPDATE)
describe('AlbumCard', () => {
    const mockProps = {
        albumId: 'album-1',
        name: 'Beach Trip',
        start: '2026-07-15',
        end: '2026-07-22',
        totalCount: 80,  // 80 photos / 8 days = 10/day (medium)
        ownerId: 'owner-1',
        folderName: 'beach-trip',
        mediaIds: ['media-1', 'media-2', 'media-3', 'media-4'],
        owner: {
            userId: 'owner-1',
            name: 'John Doe',
            profilePicture: 'https://example.com/john.jpg',
        },
        sharedWith: [
            {
                userId: 'user-2',
                name: 'Jane Smith',
                email: 'jane@example.com',
                profilePicture: 'https://example.com/jane.jpg',
            },
        ],
    }
    
    it('renders 4 photo thumbnails', () => {
        render(<AlbumCard {...mockProps} />)
        const images = screen.getAllByRole('img')
        expect(images.length).toBeGreaterThanOrEqual(4)  // 4 photos + avatars
    })
    
    it('displays density indicator with medium density color', () => {
        render(<AlbumCard {...mockProps} />)
        const indicator = screen.getByTestId('density-indicator')  // Add testId in component
        expect(indicator).toHaveStyle({ backgroundColor: '#185986' })
    })
    
    it('displays owner avatar', () => {
        render(<AlbumCard {...mockProps} />)
        expect(screen.getByAltText('John Doe')).toBeInTheDocument()
    })
    
    it('displays shared user avatars', () => {
        render(<AlbumCard {...mockProps} />)
        expect(screen.getByAltText('Jane Smith')).toBeInTheDocument()
    })
    
    it('shows tooltip with user name on avatar hover', async () => {
        render(<AlbumCard {...mockProps} />)
        const avatar = screen.getByAltText('John Doe')
        fireEvent.mouseOver(avatar)
        await waitFor(() => {
            expect(screen.getByText('John Doe')).toBeInTheDocument()
        })
    })
    
    it('hides sharing avatars when album not shared', () => {
        const propsWithoutSharing = { ...mockProps, owner: undefined, sharedWith: undefined }
        render(<AlbumCard {...propsWithoutSharing} />)
        expect(screen.queryByTestId('sharing-avatars')).not.toBeInTheDocument()
    })
})

// DensityIndicator.test.tsx
describe('DensityIndicator', () => {
    it('renders with high density color', () => {
        render(<DensityIndicator density={15} color="#e57373" />)
        const indicator = screen.getByTestId('density-indicator')
        expect(indicator).toHaveStyle({ backgroundColor: '#e57373' })
    })
    
    it('renders with medium density color', () => {
        render(<DensityIndicator density={5} color="#185986" />)
        const indicator = screen.getByTestId('density-indicator')
        expect(indicator).toHaveStyle({ backgroundColor: '#185986' })
    })
    
    it('renders with low density color', () => {
        render(<DensityIndicator density={1} color="#4a9ece" />)
        const indicator = screen.getByTestId('density-indicator')
        expect(indicator).toHaveStyle({ backgroundColor: '#4a9ece' })
    })
})

// SharingAvatars.test.tsx
describe('SharingAvatars', () => {
    const mockOwner = {
        userId: 'owner-1',
        name: 'John Doe',
        profilePicture: 'https://example.com/john.jpg',
    }
    
    const mockSharedWith = [
        { userId: 'user-2', name: 'Jane Smith', email: 'jane@example.com', profilePicture: 'https://example.com/jane.jpg' },
        { userId: 'user-3', name: 'Bob Wilson', email: 'bob@example.com' },
    ]
    
    it('renders owner avatar', () => {
        render(<SharingAvatars owner={mockOwner} />)
        expect(screen.getByAltText('John Doe')).toBeInTheDocument()
    })
    
    it('renders shared user avatars', () => {
        render(<SharingAvatars owner={mockOwner} sharedWith={mockSharedWith} />)
        expect(screen.getByAltText('Jane Smith')).toBeInTheDocument()
        expect(screen.getByAltText('Bob Wilson')).toBeInTheDocument()
    })
    
    it('displays initials when profile picture missing', () => {
        render(<SharingAvatars sharedWith={mockSharedWith} />)
        expect(screen.getByText('BW')).toBeInTheDocument()  // Bob Wilson initials
    })
    
    it('returns null when no owner or shared users', () => {
        const { container } = render(<SharingAvatars />)
        expect(container.firstChild).toBeNull()
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
│   ├── page.tsx                    # ✅ Already exists (Story 1.3) - NO CHANGES
│   └── _components/                # ✅ Already exists (Story 1.3)
│       ├── AlbumsPage.tsx          # ✅ Already exists - NO CHANGES
│       ├── AlbumsGrid.tsx          # ✅ Already exists - UPDATE: pass owner, sharedWith
│       ├── AlbumCard.tsx           # ✅ Already exists - ENHANCE: add density, avatars
│       ├── DensityIndicator.tsx    # 🆕 Story 1.4
│       ├── SharingAvatars.tsx      # 🆕 Story 1.4
│       ├── LoadingSkeleton.tsx     # ✅ Already exists - NO CHANGES
│       ├── EmptyState.tsx          # ✅ Already exists - NO CHANGES
│       ├── ErrorDisplay.tsx        # ✅ Already exists - NO CHANGES
│       └── __tests__/              # ✅ Already exists
│           ├── AlbumCard.test.tsx  # UPDATE: add tests for density, avatars
│           ├── DensityIndicator.test.tsx  # 🆕 Story 1.4
│           └── SharingAvatars.test.tsx    # 🆕 Story 1.4

domains/
└── catalog/                        # ✅ Already exists (Story 1.2)
    ├── language/
    │   ├── Album.ts                # UPDATE: ensure owner, sharedWith fields
    │   ├── OwnerDetails.ts         # 🆕 Story 1.4 (if not exists)
    │   └── UserDetails.ts          # 🆕 Story 1.4 (if not exists)
    └── adapters/
        └── fetch-adapter.ts        # CHECK: ensure owner, sharedWith returned
```

### Fetch Adapter Update (If Needed)

**Check fetch-adapter.ts from Story 1.2:**

If the adapter doesn't return owner and sharedWith, update it:

```typescript
// domains/catalog/adapters/fetch-adapter.ts
export async function fetchAlbumsAdapter(): Promise<Album[]> {
    const response = await fetch('/api/v1/albums', {
        credentials: 'include',
    })
    
    if (!response.ok) {
        throw new CatalogError(`Failed to fetch albums: ${response.status}`)
    }
    
    const data = await response.json()
    
    return data.albums.map(album => ({
        albumId: album.albumId,
        name: album.name,
        start: album.start,
        end: album.end,
        folderName: album.folderName,
        ownerId: album.ownerId,
        totalCount: album.totalCount,
        mediaIds: album.mediaIds || [],
        owner: album.owner ? {                    // 🆕 Story 1.4
            userId: album.owner.userId,
            name: album.owner.name,
            profilePicture: album.owner.profilePicture,
        } : undefined,
        sharedWith: album.sharedWith?.map(user => ({  // 🆕 Story 1.4
            userId: user.userId,
            name: user.name,
            email: user.email,
            profilePicture: user.profilePicture,
        })) || [],
    }))
}
```

**If the API already returns owner and sharedWith, no changes needed.**

### Success Validation

**Before Raising PR:**

1. ✅ All tests pass: `npm run test`
2. ✅ Build succeeds: `npm run build`
3. ✅ Visual tests pass: `npm run test:visual` (if applicable)
4. ✅ No TypeScript errors
5. ✅ No missing imports
6. ✅ Album cards display 4 photo thumbnails edge-to-edge
7. ✅ Density indicators appear with correct colors
8. ✅ Sharing avatars display in top-right corner
9. ✅ Tooltips show user names on avatar hover
10. ✅ Responsive layout works on mobile/tablet/desktop
11. ✅ Hover effects still work from Story 1.3
12. ✅ All acceptance criteria met

### Critical Guardrails

❌ **DO NOT:**

1. Recreate AlbumCard component (enhance existing from Story 1.3)
2. Change Server Component architecture (no new client state)
3. Make new API calls (data already in albums)
4. Break existing hover effects from Story 1.3
5. Change edge-to-edge photo design from Story 1.3
6. Use inline styles (MUI sx prop only)
7. Modify page.tsx (no changes to Server Component)
8. Add useReducer or client state (pure UI only)
9. Change image loader configuration
10. Forget to test all density ranges (high/medium/low)

✅ **DO:**

1. Enhance existing AlbumCard from Story 1.3
2. Keep pure UI components (NO hooks, NO state)
3. Use density calculation for color-coding
4. Display sharing avatars with MUI AvatarGroup
5. Show tooltips on avatar hover
6. Use MUI sx prop for all styling
7. Maintain responsive design from Story 1.3
8. Test density calculations thoroughly
9. Handle missing profile pictures (initials fallback)
10. Pass owner and sharedWith as props from parent
11. Use brand color #185986 for medium density
12. Add data-testid attributes for testing
13. Follow existing patterns from Stories 1.1-1.3

### Anti-Patterns

**From AGENTS.md - Testing Strategy:**

❌ **Avoid:**

- Testing implementation details (how density is calculated)
- Tight coupling between tests and styling
- Skipping edge cases (missing avatars, zero photos, one-day albums)

✅ **Follow:**

- Test behavior (density indicator appears, correct color)
- Test user-facing elements (avatars visible, tooltips work)
- Test all density ranges (high/medium/low)
- Test missing data handling (no owner, no shared users)

### Commit Message Pattern

**From Git History:**

```
catalog/web - enhance album cards with density indicators and sharing avatars

- Add density calculation based on photos-per-day
- Implement color-coded density indicators (high/medium/low)
- Display sharing status with user avatars and tooltips
- Create DensityIndicator component with subtle accent bar
- Create SharingAvatars component with MUI AvatarGroup
- Handle missing profile pictures with initials fallback
- Update AlbumCard to integrate new components
- Maintain edge-to-edge photo design from Story 1.3
- All tests passing (X), build succeeds
```

**Pattern:** `catalog/web - <action description>`

### References

- **Story 1.1**: Material UI theme with brand color #185986
- **Story 1.2**: Catalog domain, fetch adapter, image loader configuration
- **Story 1.3**: Basic album list with AlbumCard component (edge-to-edge photos, text overlay, hover effects)
- **Epic 1, Story 1.4**: Album Card Enhancements requirements
- **Architecture.md**: Material UI patterns, responsive design, MUI sx prop
- **UX Design Specification**: Density color-coding, sharing status display
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
