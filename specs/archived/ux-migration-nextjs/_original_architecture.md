---
stepsCompleted: [ 1, 2, 3, 4, 5, 6, 7, 8 ]
workflowType: 'architecture'
lastStep: 8
status: 'complete'
completedAt: '2026-01-31'
inputDocuments:
  - '/home/dush/dev/git/dphoto/specs/designs/prd.md'
  - '/home/dush/dev/git/dphoto/specs/designs/ux-design-specification.md'
  - '/home/dush/dev/git/dphoto/specs/2026-01-ux-functionnal.md'
  - '/home/dush/dev/git/dphoto/AGENTS.md'
project_name: 'dphoto'
user_name: 'Arch'
date: '2026-01-31'
scope: 'NextJS Web UI (web-nextjs) - Frontend architecture only'
outOfScope: 'Backend (pkg/), Infrastructure (deployments/cdk/), API (api/lambdas/), CLI (cmd/dphoto/), Data Model (DATA_MODEL.md)'
---

# Architecture Decision Document - DPhoto NextJS Web UI

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

The NextJS web UI must implement **50 functional requirements** organized into 9 capability areas:

1. **Album Discovery & Browsing (FR1-FR9)** - View owned and shared albums, filter by owner, see metadata (name, dates, count, owner info, sharing status),
   chronological ordering, visual activity indicators, random photo discovery with navigation to source albums

2. **Photo Viewing & Navigation (FR10-FR17)** - Day-grouped photo display, full-screen view, keyboard controls (arrows/ESC/ENTER), touch gestures (swipe),
   progressive quality loading, zoom capability, clear navigation back to album list, date headers for context

3. **Album Management (FR18-FR25)** - Create albums with name and date range, custom folder names (optional), edit names and dates, delete owned albums,
   validate dates don't orphan media, re-index on date changes, clear success/failure feedback

4. **Sharing & Access Control (FR26-FR32)** - Share albums via email, revoke access, view shared users, display owner information, distinguish owner vs viewer
   capabilities in UI, validate emails, load user profiles (name, picture)

5. **Authentication & User Management (FR33-FR37)** - Google OAuth authentication (backend-provided), session management, profile display, owner role
   identification, all pages require authentication

6. **Visual Presentation & UI State (FR38-FR44)** - Loading states during fetches, empty states when no albums exist, error messages on failures, selected album
   and active filter indicators, smooth view transitions, sharing status with avatars, responsive layouts (mobile/tablet/desktop)

7. **Error Handling & Validation (FR45-FR50)** - Validate date ranges, non-empty album names, handle failed operations, album not found errors, recovery options

**Non-Functional Requirements:**

**Performance:**

- **Interaction responsiveness** - User interactions (clicks, taps, swipes) must provide immediate visual feedback (<100ms)
- **Image loading** - Progressive loading displays thumbnail-quality images immediately with blur-up to full resolution within 3 seconds on slow networks
- **Page performance** - Lighthouse Performance score ≥90 on mobile, optimized initial load and time-to-interactive through code splitting and lazy loading
- **Responsive image sizing** - Request minimum appropriate size for current screen (mobile/tablet/desktop) using existing API quality parameters

**Integration:**

- **Backend API compatibility** - Consume existing REST API endpoints without requiring backend modifications
- **No breaking changes** - Maintain data contract compatibility with existing API response formats
- **Graceful API handling** - Handle response times with loading states, errors with clear feedback and retry mechanisms

**Security:**

- **Backend authentication** - Relies on existing Google OAuth mechanism
- **Session security** - Maintained by backend
- **Authenticated pages** - All routes restricted to authenticated users (enforced by backend)
- **No additional security** - Frontend consumes authenticated API endpoints only

**Usability:**

- **Keyboard navigation** - Core flows fully keyboard-accessible (arrows, ESC, ENTER), clear focus management, no browser default conflicts
- **Mobile & Responsive** - Touch interactions with visual feedback, layouts adapt across breakpoints (mobile <600px, tablet 600-960px, desktop >960px), natural
  mobile gestures (swipe, pinch-to-zoom), no performance degradation on mobile
- **Browser compatibility** - Latest 2 versions of Chrome, Firefox, Safari, Edge (evergreen browsers), may use modern features (CSS Grid, Flexbox, ES2020+)
  without polyfills, no IE11 or legacy browser support

**Reliability:**

- **Best-effort availability** - No uptime SLA required
- **Manual refresh recovery** - Acceptable for transient errors
- **Clear error messaging** - Provide guidance on recovery
- **No data loss** - Network failures must not lose in-progress operations

### Scale & Complexity

**Project Scale:**

- **Primary domain:** Web application (NextJS/React/TypeScript)
- **Complexity level:** Medium
- **Architecture type:** Multi-page application structure with NextJS App Router
- **Estimated components:** 20-30 reusable components

**Complexity Indicators:**

- **High complexity:** Progressive image loading with responsive sizing and caching, multi-device optimization with different interaction models per platform
- **Medium complexity:** Visual date selection with real-time photo previews, Material UI customization (theming, component extensions), keyboard navigation
  throughout application
- **Low complexity:** No real-time features (manual refresh model), no offline support or PWA capabilities, no SEO requirements (all behind authentication),
  authentication handled entirely by backend

### Technical Constraints & Dependencies

**Must Use:**

- NextJS App Router (already initiated in `web-nextjs/`)
- Material UI component library and design system
- TypeScript throughout
- Existing REST API without modifications
- Existing authentication mechanism (Google OAuth via backend)

**Browser Support:**

- Modern evergreen browsers only (latest 2 versions)
- No polyfills required - can use ES2020+, CSS Grid, Flexbox directly
- No legacy browser support (IE11, older versions)

**Performance Targets:**

- Lighthouse Performance score ≥90 on mobile
- <100ms response to user interactions
- Progressive image loading (blur-up within 3 seconds on slow networks)

**Design System:**

