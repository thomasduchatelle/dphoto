---
stepsCompleted: [ 'step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation' ]
workflowComplete: true
completedDate: '2026-02-01'
inputDocuments:
  - '/home/dush/dev/git/dphoto/specs/designs/prd.md'
  - '/home/dush/dev/git/dphoto/specs/designs/architecture.md'
  - '/home/dush/dev/git/dphoto/specs/designs/ux-design-specification.md'
  - '/home/dush/dev/git/dphoto/AGENTS.md'
---

# dphoto - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for dphoto, decomposing the requirements from the PRD, UX Design, and Architecture requirements
into implementable stories.

## Requirements Inventory

### Functional Requirements

**Album Discovery & Browsing:**

- **FR1:** Album owners can view all albums they own
- **FR2:** Shared viewers can view all albums shared with them
- **FR3:** Users can filter albums by owner (my albums, all albums, specific owner)
- **FR4:** Users can view album metadata (name, date range, media count, owner information)
- **FR5:** Users can see which users an album is shared with
- **FR6:** Users can view albums in chronological order
- **FR7:** Users can see visual indicators of album activity density (temperature)
- **FR8:** Users can discover random photos from across their accessible collection
- **FR9:** Users can navigate from random photo highlights to the source album

**Photo Viewing & Navigation:**

- **FR10:** Users can view photos grouped by capture date within an album
- **FR11:** Users can open individual photos in full-screen view
- **FR12:** Users can navigate between photos using keyboard controls (arrow keys, esc, enter)
- **FR13:** Users can navigate between photos using touch gestures (swipe on mobile/tablet)
- **FR14:** Users can view photos with progressive quality loading (low to high resolution)
- **FR15:** Users can zoom into photo details
- **FR16:** Users can navigate back to album list from photo view
- **FR17:** Users can see date headers separating photos by day

**Album Management:**

- **FR18:** Album owners can create new albums by specifying name and date range
- **FR19:** Album owners can specify custom folder names for albums (optional)
- **FR20:** Album owners can edit album names
- **FR21:** Album owners can edit album date ranges
- **FR22:** Album owners can delete albums they own
- **FR23:** System validates that date edits don't orphan media
- **FR24:** System re-indexes photos when album date ranges change
- **FR25:** System provides feedback when album operations succeed or fail

**Sharing & Access Control:**

- **FR26:** Album owners can share albums with other users via email address
- **FR27:** Album owners can revoke access from users they previously shared with
- **FR28:** Album owners can view who has access to their albums
- **FR29:** Shared viewers can see who owns albums shared with them
- **FR30:** System distinguishes between owner capabilities and viewer capabilities in UI
- **FR31:** System validates email addresses when granting access
- **FR32:** System loads user profile information (name, picture) for display

**Authentication & User Management:**

- **FR33:** Users can authenticate using Google OAuth
- **FR34:** System maintains user session across page navigation
- **FR35:** Users can view their profile information (name, picture)
- **FR36:** System identifies if authenticated user is an owner (can create albums)
- **FR37:** System restricts all pages to authenticated users only

**Visual Presentation & UI State:**

- **FR38:** System displays loading states while fetching albums or photos
- **FR39:** System displays empty states when no albums exist
- **FR40:** System displays error messages when operations fail
- **FR41:** System indicates selected albums and active filters
- **FR42:** System provides visual transitions when navigating between views
- **FR43:** System displays album sharing status with user avatars
- **FR44:** System provides responsive layouts for mobile, tablet, and desktop devices

**Error Handling & Validation:**

- **FR45:** System validates album date ranges (end date after start date)
- **FR46:** System validates album names are not empty
- **FR47:** System handles and displays errors for failed album operations
- **FR48:** System handles and displays errors for failed sharing operations
- **FR49:** System handles and displays errors when albums are not found
- **FR50:** System provides recovery options when operations fail

### Non-Functional Requirements

**Performance:**

- **NFR1:** Visual transitions and animations must maintain 60fps on modern mobile and desktop devices
- **NFR2:** User interactions (clicks, taps, swipes) must provide immediate visual feedback (<100ms)
- **NFR3:** Keyboard navigation must respond without perceptible lag
- **NFR4:** Page transitions and photo zoom animations must feel smooth and purposeful
- **NFR5:** System must display thumbnail-quality images immediately while full-resolution images load in background
- **NFR6:** System must request minimum image size appropriate for current screen dimensions (mobile, tablet, desktop)
- **NFR7:** Progressive image loading (blur-up from low to high quality) must complete within 3 seconds on slow network conditions
- **NFR8:** System must use existing API quality parameters to optimize image delivery
- **NFR9:** Initial page load must achieve Lighthouse Performance score ≥90 on mobile devices
- **NFR10:** Time to interactive must be optimized through efficient code splitting and lazy loading
- **NFR11:** Subsequent page navigations must feel instantaneous through appropriate caching strategies

**Integration:**

- **NFR12:** Frontend must consume existing REST API endpoints without requiring backend modifications
- **NFR13:** System must handle API response times gracefully with appropriate loading states
- **NFR14:** System must handle API errors with clear user feedback and retry mechanisms
- **NFR15:** System must maintain data contract compatibility with existing API response formats
- **NFR16:** System must leverage existing API image quality parameters for progressive loading

**Security:**

- **NFR17:** System relies on existing backend Google OAuth authentication mechanism
- **NFR18:** System must maintain user session security as provided by backend
- **NFR19:** System must restrict all pages to authenticated users only (enforced by backend)
- **NFR20:** Frontend requires no additional security measures beyond consuming authenticated API endpoints

**Usability:**