- Material UI with dark theme
- Brand color integration (#185986 blue throughout)
- Responsive breakpoints: Mobile (<600px), Tablet (600-960px), Desktop (>960px)

### Cross-Cutting Concerns Identified

**State Management:**

- Album list state (owned and shared albums)
- Selected album and active filters
- Photo data for current album
- Sharing state (users with access)
- User authentication state and profile
- Loading and error states across all operations

**Image Optimization:**

- Progressive loading strategy (blur-up from low to high quality)
- Responsive image sizing based on device/viewport
- Caching strategies for thumbnails and full-resolution images
- Leverage existing API quality parameters

**Error Handling:**

- Network failure scenarios
- Permission errors (editing shared albums)
- Photo/album loading failures
- API error responses with user-friendly messages
- Retry mechanisms and recovery paths

**Responsive Design:**

- Three distinct breakpoints with appropriate layouts
- Different interaction models: keyboard (desktop), touch (mobile/tablet), mouse (desktop)
- Progressive enhancement from mobile to desktop
- Touch-friendly controls on mobile/tablet

**Accessibility:**

- Keyboard navigation for all core flows
- Focus management in dialogs and photo viewing
- ARIA labels and semantic HTML
- Clear visual indicators for keyboard users

**Material UI Integration:**

- Dark theme configuration
- Brand color (#185986) as primary throughout
- Component customization strategy (when to extend vs build custom)
- Theming approach for consistent styling

**Innovative UX Patterns:**

- Visual date selection with live photo previews (real-time thumbnails as dates adjust)
- Random photo discovery (home page highlights and album card samples)
- Contextual album creation from media list (smart date suggestions)
- Timeline navigation with visual density indicators

## Starter Template Evaluation

### Primary Technology Domain

**Web Application (NextJS/React/TypeScript)** - Frontend for authenticated photo management application with progressive image loading, responsive design, and
Material UI component library integration.

### Current Foundation Analysis

The `web-nextjs/` project has already been initialized and is production-ready with core infrastructure in place. Rather than using a third-party starter
template, the project was custom-initialized with specific requirements for the DPhoto use case.

**Project Status:** ✅ **Already Initialized** - Architecture decisions focus on completing the feature implementation on top of this established foundation.

### Existing Technical Stack

**Language & Runtime:**

- **TypeScript 5.x** with strict mode enabled
- **ES2017 target** with modern library support (DOM, ESNext)
- Path aliases configured (`@/*` for root imports)
- Strict type checking enabled for code quality

**Framework & Core Libraries:**

- **Next.js 16.1.1** with App Router architecture
- **React 19.2.3** with React Compiler plugin enabled
- **Server-only components** support for optimized server rendering
- Standalone output mode configured for deployment

**Styling Solution:**

- **Tailwind CSS 4.x** with PostCSS integration
- Dark theme and brand color (#185986) integration needed (architecture decision)
- Material UI integration required per UX specification (not yet installed)

**Authentication & Security:**

- **OpenID Client 6.8.1** for OAuth integration
- Google OAuth authentication already implemented
- Session management utilities in `libs/security/`
- Access token service with JWT utilities
- Backend store for secure token handling
- Cookie management utilities

**Build & Development Tools:**

- **ESLint** with Next.js configuration for code quality
- **Open Next 3.9.7** for AWS Lambda deployment
- Next.js development server with hot reloading
- Path-to-regexp for route matching utilities

**Testing Infrastructure:**

- **Vitest 4.0.16** with jsdom environment for React component testing
- **MSW (Mock Service Worker) 2.12.7** for API mocking in tests
- Vitest UI for interactive test running
- Test utilities in `__tests__/helpers/` for OIDC and assertions

**Request & API Layer:**

- Request utilities in `libs/requests/` for API communication
- Base path configuration (`/nextjs`) for deployment
- Remote image patterns configured for external image sources

**Project Structure:**

```
web-nextjs/
├── app/
│   ├── (authenticated)/    # Protected routes requiring auth
│   │   ├── layout.tsx      # Authenticated layout wrapper
│   │   └── page.tsx        # Main authenticated home page
│   ├── auth/               # Authentication flows
│   │   ├── callback/       # OAuth callback handler
│   │   ├── error/          # Auth error page
│   │   └── logout/         # Logout handler
│   ├── layout.tsx          # Root layout
│   └── globals.css         # Global styles
├── components/
│   └── UserInfo/           # User profile display component
├── libs/
│   ├── nextjs-cookies/     # Cookie utilities
│   ├── requests/           # API request layer
│   └── security/           # Auth & session management
├── __tests__/              # Test infrastructure
└── public/                 # Static assets
```

**Deployment Configuration:**

- **AWS deployment** via SST (Serverless Stack) integration
- Open Next adapter for Lambda@Edge optimization
- Base path `/nextjs` for deployment routing
- Standalone output for containerization

### Architectural Decisions Provided by Current Setup

**✅ Already Decided:**

1. **Framework Choice:** NextJS 16 with App Router - provides file-based routing, server components, and optimized builds
2. **Language:** TypeScript with strict mode - ensures type safety throughout
3. **Authentication:** Google OAuth via backend with secure session management - production-ready implementation
4. **Testing:** Vitest with MSW for mocked API tests - modern, fast test infrastructure
5. **Deployment:** AWS via Open Next and SST - serverless architecture for cost efficiency
6. **Build Process:** Next.js compiler with React Compiler plugin - optimized production builds
7. **Code Quality:** ESLint with Next.js standards - enforces consistency

**🔧 Architecture Decisions Still Needed:**

1. **Material UI Integration** - Theme configuration (dark mode, brand color #185986), component library setup
2. **State Management** - Approach for managing album/photo/filter/sharing state (React Context, Zustand, or other)
3. **Component Architecture** - Structure and organization for 20-30 feature components
4. **Image Optimization Strategy** - Progressive loading implementation, responsive sizing, caching
5. **Routing Patterns** - App Router structure for albums, photos, and management dialogs
6. **Layout Architecture** - Responsive patterns for mobile/tablet/desktop breakpoints
7. **Error Boundaries** - Error handling strategy for graceful degradation
8. **Data Fetching Patterns** - Server Component vs Client Component strategy, caching approach

### Dependencies to Add

Based on UX specification requirements, the following dependencies need to be added:

**Required:**

- `@mui/material` - Material UI component library
- `@mui/icons-material` - Material UI icon set
- `@emotion/react` - Material UI peer dependency for styling
- `@emotion/styled` - Material UI peer dependency for styled components

**Potential (Architecture Decisions):**

- State management library (if not using React Context) - Zustand, Jotai, or Redux Toolkit
- Form management library - React Hook Form (if complex forms need validation)
- Image optimization utilities - if native Next.js Image component insufficient

### Development Workflow

**Current Commands:**

```bash
npm run dev          # Start development server (localhost:3000)
npm run build        # Production build with Open Next
npm run build:next   # Standard Next.js build
npm run test         # Run tests with Vitest
npm run test:watch   # Watch mode for test development
npm run lint         # Run ESLint
npm run clean        # Clean build artifacts
```

**Package Manager:** npm 11.6.4 (configured via packageManager field)

### Next Steps

With this foundation established, architectural decisions will focus on:

1. Completing Material UI integration and theming
2. Defining component structure and patterns
3. Establishing state management approach
4. Implementing progressive image loading strategy
5. Creating routing patterns for the 6 core capabilities
6. Defining responsive layout architecture
7. Establishing error handling patterns

The existing authentication, testing, and deployment infrastructure provides a solid base for building the feature implementation.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**

1. Material UI Integration - Remove Tailwind, use MUI exclusively
2. State Management - React Context API with well-designed patterns
3. Component Architecture - Colocation principle (components live with pages unless shared)
4. Routing Structure - Owner-based paths with parallel routes for photo modal
5. Image Optimization - Next.js Image with custom loader for backend API integration
6. Layout Architecture - Material UI sx prop with theme breakpoints
7. Error Boundaries - NextJS page-level error.tsx files
8. Data Fetching - Server Components for initial data, Client Components for interactions

**Important Decisions (Shape Architecture):**

- Backend API image quality requirements documented
- Responsive image sizing by context defined
- Context organization (split by concern to prevent re-renders)

**Deferred Decisions (Post-Launch):**

- Advanced caching strategies beyond browser defaults
- Performance monitoring and optimization tools
- Progressive Web App capabilities

---

### Material UI Integration

**Decision:** Remove Tailwind CSS, use Material UI exclusively

**Rationale:**

- UX specification explicitly requires Material UI with dark theme
- Single styling system reduces complexity and bundle size
- Material UI's `sx` prop and theme system provides comprehensive styling capabilities
- Eliminates potential conflicts between two CSS systems

**Implementation:**

- Remove Tailwind CSS dependencies from package.json
- Configure Material UI theme with dark mode as default
- Set brand color (#185986) as primary throughout theme
- Use Material UI's breakpoint system: `xs` (<600px), `sm` (600px), `md` (960px), `lg` (1280px)

**Dependencies to Add:**

```json
{
  "@mui/material": "^6.x",
  "@mui/icons-material": "^6.x",
  "@emotion/react": "^11.x",
  "@emotion/styled": "^11.x"
}
```

**Dependencies to Remove:**

```json
{
  "tailwindcss": "^4",
  "@tailwindcss/postcss": "^4"
}
```

---

### State Management

**Decision:** React Context API with multiple contexts split by concern

**Rationale:**

- Built-in React solution, no additional dependencies
- Sufficient for medium complexity (20-30 components)
- Well-understood pattern for all developers and AI agents
- TypeScript provides type safety for context values

**Context Organization:**

**Split contexts to prevent unnecessary re-renders:**

1. **AlbumsContext** - Album list state
    - Albums array (owned + shared)
    - Loading and error states for album list
    - Actions: fetch albums, refresh list

2. **SelectedAlbumContext** - Current album state
    - Selected album ID
    - Album details
    - Actions: select album, clear selection

3. **FilterContext** - Filter state
    - Active filter (My Albums, All Albums, By Owner)
    - Selected owner (for owner filter)
    - Actions: set filter, set owner

4. **PhotosContext** - Current album photos state
    - Photos array for selected album
    - Loading and error states for photos
    - Actions: fetch photos, refresh photos

**Context Location:**

```
components/
  contexts/
    AlbumsContext.tsx
    SelectedAlbumContext.tsx
    FilterContext.tsx
    PhotosContext.tsx
```

**Provider Structure:**

```tsx
// app/(authenticated)/layout.tsx
<AlbumsProvider>
    <FilterProvider>
        <SelectedAlbumProvider>
            <PhotosProvider>
                {children}
            </PhotosProvider>
        </SelectedAlbumProvider>
    </FilterProvider>
</AlbumsProvider>
```

**Optimization Strategy:**

- Use `useMemo` and `useCallback` for context values to prevent re-renders
- Split contexts by concern (already done above)
- Components subscribe only to contexts they need

---

### Component Architecture

**Decision:** Colocation principle - components live with pages unless truly shared

**Rationale:**

- Helps agents and developers quickly find what needs changing
- Reduces cognitive load - only read relevant code
- Prevents premature abstraction
- Follows NextJS App Router best practices

**Structure:**

```
app/
├── (authenticated)/
│   ├── page.tsx                                    # Home: Album list
│   │   └── _components/
│   │       ├── AlbumCard.tsx                       # Used only on home
│   │       ├── AlbumFilter.tsx                     # Used only on home
│   │       ├── RandomHighlights.tsx                # Used only on home
│   │       ├── CreateAlbumDialog.tsx               # Used only on home
│   │       ├── DeleteAlbumDialog.tsx               # Used only on home
│   │       └── AlbumListSkeleton.tsx               # Used only on home
│   │
│   └── owners/[ownerId]/albums/[albumId]/
│       ├── layout.tsx                              # Album layout (for @modal slot)
│       │
│       ├── @modal/(.)photos/[photoId]/
│       │   └── page.tsx                            # Photo viewer modal
│       │       └── _components/
│       │           └── PhotoViewerModal.tsx        # Used only in modal
│       │
│       ├── page.tsx                                # Album view
│       │   └── _components/
│       │       ├── PhotoGrid.tsx                   # Used only in album view
│       │       ├── DayGroupHeader.tsx              # Used only in album view
│       │       ├── PhotoThumbnail.tsx              # Used only in album view
│       │       ├── EditAlbumDialog.tsx             # Used only in album view
│       │       ├── EditAlbumDatesDialog.tsx        # Used only in album view
│       │       ├── CreateAlbumFromRangeDialog.tsx  # Used only in album view
│       │       └── SharingDialog.tsx               # Used only in album view
│       │
│       └── photos/[photoId]/
│           └── page.tsx                            # Direct photo viewer
│               └── _components/
│                   └── PhotoViewerPage.tsx         # Used only on direct access
│
components/                                         # ONLY truly shared components
  ├── contexts/                                     # Shared state (used across app)
  │   ├── AlbumsContext.tsx
  │   ├── SelectedAlbumContext.tsx
  │   ├── FilterContext.tsx
  │   └── PhotosContext.tsx
  │
  ├── shared/                                       # Shared UI components
  │   ├── ErrorBoundary.tsx                         # Used across multiple pages
  │   ├── LoadingSkeleton.tsx                       # Used across multiple pages
  │   └── EmptyState.tsx                            # Used across multiple pages
  │
  └── UserInfo/                                     # Already exists
      └── index.tsx
```

**Rules:**

1. **Component lives in page folder if used by single page** - Use `_components/` subfolder
2. **Component moves to `components/shared/` only when used by 2+ pages** - Document which pages use it
3. **Contexts always in `components/contexts/`** - Used across application
4. **No premature abstraction** - Don't extract until duplication pain is real

---

### Routing Structure

**Decision:** Owner-based paths with NextJS parallel routes for photo modal interception

**Routes:**

```
/                                                   # Home: All albums (owned + shared)
/owners/[ownerId]/albums/[albumId]                 # Album view: Photo grid
/owners/[ownerId]/albums/[albumId]/photos/[photoId] # Photo viewer (modal or full page)
```

**Parallel Route for Photo Modal:**

When user clicks photo from album page:

- URL: `/owners/[ownerId]/albums/[albumId]/photos/[photoId]`
- Behavior: Opens as **modal** over album page (intercepted route)
- ESC key or close button dismisses modal, returns to album
- Browser back button closes modal

When user refreshes or direct link:

- Same URL: `/owners/[ownerId]/albums/[albumId]/photos/[photoId]`
- Behavior: Loads **full page** photo viewer (fallback route)

**Implementation:**

```
app/
└── (authenticated)/
    └── owners/[ownerId]/albums/[albumId]/
        ├── layout.tsx              # Renders {children} and {modal}
        ├── @modal/                 # Parallel route slot
        │   └── (.)photos/[photoId]/
        │       └── page.tsx        # Intercepted modal
        ├── page.tsx                # Album page
        └── photos/[photoId]/
            └── page.tsx            # Fallback full page
```

**Layout Implementation:**

```tsx
// owners/[ownerId]/albums/[albumId]/layout.tsx
export default function AlbumLayout({
                                        children,
                                        modal,
                                    }: {
    children: React.ReactNode
    modal: React.ReactNode
}) {
    return (
        <>
            {children} {/* Album page with photo grid */}
            {modal} {/* Photo viewer modal when intercepted */}
        </>
    )
}
```

**User Flow:**

1. User on `/` (home) sees album list
2. Clicks album → navigates to `/owners/123/albums/456`
3. Sees photo grid grouped by day
4. Clicks photo → URL becomes `/owners/123/albums/456/photos/789`
    - NextJS intercepts and shows modal
    - Album page stays visible behind modal
5. User presses ESC or clicks close → URL returns to `/owners/123/albums/456`
6. User refreshes on photo URL → Loads full page photo viewer

**Owner ID in URL:**

- All users can view albums from multiple owners (shared albums)
- Owner ID distinguishes albums with same folder name from different owners
- URL is shareable - recipient loads correct owner's album

---

### Image Optimization

**Decision:** Next.js Image component with custom loader for backend API integration

**Rationale:**

- Leverage Next.js automatic optimization and caching
- Map Next.js parameters to existing backend API quality system
- Progressive loading: blur → medium → high quality
- Responsive sizing handled automatically

**Custom Loader Implementation:**

```typescript
// libs/image-loader.ts
type ImageLoaderProps = {
    src: string      // mediaId
    width: number
    quality?: number
}

export function dphotoImageLoader({src, width, quality}: ImageLoaderProps) {
    const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL

    // Map Next.js width to backend quality levels
    let backendQuality: string
    if (width <= 40) {
        backendQuality = 'blur'      // Ultra-low for placeholder
    } else if (width <= 500) {
        backendQuality = 'medium'    // Reasonable quality for grid
    } else if (width <= 1200) {
        backendQuality = 'high'      // High quality for viewer
    } else {
        backendQuality = 'full'      // Original quality for zoom
    }

    return `${apiBaseUrl}/api/v1/media/${src}/image?quality=${backendQuality}&width=${width}`
}
```

**Usage:**

```tsx
import Image from 'next/image'
import {dphotoImageLoader} from '@/libs/image-loader'

<
Image
loader = {dphotoImageLoader}
src = {mediaId}
alt = {photoAlt}
placeholder = "blur"
blurDataURL = {blurThumbnailUrl}
sizes = "(max-width: 600px) 50vw, (max-width: 960px) 33vw, 25vw"
fill
/ >
```

**Backend API Requirements:**

The backend REST API must support image quality endpoints:

**Required Quality Levels:**

1. **`blur` / `thumbnail`**
    - Purpose: Instant blur placeholder
    - Size: 20-40px longest dimension
    - File size: <2KB
    - Format: JPEG heavy compression or data URL

2. **`medium`**
    - Purpose: First clear image (grid display)
    - File size: 50-150KB
    - Loading target: <500ms on 3G

3. **`high`**
    - Purpose: Full-screen viewer
    - Quality: Near-original
    - Loaded on-demand

4. **`full`**
    - Purpose: Zoom interactions
    - Quality: Original
    - Loaded on-demand

**API Endpoint Pattern:**

```
GET /api/v1/media/{mediaId}/image?quality={quality}&width={width}

Parameters:
- quality: "blur" | "medium" | "high" | "full"
- width: number (optional, for responsive sizing)
```

**Image Sizes by Context:**

| Context                      | Display Size     | Quality | Requested Width | Progressive Loading |
|------------------------------|------------------|---------|-----------------|---------------------|
| Album Card Samples           | 80-120px square  | medium  | 240px (2x)      | blur → medium       |
| Random Highlights            | 150-200px height | medium  | 600px (2x)      | blur → medium       |
| Photo Grid (Mobile)          | ~180px width     | medium  | 360px (2x)      | blur → medium       |
| Photo Grid (Tablet)          | ~220px width     | medium  | 440px (2x)      | blur → medium       |
| Photo Grid (Desktop)         | ~250px width     | medium  | 500px (2x)      | blur → medium       |
| Full-Screen Viewer (Mobile)  | Viewport         | high    | 1080px max      | medium → high       |
| Full-Screen Viewer (Tablet)  | Viewport         | high    | 1536px max      | medium → high       |
| Full-Screen Viewer (Desktop) | Viewport         | high    | 2560px max      | medium → high       |
| Zoom Interaction             | Full resolution  | full    | Original        | high → full         |

**Caching Strategy:**

- Browser caching handled by Next.js Image component
- Backend serves images with `Cache-Control: max-age=31536000, immutable`
- Images are immutable (mediaId changes if content changes)
- No additional client-side caching needed

**Image Configuration:**

```typescript
// next.config.ts
images: {
    loader: 'custom',
        loaderFile
:
    './libs/image-loader.ts',
        remotePatterns
:
    [
        {
            protocol: 'https',
            hostname: '**', // Configured API domain
        }
    ]
}
```

---

### Layout Architecture

**Decision:** Material UI sx prop with theme breakpoints for responsive layouts

**Rationale:**

- Consistent with Material UI ecosystem
- Type-safe breakpoint system
- Integrated with theme configuration
- Works seamlessly with all MUI components

**Breakpoint Configuration:**

```typescript
// Material UI theme breakpoints (default)
breakpoints: {
    values: {
        xs: 0,      // Mobile
            sm
    :
        600,    // Tablet start
            md
    :
        960,    // Desktop start
            lg
    :
        1280,   // Large desktop
            xl
    :
        1920    // Extra large
    }
}
```

**Usage Patterns:**

**Photo Grid Example:**

```tsx
<Box sx={{
    display: 'grid',
    gridTemplateColumns: {
        xs: 'repeat(2, 1fr)',    // Mobile: 2 columns
        sm: 'repeat(3, 1fr)',    // Tablet: 3 columns
        md: 'repeat(4, 1fr)',    // Desktop: 4 columns
    },
    gap: {xs: 1, sm: 2, md: 3},
    padding: {xs: 1, sm: 2, md: 3},
}}>
    {photos.map(photo => <PhotoThumbnail key={photo.id} photo={photo}/>)}
</Box>
```

**Album Card List Example:**

```tsx
<Box sx={{
    display: 'grid',
    gridTemplateColumns: {
        xs: '1fr',              // Mobile: 1 column (stacked)
        sm: 'repeat(2, 1fr)',   // Tablet: 2 columns
        md: 'repeat(3, 1fr)',   // Desktop: 3 columns
    },
    gap: 2,
}}>
    {albums.map(album => <AlbumCard key={album.id} album={album}/>)}
</Box>
```

**Responsive Typography:**

```tsx
<Typography variant="h4" sx={{
    fontSize: {
        xs: '1.5rem',   // Mobile
        sm: '2rem',     // Tablet
        md: '2.5rem',   // Desktop
    }
}}>
    Album Title
</Typography>
```

**Container Width:**

```tsx
<Container maxWidth="lg" sx={{
    px: {xs: 2, sm: 3, md: 4}  // Responsive padding
}}>
    {/* Content */}
</Container>
```

---

### Error Boundaries

**Decision:** NextJS page-level error boundaries using error.tsx files

**Rationale:**

- Built into NextJS App Router
- Automatic error isolation by route
- Simple to implement and understand for agents
- Sufficient granularity for application needs

**Error Boundary Structure:**

```
app/
├── error.tsx                                   # Root catch-all
└── (authenticated)/
    ├── error.tsx                               # Authenticated routes catch-all
    ├── page.tsx                                # Home page
    └── owners/[ownerId]/albums/[albumId]/
        ├── error.tsx                           # Album-specific errors
        ├── page.tsx
        └── photos/[photoId]/
            └── error.tsx                       # Photo viewer errors
```

**Error Component Implementation:**

```tsx
// app/(authenticated)/error.tsx
'use client'

export default function AuthenticatedError({
                                               error,
                                               reset,
                                           }: {
    error: Error & { digest?: string }
    reset: () => void
}) {
    return (
        <Box sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '50vh',
            gap: 2,
        }}>
            <Typography variant="h5">Something went wrong</Typography>
            <Typography color="text.secondary">
                {error.message || 'An unexpected error occurred'}
            </Typography>
            <Button variant="contained" onClick={reset}>
                Try Again
            </Button>
            <Button variant="text" href="/">
                Return Home
            </Button>
        </Box>
    )
}
```

**Error Types Handled:**

1. **Network Errors** - API failures, timeouts
2. **Permission Errors** - Accessing albums without permission
3. **Not Found Errors** - Album or photo doesn't exist (use not-found.tsx)
4. **Rendering Errors** - Component failures

**Not Found Handling:**

```
app/
└── (authenticated)/
    └── owners/[ownerId]/albums/[albumId]/
        ├── not-found.tsx               # Album not found
        └── photos/[photoId]/
            └── not-found.tsx           # Photo not found
```

```tsx
// app/(authenticated)/owners/[ownerId]/albums/[albumId]/not-found.tsx
export default function AlbumNotFound() {
    return (
        <Box sx={{textAlign: 'center', py: 8}}>
            <Typography variant="h5">Album Not Found</Typography>
            <Typography color="text.secondary">
                This album doesn't exist or you don't have access to it.
            </Typography>
            <Button href="/" sx={{mt: 2}}>
                Return to Albums
            </Button>
        </Box>
    )
}
```

**Error Recovery:**

- **Try Again** button calls `reset()` to re-render component
- **Return Home** link navigates to root
- Error boundaries prevent entire app from crashing
- User sees contextual error message based on route

---

### Data Fetching

**Decision:** Server Components for initial data, Client Components for interactions

**Rationale:**

- Fast initial page load without loading spinners
- Reduced client bundle size
- SEO-friendly approach (good practice even if not required)
- Leverage NextJS App Router server/client component split

**Pattern:**

**Server Component fetches initial data:**

```tsx
// app/(authenticated)/page.tsx - Server Component
import {AlbumListClient} from './_components/AlbumListClient'

async function fetchAlbums() {
    const res = await fetch(`${process.env.API_BASE_URL}/api/v1/albums`, {
        headers: {
            // Include auth from server-side context
        }
    })
    return res.json()
}

export default async function HomePage() {
    const albums = await fetchAlbums()

    return <AlbumListClient initialAlbums={albums}/>
}
```

**Client Component manages interactions:**

```tsx
// app/(authenticated)/_components/AlbumListClient.tsx
'use client'

import {useState} from 'react'
import {useAlbums} from '@/components/contexts/AlbumsContext'

export function AlbumListClient({initialAlbums}) {
    const {albums, setAlbums} = useAlbums()

    // Initialize context with server data
    useEffect(() => {
        setAlbums(initialAlbums)
    }, [initialAlbums])

    // Handle filtering, refreshing on client
    const handleFilter = (filter) => {
        // Client-side filtering or refetch
    }

    return (
        <>
            <AlbumFilter onFilterChange={handleFilter}/>
            {albums.map(album => <AlbumCard key={album.id} album={album}/>)}
        </>
    )
}
```

**Data Flow:**

1. **Initial Load (Server)**
    - Server Component fetches albums from API
    - Passes initial data as props to Client Component
    - Fast initial render, no loading spinner

2. **User Interactions (Client)**
    - Filtering: Client-side or refetch from API
    - Album selection: Update context state
    - Creating album: POST to API, update context
    - Deleting album: DELETE to API, update context

3. **Context Integration**
    - Server data initializes Context state
    - Context manages updates from user interactions
    - Context provides state to all components

**Fetching Strategy by Route:**

| Route                     | Component Type  | Data Fetch             | Update Strategy          |
|---------------------------|-----------------|------------------------|--------------------------|
| `/` (Home)                | Server → Client | Server: initial albums | Client: filters, refresh |
| `/owners/.../albums/[id]` | Server → Client | Server: album + photos | Client: updates, modals  |
| `/owners/.../photos/[id]` | Server → Client | Server: photo details  | Client: navigation       |

**API Request Utilities:**

```typescript
// libs/requests/api-client.ts
export async function fetchFromAPI(endpoint: string, options?: RequestInit) {
    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL
    const res = await fetch(`${baseUrl}${endpoint}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers,
        },
    })

    if (!res.ok) {
        throw new Error(`API Error: ${res.statusText}`)
    }

    return res.json()
}
```

**Caching:**

- Server Component fetches are cached by NextJS by default
- Use `revalidate` or `cache: 'no-store'` for dynamic data
- Client-side updates don't trigger server refetch
- Manual refresh (F5) re-runs server fetch

**Example with Revalidation:**

```tsx
// Revalidate albums every 5 minutes
export const revalidate = 300

async function fetchAlbums() {
    const res = await fetch(...)
    return res.json()
}
```

---

### Decision Impact Analysis

**Implementation Sequence:**

1. **Material UI Setup** (First - affects all components)
    - Remove Tailwind dependencies
    - Install Material UI packages
    - Configure theme with dark mode and brand color (#185986)
    - Set up breakpoint system

2. **Context Architecture** (Second - affects state management)
    - Create AlbumsContext, FilterContext, SelectedAlbumContext, PhotosContext
    - Set up providers in authenticated layout
    - Define context interfaces and actions

3. **Routing Structure** (Third - defines page organization)
    - Create authenticated layout with nav/logo/profile
    - Set up home page route
    - Create owner/album/photo route structure
    - Implement parallel route for photo modal

4. **Image Loader** (Fourth - needed for all image displays)
    - Implement custom image loader
    - Configure Next.js image settings
    - Document backend API requirements

5. **Layout Patterns** (Fifth - defines responsive structure)
    - Create base layouts with MUI sx breakpoints
    - Define grid systems for album cards and photo grids
    - Implement responsive patterns

6. **Error Boundaries** (Sixth - wrap routes)
    - Create error.tsx files for each route level
    - Implement not-found.tsx for missing resources
    - Test error recovery flows

7. **Data Fetching** (Seventh - populate UI with data)
    - Implement Server Components for initial loads
    - Create Client Components for interactions
    - Wire up Context with server data
    - Build API client utilities

8. **Feature Components** (Eighth - build capabilities)
    - Album list and filtering
    - Photo grid and viewer
    - Album management dialogs
    - Sharing dialogs

**Cross-Component Dependencies:**

```
Material UI Theme
  ↓
  ├─→ Layout Architecture (uses theme breakpoints)
  ├─→ Error Boundaries (uses theme components)
  └─→ All Feature Components (use MUI components)

Context Architecture
  ↓
  ├─→ Data Fetching (initializes contexts with server data)
  └─→ Feature Components (subscribe to contexts)

Image Loader
  ↓
  └─→ All Image-Displaying Components (album cards, photo grid, viewer)

Routing Structure
  ↓
  ├─→ Layout Architecture (defines where layouts apply)
  ├─→ Error Boundaries (defines error boundary scope)
  ├─→ Data Fetching (defines server vs client components)
  └─→ Component Organization (defines colocation structure)
```

**Technology Dependencies:**

- Next.js 16.1.1 (already installed)
- React 19.2.3 (already installed)
- TypeScript 5.x (already installed)
- Material UI 6.x (to be installed)
- @emotion/react 11.x (to be installed)
- @emotion/styled 11.x (to be installed)

**Backend Dependencies:**

- REST API must support image quality parameters: `blur`, `medium`, `high`, `full`
- REST API must support width parameter for responsive sizing
- Images must be served with appropriate cache headers (immutable)
- Authentication already implemented (Google OAuth via backend)

## Implementation Patterns & Consistency Rules

_These patterns ensure consistent code structure across all AI agent implementations. Following these rules prevents conflicts and maintains architectural
integrity._

### File & Component Naming

**Rule:** Use PascalCase for all React components and their files

**Rationale:**

- Standard React convention for component names
- Immediate visual distinction between components and utilities
- TypeScript/IDE autocomplete works better with consistent casing

**Examples:**

```
✅ Correct:
AlbumCard.tsx
PhotoGrid.tsx
CreateAlbumDialog.tsx
AlbumsContext.tsx

❌ Incorrect:
albumCard.tsx
photo-grid.tsx
createAlbumDialog.tsx
```

**Non-Component Files:**

```
✅ Correct:
api-client.ts          # Utility files use kebab-case
image-loader.ts
types/api.ts

❌ Incorrect:
ApiClient.ts
imageLoader.ts
```

---

### Plural vs Singular Naming

**Rule:** Use semantic naming - plural for collections, singular for single items

**Rationale:**

- Semantic clarity about what the variable/context contains
- Prevents confusion when reading code
- Natural language patterns

**Context Naming:**

```tsx
✅ Correct:
    AlbumsContext          // Contains array of albums
SelectedAlbumContext   // Contains single album
PhotosContext          // Contains array of photos
FilterContext          // Contains filter state (not a collection)

❌ Incorrect:
    AlbumContext           // Ambiguous - one or many?
AlbumListContext       // Redundant
PhotoContext           // Unclear if single or multiple
```

**Variable Naming:**

```tsx
✅ Correct:
    const albums = []
const selectedAlbum = {}
const photos = []
const activeFilter = 'my-albums'

❌ Incorrect:
    const album = []        // Array should be plural
const selectedAlbums = {} // Single item should be singular
```

---

### Function Naming Conventions

**Rule:** Prefix functions by purpose: `fetch*` for API calls, `get*` for state access, `handle*` for event handlers

**Rationale:**

- Clear intent from function name alone
- Easy to identify async operations
- Distinguishes data sources (API vs state)

**API Functions:**

```typescript
✅ Correct:
    async function fetchAlbums() { ...
    }

async function fetchPhotos(albumId: string) { ...
}

async function fetchAlbumDetails(albumId: string) { ...
}

❌ Incorrect:
    async function getAlbums() { ...
    }      // Ambiguous - API or state?
async function loadAlbums() { ...
}     // Inconsistent prefix
async function albums() { ...
}         // Missing action verb
```

**State Access Functions:**

```typescript
✅ Correct:
    function getAlbums() {
        return albums
    }

function getSelectedAlbum() {
    return selectedAlbum
}

function getFilterState() {
    return filter
}

❌ Incorrect:
    function albums() { ...
    }               // Missing verb
function fetchAlbums() { ...
}          // Implies API call
function retrieveAlbums() { ...
}       // Inconsistent verb
```

**Event Handlers:**

```typescript
✅ Correct:
    function handleFilterChange(filter: string) { ...
    }

function handleAlbumClick(albumId: string) { ...
}

function handlePhotoDelete(photoId: string) { ...
}

❌ Incorrect:
    function onFilterChange() { ...
    }       // Reserve 'on' for props
function filterChange() { ...
}         // Missing 'handle'
function changeFilter() { ...
}         // Verb-first (less React-idiomatic)
```

**Component Props (callbacks):**

```typescript
✅ Correct:
    interface AlbumCardProps {
        onAlbumClick: (id: string) => void     // Use 'on' prefix for props
        onDeleteClick: () => void
    }

// Usage:
<AlbumCard
    onAlbumClick = {handleAlbumClick}        // Pass handler to 'on' prop
onDeleteClick = {handleDeleteClick}
/>
```

---

### Test File Location

**Rule:** Co-locate tests in `__tests__/` subfolder within component directory

**Rationale:**

- Tests stay close to components (easy to find)
- Organized in dedicated folder (clear separation)
- Standard pattern recognized by test runners
- Works with both page-level and shared components

**Structure:**

```
app/(authenticated)/page.tsx
app/(authenticated)/_components/
  AlbumCard.tsx
  AlbumFilter.tsx
  __tests__/
    AlbumCard.test.tsx
    AlbumFilter.test.tsx
    test-utils.ts

components/shared/
  EmptyState.tsx
  LoadingSkeleton.tsx
  __tests__/
    EmptyState.test.tsx
    LoadingSkeleton.test.tsx
```

**Test File Naming:**

```tsx
✅ Correct:
    AlbumCard.test.tsx            // Matches component name
PhotoGrid.test.tsx

❌ Incorrect:
    AlbumCard.spec.tsx            // Inconsistent suffix
albumCard.test.tsx            // Wrong casing
AlbumCardTest.tsx             // 'Test' should be after dot
```

---

### Type Definitions

**Rule:** Hybrid approach - shared API types centralized, component-specific types co-located

**Rationale:**

- API types are shared contracts (central location)
- Component-specific types are implementation details (co-locate)
- Balance between organization and proximity

**Shared API Types:**

```typescript
// types/api.ts
export interface Album {
    id: string
    name: string
    start: string
    end: string
    ownerId: string
    folderName?: string
    count: number
}

export interface Photo {
    id: string
    mediaId: string
    albumId: string
    capturedAt: string
}

export interface Owner {
    id: string
    name: string
    email: string
    pictureUrl?: string
}

export interface SharingInfo {
    albumId: string
    sharedWith: Array<{
        email: string
        name?: string
        pictureUrl?: string
    }>
}
```

**Component-Specific Types:**

```typescript
// app/(authenticated)/_components/AlbumCard.tsx
import {Album, Owner} from '@/types/api'

// Component-specific props and state
interface AlbumCardProps {
    album: Album
    owner: Owner
    isSelected: boolean
    onAlbumClick: (id: string) => void
    onDeleteClick?: (id: string) => void
}

// Internal component state types
interface AlbumCardState {
    isHovered: boolean
    showActions: boolean
}
```

**Context Types:**

```typescript
// components/contexts/AlbumsContext.tsx
import {Album} from '@/types/api'

interface AlbumsContextValue {
    albums: Album[]
    loading: boolean
    error: Error | null
    actions: {
        fetchAlbums: () => Promise<void>
        refreshAlbums: () => Promise<void>
        addAlbum: (album: Album) => void
        removeAlbum: (id: string) => void
    }
}
```

**File Organization:**

```
types/
  api.ts              # Shared API response types
  
components/contexts/
  AlbumsContext.tsx   # Context-specific types inline
  
app/(authenticated)/_components/
  AlbumCard.tsx       # Component-specific types inline
```

---

### Loading & Error State Patterns

**Rule:** Use object pattern `{ loading, error, data }` for async state

**Rationale:**

- Single source of truth for async operation state
- Prevents impossible states (loading + error simultaneously)
- Consistent pattern across all data fetching
- TypeScript discriminated unions for type safety

**Context Pattern:**

```typescript
// components/contexts/AlbumsContext.tsx
interface AlbumsContextValue {
    albums: Album[]
    loading: boolean
    error: Error | null
    actions: {
        fetchAlbums: () => Promise<void>
    }
}

function AlbumsProvider({children}: { children: React.ReactNode }) {
    const [albums, setAlbums] = useState<Album[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<Error | null>(null)

    const fetchAlbums = async () => {
        setLoading(true)
        setError(null)
        try {
            const data = await fetchFromAPI('/api/v1/albums')
            setAlbums(data)
        } catch (err) {
            setError(err as Error)
        } finally {
            setLoading(false)
        }
    }

    return (
        <AlbumsContext.Provider value = {
    {
        albums, loading, error, actions
    :
        {
            fetchAlbums
        }
    }
}>
    {
        children
    }
    </AlbumsContext.Provider>
)
}
```

**Component Usage:**

```tsx
function AlbumList() {
    const {albums, loading, error} = useAlbums()

    if (loading) return <LoadingSkeleton/>
    if (error) return <ErrorMessage error={error}/>
    if (albums.length === 0) return <EmptyState/>

    return (
        <Box>
            {albums.map(album => <AlbumCard key={album.id} album={album}/>)}
        </Box>
    )
}
```

**TypeScript Discriminated Union (Advanced):**

```typescript
// Alternative: Discriminated union for mutually exclusive states
type AsyncState<T> =
    | { status: 'idle' }
    | { status: 'loading' }
    | { status: 'success'; data: T }
    | { status: 'error'; error: Error }

// Usage:
const [albumsState, setAlbumsState] = useState<AsyncState<Album[]>>({status: 'idle'})

// TypeScript narrows type based on status
if (albumsState.status === 'success') {
    // albumsState.data is available here
    albumsState.data.map(...)
}
```

**Recommendation:** Use simple `{ loading, error, data }` pattern for this project. Discriminated unions add complexity without significant benefit for our use
case.

---

### Context Update Patterns

**Rule:** Use named action methods, not direct setters

**Rationale:**

- Encapsulates update logic in context
- Prevents invalid state updates
- Clear intent from action names
- Easy to add side effects later

**❌ Incorrect - Direct Setters:**

```typescript
// DON'T expose raw setters
interface AlbumsContextValue {
    albums: Album[]
    setAlbums: (albums: Album[]) => void     // Direct setter exposed
}

// Components can cause invalid states:
setAlbums([])  // Accidentally clear albums without loading state
```

**✅ Correct - Named Actions:**

```typescript
// DO provide named action methods
interface AlbumsContextValue {
    albums: Album[]
    loading: boolean
    error: Error | null
    actions: {
        fetchAlbums: () => Promise<void>
        refreshAlbums: () => Promise<void>
        addAlbum: (album: Album) => void
        updateAlbum: (id: string, updates: Partial<Album>) => void
        removeAlbum: (id: string) => void
    }
}

// Usage in components:
const {albums, actions} = useAlbums()

// Clear intent and encapsulated logic
await actions.fetchAlbums()
actions.addAlbum(newAlbum)
actions.removeAlbum(albumId)
```

**Implementation:**

```typescript
function AlbumsProvider({children}: { children: React.ReactNode }) {
    const [albums, setAlbums] = useState<Album[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<Error | null>(null)

    const actions = useMemo(() => ({
        fetchAlbums: async () => {
            setLoading(true)
            setError(null)
            try {
                const data = await fetchFromAPI('/api/v1/albums')
                setAlbums(data)
            } catch (err) {
                setError(err as Error)
            } finally {
                setLoading(false)
            }
        },

        addAlbum: (album: Album) => {
            setAlbums(prev => [...prev, album])
        },

        updateAlbum: (id: string, updates: Partial<Album>) => {
            setAlbums(prev => prev.map(album =>
                album.id === id ? {...album, ...updates} : album
            ))
        },

        removeAlbum: (id: string) => {
            setAlbums(prev => prev.filter(album => album.id !== id))
        },

        refreshAlbums: async () => {
            // Same as fetchAlbums, but doesn't show loading (silent refresh)
            try {
                const data = await fetchFromAPI('/api/v1/albums')
                setAlbums(data)
            } catch (err) {
                // Silent error or toast notification
                console.error('Failed to refresh albums:', err)
            }
        },
    }), [])

    const value = useMemo(
        () => ({albums, loading, error, actions}),
        [albums, loading, error, actions]
    )

    return (
        <AlbumsContext.Provider value = {value} >
            {children}
            < /AlbumsContext.Provider>
    )
}
```

**Benefits:**

1. **Encapsulation** - Update logic lives in context, not scattered across components
2. **Consistency** - All updates follow same pattern
3. **Side Effects** - Easy to add logging, analytics, optimistic updates
4. **Type Safety** - Action parameters are typed and validated
5. **Testing** - Mock actions object instead of tracking setter calls

---

### Error Object Structure

**Rule:** Use custom error objects with `{ message, code?, details? }` structure

**Rationale:**

- Structured error information for UI rendering
- Error codes enable specific handling (retry, redirect, etc.)
- Additional details for debugging without cluttering message
- Consistent error shape across application

**Error Type Definition:**

```typescript
// types/errors.ts
export interface AppError {
    message: string              // User-friendly message
    code?: string                // Machine-readable code
    details?: Record<string, unknown>  // Additional context
    originalError?: Error        // Original error for debugging
}

export class DPhotoError extends Error implements AppError {
    code?: string
    details?: Record<string, unknown>
    originalError?: Error

    constructor(message: string, options?: {
        code?: string
        details?: Record<string, unknown>
        originalError?: Error
    }) {
        super(message)
        this.name = 'DPhotoError'
        this.code = options?.code
        this.details = options?.details
        this.originalError = options?.originalError
    }
}
```

**Error Codes:**

```typescript
// types/errors.ts
export const ErrorCode = {
    // Network errors
    NETWORK_ERROR: 'NETWORK_ERROR',
    TIMEOUT: 'TIMEOUT',

    // Permission errors
    UNAUTHORIZED: 'UNAUTHORIZED',
    FORBIDDEN: 'FORBIDDEN',

    // Not found errors
    ALBUM_NOT_FOUND: 'ALBUM_NOT_FOUND',
    PHOTO_NOT_FOUND: 'PHOTO_NOT_FOUND',

    // Validation errors
    INVALID_DATE_RANGE: 'INVALID_DATE_RANGE',
    INVALID_ALBUM_NAME: 'INVALID_ALBUM_NAME',

    // Server errors
    SERVER_ERROR: 'SERVER_ERROR',
    UNKNOWN_ERROR: 'UNKNOWN_ERROR',
} as const

export type ErrorCodeType = typeof ErrorCode[keyof typeof ErrorCode]
```

**API Client Error Handling:**

```typescript
// libs/requests/api-client.ts
export async function fetchFromAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}${endpoint}`, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options?.headers,
            },
        })

        if (!res.ok) {
            const errorData = await res.json().catch(() => ({}))

            throw new DPhotoError(
                errorData.message || `API request failed: ${res.statusText}`,
                {
                    code: mapHTTPStatusToErrorCode(res.status),
                    details: {
                        status: res.status,
                        endpoint,
                        ...errorData,
                    },
                }
            )
        }

        return res.json()
    } catch (error) {
        if (error instanceof DPhotoError) {
            throw error
        }

        // Network errors or other exceptions
        throw new DPhotoError(
            'Failed to connect to server. Please check your connection.',
            {
                code: ErrorCode.NETWORK_ERROR,
                originalError: error as Error,
            }
        )
    }
}

function mapHTTPStatusToErrorCode(status: number): ErrorCodeType {
    switch (status) {
        case 401:
            return ErrorCode.UNAUTHORIZED
        case 403:
            return ErrorCode.FORBIDDEN
        case 404:
            return ErrorCode.ALBUM_NOT_FOUND
        case 408:
            return ErrorCode.TIMEOUT
        case 500:
            return ErrorCode.SERVER_ERROR
        default:
            return ErrorCode.UNKNOWN_ERROR
    }
}
```

**Error Display Component:**

```tsx
// components/shared/ErrorMessage.tsx
import {AppError, ErrorCode} from '@/types/errors'

interface ErrorMessageProps {
    error: AppError | Error
    onRetry?: () => void
}

export function ErrorMessage({error, onRetry}: ErrorMessageProps) {
    const appError = error as AppError

    // Customize message based on error code
    const getActionButton = () => {
        if (appError.code === ErrorCode.NETWORK_ERROR && onRetry) {
            return <Button onClick={onRetry}>Retry</Button>
        }
        if (appError.code === ErrorCode.UNAUTHORIZED) {
            return <Button href="/auth/logout">Sign In Again</Button>
        }
        if (appError.code === ErrorCode.ALBUM_NOT_FOUND) {
            return <Button href="/">Return to Albums</Button>
        }
        return onRetry ? <Button onClick={onRetry}>Try Again</Button> : null
    }

    return (
        <Box sx={{textAlign: 'center', py: 4}}>
            <Typography variant="h6" color="error">
                {appError.message || 'An error occurred'}
            </Typography>
            {appError.details?.status && (
                <Typography variant="caption" color="text.secondary">
                    Error code: {appError.code}
                </Typography>
            )}
            <Box sx={{mt: 2}}>
                {getActionButton()}
            </Box>
        </Box>
    )
}
```

**Usage in Context:**

```typescript
// components/contexts/AlbumsContext.tsx
const [error, setError] = useState<AppError | null>(null)

const fetchAlbums = async () => {
    setLoading(true)
    setError(null)
    try {
        const data = await fetchFromAPI<Album[]>('/api/v1/albums')
        setAlbums(data)
    } catch (err) {
        setError(err as AppError)
    } finally {
        setLoading(false)
    }
}
```

**Benefits:**

1. **User-Friendly Messages** - Clear error messages for users
2. **Specific Handling** - Error codes enable custom recovery flows
3. **Debugging** - Original error and details preserved
4. **Consistency** - All errors follow same structure
5. **Type Safety** - TypeScript ensures correct error handling

---

### Pattern Summary

**Quick Reference for AI Agents:**

| Pattern             | Rule                   | Example                                       |
|---------------------|------------------------|-----------------------------------------------|
| **Component Files** | PascalCase             | `AlbumCard.tsx`, `PhotoGrid.tsx`              |
| **Utility Files**   | kebab-case             | `api-client.ts`, `image-loader.ts`            |
| **Collections**     | Plural                 | `albums`, `photos`, `AlbumsContext`           |
| **Single Items**    | Singular               | `selectedAlbum`, `SelectedAlbumContext`       |
| **API Functions**   | `fetch*` prefix        | `fetchAlbums()`, `fetchPhotos()`              |
| **State Access**    | `get*` prefix          | `getAlbums()`, `getSelectedAlbum()`           |
| **Event Handlers**  | `handle*` prefix       | `handleClick()`, `handleFilterChange()`       |
| **Callback Props**  | `on*` prefix           | `onAlbumClick`, `onDeleteClick`               |
| **Tests**           | `__tests__/` subfolder | `__tests__/AlbumCard.test.tsx`                |
| **API Types**       | `types/api.ts`         | `Album`, `Photo`, `Owner`                     |
| **Component Types** | Co-located inline      | `AlbumCardProps`, `AlbumCardState`            |
| **Async State**     | Object pattern         | `{ loading, error, data }`                    |
| **Context Updates** | Named actions          | `actions.fetchAlbums()`, `actions.addAlbum()` |
| **Errors**          | Custom structure       | `{ message, code?, details? }`                |

**Consistency Checklist:**

- [ ] Component files use PascalCase
- [ ] Utility files use kebab-case
- [ ] Collections are plural, single items are singular
- [ ] API calls use `fetch*` prefix
- [ ] Event handlers use `handle*` prefix
- [ ] Tests are in `__tests__/` subfolders
- [ ] API types are in `types/api.ts`
- [ ] Component types are co-located
- [ ] Async state uses `{ loading, error, data }` pattern
- [ ] Context provides named actions, not setters
- [ ] Errors use `{ message, code, details }` structure

## Project Structure & Boundaries

### Complete Project Directory Structure

```
web-nextjs/
├── README.md
├── package.json
├── package-lock.json
├── next.config.ts
├── tsconfig.json
├── .env.local
├── .env.example
├── .gitignore
├── .eslintrc.json
├── vitest.config.ts
├── playwright.config.ts
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
│
├── app/
│   ├── layout.tsx                              # Root layout with MUI theme provider
│   ├── globals.css                             # Minimal global styles
│   ├── error.tsx                               # Root error boundary
│   ├── not-found.tsx                           # Root 404 page
│   │
│   ├── auth/                                   # Authentication flows (already exists)
│   │   ├── callback/
│   │   │   └── page.tsx
│   │   ├── error/
│   │   │   └── page.tsx
│   │   └── logout/
│   │       └── page.tsx
│   │
│   └── (authenticated)/                        # Protected routes group
│       ├── layout.tsx                          # Authenticated layout with nav + context providers
│       ├── error.tsx                           # Authenticated routes error boundary
│       │
│       ├── page.tsx                            # HOME: Album list (Server Component)
│       └── _components/                        # Home page components (colocation)
│           ├── AlbumListClient.tsx             # Client wrapper for server data
│           ├── AlbumCard.tsx                   # Individual album card
│           ├── AlbumFilter.tsx                 # Filter UI (My/All/By Owner)
│           ├── RandomHighlights.tsx            # Random photo highlights section
│           ├── CreateAlbumDialog.tsx           # Dialog to create new album
│           ├── DeleteAlbumDialog.tsx           # Confirmation dialog for delete
│           ├── AlbumListSkeleton.tsx           # Loading skeleton for album list
│           ├── EmptyAlbumList.tsx              # Empty state when no albums
│           └── __tests__/                      # Co-located tests
│               ├── AlbumCard.test.tsx
│               ├── AlbumFilter.test.tsx
│               ├── CreateAlbumDialog.test.tsx
│               └── test-utils.ts
│
│       └── owners/[ownerId]/albums/[albumId]/  # Album detail route
│           ├── layout.tsx                      # Album layout (renders {children} + {modal})
│           ├── error.tsx                       # Album-specific errors
│           ├── not-found.tsx                   # Album not found page
│           │
│           ├── page.tsx                        # ALBUM VIEW: Photo grid (Server Component)
│           ├── _components/                    # Album page components (colocation)
│           │   ├── AlbumViewClient.tsx         # Client wrapper for server data
│           │   ├── PhotoGrid.tsx               # Photo grid layout
│           │   ├── PhotoThumbnail.tsx          # Individual photo in grid
│           │   ├── DayGroupHeader.tsx          # Date header for photo groups
│           │   ├── EditAlbumDialog.tsx         # Dialog to edit album name
│           │   ├── EditAlbumDatesDialog.tsx    # Dialog with visual date picker
│           │   ├── CreateAlbumFromRangeDialog.tsx  # Create album from media range
│           │   ├── SharingDialog.tsx           # Share/revoke access dialog
│           │   ├── PhotoGridSkeleton.tsx       # Loading skeleton for photos
│           │   └── __tests__/
│           │       ├── PhotoGrid.test.tsx
│           │       ├── PhotoThumbnail.test.tsx
│           │       ├── EditAlbumDatesDialog.test.tsx
│           │       └── SharingDialog.test.tsx
│           │
│           ├── @modal/                         # Parallel route slot for modal
│           │   └── (.)photos/[photoId]/
│           │       └── page.tsx                # PHOTO MODAL: Intercepted route
│           │           └── _components/
│           │               ├── PhotoViewerModal.tsx      # Modal photo viewer
│           │               ├── PhotoNavigation.tsx       # Prev/Next controls
│           │               ├── PhotoMetadata.tsx         # Photo info display
│           │               └── __tests__/
│           │                   └── PhotoViewerModal.test.tsx
│           │
│           └── photos/[photoId]/
│               ├── page.tsx                    # PHOTO VIEWER: Full page (fallback)
│               ├── error.tsx                   # Photo-specific errors
│               ├── not-found.tsx               # Photo not found page
│               └── _components/
│                   ├── PhotoViewerPage.tsx     # Full-page photo viewer
│                   ├── PhotoZoomControls.tsx   # Zoom in/out controls
│                   └── __tests__/
│                       └── PhotoViewerPage.test.tsx
│
├── components/                                 # ONLY truly shared components
│   │
│   ├── contexts/                               # Shared React contexts
│   │   ├── AlbumsContext.tsx                   # Albums list state + actions
│   │   ├── SelectedAlbumContext.tsx            # Current album state
│   │   ├── FilterContext.tsx                   # Filter state (My/All/Owner)
│   │   ├── PhotosContext.tsx                   # Photos for selected album
│   │   └── __tests__/
│   │       ├── AlbumsContext.test.tsx
│   │       ├── FilterContext.test.tsx
│   │       └── PhotosContext.test.tsx
│   │
│   ├── shared/                                 # Shared UI components (used 2+ places)
│   │   ├── ErrorMessage.tsx                    # Error display with recovery actions
│   │   ├── LoadingSkeleton.tsx                 # Generic loading skeleton
│   │   ├── EmptyState.tsx                      # Generic empty state display
│   │   └── __tests__/
│   │       ├── ErrorMessage.test.tsx
│   │       └── EmptyState.test.tsx
│   │
│   ├── UserInfo/                               # Already exists
│   │   └── index.tsx                           # User profile display
│   │
│   └── theme/                                  # Material UI theme configuration
│       ├── ThemeProvider.tsx                   # MUI theme provider wrapper
│       ├── theme.ts                            # Theme configuration (dark + #185986)
│       └── __tests__/
│           └── theme.test.ts
│
├── libs/                                       # Utility libraries
│   ├── image-loader.ts                         # Custom Next.js image loader
│   │
│   ├── requests/                               # API request layer (already exists)
│   │   ├── api-client.ts                       # Enhanced with error handling
│   │   └── __tests__/
│   │       └── api-client.test.ts
│   │
│   ├── security/                               # Auth & session management (already exists)
│   │   ├── session.ts
│   │   ├── access-token.ts
│   │   ├── backend-store.ts
│   │   └── __tests__/
│   │
│   └── nextjs-cookies/                         # Cookie utilities (already exists)
│       └── index.ts
│
├── types/                                      # Shared TypeScript types
│   ├── api.ts                                  # API response types (Album, Photo, Owner, etc.)
│   └── errors.ts                               # Error types and codes
│
├── __tests__/                                  # Global test infrastructure
│   ├── setup.ts                                # Vitest global setup
│   ├── helpers/                                # Test helpers (already exists)
│   │   ├── oidc.ts
│   │   └── assertions.ts
│   ├── mocks/                                  # MSW mocks
│   │   ├── handlers.ts                         # API mock handlers
│   │   └── server.ts                           # MSW server setup
│   └── fixtures/                               # Test data fixtures
│       ├── albums.ts                           # Sample album data
│       ├── photos.ts                           # Sample photo data
│       └── users.ts                            # Sample user data
│
├── public/                                     # Static assets
│   ├── favicon.ico
│   ├── logo.svg
│   └── assets/
│       └── placeholder-image.svg               # Placeholder for failed images
│
└── playwright/                                 # E2E tests (if using Playwright)
    ├── tests/
    │   ├── album-list.spec.ts
    │   ├── album-view.spec.ts
    │   └── photo-viewer.spec.ts
    └── playwright.config.ts
```

---

### Architectural Boundaries

#### API Boundaries

**Backend REST API:**

- Base URL: `process.env.NEXT_PUBLIC_API_BASE_URL`
- Authentication: Google OAuth tokens in cookies (handled by backend)
- All API calls go through `libs/requests/api-client.ts`

**API Endpoints Used:**

```
GET  /api/v1/albums                          # List all albums (owned + shared)
GET  /api/v1/albums/:albumId                 # Get album details
POST /api/v1/albums                          # Create new album
PUT  /api/v1/albums/:albumId                 # Update album (name, dates)
DELETE /api/v1/albums/:albumId               # Delete album

GET  /api/v1/albums/:albumId/photos          # Get photos in album
GET  /api/v1/media/:mediaId/image            # Get image with quality param
     ?quality=blur|medium|high|full
     &width={number}

GET  /api/v1/albums/:albumId/sharing         # Get sharing info
POST /api/v1/albums/:albumId/sharing         # Share album with user
DELETE /api/v1/albums/:albumId/sharing/:email # Revoke access

GET  /api/v1/owners/:ownerId                 # Get owner profile
GET  /api/v1/users/me                        # Get current user profile
```

**API Client Boundary:**

```typescript
// libs/requests/api-client.ts
export async function fetchFromAPI<T>(endpoint: string, options?: RequestInit): Promise<T>

// All API calls use this single function
// Handles: authentication, error mapping, DPhotoError creation
```

---

#### Component Boundaries

**Server vs Client Component Split:**

**Server Components (data fetching):**

- `app/(authenticated)/page.tsx` - Fetches initial album list
- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/page.tsx` - Fetches album + photos
- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/photos/[photoId]/page.tsx` - Fetches photo details

**Client Components (interactivity):**

- All `_components/` - Handle user interactions
- All dialog components - Manage modal state
- Context providers - Manage application state

**Context Provider Nesting (in authenticated layout):**

```tsx
// app/(authenticated)/layout.tsx
<ThemeProvider>          {/* Material UI theme */}
    <AlbumsProvider>       {/* Albums list */}
        <FilterProvider>     {/* Active filter */}
            <SelectedAlbumProvider>  {/* Current album */}
                <PhotosProvider>       {/* Photos in album */}
                    {children}
                </PhotosProvider>
            </SelectedAlbumProvider>
        </FilterProvider>
    </AlbumsProvider>
</ThemeProvider>
```

**Component Communication:**

1. **Server → Client:** Props (initial data passed from Server Component to Client wrapper)
2. **Context → Components:** React Context hooks (`useAlbums`, `usePhotos`, etc.)
3. **Parent → Child:** Props (callbacks via `onAlbumClick`, `onDeleteClick`, etc.)
4. **Child → Parent:** Event callbacks (handled via `handle*` functions passed as props)

---

#### Service Boundaries

**No Backend Services in This Project:**
This is frontend-only architecture. Backend services (Golang packages in `pkg/`) are out of scope.

**Frontend "Services" (Utility Functions):**

- `libs/requests/api-client.ts` - API communication boundary
- `libs/image-loader.ts` - Image URL generation boundary
- `libs/security/` - Authentication boundary (already implemented)

**External Service Integrations:**

- **Google OAuth** - Handled by backend, frontend receives session cookies
- **AWS S3 (images)** - Accessed via backend API endpoints, not directly

---

#### Data Boundaries

**State Management Layers:**

1. **Server Data (NextJS Cache):**
    - Initial page loads cached by Next.js
    - Revalidation controlled via `revalidate` export

2. **Client Context State:**
    - `AlbumsContext` - Albums array, loading, error
    - `FilterContext` - Active filter, selected owner
    - `SelectedAlbumContext` - Current album details
    - `PhotosContext` - Photos for current album

3. **Local Component State:**
    - Dialog open/close state
    - Form input state
    - Hover states
    - Loading indicators for actions

**Data Flow:**

```
Backend API
    ↓
Server Component (fetch)
    ↓
Client Component (initial props)
    ↓
Context Provider (initialize)
    ↓
Consumer Components (useContext)
    ↓
User Interactions
    ↓
Context Actions (update state)
    ↓
Backend API (POST/PUT/DELETE)
    ↓
Context State (optimistic or refetch)
```

**Cache Boundaries:**

- **Browser HTTP Cache:** Images (immutable, 1 year cache)
- **Next.js Server Cache:** Initial page data (configurable revalidation)
- **No Client Cache:** React Context holds current state, cleared on refresh

---

### Requirements to Structure Mapping

#### Feature: Album Discovery & Browsing (FR1-FR9)

**UI Components:**

- `app/(authenticated)/page.tsx` - Server Component for initial album list fetch
- `app/(authenticated)/_components/AlbumListClient.tsx` - Client wrapper with interactions
- `app/(authenticated)/_components/AlbumCard.tsx` - Individual album display
- `app/(authenticated)/_components/AlbumFilter.tsx` - Filter UI (My/All/By Owner)
- `app/(authenticated)/_components/RandomHighlights.tsx` - Random photo highlights

**State Management:**

- `components/contexts/AlbumsContext.tsx` - Albums list, loading, error, actions
- `components/contexts/FilterContext.tsx` - Active filter state

**API Integration:**

- `libs/requests/api-client.ts` → `GET /api/v1/albums`

**Tests:**

- `app/(authenticated)/_components/__tests__/AlbumCard.test.tsx`
- `app/(authenticated)/_components/__tests__/AlbumFilter.test.tsx`

---

#### Feature: Photo Viewing & Navigation (FR10-FR17)

**UI Components:**

- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/_components/PhotoGrid.tsx` - Photo grid layout
- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/_components/PhotoThumbnail.tsx` - Individual photo
- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/@modal/(.)photos/[photoId]/page.tsx` - Modal viewer
- `app/(authenticated)/owners/[ownerId]/albums/[albumId]/photos/[photoId]/page.tsx` - Full-page viewer

**State Management:**

- `components/contexts/PhotosContext.tsx` - Photos array for current album
- `components/contexts/SelectedAlbumContext.tsx` - Current album details

**API Integration:**

- `libs/requests/api-client.ts` → `GET /api/v1/albums/:albumId/photos`
- `libs/image-loader.ts` → `GET /api/v1/media/:mediaId/image?quality=...`

**Tests:**

- `_components/__tests__/PhotoGrid.test.tsx`
- `_components/__tests__/PhotoThumbnail.test.tsx`
- `@modal/(.)photos/[photoId]/_components/__tests__/PhotoViewerModal.test.tsx`

---

#### Feature: Album Management (FR18-FR25)

**UI Components:**

- `app/(authenticated)/_components/CreateAlbumDialog.tsx` - Create new album
- `app/(authenticated)/owners/.../albums/[albumId]/_components/EditAlbumDialog.tsx` - Edit album name
- `app/(authenticated)/owners/.../albums/[albumId]/_components/EditAlbumDatesDialog.tsx` - Edit dates with visual picker
- `app/(authenticated)/_components/DeleteAlbumDialog.tsx` - Confirm deletion

**State Management:**

- `components/contexts/AlbumsContext.tsx` - Actions: addAlbum, updateAlbum, removeAlbum

**API Integration:**

- `POST /api/v1/albums` - Create album
- `PUT /api/v1/albums/:albumId` - Update album
- `DELETE /api/v1/albums/:albumId` - Delete album

**Tests:**

- `_components/__tests__/CreateAlbumDialog.test.tsx`
- `_components/__tests__/EditAlbumDatesDialog.test.tsx`

---

#### Feature: Sharing & Access Control (FR26-FR32)

**UI Components:**

- `app/(authenticated)/owners/.../albums/[albumId]/_components/SharingDialog.tsx` - Share/revoke UI

**State Management:**

- Local component state + refetch on changes

**API Integration:**

- `GET /api/v1/albums/:albumId/sharing` - Get sharing info
- `POST /api/v1/albums/:albumId/sharing` - Share with user
- `DELETE /api/v1/albums/:albumId/sharing/:email` - Revoke access

**Tests:**

- `_components/__tests__/SharingDialog.test.tsx`

---

#### Feature: Authentication (FR33-FR37)

**Already Implemented:**

- `app/auth/` - OAuth callback, error, logout pages
- `libs/security/` - Session and token management

**No New Structure Needed**

---

#### Cross-Cutting: Error Handling (FR45-FR50)

**Error Types:**

- `types/errors.ts` - DPhotoError class, ErrorCode constants

**Error Boundaries:**

- `app/error.tsx` - Root error boundary
- `app/(authenticated)/error.tsx` - Authenticated routes error boundary
- Page-level `error.tsx` files for album and photo routes
- Page-level `not-found.tsx` files for 404 handling

**Error Display:**

- `components/shared/ErrorMessage.tsx` - Reusable error component with recovery actions

**API Integration:**

- `libs/requests/api-client.ts` - Maps HTTP status codes to DPhotoError

---

#### Cross-Cutting: Loading States (FR38-FR44)

**Pattern:**

- All contexts use `{ loading, error, data }` pattern

**Components:**

- `components/shared/LoadingSkeleton.tsx` - Generic skeleton
- `app/(authenticated)/_components/AlbumListSkeleton.tsx` - Album list skeleton
- `_components/PhotoGridSkeleton.tsx` - Photo grid skeleton

**Tests:**

- Verify loading states in all data-fetching component tests

---

#### Cross-Cutting: Responsive Design (FR38-FR44)

**Theme:**

- `components/theme/theme.ts` - MUI breakpoint configuration

**Breakpoints:**

- xs: 0px (Mobile)
- sm: 600px (Tablet)
- md: 960px (Desktop)
- lg: 1280px (Large desktop)

**Usage:**

- All layout components use `sx` prop with responsive values
- Grid columns adjust per breakpoint
- Typography sizes scale per breakpoint

---

### Integration Points

#### Internal Communication

**Context-to-Context:**

- `FilterContext` changes trigger `AlbumsContext` updates (mediated by components)
- `SelectedAlbumContext` selection triggers `PhotosContext` fetch (mediated by components)
- No direct context-to-context dependencies

**Page-to-Page Navigation:**

- Next.js Link or router.push for navigation
- Context state preserved during navigation
- Parallel routes for modal interception (photo viewer)

**Component-to-API Flow:**

```
User Interaction (Component)
    ↓
Event Handler (handle*)
    ↓
Context Action (actions.fetch*, actions.add*, etc.)
    ↓
API Client (fetchFromAPI)
    ↓
Backend REST API
    ↓
Response / Error
    ↓
Context State Update
    ↓
Consumer Components Re-render
```

---

#### External Integrations

**Backend REST API:**

- **Integration Point:** `libs/requests/api-client.ts`
- **Authentication:** Cookies managed by browser, sent automatically
- **Error Handling:** HTTP status → DPhotoError mapping

**Google OAuth (via Backend):**

- **Integration Point:** `libs/security/` (already implemented)
- **Flow:** Frontend redirects → Backend handles OAuth → Cookies set → Frontend accesses protected pages

**AWS S3 Images (via Backend API):**

- **Integration Point:** `libs/image-loader.ts` + Next.js Image component
- **URL Pattern:** `/api/v1/media/:mediaId/image?quality=...&width=...`
- **Caching:** Browser handles via immutable cache headers from backend

---

#### Data Flow Diagram

```
User Action (Component)
    ↓
Event Handler (handle*)
    ↓
Context Action (actions.fetch*, actions.add*, etc.)
    ↓
API Client (fetchFromAPI)
    ↓
Backend REST API
    ↓
Response / Error
    ↓
Context State Update (setAlbums, setError, etc.)
    ↓
Context Consumers Re-render
    ↓
UI Updates
```

---

### File Organization Patterns

#### Configuration Files (Root Level)

- `package.json` - Dependencies, scripts, package manager config
- `next.config.ts` - Next.js configuration (image loader, base path, standalone output)
- `tsconfig.json` - TypeScript compiler options (strict mode, path aliases)
- `vitest.config.ts` - Test configuration (jsdom environment, MSW setup)
- `.env.local` - Local environment variables (gitignored)
- `.env.example` - Template for environment variables (committed)
- `.eslintrc.json` - Linting rules (Next.js standards)

#### Source Organization

**app/ Directory:**

- NextJS App Router pages and layouts
- `(authenticated)/` - Route group for protected pages
- File naming: `page.tsx`, `layout.tsx`, `error.tsx`, `not-found.tsx`
- `_components/` - Colocation for page-specific components (underscore prefix)

**components/ Directory:**

- ONLY truly shared components (used by 2+ pages)
- `contexts/` - React Context providers
- `shared/` - Shared UI components
- `theme/` - Material UI theme configuration

**libs/ Directory:**

- Utility functions and integrations
- Each lib is self-contained with its own purpose

**types/ Directory:**

- Shared TypeScript types
- `api.ts` - API response types
- `errors.ts` - Error types and codes

#### Test Organization

- **Unit tests:** `__tests__/` subfolder within component directory
- **Test file naming:** `ComponentName.test.tsx` (matches component name)
- **Integration tests:** `__tests__/mocks/` for MSW handlers
- **E2E tests:** `playwright/tests/` (if using Playwright)
- **Test utilities:** `__tests__/helpers/` at project root
- **Test fixtures:** `__tests__/fixtures/` for sample data

#### Asset Organization

- **Static assets:** `public/` directory (served as-is)
- **Images:** `public/assets/` for logos, icons, placeholders
- **Favicon:** `public/favicon.ico`
- **Backend media:** Fetched via API (NOT in public folder)

---

### Development Workflow Integration

#### Development Server Structure

```bash
npm run dev
# Serves from app/ directory
# Hot reloading enabled for all files
# Base path: / (no /nextjs prefix in development)
# API calls: NEXT_PUBLIC_API_BASE_URL from .env.local
# Port: 3000 (default)
```

#### Build Process Structure

```bash
npm run build
# Next.js builds to .next/ directory
# Open Next adapter processes for Lambda deployment
# Standalone mode: copies dependencies to .next/standalone/
# Static assets: .next/static/
# Server functions: .next/server/
```

#### Deployment Structure

```bash
# SST deployment (handled by deployments/cdk/)
# Serves from .next/standalone/
# Base path: /nextjs (configured in next.config.ts)
# Static assets: CloudFront CDN
# Lambda@Edge for SSR pages
# API Gateway integration for backend API
```

#### Testing Workflow

```bash
npm run test              # Unit tests (Vitest) - all __tests__/**/*.test.tsx
npm run test:watch        # Watch mode for development
npm run test:coverage     # Coverage report generation
npm run test:e2e          # Playwright E2E tests (if configured)
npm run lint              # ESLint code quality checks
```

---

### Project Structure Summary

**Total Estimated Files:**

- **Pages:** 12 (routes + error/not-found pages)
- **Page Components:** 25 (_components/ colocation)
- **Shared Components:** 8 (contexts + shared UI + theme)
- **Utilities:** 5 (libs/ directory)
- **Types:** 2 (types/ directory)
- **Tests:** 30+ (one per component + integration)
- **Config:** 7 (root config files)
- **Total:** ~90 source files, ~30 test files

**Key Directories by Purpose:**

- `app/` - 40 files (pages, layouts, errors, page components)
- `components/` - 15 files (contexts, shared UI, theme)
- `libs/` - 5 files (utilities for API, images, security)
- `types/` - 2 files (API types + error types)
- `__tests__/` - 30+ files (unit + integration tests)

**Colocation Strategy:**

- Page-specific components live in `_components/` next to their page
- Only move to `components/shared/` when used by 2+ pages
- Tests co-located with components in `__tests__/` subfolder
- Types co-located inline unless shared across pages (then in `types/`)

## Architecture Validation Results

### Coherence Validation ✅

#### Decision Compatibility

**Technology Stack Compatibility:**
All technology choices work together without conflicts:

- **Next.js 16.1.1 + React 19.2.3 + TypeScript 5.x** - Fully compatible versions
- **Material UI 6.x + Emotion 11.x** - Required peer dependencies specified
- **Vitest 4.x + MSW 2.x** - Compatible testing infrastructure
- **Next.js Image + Custom Loader** - Supported pattern for backend integration

**Pattern-Technology Alignment:**
All patterns align with chosen technologies:

- React Context API works seamlessly with React 19.2.3
- Material UI `sx` prop pattern aligns with styling approach
- NextJS App Router natively supports parallel routes decision
- Server/Client component split is core App Router feature

**No Conflicts Detected:** All 8 major architectural decisions work together harmoniously.

---

#### Pattern Consistency

**Naming Conventions:**
All naming patterns are consistent and comprehensive:

- PascalCase for components (AlbumCard.tsx, PhotoGrid.tsx)
- kebab-case for utilities (api-client.ts, image-loader.ts)
- Function prefixes cover all cases: `fetch*` (API), `get*` (state), `handle*` (events), `on*` (props)
- Semantic plural/singular rules clear and consistent

**Structure Patterns:**
Organization patterns support architectural decisions:

- Colocation principle (`_components/`) supports component architecture
- Context organization (`components/contexts/`) matches state management
- Test location pattern (`__tests__/`) aligns with component organization

**Communication Patterns:**
Data flow patterns are coherent across all layers:

- Server → Client → Context → Components flow is well-defined
- Single API boundary through `api-client.ts` ensures consistency
- Context actions pattern prevents direct state manipulation issues

---

#### Structure Alignment

**Structure Supports All Decisions:**

- `app/` directory structure matches NextJS App Router requirements
- `_components/` colocation implements component architecture decision
- `components/contexts/` location supports state management approach
- Parallel route `@modal/` structure enables photo modal interception

**Boundaries Properly Defined:**

- **API Boundary:** Single entry point at `libs/requests/api-client.ts`
- **Component Boundary:** Clear Server/Client component split
- **Data Boundary:** Layered through Context providers
- **Error Boundary:** Hierarchical page-level error.tsx files

**Integration Points Well-Structured:**

- Context provider nesting order defined in authenticated layout
- Image loader integration point specified with backend API contract
- Error boundary hierarchy matches route structure
- Component communication patterns clearly mapped

---

### Requirements Coverage Validation ✅

#### Functional Requirements Coverage (50 FRs)

**Album Discovery & Browsing (FR1-FR9):** ✅ Fully Covered

- **Structure:** `app/(authenticated)/page.tsx` + `_components/` (AlbumCard, AlbumFilter, RandomHighlights)
- **State:** AlbumsContext, FilterContext
- **API:** GET /api/v1/albums endpoint documented
- **Tests:** Component tests specified in `__tests__/`

**Photo Viewing & Navigation (FR10-FR17):** ✅ Fully Covered

- **Structure:** Photo grid in album page, modal viewer via parallel route, full-page fallback
- **Progressive Loading:** Custom image loader with blur/medium/high/full quality levels
- **Interactions:** Keyboard/touch handlers in viewer components
- **State:** PhotosContext, SelectedAlbumContext
- **API:** Image quality endpoints with parameters documented

**Album Management (FR18-FR25):** ✅ Fully Covered

- **Structure:** Create/Edit/Delete dialogs specified in structure
- **State:** Context actions (addAlbum, updateAlbum, removeAlbum) defined
- **API:** POST/PUT/DELETE endpoints documented
- **Validation:** Error types and validation patterns defined

**Sharing & Access Control (FR26-FR32):** ✅ Fully Covered

- **Structure:** SharingDialog.tsx component specified
- **API:** Sharing endpoints (GET/POST/DELETE) documented
- **Display:** Owner information integrated in album cards

**Authentication & User Management (FR33-FR37):** ✅ Already Implemented

- **Structure:** `app/auth/` flows exist, `libs/security/` implemented
- **No gaps:** Existing implementation meets all requirements

**Visual Presentation & UI State (FR38-FR44):** ✅ Fully Covered

- **Loading States:** Pattern defined (`{ loading, error, data }`), skeleton components specified
- **Empty States:** EmptyState component + page-specific empty states
- **Error Messages:** ErrorMessage component + page-level error boundaries
- **Responsive:** Material UI breakpoints configured (xs/sm/md/lg)

**Error Handling & Validation (FR45-FR50):** ✅ Fully Covered

- **Validation:** Error types, error codes (ErrorCode enum) defined
- **Display:** ErrorMessage component with recovery actions
- **Boundaries:** Page-level error.tsx hierarchy
- **Recovery:** Error-specific action buttons (retry, sign in, return home)

**Functional Requirements Coverage: 50/50 (100%) ✅**

---

#### Non-Functional Requirements Coverage

**Performance Requirements:** ✅ Fully Addressed

- **Interaction Responsiveness (<100ms):** React optimizations (useMemo, useCallback) specified in patterns
- **Image Loading (3s on slow networks):** Progressive loading with blur-up defined
- **Page Performance (Lighthouse ≥90):** Server Components for initial load, Next.js automatic code splitting
- **Responsive Sizing:** Custom image loader maps viewport to appropriate quality/width

**Integration Requirements:** ✅ Fully Addressed

- **Backend API Compatibility:** API client uses existing endpoints without requiring modifications
- **No Breaking Changes:** Frontend-only, read-only consumption of backend API
- **Graceful API Handling:** Loading states, structured error handling, retry mechanisms defined

**Security Requirements:** ✅ Fully Addressed

- **Backend Authentication:** Leverages existing Google OAuth implementation
- **Session Security:** Handled by backend with secure cookies
- **Authenticated Pages:** Route group `(authenticated)/` structure enforces authentication

**Usability Requirements:** ✅ Fully Addressed

- **Keyboard Navigation:** Specified in photo viewer component requirements
- **Mobile & Responsive:** MUI breakpoints configured, responsive grid patterns defined
- **Browser Compatibility:** Modern evergreen browsers only (ES2020+ allowed, no polyfills)

**Reliability Requirements:** ✅ Fully Addressed

- **Best-Effort Availability:** Acceptable per requirements, no uptime SLA needed
- **Manual Refresh Recovery:** Context pattern supports state refresh
- **Clear Error Messaging:** ErrorMessage component with code-specific recovery actions
- **No Data Loss:** Context actions are atomic, state updates properly managed

**All NFRs Architecturally Supported ✅**

---

### Implementation Readiness Validation ✅

#### Decision Completeness

**All Critical Decisions Documented with Versions:**

1. **Material UI Integration** - v6.x specified, dependencies listed, Tailwind removal documented
2. **State Management** - React Context pattern, provider nesting order defined
3. **Component Architecture** - Colocation principle with clear rules
4. **Routing Structure** - Owner-based paths, parallel routes for modal
5. **Image Optimization** - Custom loader implementation, quality levels, backend API contract
6. **Layout Architecture** - Material UI `sx` prop, responsive breakpoints
7. **Error Boundaries** - Page-level error.tsx hierarchy
8. **Data Fetching** - Server/Client component split with initialization pattern

**Implementation Patterns Comprehensive:**

- 8 pattern categories defined with detailed examples
- Quick reference table for AI agents
- Consistency checklist (11 items)
- Rationale provided for each pattern choice

**Examples Provided for All Major Patterns:**

- Context implementation with actions
- Component communication (Server → Client → Context)
- Error handling (types, codes, display, recovery)
- Loading state patterns
- API client with error mapping
- Responsive layouts with MUI breakpoints

---

#### Structure Completeness

**Project Structure Complete and Specific:**

- **Complete directory tree:** Root to leaves, all files specified
- **Purpose comments:** Each file/directory documented
- **Test structure:** Co-located `__tests__/` pattern defined
- **File estimates:** 90 source files, 30+ test files documented

**All Integration Points Clearly Specified:**

- **API Endpoints:** 12 endpoints listed with parameters
- **Context Nesting:** Provider order documented
- **Data Flow:** Diagram provided (User → Handler → Context → API → State → UI)
- **Component Communication:** 4 patterns explained

**Component Boundaries Well-Defined:**

- **Server Components:** 3 pages for initial data fetching
- **Client Components:** All `_components/` for interactivity
- **Colocation Rules:** Move to shared only when used 2+ places
- **Context Scope:** 4 contexts with clear responsibilities

---

#### Pattern Completeness

**All Potential Conflict Points Addressed:**

- **File Naming:** PascalCase vs kebab-case prevents conflicts
- **Context Re-renders:** Split by concern (4 contexts instead of 1)
- **Premature Abstraction:** Colocation principle prevents over-engineering

**Naming Conventions Comprehensive:**

- **Components:** PascalCase (AlbumCard.tsx)
- **Utilities:** kebab-case (api-client.ts)
- **Tests:** ComponentName.test.tsx pattern
- **Functions:** fetch*/get*/handle*/on* prefixes cover all types
- **Variables:** Semantic plural/singular rules

**Communication Patterns Fully Specified:**

- **Server → Client:** Props with initial data
- **Context → Components:** React hooks (useAlbums, usePhotos)
- **Parent → Child:** Callback props (onAlbumClick, onDeleteClick)
- **Child → Parent:** Event handlers (handleClick functions)

**Process Patterns Complete:**

- **Error Handling:** DPhotoError class, ErrorCode enum, mapHTTPStatusToErrorCode, ErrorMessage component
- **Loading States:** `{ loading, error, data }` pattern in all contexts
- **Data Fetching:** Server fetch → Client initialize Context → Components consume

---

### Gap Analysis Results

#### Critical Gaps: ✅ NONE

All blocking architectural decisions are documented with sufficient detail for consistent implementation by AI agents.

---

#### Important Gaps: ⚠️ 1 ASSUMPTION TO VERIFY

**Backend API Quality Endpoint Parameters:**

- **Assumption:** Architecture assumes backend supports exact parameters: `?quality=blur|medium|high|full&width={number}`
- **Impact:** If backend API uses different parameter names or values, custom image loader requires adjustment
- **Mitigation Strategy:**
    - Verify backend API endpoint signature before implementing image loader
    - If parameters differ, update `libs/image-loader.ts` mapping logic
    - Document actual API parameters in code comments
- **Severity:** Important but not blocking (image loader can be adapted to actual API)
- **Documentation:** Backend API requirements clearly stated in Image Optimization section

---

#### Nice-to-Have Gaps: 💡 3 SUGGESTIONS FOR FUTURE ENHANCEMENT

1. **Material UI Component Customization Guidelines:**
    - When to extend MUI components vs build custom
    - Where MUI theme overrides go vs custom components
    - Impact: Minor - can be determined during implementation
    - Recommendation: Add guidelines when patterns emerge during development

2. **Performance Monitoring Strategy:**
    - Lighthouse testing frequency and automation
    - Performance budget enforcement
    - Impact: Minor - good practice but not architectural decision
    - Recommendation: Define in development workflow, not architecture

3. **Accessibility Testing Approach:**
    - Keyboard navigation testing strategy
    - ARIA label patterns and testing
    - Impact: Minor - implementation detail rather than architecture
    - Recommendation: Document as implementation patterns emerge

**None of these gaps block implementation. All can be addressed during or after initial development.**

---

### Validation Issues Addressed

**No Critical or Important Issues Found** ✅

The architecture is coherent, complete, and ready for AI agent implementation.

**Backend API Assumption Documented:**
The one important assumption (image quality API parameters) is clearly documented in the "Backend API Requirements" section of the Image Optimization decision.
Implementation teams must verify actual backend API signature and adapt the custom image loader accordingly.

---

### Architecture Completeness Checklist

#### ✅ Requirements Analysis

- [x] Project context thoroughly analyzed (50 FRs, 34 NFRs across 5 categories)
- [x] Scale and complexity assessed (Medium complexity, 20-30 components, 90 source files)
- [x] Technical constraints identified (Modern browsers, existing backend API, AWS deployment)
- [x] Cross-cutting concerns mapped (State, images, errors, responsive, Material UI, UX patterns)

#### ✅ Architectural Decisions

- [x] Critical decisions documented with versions (8 major decisions, all with specific versions)
- [x] Technology stack fully specified (Next.js 16.1.1, React 19.2.3, MUI 6.x, TypeScript 5.x)
- [x] Integration patterns defined (API client boundary, Context nesting, Server/Client split)
- [x] Performance considerations addressed (Progressive loading, Server Components, responsive sizing)

#### ✅ Implementation Patterns

- [x] Naming conventions established (PascalCase, kebab-case, function prefixes, plural/singular)
- [x] Structure patterns defined (Colocation, `_components/`, `__tests__/`, shared criteria)
- [x] Communication patterns specified (Server→Client→Context, props, callbacks, actions)
- [x] Process patterns documented (Error handling, loading states, data fetching, context updates)

#### ✅ Project Structure

- [x] Complete directory structure defined (Root to leaves, 90 source + 30 test files)
- [x] Component boundaries established (Server vs Client, colocation rules, shared criteria)
- [x] Integration points mapped (12 API endpoints, Context nesting, data flow diagram)
- [x] Requirements to structure mapping complete (All 50 FRs mapped to specific files/directories)

---

### Architecture Readiness Assessment

**Overall Status:** ✅ **READY FOR IMPLEMENTATION**

**Confidence Level:** **HIGH**

**Justification:**

- All 50 functional requirements architecturally supported
- All non-functional requirements addressed
- 8 major architectural decisions fully documented with versions
- 8 implementation pattern categories with examples and consistency rules
- Complete project structure (90 files) specified with purpose
- All integration points and boundaries clearly defined
- No critical gaps, only 1 minor assumption to verify
- Comprehensive validation completed with 100% requirements coverage

---

**Key Strengths:**

1. **Coherent Technology Choices:** Material UI-only approach eliminates CSS conflicts, React Context avoids state library complexity
2. **Clear Boundaries:** Single API client, layered Context architecture, Server/Client split prevents confusion
3. **Implementation Patterns:** 8 comprehensive pattern categories with examples prevent AI agent conflicts
4. **Colocation Principle:** Components near pages speeds up development and understanding
5. **Progressive Image Loading:** Custom loader with quality levels ensures optimal performance
6. **Complete Structure:** 90-file project tree leaves no organizational ambiguity
7. **Testing Strategy:** Co-located tests, MSW mocks, Vitest infrastructure ready
8. **Error Handling:** Structured errors with codes enable specific recovery flows

---

**Areas for Future Enhancement:**

1. **Performance Monitoring:** Add Lighthouse CI integration to enforce performance budgets
2. **Component Library Guidelines:** Document MUI customization patterns as they emerge
3. **Accessibility Testing:** Formalize keyboard navigation and ARIA testing strategy
4. **E2E Test Coverage:** Expand Playwright tests beyond initial setup
5. **Image Quality Optimization:** Fine-tune quality levels based on actual usage data
6. **Context Optimization:** Monitor for unnecessary re-renders and add selective subscriptions if needed
7. **Backend API Caching:** Explore Next.js ISR or other caching strategies for album data

**None of these enhancements are required for initial implementation. All represent iterative improvements.**

---

### Implementation Handoff

#### AI Agent Guidelines

**Core Principles:**

1. **Follow all architectural decisions exactly as documented** - No deviations without explicit approval
2. **Use implementation patterns consistently across all components** - Check pattern summary before writing code
3. **Respect project structure and boundaries** - Files go in specified locations, not elsewhere
4. **Refer to this document for all architectural questions** - Architecture is single source of truth

**Before Writing Code:**

1. Read relevant architectural decision section
2. Review implementation patterns that apply
3. Check project structure for file locations
4. Verify requirements mapping for feature

**During Implementation:**

1. Follow naming conventions (PascalCase components, kebab-case utilities)
2. Use specified patterns (Context actions, error handling, loading states)
3. Place components in correct locations (colocation vs shared)
4. Write co-located tests in `__tests__/` subfolders

**Quality Checks:**

1. Run consistency checklist (11 items) on your code
2. Verify component communication follows patterns
3. Ensure error handling uses DPhotoError structure
4. Test responsive behavior at all breakpoints

---

#### First Implementation Priorities

**Phase 1: Foundation (Infrastructure)**

1. **Remove Tailwind CSS:** Uninstall tailwindcss, @tailwindcss/postcss from package.json
2. **Install Material UI:** Add @mui/material ^6.x, @mui/icons-material ^6.x, @emotion/react ^11.x, @emotion/styled ^11.x
3. **Create Theme:** Implement `components/theme/theme.ts` with dark mode + brand color #185986
4. **Setup Root Layout:** Add ThemeProvider to `app/layout.tsx`

**Phase 2: Core Utilities**

1. **Enhance API Client:** Add DPhotoError handling to `libs/requests/api-client.ts`
2. **Create Error Types:** Implement `types/errors.ts` (DPhotoError, ErrorCode)
3. **Create API Types:** Implement `types/api.ts` (Album, Photo, Owner, SharingInfo)
4. **Create Image Loader:** Implement `libs/image-loader.ts` with quality mapping

**Phase 3: State Management**

1. **AlbumsContext:** Implement with actions (fetchAlbums, addAlbum, updateAlbum, removeAlbum)
2. **FilterContext:** Implement with filter state management
3. **SelectedAlbumContext:** Implement with current album state
4. **PhotosContext:** Implement with photos array for selected album

**Phase 4: Shared Components**

1. **ErrorMessage:** Implement with code-specific recovery actions
2. **LoadingSkeleton:** Generic loading skeleton component
3. **EmptyState:** Generic empty state component

**Phase 5: Feature Implementation**

1. **Home Page:** Album list with filtering (start here - most visible feature)
2. **Album View:** Photo grid with day grouping
3. **Photo Viewer:** Modal + full-page viewer with parallel routes
4. **Album Management:** Create/Edit/Delete dialogs
5. **Sharing:** SharingDialog implementation

**Recommended Start:** Begin with Phase 1-2 (Foundation + Core Utilities) to establish the architectural foundation, then implement Phase 3 (State Management)
before building features in Phase 5.

---

### Architecture Document Complete ✅

This architecture document provides comprehensive guidance for consistent implementation by AI agents. All architectural decisions, implementation patterns,
project structure, and validation results are documented and ready for development.

**Document Sections:**

1. ✅ Project Context Analysis
2. ✅ Starter Template Evaluation
3. ✅ Core Architectural Decisions (8 decisions)
4. ✅ Implementation Patterns & Consistency Rules (8 patterns)
5. ✅ Project Structure & Boundaries (90-file tree)
6. ✅ Architecture Validation Results (100% coverage verified)

**Next Steps:**

- Proceed to implementation following phased priorities above
- Refer to this document for all architectural questions
- Update architecture if significant new requirements emerge

---

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** 2026-01-31
**Document Location:** `/home/dush/dev/git/dphoto/specs/designs/architecture.md`

---

### Final Architecture Deliverables

#### 📋 Complete Architecture Document

- **All architectural decisions documented** with specific versions (Next.js 16.1.1, React 19.2.3, Material UI 6.x, TypeScript 5.x)
- **Implementation patterns** ensuring AI agent consistency (8 pattern categories with examples)
- **Complete project structure** with all files and directories (90 source files, 30+ test files)
- **Requirements to architecture mapping** (50 FRs + all NFRs mapped to specific locations)
- **Validation** confirming coherence and completeness (100% requirements coverage)

---

#### 🏗️ Implementation Ready Foundation

- **8 architectural decisions** made (Material UI, React Context, colocation, routing, images, layouts, errors, data fetching)
- **8 implementation patterns** defined (naming, plural/singular, functions, tests, types, loading, context, errors)
- **6 architectural components** specified (app/, components/, libs/, types/, __tests__/, public/)
- **50 functional requirements** fully supported + all non-functional requirements addressed

---

#### 📚 AI Agent Implementation Guide

- **Technology stack** with verified compatible versions
- **Consistency rules** that prevent implementation conflicts (11-item checklist)
- **Project structure** with clear boundaries (API, Component, Service, Data)
- **Integration patterns** and communication standards (Server→Client→Context→Component)

---

### Implementation Handoff

#### For AI Agents:

This architecture document is your complete guide for implementing the DPhoto NextJS Web UI. Follow all decisions, patterns, and structures exactly as
documented.

**First Implementation Priority:**

**Phase 1: Foundation (Start Here)**

1. Remove Tailwind CSS dependencies
2. Install Material UI packages (@mui/material, @mui/icons-material, @emotion/react, @emotion/styled)
3. Create theme with dark mode + brand color #185986
4. Setup ThemeProvider in root layout

**Phase 2: Core Utilities**

1. Enhance API client with DPhotoError handling
2. Create error types (types/errors.ts)
3. Create API types (types/api.ts)
4. Create custom image loader (libs/image-loader.ts)

**Phase 3: State Management**

1. Implement AlbumsContext with actions
2. Implement FilterContext
3. Implement SelectedAlbumContext
4. Implement PhotosContext

**Phase 4: Shared Components**

1. ErrorMessage component
2. LoadingSkeleton component
3. EmptyState component

**Phase 5: Features**

1. Home page (album list with filtering)
2. Album view (photo grid)
3. Photo viewer (modal + full-page)
4. Album management (dialogs)
5. Sharing functionality

---

#### Development Sequence:

1. **Initialize project** using existing Next.js setup in `web-nextjs/`
2. **Set up development environment** per architecture (remove Tailwind, add MUI)
3. **Implement core architectural foundations** (theme, contexts, utilities)
4. **Build features** following established patterns (colocation, naming conventions)
5. **Maintain consistency** with documented rules (consistency checklist)

---

### Quality Assurance Checklist

#### ✅ Architecture Coherence

- [x] All decisions work together without conflicts
- [x] Technology choices are compatible (Next.js 16 + React 19 + MUI 6 + TypeScript 5)
- [x] Patterns support the architectural decisions (Context supports state, colocation supports organization)
- [x] Structure aligns with all choices (App Router structure, Material UI theming, Context layers)

#### ✅ Requirements Coverage

- [x] All functional requirements are supported (50/50 FRs mapped to structure)
- [x] All non-functional requirements are addressed (Performance, Integration, Security, Usability, Reliability)
- [x] Cross-cutting concerns are handled (Error handling, loading states, responsive design, image optimization)
- [x] Integration points are defined (API client, Context nesting, image loader, authentication)

#### ✅ Implementation Readiness

- [x] Decisions are specific and actionable (8 decisions with versions and examples)
- [x] Patterns prevent agent conflicts (Naming conventions, colocation rules, consistency checklist)
- [x] Structure is complete and unambiguous (90-file tree specified from root to leaves)
- [x] Examples are provided for clarity (Code examples for all patterns, data flow diagrams)

---

### Project Success Factors

#### 🎯 Clear Decision Framework

Every technology choice was made collaboratively with clear rationale (Material UI-only for simplicity, React Context for medium complexity, colocation for
developer speed), ensuring all stakeholders understand the architectural direction.

#### 🔧 Consistency Guarantee

Implementation patterns and rules ensure that multiple AI agents will produce compatible, consistent code that works together seamlessly (PascalCase components,
Context actions not setters, error structure, loading pattern).

#### 📋 Complete Coverage

All project requirements are architecturally supported, with clear mapping from business needs to technical implementation (Album discovery → home page +
contexts, Photo viewing → grid + modal + parallel routes, etc.).

#### 🏗️ Solid Foundation

The existing Next.js 16 setup with TypeScript provides a production-ready foundation. Architecture decisions enhance it with Material UI theming, React Context
state management, and progressive image loading following current best practices.

---

**Architecture Status:** ✅ **READY FOR IMPLEMENTATION**

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation. Keep the architecture document as the single
source of truth for all technical decisions.