- **NFR21:** Core photo browsing flows must be fully keyboard-accessible (arrow keys, esc, enter)
- **NFR22:** Focus management must be clear and logical during navigation and in dialogs
- **NFR23:** Keyboard shortcuts must not conflict with browser defaults
- **NFR24:** Touch interactions must feel responsive with appropriate visual feedback
- **NFR25:** Layouts must adapt appropriately across mobile (<600px), tablet (600-960px), and desktop (>960px) breakpoints
- **NFR26:** Mobile gestures (swipe, pinch-to-zoom) must feel natural and performant
- **NFR27:** Mobile performance must not degrade below acceptable animation frame rates
- **NFR28:** System must support latest 2 versions of Chrome, Firefox, Safari, and Edge (evergreen browsers)
- **NFR29:** System may use modern web features (CSS Grid, Flexbox, ES2020+) without polyfills for legacy browsers
- **NFR30:** No support required for Internet Explorer or older browser versions

**Reliability:**

- **NFR31:** System operates on best-effort availability basis (no uptime SLA required)
- **NFR32:** Manual page refresh is acceptable recovery mechanism for transient errors
- **NFR33:** System must provide clear error messages when operations fail with guidance on recovery
- **NFR34:** Network failures must not cause data loss for in-progress operations

### Additional Requirements

**From Architecture Document:**

**Material UI Integration:**

- Remove Tailwind CSS dependencies completely
- Configure MUI theme with dark mode as default
- Set brand color (#185986) as primary throughout theme
- Use MUI breakpoint system: `xs` (<600px), `sm` (600px), `md` (960px), `lg` (1280px)

**State Management (Lift-and-Shift from existing web/):**

- Migrate battle-tested state management from `web/src/core/catalog/` (90+ tests)
- Copy state definition (CatalogViewerState), actions, reducer, thunks
- Replace axios adapter with fetch (server + client compatible)
- Initialize state server-side in Server Components, pass to Client Components as props
- Instantiate thunks client-side for user interactions
- Pure UI components receive state and handlers as props (NO internal state management)
- Add `router.refresh()` for NextJS cache invalidation after mutations

**Component Architecture:**

- Colocation principle - components live with pages unless used by 2+ pages
- Place page-specific components in `_components/` subfolder next to page
- Move to `components/shared/` only when used by 2+ pages
- Keep contexts in `components/contexts/`

**Routing Structure:**

- Owner-based paths: `/owners/[ownerId]/albums/[albumId]`
- NextJS parallel routes for photo modal interception: `@modal/(.)photos/[photoId]/`
- Fallback route `photos/[photoId]/page.tsx` for direct access
- Render `{children}` and `{modal}` in album layout

**Image Optimization:**

- Next.js Image component with custom loader in `libs/image-loader.ts`
- Map width to backend quality levels:
    - width ≤40: `blur` (ultra-low placeholder)
    - width ≤500: `medium` (grid display)
    - width ≤1200: `high` (full-screen viewer)
    - width >1200: `full` (original for zoom)
- Use `placeholder="blur"`, `blurDataURL`, responsive `sizes`
- Configure image loader in `next.config.ts`

**Layout Architecture:**

- Material UI `sx` prop with theme breakpoints for responsive layouts
- Responsive object syntax for grid columns, gap, padding, typography
- No inline styles or custom breakpoint systems

**Error Boundaries:**

- NextJS page-level error boundaries using error.tsx files
- Create `error.tsx` at each route level (root, authenticated, album, photo)
- Create `not-found.tsx` for 404 handling
- Provide "Try Again" button calling `reset()` function

**Data Fetching:**

- Server Components fetch initial data in page.tsx files
- Pass initial data as props to Client Component wrappers
- Initialize Context state with server data via `useEffect`
- Handle user interactions on client
- Use `revalidate` or `cache: 'no-store'` for dynamic data

**Backend API Requirements:**

- Image endpoint: `GET /api/v1/media/{mediaId}/image?quality={quality}&width={width}`
- Quality levels: `blur`, `medium`, `high`, `full`
- Cache headers: `Cache-Control: max-age=31536000, immutable`
- MediaId changes if content changes (immutability)

**From UX Design Document:**

**Visual Date Selection Pattern:**

- Hybrid contextual + preview approach
- Primary: Create album from media list with suggested dates
- Secondary: Manual creation with live photo preview (~12 thumbnails)
- Photo count updates in real-time as dates adjust
- Side panel shows photos in selected range

**Random Photo Discovery:**

- Album cards display 3-4 random photo thumbnails as preview grid/carousel
- Home page highlights section with 5-8 random photos from all albums
- Each photo links to source album
- Refreshes on page reload

**Album Activity Indicators:**

- Density color-coding based on photos-per-day
- High density: Warmer colors or bolder display
- Low density: Cooler colors or lighter display
- Integrates with album card layout

**Dark Theme with Brand Color:**

- Background: #121212
- Surface/Cards: #1e1e1e
- Brand Blue (#185986) for primary actions, focus states, links
- Accent lines: Light cyan/blue (#4a9ece) or desaturated white (#6a7a8a)
- Text primary: #ffffff
- Text secondary: rgba(255, 255, 255, 0.7)

**Performance Targets:**

- 60fps animations
- <100ms response to user interactions
- Progressive image loading (blur-up within 3s on slow networks)
- Lighthouse Performance ≥90 on mobile

**From AGENTS.md (Project Context):**

**Domain Architecture:**

- Frontend only affects `web-nextjs/` (no backend, infrastructure, API, CLI, or data model changes)
- Must respect domain separation: catalog, archive, backup, ACL
- Frontend coding standards in `.github/instructions/nextjs.instructions.md`

**Project Structure:**

- NextJS App Router file structure
- Best practices from NextJS
- Build: `npm run test`, `npm run test:visual`, `npm run laddle`
- Always run `npm install` before other commands

**Testing Strategy:**

- Unit tests (`npm run test`)
- Visual tests (`npm run test:visual`)
- Laddle component viewer on :61000

**No Backend Changes:**

- All backend REST API endpoints already exist
- No modifications to `pkg/`, `api/lambdas/`, `deployments/cdk/`, or `DATA_MODEL.md`
- Frontend consumes existing API contracts

### FR Coverage Map

**Epic 1 (Album List Home Page):**

- FR1: Album owners can view all albums they own
- FR2: Shared viewers can view all albums shared with them
- FR3: Users can filter albums by owner (my albums, all albums, specific owner)
- FR4: Users can view album metadata (name, date range, media count, owner information)
- FR5: Users can see which users an album is shared with
- FR6: Users can view albums in chronological order
- FR7: Users can see visual indicators of album activity density (temperature)
- FR33: Users can authenticate using Google OAuth
- FR34: System maintains user session across page navigation
- FR35: Users can view their profile information (name, picture)
- FR36: System identifies if authenticated user is an owner (can create albums)
- FR37: System restricts all pages to authenticated users only
- FR38: System displays loading states while fetching albums or photos
- FR39: System displays empty states when no albums exist
- FR40: System displays error messages when operations fail
- FR41: System indicates selected albums and active filters
- FR43: System displays album sharing status with user avatars
- FR44: System provides responsive layouts for mobile, tablet, and desktop devices

**Epic 2 (Photo Viewing & Navigation):**

- FR10: Users can view photos grouped by capture date within an album
- FR11: Users can open individual photos in full-screen view
- FR12: Users can navigate between photos using keyboard controls (arrow keys, esc, enter)
- FR13: Users can navigate between photos using touch gestures (swipe on mobile/tablet)
- FR14: Users can view photos with progressive quality loading (low to high resolution)
- FR15: Users can zoom into photo details
- FR16: Users can navigate back to album list from photo view
- FR17: Users can see date headers separating photos by day
- FR42: System provides visual transitions when navigating between views

**Epic 3 (Album Management):**

- FR18: Album owners can create new albums by specifying name and date range
- FR19: Album owners can specify custom folder names for albums (optional)
- FR20: Album owners can edit album names
- FR21: Album owners can edit album date ranges
- FR22: Album owners can delete albums they own
- FR23: System validates that date edits don't orphan media
- FR24: System re-indexes photos when album date ranges change
- FR25: System provides feedback when album operations succeed or fail
- FR45: System validates album date ranges (end date after start date)
- FR46: System validates album names are not empty
- FR47: System handles and displays errors for failed album operations
- FR49: System handles and displays errors when albums are not found
- FR50: System provides recovery options when operations fail

**Epic 4 (Sharing & Access Control):**

- FR26: Album owners can share albums with other users via email address
- FR27: Album owners can revoke access from users they previously shared with
- FR28: Album owners can view who has access to their albums
- FR29: Shared viewers can see who owns albums shared with them
- FR30: System distinguishes between owner capabilities and viewer capabilities in UI
- FR31: System validates email addresses when granting access
- FR32: System loads user profile information (name, picture) for display
- FR48: System handles and displays errors for failed sharing operations

**Epic 5 (Random Photo Discovery):**

- FR8: Users can discover random photos from across their accessible collection
- FR9: Users can navigate from random photo highlights to the source album

**Total FRs Covered:** 50 / 50 ✅

## Epic List

### Epic 1: Album List Home Page

Users can authenticate, view their album list on the home page with album cards showing random photo samples and density indicators, and filter albums by owner
with a fully functional filter.

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR33, FR34, FR35, FR36, FR37, FR38, FR39, FR40, FR41, FR43, FR44 (18 FRs)

**NFRs addressed:** NFR9 (Lighthouse ≥90), NFR17-NFR20 (auth security), NFR25 (responsive breakpoints), NFR38 (loading states)

**Additional Requirements:**

- Remove Tailwind, install Material UI with dark theme (#185986)
- Migrate state management from `web/src/core/catalog/` (90+ tests)
- Replace axios with fetch adapter
- Set up NextJS App Router with Server/Client component pattern
- Configure custom image loader for progressive loading
- Establish MUI theme, breakpoints, error boundaries
- Album cards display 3-4 random photo samples (UX requirement)
- Density color-coding for album activity
- Colocation principle for components

### Epic 2: Photo Viewing & Navigation

Users can view photos within albums grouped by day, navigate smoothly with keyboard and touch controls, see progressive image loading, zoom into details, and
experience smooth transitions.

**FRs covered:** FR10, FR11, FR12, FR13, FR14, FR15, FR16, FR17, FR42 (9 FRs)

**NFRs addressed:** NFR1-NFR4 (60fps, <100ms response), NFR5-NFR8 (progressive loading), NFR12-NFR16 (API integration), NFR21-NFR23 (keyboard nav),
NFR24-NFR27 (mobile gestures)

**Additional Requirements:**

- NextJS parallel routes for photo modal interception: `@modal/(.)photos/[photoId]/`
- Owner-based paths: `/owners/[ownerId]/albums/[albumId]`
- Next.js Image component with custom loader (blur/medium/high/full quality mapping)
- MUI sx prop for responsive layouts

### Epic 3: Album Management

Album owners can create, edit, and delete albums with visual date selection showing real-time photo previews and photo counts that update as dates adjust.

**FRs covered:** FR18, FR19, FR20, FR21, FR22, FR23, FR24, FR25, FR45, FR46, FR47, FR49, FR50 (13 FRs)

**NFRs addressed:** NFR13-NFR14 (error handling), NFR33-NFR34 (clear error messages, no data loss)

**Additional Requirements:**

- Visual date picker with live photo preview (~12 thumbnails)
- Contextual creation from media list with suggested dates
- Photo count updates in real-time as dates adjust
- Validation, orphan photo detection, error handling

### Epic 4: Sharing & Access Control

Album owners can share albums with users via email, revoke access, view who has access with avatars and names; shared viewers see ownership information and
clear UI distinction between owner and viewer capabilities.

**FRs covered:** FR26, FR27, FR28, FR29, FR30, FR31, FR32, FR48 (8 FRs)

**NFRs addressed:** NFR13-NFR14 (error handling)

**Additional Requirements:**

- User profile information loading (name, picture)
- Email validation
- Error handling for sharing operations

### Epic 5: Random Photo Discovery

Users can discover forgotten memories through random photo highlights on the home page (5-8 photos from all accessible albums) that link directly to source
albums and refresh on each page reload.

**FRs covered:** FR8, FR9 (2 FRs)

**NFRs addressed:** None specific (leverages existing performance from Epic 1)

**Additional Requirements:**

- Home page "Your Memories" or "Highlights" section at top
- 5-8 random photos from all accessible albums
- Each photo links to source album
- Refreshes on page reload

## Epic 1: Album List Home Page

Users can authenticate, view their album list on the home page with album cards showing random photo samples and density indicators, and filter albums by owner
with a fully functional filter.

### Story 1.1: Project Foundation Setup

As a developer,
I want to set up Material UI with the dark theme and remove Tailwind CSS,
So that the project has a consistent design system foundation.

**Acceptance Criteria:**

**Given** the project currently uses Tailwind CSS
**When** I set up the Material UI foundation
**Then** Tailwind CSS dependencies are completely removed from package.json and configuration files
**And** Material UI (@mui/material ^6.x, @mui/icons-material ^6.x) is installed
**And** Emotion dependencies (@emotion/react ^11.x, @emotion/styled ^11.x) are installed
**And** MUI theme is configured in `components/theme/theme.ts` with:

- Dark mode as default (`mode: 'dark'`)
- Brand color #185986 as primary color
- Background #121212 and surface #1e1e1e
- Text colors (primary: #ffffff, secondary: rgba(255,255,255,0.7))
  **And** MUI breakpoint system is configured: xs (<600px), sm (600px), md (960px), lg (1280px)
  **And** ThemeProvider wraps the application in root layout
  **And** No Tailwind classes remain in any component files
  **And** The application builds successfully with `npm run build`
  **And** Unit tests pass with `npm run test`

---

### Story 1.2: State Management Migration

As a developer,
I want to migrate the catalog state management from the existing web app,
So that I can reuse battle-tested state logic with 90+ tests.

**Acceptance Criteria:**

**Given** the existing state management in `web/src/core/catalog/` has 90+ passing tests
**When** I migrate the state management to web-nextjs
**Then** state management code is copied to `web-nextjs/domains/catalog/` maintaining folder structure:

- `language/` (state types - CatalogViewerState)
- `actions.ts` (reducer and actions)
- `album-create/`, `album-edit-*/`, `album-delete/` (thunk folders)
  **And** a new fetch adapter is created at `domains/catalog/adapters/fetch-adapter.ts` replacing axios
  **And** the fetch adapter works in both Server Components and Client Components
  **And** the fetch adapter implements the same interface as the original CatalogAPIAdapter
  **And** custom image loader is created at `libs/image-loader.ts` mapping:
- width ≤40 → quality=blur
- width ≤500 → quality=medium
- width ≤1200 → quality=high
- width >1200 → quality=full
  **And** image loader is configured in `next.config.ts`
  **And** error boundaries are created: `app/error.tsx` and `app/(authenticated)/error.tsx`
  **And** not-found pages are created: `app/not-found.tsx`
  **And** all migrated tests continue to pass (run with `npm run test`)

---

### Story 1.3: Basic Album List Display

As a user,
I want to view my album list on the authenticated home page,
So that I can see all albums I own and albums shared with me in chronological order.

**Acceptance Criteria:**

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

---

### Story 1.4: Album Card Enhancements

As a user,
I want to see preview photos and activity indicators on each album card,
So that I can quickly understand what's in each album before clicking.

**Acceptance Criteria:**

**Given** the basic album list is displayed
**When** I view an album card
**Then** each album card displays 3-4 random photo thumbnails in a preview grid
**And** random photos are fetched from the existing API endpoint for random media
**And** thumbnails use the Next.js Image component with the custom loader
**And** thumbnails request quality=medium (appropriate for card preview size)
**And** thumbnail images have blur placeholders while loading (quality=blur)
**And** density color-coding is applied based on photos-per-day calculation:

- High density (>10 photos/day): Warmer color accent or bolder display
- Medium density (3-10 photos/day): Neutral color
- Low density (<3 photos/day): Cooler color or lighter display
  **And** visual density indicator appears on the card (FR7)
  **And** sharing status is displayed with user avatars when album is shared (FR43, FR5)
  **And** user avatars and names are loaded from API (FR32)
  **And** the card maintains responsive layout on mobile, tablet, and desktop
  **And** images load progressively without blocking card rendering
  **And** card appearance uses MUI sx prop for styling (no inline styles)
  **And** brand color #185986 is used for primary interactive elements

---

### Story 1.5: Album Filtering

As a user,
I want to filter albums by owner,
So that I can focus on my own albums or view all albums including shared ones.

**Acceptance Criteria:**

**Given** I am viewing the album list with multiple albums from different owners
**When** I use the owner filter
**Then** a filter control is displayed above the album grid with options:

- "All Albums" - shows all accessible albums (owned + shared)
- "My Albums" - shows only albums I own
- Specific owner names - shows albums owned by that specific owner
  **And** the filter uses MUI Select or ToggleButtonGroup component
  **And** the active filter state is visually indicated (FR41)
  **And** filtering is handled client-side using the migrated catalog state reducer
  **And** the filter action dispatches to the reducer to update displayed albums (FR3)
  **And** the filtered album list updates immediately without API call
  **And** the filter state persists during navigation within the session
  **And** "My Albums" filter shows only albums where I am the owner (FR1)
  **And** "All Albums" shows owned albums (FR1) and shared albums (FR2)
  **And** the filter control is responsive and works on mobile, tablet, and desktop
  **And** keyboard navigation works for the filter control (tab, arrow keys, enter)
  **And** the filter state is cleared when logging out
  **And** changing filter maintains scroll position in the album list

---

## Epic 2: Photo Viewing & Navigation

Users can view photos within albums grouped by day, navigate smoothly with keyboard and touch controls, see progressive image loading, zoom into details, and
experience smooth transitions.

### Story 2.1: Album Photo Grid View

As a user,
I want to view all photos in an album grouped by the day they were captured,
So that I can browse through the album chronologically and see the story of that time period.

**Acceptance Criteria:**

**Given** I have clicked on an album from the home page
**When** I view the album page
**Then** the page is located at `app/(authenticated)/owners/[ownerId]/albums/[albumId]/page.tsx` (Server Component)
**And** the Server Component fetches album details and photos from the existing REST API
**And** photos are passed as initial data to Client Component wrapper
**And** photos are displayed in a responsive grid:

- Mobile (xs): 2 columns
- Tablet (sm): 3 columns
- Desktop (md): 4 columns
- Large Desktop (lg): 5 columns
  **And** photos are grouped by capture date with date headers separating each day (FR10, FR17)
  **And** date headers display in format "July 15, 2026" or similar
  **And** photos within each day are ordered chronologically
  **And** each photo thumbnail uses Next.js Image component with custom loader
  **And** thumbnails request quality=medium (appropriate for grid display)
  **And** progressive loading displays blur placeholder immediately (quality=blur) (FR14)
  **And** full quality loads in background after blur is displayed
  **And** responsive `sizes` attribute optimizes image requests for screen dimensions (NFR6)
  **And** a back button/link navigates to the album list home page (FR16)
  **And** album name and metadata are displayed in the page header
  **And** loading skeleton is shown while fetching photos
  **And** empty state message is displayed if album has no photos
  **And** error handling displays message with "Try Again" button if fetch fails
  **And** the grid layout uses MUI sx prop with responsive breakpoints
  **And** clicking a photo thumbnail navigates to the photo viewer

---

### Story 2.2: Full-Screen Photo Viewer with Modal Interception

As a user,
I want to view photos in full-screen mode with smooth transitions,
So that I can see photo details clearly and enjoy a focused viewing experience.

**Acceptance Criteria:**

**Given** I am viewing the album photo grid
**When** I click on a photo thumbnail
**Then** NextJS parallel route intercepts the navigation using `@modal/(.)photos/[photoId]/`
**And** a modal opens overlaying the photo grid (grid remains visible in background)
**And** the modal displays the photo in full-screen view (FR11)
**And** the URL updates to `/owners/[ownerId]/albums/[albumId]/photos/[photoId]`
**And** the photo loads progressively using Next.js Image with custom loader:

- Blur placeholder displays immediately (quality=blur, <2KB)
- Medium quality loads next (quality=medium, <500ms target)
- High quality loads for full-screen viewing (quality=high)
- Full original quality available for zoom (quality=full)
  **And** progressive loading completes within 3 seconds on slow networks (NFR7)
  **And** the photo is optimized for current screen size (mobile/tablet/desktop) (NFR6)
  **And** smooth transition animation displays when opening photo (if MUI provides) (FR42, NFR4)
  **And** photo metadata is displayed (date captured, album name)
  **And** a close button (X) is visible in the modal
  **And** clicking the close button or clicking outside returns to album grid
  **And** the modal maintains the album grid scroll position after closing
  **And** fallback route `photos/[photoId]/page.tsx` handles direct URL access (shows full page viewer)
  **And** refreshing on photo URL loads the full page viewer correctly
  **And** the layout component renders both `{children}` and `{modal}` slots
  **And** error boundary handles photo not found (displays not-found.tsx)
  **And** visual transitions maintain 60fps performance (NFR1)

---

### Story 2.3: Keyboard Navigation

As a user,
I want to navigate photos using keyboard shortcuts,
So that I can efficiently browse through albums without using the mouse.

**Acceptance Criteria:**

**Given** I am viewing a photo in full-screen mode
**When** I use keyboard controls
**Then** pressing RIGHT ARROW navigates to the next photo in sequence (FR12)
**And** pressing LEFT ARROW navigates to the previous photo in sequence (FR12)
**And** pressing ESC closes the photo viewer and returns to album grid (FR12)
**And** keyboard navigation responds without perceptible lag (<100ms) (NFR2, NFR3)
**And** navigation wraps to first photo when on last photo and next is pressed
**And** navigation wraps to last photo when on first photo and previous is pressed
**And** photo position indicator shows current position (e.g., "5 of 47")
**And** each keyboard action updates the URL to reflect current photo
**And** the photo transitions smoothly when navigating (if MUI provides naturally)
**And** progressive loading continues to work during keyboard navigation
**And** keyboard focus is managed correctly (visible focus indicators)
**And** keyboard shortcuts don't conflict with browser defaults (NFR23)
**And** pressing ENTER from the grid opens the focused photo

**Given** I am viewing the album grid
**When** I use keyboard controls
**Then** TAB navigates between photo thumbnails
**And** ENTER opens the focused photo in full-screen view
**And** focus indicators are clearly visible using brand color #185986
**And** focus management follows logical order (left to right, top to bottom)
**And** ESC from grid navigates back to album list

---

### Story 2.4: Touch Gestures and Mobile Navigation

As a mobile or tablet user,
I want to navigate photos using natural touch gestures,
So that I can browse albums comfortably on my device.

**Acceptance Criteria:**

**Given** I am viewing photos on a mobile or tablet device
**When** I use touch gestures in full-screen photo view
**Then** swiping LEFT navigates to the next photo (FR13)
**And** swiping RIGHT navigates to the previous photo (FR13)
**And** tap on the photo or swipe DOWN closes the viewer and returns to grid
**And** pinch-to-zoom gesture zooms into photo details (FR15)
**And** zoomed photo can be panned by dragging
**And** double-tap zooms to fit or zooms in (standard behavior)
**And** touch interactions feel responsive with immediate visual feedback (<100ms) (NFR24)
**And** gestures feel natural and performant (NFR26)
**And** animations maintain acceptable frame rates on mobile (NFR27)
**And** mobile performance doesn't degrade below 30fps minimum
**And** swipe gestures have appropriate threshold to prevent accidental navigation
**And** smooth transitions occur between photos during swipe (FR42)

**Given** I am viewing the album grid on mobile
**When** I interact with the grid
**Then** tapping a photo opens it in full-screen view
**And** scrolling the grid feels smooth and responsive
**And** touch targets are appropriately sized for mobile (minimum 44px)
**And** pull-to-refresh reloads the album (if supported by browser)
**And** the grid layout adapts correctly to mobile breakpoints (<600px)
**And** photo thumbnails load progressively on mobile networks
**And** the back button navigation works correctly

**Given** I am on an older mobile device (like Marie's iPad)
**When** I use the photo viewer
**Then** the interface works smoothly without performance degradation
**And** touch gestures remain responsive
**And** images load progressively to accommodate slower networks
**And** the interface recovers gracefully from errors

---

## Epic 3: Album Management

Album owners can create, edit, and delete albums with visual date selection showing real-time photo previews and photo counts that update as dates adjust.

### Story 3.1: Create Album with Visual Date Selection

As an album owner,
I want to create a new album with visual date selection that shows me which photos will be included,
So that I can confidently organize photos without guessing at dates.

**Acceptance Criteria:**

**Given** I am an authenticated user identified as an owner (FR36)
**When** I create a new album
**Then** a "Create Album" button is visible on the home page for owners
**And** clicking the button opens a create album dialog using MUI Dialog component
**And** the dialog contains:

- Album name text field (required) (FR18)
- Start date picker (MUI DatePicker)
- End date picker (MUI DatePicker)
- Optional toggle for custom folder name (FR19)
- Custom folder name text field (shown only if toggle enabled)
- Photo preview panel showing ~12 thumbnail samples
- Photo count display (e.g., "47 photos will be included")
- Cancel and Create buttons
  **And** when I select or adjust start/end dates, the photo preview updates in real-time
  **And** photo count updates dynamically as dates change (e.g., "47 photos" → "52 photos")
  **And** preview thumbnails are fetched from API based on selected date range
  **And** thumbnails use Next.js Image with quality=medium
  **And** date range validation ensures end date is after start date (FR45)
  **And** album name validation ensures name is not empty (FR46)
  **And** validation errors display inline below the relevant field
  **And** warning is shown if date selection would orphan existing photos (FR23)
  **And** Create button is disabled until validation passes
  **And** clicking Create button:
- Calls existing API endpoint to create album
- Uses migrated catalog thunk for album creation
- Shows loading state on button during API call
- Closes dialog on success
- Displays success message/toast (FR25)
- Refreshes album list to show new album
- Calls router.refresh() for NextJS cache invalidation
  **And** clicking Cancel button closes dialog without creating album
  **And** error handling displays clear message if creation fails (FR47)
  **And** recovery option (retry button) is provided on error (FR50)
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, Enter to submit, ESC to cancel) (NFR22)
  **And** focus management is clear and logical within dialog
  **And** new album appears in chronological order in the album list

---

### Story 3.2: Edit Album Name

As an album owner,
I want to edit the name of my album,
So that I can fix typos or better describe the album content.

**Acceptance Criteria:**

**Given** I am viewing an album I own
**When** I edit the album name
**Then** an "Edit Name" button/icon is visible only on albums I own (FR30)
**And** the edit option is not visible on shared albums I don't own
**And** clicking edit opens an edit name dialog using MUI Dialog
**And** the dialog contains:

- Album name text field pre-filled with current name
- Current album metadata displayed for context
- Cancel and Save buttons
  **And** album name validation ensures name is not empty (FR46)
  **And** validation error displays inline if name is cleared
  **And** Save button is disabled if validation fails
  **And** clicking Save button:
- Calls existing API endpoint to update album name
- Uses migrated catalog thunk for album edit
- Shows loading state on button during API call
- Closes dialog on success
- Displays success message (FR25)
- Updates album name in the album list immediately
- Calls router.refresh() for NextJS cache invalidation
  **And** clicking Cancel button closes dialog without saving changes
  **And** error handling displays clear message if update fails (FR47)
  **And** recovery option (retry button) is provided on error (FR50)
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, Enter to submit, ESC to cancel)
  **And** focus is set to the name field when dialog opens
  **And** the updated name appears throughout the UI (album list, album page header)

---

### Story 3.3: Edit Album Date Range

As an album owner,
I want to edit the date range of my album with a live photo preview,
So that I can adjust which photos are included and see the results before saving.

**Acceptance Criteria:**

**Given** I am viewing an album I own
**When** I edit the album date range
**Then** an "Edit Dates" button/icon is visible only on albums I own (FR30)
**And** the edit option is not visible on shared albums I don't own
**And** clicking edit opens an edit date range dialog using MUI Dialog
**And** the dialog contains:

- Current album name displayed for context
- Start date picker pre-filled with current start date
- End date picker pre-filled with current end date
- Photo preview panel showing ~12 thumbnail samples from current range
- Photo count display (e.g., "47 photos in this album")
- Cancel and Save buttons
  **And** when I adjust start or end dates, the photo preview updates in real-time (FR21)
  **And** photo count updates dynamically as dates change
  **And** preview thumbnails re-fetch from API based on new date range
  **And** thumbnails use Next.js Image with quality=medium
  **And** date range validation ensures end date is after start date (FR45)
  **And** validation error displays if end date is before or equal to start date
  **And** warning is shown if new date range would orphan existing photos (FR23)
  **And** orphan warning displays which photos would be affected
  **And** Save button is disabled until validation passes
  **And** clicking Save button:
- Calls existing API endpoint to update album dates
- Uses migrated catalog thunk for album edit
- Shows loading state on button during API call
- Triggers backend re-indexing of photos (FR24)
- Closes dialog on success
- Displays success message (FR25)
- Refreshes album to show updated photo set
- Calls router.refresh() for NextJS cache invalidation
  **And** clicking Cancel button closes dialog without saving changes
  **And** error handling displays clear message if update fails (FR47)
  **And** recovery option (retry button) is provided on error (FR50)
  **And** network failures do not cause data loss (NFR34)
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, arrow keys to adjust dates, Enter to save, ESC to cancel)
  **And** focus management is clear within dialog
  **And** the album page refreshes to show correct photos after save

---

### Story 3.4: Delete Album

As an album owner,
I want to delete an album I own,
So that I can remove albums I no longer need.

**Acceptance Criteria:**

**Given** I am viewing an album I own
**When** I delete the album
**Then** a "Delete Album" button/icon is visible only on albums I own (FR30)
**And** the delete option is not visible on shared albums I don't own
**And** clicking delete opens a confirmation dialog using MUI Dialog
**And** the confirmation dialog contains:

- Warning message explaining deletion is permanent
- Album name displayed for confirmation
- Number of photos that will remain (not deleted, just album removed)
- Cancel and Delete buttons
  **And** Delete button uses warning/error color (red) to indicate destructive action
  **And** clicking Delete button:
- Calls existing API endpoint to delete album
- Uses migrated catalog thunk for album deletion
- Shows loading state on button during API call
- Closes dialog on success
- Displays success message (FR25)
- Removes album from the album list immediately
- Navigates to home page if currently viewing the deleted album
- Calls router.refresh() for NextJS cache invalidation
  **And** clicking Cancel button closes dialog without deleting album
  **And** error handling displays clear message if deletion fails (FR47)
  **And** recovery option is provided on error (FR50)
  **And** if album is not found during deletion, appropriate error displays (FR49)
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, Enter on Cancel as default, must explicitly select Delete)
  **And** focus defaults to Cancel button (safer default for destructive action)
  **And** the album is removed from all views after successful deletion

---

## Epic 4: Sharing & Access Control

Album owners can share albums with users via email, revoke access, view who has access with avatars and names; shared viewers see ownership information and
clear UI distinction between owner and viewer capabilities.

### Story 4.1: Share Album with Users

As an album owner,
I want to share my album with specific users via their email addresses,
So that family members can view and enjoy the photos.

**Acceptance Criteria:**

**Given** I am viewing an album I own
**When** I share the album
**Then** a "Share" button/icon is visible only on albums I own (FR30)
**And** the share option is not visible on shared albums I don't own
**And** clicking share opens a sharing dialog using MUI Dialog
**And** the dialog contains:

- Album name displayed for context
- Email address input field with validation
- "Grant Access" button
- List of currently shared users with their avatars and names
- Close button
  **And** the email input field validates email format on blur (FR31)
  **And** validation error displays inline if email format is invalid
  **And** Grant Access button is disabled if email is invalid or empty
  **And** when I enter a valid email and click Grant Access:
- Calls existing API endpoint to grant album access
- Uses migrated catalog thunk for sharing operation
- Shows loading state during API call
- Fetches user profile information (name, picture) from API (FR32)
- Adds user to the shared users list with avatar and name
- Displays success message (FR25)
- Clears email input field for next entry
- Calls router.refresh() for NextJS cache invalidation
  **And** the shared users list displays:
- User avatar image
- User name
- User email
- Remove/revoke button (X icon) next to each user
  **And** shared users are loaded from API when dialog opens (FR28)
  **And** I can see all users who currently have access to the album (FR28)
  **And** error handling displays clear message if sharing fails (FR48)
  **And** recovery option (retry button) is provided on error (FR50)
  **And** network failures do not cause data loss (NFR34)
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, Enter to grant access, ESC to close)
  **And** focus is set to email input when dialog opens
  **And** shared user avatars appear on the album card after sharing (FR43, FR5)

**Given** I am a shared viewer of an album
**When** I view the album
**Then** I can see who owns the album (name and avatar) (FR29)
**And** I cannot see the Share button (not an owner) (FR30)
**And** the UI clearly distinguishes that I am viewing, not owning

---

### Story 4.2: Revoke Album Access

As an album owner,
I want to revoke access from users I previously shared with,
So that I can control who continues to have access to my albums.

**Acceptance Criteria:**

**Given** I have an album shared with multiple users
**When** I revoke access from a user
**Then** the sharing dialog shows all currently shared users (FR28)
**And** each shared user in the list has a remove/revoke button (X icon)
**And** clicking the revoke button opens a confirmation dialog
**And** the confirmation dialog contains:

- User name and avatar being removed
- Album name for context
- Warning that user will lose access to the album
- Cancel and Revoke Access buttons
  **And** Revoke Access button uses warning color to indicate action
  **And** clicking Revoke Access button:
- Calls existing API endpoint to revoke album access
- Uses migrated catalog thunk for revoke operation
- Shows loading state during API call
- Removes user from the shared users list on success
- Displays success message (FR25, FR27)
- Updates album card to remove user's avatar if no longer shared
- Calls router.refresh() for NextJS cache invalidation
  **And** clicking Cancel button closes confirmation without revoking
  **And** error handling displays clear message if revoke fails (FR48)
  **And** recovery option (retry button) is provided on error (FR50)
  **And** the revoked user no longer sees the album in their album list
  **And** dialog is responsive and works on mobile, tablet, desktop
  **And** keyboard navigation works (Tab, Enter on Cancel as default, must explicitly select Revoke)
  **And** focus defaults to Cancel button (safer default)
  **And** the sharing dialog remains open after revoke to allow additional changes
  **And** I can close the sharing dialog when done managing access
  **And** album sharing indicators update across the UI after revoke

---

## Epic 5: Random Photo Discovery

Users can discover forgotten memories through random photo highlights on the home page (5-8 photos from all accessible albums) that link directly to source
albums and refresh on each page reload.

### Story 5.1: Random Photo Highlights UI with Mock Data

As a user,
I want to see random photo highlights on the home page,
So that I can discover and rediscover forgotten memories from my collection.

**Acceptance Criteria:**

**Given** I am viewing the authenticated home page with albums
**When** the page loads
**Then** a "Your Memories" or "Highlights" section is displayed at the top of the page above the album list
**And** the section displays 5-8 random photos in a horizontal layout
**And** the mock implementation selects random photos from the album cards displayed below (using the 3-4 random samples already loaded per album card)
**And** each photo is displayed using Next.js Image component with custom loader
**And** photos use quality=medium for the highlights display
**And** progressive loading displays blur placeholder (quality=blur) immediately
**And** the layout is responsive:

- Mobile (xs): Horizontal scroll with 2-3 photos visible
- Tablet (sm): 4-5 photos visible
- Desktop (md): 6-8 photos visible
  **And** each photo is clickable and links to its source album (FR9)
  **And** clicking a photo navigates to `/owners/[ownerId]/albums/[albumId]`
  **And** the section has a clear heading: "Your Memories" or "Highlights"
  **And** photos are displayed with subtle spacing and styling (MUI sx prop)
  **And** the section uses brand color #185986 for accent elements
  **And** loading skeleton is shown while photos are being selected
  **And** the section works on mobile, tablet, and desktop devices
  **And** the random selection changes on page reload (different subset from album cards)
  **And** the feature is demoable and validates the UX flow (FR8)
  **And** if no albums exist, the highlights section is hidden
  **And** the mock implementation is clearly commented as temporary

---

### Story 5.2: Backend API for Random Photo Discovery

As a backend developer,
I want to implement a REST API endpoint that returns random photos from a user's accessible albums,
So that the frontend can display diverse memory highlights without relying on mock data.

**Acceptance Criteria:**

**Given** the backend API infrastructure exists
**When** implementing the random photos endpoint
**Then** a new REST API endpoint is created: `GET /api/v1/photos/random`
**And** the endpoint accepts an optional query parameter `count` (default: 8, max: 20)
**And** the endpoint returns random photos from all albums accessible to the authenticated user (owned + shared)
**And** the response includes for each photo:

- Photo/media ID
- Album ID
- Owner ID
- Photo metadata (date captured, if available)
- Image URL template or media ID for frontend image loader
  **And** the endpoint respects user permissions (only photos from accessible albums)
  **And** randomization algorithm ensures variety across albums (not all from one album)
  **And** the endpoint handles edge cases:
- User has no albums (returns empty array)
- User has fewer photos than requested (returns all available)
- Authentication failures (returns 401)
  **And** the endpoint follows existing API conventions and patterns
  **And** appropriate error responses are returned with clear messages
  **And** the endpoint is tested with unit tests
  **And** the endpoint is deployed to the backend environment

**NOTE:** This story is backend work in `api/lambdas/` (out of scope for web-nextjs frontend, but required for Epic 5 completion)

---

### Story 5.3: Connect UI to Random Photos API

As a user,
I want the random photo highlights to show diverse photos from across my entire collection,
So that I can discover memories I haven't seen in a while.

**Acceptance Criteria:**

**Given** the backend random photos API is deployed and available
**When** the home page loads
**Then** the highlights section fetches random photos from `GET /api/v1/photos/random?count=8`
**And** the API call is made from the Server Component during initial page load
**And** random photos data is passed to the Client Component as props
**And** the mock implementation from Story 5.1 is completely removed
**And** photos are displayed using the same UI layout from Story 5.1
**And** each photo uses Next.js Image component with custom loader
**And** photo URLs are constructed using mediaId from API response
**And** clicking a photo navigates to `/owners/[ownerId]/albums/[albumId]` using IDs from API
**And** the highlights refresh on each page reload (new API call fetches different randoms) (FR8)
**And** progressive loading works (blur → medium quality)
**And** loading state is shown while API call is in progress
**And** error handling displays graceful fallback if API fails (hide section or show friendly message)
**And** if API returns empty array (no photos), the section is hidden
**And** the feature works across mobile, tablet, and desktop
**And** the implementation uses the migrated fetch adapter for consistency
**And** photos display variety across different albums (leveraging backend randomization)
**And** the feature fulfills FR8 (discover random photos) and FR9 (navigate to source album)
**And** users can successfully rediscover forgotten memories through this feature
