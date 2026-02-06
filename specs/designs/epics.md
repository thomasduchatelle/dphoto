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
- **FR2:** All users can view albums shared with them
- **FR3:** Users can filter albums by owner (my albums, all albums, specific owner)
- **FR4:** Users can view album metadata (name, date range, media count, owner information)
- **FR5:** Users can see which users an album is shared with
- **FR6:** Users can view albums in chronological order
- **FR7:** Users can see visual indicators of album activity density (temperature)
- **FR8:** Users can discover random photos from across their accessible collection
- **FR9:** Users can navigate from random photo highlights to the source album
- **FR10:** System indicates selected albums and active filters
- **FR11:** System displays album sharing status with user avatars

**Photo Viewing & Navigation:**

- **FR12:** Users can view photos grouped by capture date within an album
- **FR13:** Users can open individual photos in full-screen view
- **FR14:** Users can navigate between photos using keyboard controls (arrow keys, esc, enter)
- **FR15:** Users can navigate between photos using touch gestures (swipe on mobile/tablet)
- **FR16:** Users can zoom into photo details
- **FR17:** Users can navigate back to album list from photo view
- **FR18:** Users can see date headers separating photos by day

**Album Management:**

- **FR19:** Only owners can create albums
- **FR20:** Albums can only be managed by their owner
- **FR21:** Album owners can create new albums by specifying name and date range
- **FR22:** Album owners can specify custom folder names for albums (optional)
- **FR23:** Album owners can edit album names
- **FR24:** Album owners can edit album date ranges
- **FR25:** Album owners can delete albums they own
- **FR26:** System validates that date edits don't orphan media
- **FR27:** System re-indexes photos when album date ranges change
- **FR28:** System provides feedback when album operations succeed or fail

**Sharing & Access Control:**

- **FR29:** Album owners can share albums with other users via email address
- **FR30:** Album owners can revoke access from users they previously shared with
- **FR31:** Album owners can view who has access to their albums
- **FR32:** Shared viewers can see who owns albums shared with them
- **FR33:** System distinguishes between owner capabilities and viewer capabilities in UI
- **FR34:** System validates email addresses when granting access
- **FR35:** System loads user profile information (name, picture) for display

**User Profile:**

- **FR36:** Users can view their profile information (name, picture)

**Visual Presentation & UI State:**

- **FR37:** System displays discrete loading feedback when a page is loaded in the background (no skeleton or big spinner)
- **FR38:** System invites user to the next step when no albums exist (create album) or no medias exist (upload medias or update album)
- **FR39:** System shows explicit errors (even if technical) when something goes wrong
- **FR40:** System provides responsive layouts for mobile, tablet, and desktop devices

**Error Handling & Validation:**

- **FR41:** System validates album date ranges (end date after start date)
- **FR42:** System validates album names are not empty
- **FR43:** System handles and displays errors for failed album operations
- **FR44:** System handles and displays errors for failed sharing operations
- **FR45:** System handles and displays errors when albums are not found
- **FR46:** System provides recovery options when operations fail

### Non-Functional Requirements

**Integration:**

- **NFR1:** Frontend must consume existing REST API endpoints without requiring backend modifications
- **NFR2:** System must handle API response times gracefully with appropriate loading states
- **NFR3:** System must handle API errors with clear user feedback and retry mechanisms
- **NFR4:** System must maintain data contract compatibility with existing API response formats

**Usability:**

- **NFR5:** Focus management must be clear and logical during navigation and in dialogs
- **NFR6:** Keyboard shortcuts must not conflict with browser defaults
- **NFR7:** Touch interactions must feel responsive with appropriate visual feedback
- **NFR8:** Layouts must adapt appropriately across mobile (<600px), tablet (600-960px), and desktop (>960px) breakpoints
- **NFR9:** Mobile gestures (swipe, pinch-to-zoom) must feel natural and performant
- **NFR10:** System must use Material UI components; modern web features (CSS Grid, Flexbox, ES2020+) may be used only when Material UI doesn't provide the
  required capabilities
- **NFR11:** No support required for Internet Explorer or older browser versions

**Reliability:**

- **NFR12:** System operates on best-effort availability basis (no uptime SLA required)
- **NFR13:** Manual page refresh is acceptable recovery mechanism for transient errors
- **NFR14:** System must provide clear error messages when operations fail with guidance on recovery

### Additional Requirements

Refer on @specs/designs/architecture.md for additional requirements and technical directions.

### FR Coverage Map

**Epic 1 (Album List Home Page):**

- FR1: Album owners can view all albums they own
- FR2: All users can view albums shared with them
- FR3: Users can filter albums by owner (my albums, all albums, specific owner)
- FR4: Users can view album metadata (name, date range, media count, owner information)
- FR5: Users can see which users an album is shared with
- FR6: Users can view albums in chronological order
- FR7: Users can see visual indicators of album activity density (temperature)
- FR10: System indicates selected albums and active filters
- FR11: System displays album sharing status with user avatars
- FR36: Users can view their profile information (name, picture)
- FR37: System displays discrete loading feedback when a page is loaded in the background
- FR38: System invites user to the next step when no albums exist or no medias exist
- FR39: System shows explicit errors when something goes wrong
- FR40: System provides responsive layouts for mobile, tablet, and desktop devices

**Epic 2 (Photo Viewing & Navigation):**

- FR12: Users can view photos grouped by capture date within an album
- FR13: Users can open individual photos in full-screen view
- FR14: Users can navigate between photos using keyboard controls (arrow keys, esc, enter)
- FR15: Users can navigate between photos using touch gestures (swipe on mobile/tablet)
- FR16: Users can zoom into photo details
- FR17: Users can navigate back to album list from photo view
- FR18: Users can see date headers separating photos by day

**Epic 3 (Album Management):**

- FR19: Only owners can create albums
- FR20: Albums can only be managed by their owner
- FR21: Album owners can create new albums by specifying name and date range
- FR22: Album owners can specify custom folder names for albums (optional)
- FR23: Album owners can edit album names
- FR24: Album owners can edit album date ranges
- FR25: Album owners can delete albums they own
- FR26: System validates that date edits don't orphan media
- FR27: System re-indexes photos when album date ranges change
- FR28: System provides feedback when album operations succeed or fail
- FR41: System validates album date ranges (end date after start date)
- FR42: System validates album names are not empty
- FR43: System handles and displays errors for failed album operations
- FR45: System handles and displays errors when albums are not found
- FR46: System provides recovery options when operations fail

**Epic 4 (Sharing & Access Control):**

- FR29: Album owners can share albums with other users via email address
- FR30: Album owners can revoke access from users they previously shared with
- FR31: Album owners can view who has access to their albums
- FR32: Shared viewers can see who owns albums shared with them
- FR33: System distinguishes between owner capabilities and viewer capabilities in UI
- FR34: System validates email addresses when granting access
- FR35: System loads user profile information (name, picture) for display
- FR44: System handles and displays errors for failed sharing operations

**Epic 5 (Random Photo Discovery):**

- FR8: Users can discover random photos from across their accessible collection
- FR9: Users can navigate from random photo highlights to the source album

**Epic 6 (Album Preview):**

- FR8: Users can discover random photos from across their accessible collection (partial)

**Total FRs Covered:** 46 / 46 ✅

## Epic List

### Epic 1: Album List Home Page

Users can view their album list on the home page with album cards showing metadata and density indicators, and filter albums by owner with a fully functional
filter.

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR10, FR11, FR36, FR37, FR38, FR39, FR40 (14 FRs)

**NFRs addressed:** NFR8 (responsive breakpoints), NFR2 (loading states)

**Additional Requirements:**

- Remove Tailwind, install Material UI with dark theme (#185986)
- Migrate state management from `web/src/core/catalog/` (90+ tests)
- Replace axios with fetch adapter
- Set up NextJS App Router with Server/Client component pattern
- Configure custom image loader
- Establish MUI theme, breakpoints, error boundaries
- Density color-coding for album activity
- Colocation principle for components
- Album paths: `/albums/[ownerId]/[albumId]`
- Redirect `/albums` and `/albums/[ownerId]` to `/`

### Epic 2: Photo Viewing & Navigation

Users can view photos within albums grouped by day, navigate smoothly with keyboard and touch controls, and zoom into details.

**FRs covered:** FR12, FR13, FR14, FR15, FR16, FR17, FR18 (7 FRs)

**NFRs addressed:** NFR1-NFR4 (API integration), NFR5-NFR6 (keyboard and focus), NFR7 (touch), NFR8-NFR9 (mobile gestures, responsive layouts)

**Additional Requirements:**

- NextJS parallel routes for photo modal interception: `@modal/(.)photos/[photoId]/`
- Album paths: `/albums/[ownerId]/[albumId]`
- Redirect `/albums` and `/albums/[ownerId]` to `/`
- Next.js Image component with custom loader
- MUI sx prop for responsive layouts

### Epic 3: Album Management

Album owners can create, edit, and delete albums with visual date selection showing photo previews.

**FRs covered:** FR19, FR20, FR21, FR22, FR23, FR24, FR25, FR26, FR27, FR28, FR41, FR42, FR43, FR45, FR46 (15 FRs)

**NFRs addressed:** NFR3 (error handling), NFR14 (clear error messages), NFR5 (keyboard navigation)

**Additional Requirements:**

- Visual date picker with live photo preview (~12 thumbnails)
- Multi-step album creation UX for better space utilization
- Validation, orphan photo detection, error handling

### Epic 4: Sharing & Access Control

Album owners can share albums with users via email, revoke access, view who has access with avatars and names; shared viewers see ownership information and
clear UI distinction between owner and viewer capabilities.

**FRs covered:** FR29, FR30, FR31, FR32, FR33, FR34, FR35, FR44 (8 FRs)

**NFRs addressed:** NFR3 (error handling)

**Additional Requirements:**

- User profile information loading (name, picture)
- Email validation
- Error handling for sharing operations

## Epic 1: Album List Home Page

Users can view their album list on the home page with album cards showing metadata and density indicators, and filter albums by owner with a fully functional
filter.

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
  **And** custom image loader is created at `libs/image-loader.ts` mapping to backend-supported widths (360, 1440, 2400)
  **And** image loader is configured in `next.config.ts`
  **And** error boundaries are created: `app/error.tsx` and `app/(authenticated)/error.tsx`
  **And** not-found pages are created: `app/not-found.tsx`
  **And** all migrated tests continue to pass (run with `npm run test`)

---

### Story 1.3: Basic Album List Display

As a user,
I want to view my album list on the authenticated home page,
So that I can see all albums I own and albums shared with me.

**Acceptance Criteria:**

**Given** I am an authenticated user with access to albums
**When** I navigate to the home page at `/`
**Then** I see a list of albums displayed as clickable text links
**And** each album displays its name (FR4)
**And** albums are displayed in chronological order (newest first) (FR6)
**And** I can see albums I own (FR1) and albums shared with me (FR2)
**And** clicking an album navigates to `/albums/[ownerId]/[albumId]`
**And** when no albums exist, I see a message inviting me to create an album (FR38)
**And** if an error occurs loading albums, I see an error message with a "Try Again" button (FR39)

---

### Story 1.4: UI Components & Layout

As a developer,
I want to create reusable UI components and establish the application layout,
So that the application has a consistent visual structure and branded design system.

**Acceptance Criteria:**

**Given** the MUI theme and foundation are set up (Story 1.1)
**When** I create the UI component library
**Then** the following components are created and tested:

**Layout Components:**

- **AppLayout**: Main application layout wrapper
    - Displays app header with navigation
    - Displays user profile information (name, picture) in header
    - Provides content area for page content
    - Responsive across mobile, tablet, desktop

- **AppHeader**: Application header component
    - Shows application logo/title
    - Shows authenticated user profile (avatar, name)
    - Uses brand color #185986 for primary elements

**Album Components:**

- **AlbumCard**: Album display card component
    - Shows album name as prominent heading
    - Shows date range (e.g., "Jan 15 - Feb 28, 2026")
    - Shows media count (e.g., "47 photos")
    - Shows owner information (name, avatar) when shared
    - Shows density indicator using color-coding:
        - High density (>10 photos/day): Warmer color accent
        - Medium density (3-10 photos/day): Neutral color
        - Low density (<3 photos/day): Cooler color
    - Shows sharing status with user avatars when shared
    - Clickable to navigate to album
    - Responsive layout for mobile, tablet, desktop

- **AlbumGrid**: Responsive grid layout for album cards
    - Mobile (xs): 1 column
    - Tablet (sm): 2 columns
    - Desktop (md): 3 columns
    - Large Desktop (lg): 4 columns

**Feedback Components:**

- **PageLoadingIndicator**: Full-page loading feedback component
    - Shows when initial page is displayed
    - Discrete loading feedback
    - No skeleton or big spinner

- **NavigationLoadingIndicator**: Inline/overlay loading feedback component
    - Shows when NextJS navigates to album page after card click
    - Discrete, non-blocking loading indicator

- **ErrorDisplay**: Error message display component
    - Shows explicit error messages (including technical details)
    - Provides "Try Again" button for recovery
    - Clear and prominent display

- **EmptyState**: Empty state message component
    - Shows when no albums exist
    - Invites user to create first album
    - Shows when no media exists with invitation to upload

**User Components:**

- **UserAvatar**: User avatar display component
    - Shows user profile picture
    - Shows user initials as fallback
    - Multiple sizes (small, medium, large)

- **SharedByIndicator**: Shows sharing status for an album
    - Displays group of user avatars for users album is shared with
    - Shows multiple user avatars compactly
    - Handles overflow (e.g., "+3 more")
    - Used by AlbumCard to display sharing status

**And** all components use MUI sx prop for styling (no inline styles)
**And** all components use brand color #185986 for primary interactive elements
**And** all components have visual regression tests using existing test infrastructure
**And** error boundaries (`app/error.tsx`, `app/(authenticated)/error.tsx`) are updated to use ErrorDisplay component
**And** not-found pages (`app/not-found.tsx`) use consistent styling with ErrorDisplay
**And** AppLayout is integrated into the root layout wrapper
**And** all components are colocated with their tests following NextJS best practices

---

### Story 1.5: Style the Home Page

As a user,
I want to see a beautifully designed home page with my albums,
So that I can enjoy using the application and easily browse my collection.

**Acceptance Criteria:**

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

---

### Story 1.6: Album Filtering

As a user,
I want to filter albums by owner,
So that I can focus on my own albums or view all albums including shared ones.

**Acceptance Criteria:**

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

---

## Epic 2: Photo Viewing & Navigation

Users can view photos within albums grouped by day, navigate smoothly with keyboard and touch controls, and zoom into details.

### Story 2.1: Album Page URL Redirects

As a user,
I want invalid album URLs to redirect to the home page,
So that I don't see broken or incomplete pages.

**Acceptance Criteria:**

**Given** I navigate to an invalid album URL
**When** I access `/albums` or `/albums/[ownerId]`
**Then** the page redirects to `/` (home page)

---

### Story 2.2: Album Page UI Components

As a developer,
I want reusable UI components for the album page,
So that the album viewing experience is consistent and maintainable.

**UI Components to Create:**

1. **MediaGrid** - Displays medias grouped by day in a responsive grid
    - Props: medias grouped by day (array of `{ day: Date, medias: Media[] }`)
    - Responsive columns: Mobile (2), Tablet (3), Desktop (4), Large (5)
    - Photos are grouped by capture date with date headers separating each day (FR12, FR18)
    - Date headers display in the user's locale format
    - Does not handle empty state
    - Uses private sub-components (not tested or exposed): PhotoThumbnail, VideoThumbnail, DayGroupGrid, etc.

2. **AlbumPageHeader** - Top section of album page
    - Props: album name, metadata, back button handler
    - Shows album title and info
    - Back button to navigate home (FR17)

3. **EmptyAlbumState** - Message when album has no photos
    - Props: album name
    - Clear messaging
    - Consistent with brand styling

4. **PhotoViewer** - Full-screen photo viewer with navigation
    - Props: current photo, next media, previous media
    - Displays photo in full-screen view (FR13)
    - Close button (X) to return to album grid
    - Next/Previous navigation buttons
    - Photo position indicator (e.g., "5 of 47")
    - Photo metadata display (date captured, album name)
    - Keyboard navigation: RIGHT ARROW (next), LEFT ARROW (previous), ESC (close) (FR14)
    - Touch gestures: swipe left (next), swipe right (previous), tap/swipe down (close) (FR15)
    - Pinch-to-zoom and pan support (FR16)
    - Double-tap zoom support
    - Navigation wraps (first ↔ last)
    - URL updates with each navigation action
    - Handles photo not found errors
    - Keyboard focus management with visible indicators
    - Keyboard shortcuts don't conflict with browser defaults (NFR6)
    - Touch interactions responsive with immediate feedback (NFR7, NFR9)
    - Swipe gestures have appropriate threshold to prevent accidental navigation

**Acceptance Criteria:**

**Given** the components are implemented
**When** they are used in the album page
**Then** all 4 components are created and exported
**And** each component has clear prop interfaces
**And** components use MUI styling system
**And** components follow the brand color scheme (#185986)
**And** components have visual regression tests

---

### Story 2.3: Load and Display Album Photos

As a user,
I want to see photo thumbnails when viewing an album,
So that I can browse the photos in that album.

**Acceptance Criteria:**

**Given** I navigate to `/albums/[ownerId]/[albumId]`
**When** the page loads
**Then** the album details and photos are fetched from the REST API
**And** photos use NextJS custom image loader for backend integration
**And** thumbnails are the small version of the image, size optimised (NFR2)

---

### Story 2.4: Complete Album Page with Grouped Photos

As a user,
I want to view all photos in an album grouped by the day they were captured,
So that I can browse through the album chronologically and see the story of that time period.

**Acceptance Criteria:**

**Given** I am viewing an album at `/albums/[ownerId]/[albumId]`
**When** the page displays
**Then** the MediaGrid component is used to display medias, grouped by days
**And** a back button/link navigates to the album list home page (FR17)
**And** album name and metadata are displayed in the page header
**And** empty state message is displayed if album has no photos

---

### Story 2.5: Album Loading Progress Indicator

As a user,
I want to see a loading indicator when clicking an album link,
So that I know the page is loading and the app hasn't frozen.

**Acceptance Criteria:**

**Given** I am on the home page viewing the album list
**When** I click on an album link
**Then** the home page remains visible
**And** a thin progress bar appears at the top of the page
**And** the progress bar indicates the next page is loading
**And** once the album page loads, the page transitions to show the album
**And** the progress bar disappears

---

### Story 2.6: Full-Screen Photo Viewer with Modal Interception

As a user,
I want to view photos in full-screen mode,
So that I can see photo details clearly and enjoy a focused viewing experience.

**Acceptance Criteria:**

**Given** I am viewing the album photo grid
**When** I click on a photo thumbnail
**Then** a modal opens overlaying the photo grid (grid remains visible in background)
**And** the PhotoViewer component displays the photo in full-screen view (FR13)
**And** the URL updates to `/albums/[ownerId]/[albumId]/photos/[photoId]`
**And** the photo is optimized for current screen size (mobile/tablet/desktop) (NFR2)
**And** the modal maintains the album grid scroll position after closing
**And** refreshing on photo URL loads the photo viewer correctly

---

### Story 2.7: Photo Viewer Navigation and Interactions

As a user,
I want to navigate photos and interact with them using keyboard and touch controls,
So that I can efficiently browse through albums on any device.

**Acceptance Criteria:**

**Given** I am viewing a photo in the PhotoViewer
**When** I use navigation controls
**Then** all PhotoViewer functionality works as specified in Story 2.2 (keyboard, touch, zoom)
**And** keyboard navigation responds without perceptible lag (NFR2, NFR5)
**And** touch interactions feel responsive with immediate visual feedback (NFR7, NFR9)
**And** gestures feel natural and performant (NFR9)
**And** mobile performance doesn't degrade significantly

**Given** I am viewing the album grid on mobile
**When** I interact with the grid
**Then** tapping a photo opens it in full-screen view
**And** scrolling the grid feels smooth and responsive
**And** touch targets are appropriately sized for mobile (minimum 44px)
**And** pull-to-refresh reloads the album (if supported by browser)
**And** the grid layout adapts correctly to mobile breakpoints (<600px) (NFR8)
**And** the back button navigation works correctly

---

## Epic 3: Album Management

Album owners can create, edit, and delete albums with visual date selection showing photo previews.

### Story 3.1: Create Album with Visual Date Selection

As an album owner,
I want to create a new album with visual date selection that shows me which photos will be included,
So that I can confidently organize photos without guessing at dates.

**Acceptance Criteria:**

**Given** I am an authenticated user identified as an owner (FR19)
**When** I create a new album
**Then** a "Create Album" button is visible on the home page for owners
**And** clicking the button opens a create album dialog
**And** the dialog uses a multi-step UX to provide ample space for photo previews
**And** the dialog contains:
- Album name text field (required) (FR21)
- Start date picker
- End date picker
- Optional toggle for custom folder name (FR22)
- Custom folder name text field (shown only if toggle enabled)
- Photo preview panel showing ~12 thumbnail samples
- Cancel and Create buttons
**And** when I select or adjust start/end dates, the photo preview updates
**And** preview thumbnails are fetched from API based on selected date range
**And** thumbnails are the small version of the image, size optimised (NFR2)
**And** date range validation ensures end date is after start date (FR41)
**And** album name validation ensures name is not empty (FR42)
**And** validation errors display inline below the relevant field
**And** warning is shown if date selection would orphan existing photos (FR26)
**And** Create button is disabled until validation passes
**And** clicking Create button:
  - Shows loading state on button during API call
  - Closes dialog on success
  - Displays success message (FR28)
  - Refreshes album list to show new album
**And** clicking Cancel button closes dialog without creating album
**And** error handling displays clear message if creation fails (FR43)
**And** recovery option (retry button) is provided on error (FR46)
**And** dialog is responsive and works on mobile, tablet, desktop
**And** keyboard navigation works (Tab, Enter to submit, ESC to cancel) (NFR5)
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
**Then** an "Edit Name" button/icon is visible only on albums I own (FR20)
**And** the edit option is not visible on shared albums I don't own
**And** clicking edit opens an edit name dialog
**And** the dialog contains:
- Album name text field pre-filled with current name
- Current album metadata displayed for context
- Cancel and Save buttons
**And** album name validation ensures name is not empty (FR42)
**And** validation error displays inline if name is cleared
**And** Save button is disabled if validation fails
**And** clicking Save button:
  - Shows loading state on button during API call
  - Closes dialog on success
  - Displays success message (FR28)
  - Updates album name in the album list immediately
**And** clicking Cancel button closes dialog without saving changes
**And** error handling displays clear message if update fails (FR43)
**And** recovery option (retry button) is provided on error (FR46)
**And** dialog is responsive and works on mobile, tablet, desktop
**And** keyboard navigation works (Tab, Enter to submit, ESC to cancel) (NFR5)
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
**Then** an "Edit Dates" button/icon is visible only on albums I own (FR20)
**And** the edit option is not visible on shared albums I don't own
**And** clicking edit opens an edit date range dialog
**And** the dialog contains:
- Current album name displayed for context
- Start date picker pre-filled with current start date
- End date picker pre-filled with current end date
- Photo preview panel showing ~12 thumbnail samples from current range
- Cancel and Save buttons
**And** when I adjust start or end dates, the photo preview updates (FR27)
**And** preview thumbnails re-fetch from API based on new date range
**And** thumbnails are the small version of the image, size optimised (NFR2)
**And** date range validation ensures end date is after start date (FR41)
**And** validation error displays if end date is before or equal to start date
**And** warning is shown if new date range would orphan existing photos (FR26)
**And** orphan warning displays which photos would be affected
**And** Save button is disabled until validation passes
**And** clicking Save button:
  - Shows loading state on button during API call
  - Triggers backend re-indexing of photos (FR27)
  - Closes dialog on success
  - Displays success message (FR28)
  - Refreshes album to show updated photo set
**And** clicking Cancel button closes dialog without saving changes
**And** error handling displays clear message if update fails (FR43)
**And** recovery option (retry button) is provided on error (FR46)
**And** network failures do not cause data loss (NFR3, NFR14)
**And** dialog is responsive and works on mobile, tablet, desktop
**And** keyboard navigation works (Tab, arrow keys to adjust dates, Enter to save, ESC to cancel) (NFR5)
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
**Then** a "Delete Album" button/icon is visible only on albums I own (FR20)
**And** the delete option is not visible on shared albums I don't own
**And** clicking delete opens a confirmation dialog
**And** the confirmation dialog contains:
- Warning message explaining deletion is permanent
- Album name displayed for confirmation
- Number of photos that will remain (not deleted, just album removed)
- Cancel and Delete buttons
**And** Delete button uses warning/error color (red) to indicate destructive action
**And** clicking Delete button:
  - Shows loading state on button during API call
  - Closes dialog on success
  - Displays success message (FR28)
  - Removes album from the album list immediately
  - Navigates to home page if currently viewing the deleted album
**And** clicking Cancel button closes dialog without deleting album
**And** error handling displays clear message if deletion fails (FR43)
**And** recovery option is provided on error (FR46)
**And** if album is not found during deletion, appropriate error displays (FR45)
**And** dialog is responsive and works on mobile, tablet, desktop
**And** keyboard navigation works (Tab, Enter on Cancel as default, must explicitly select Delete) (NFR5)
**And** focus defaults to Cancel button (safer default for destructive action)
**And** the album is removed from all views after successful deletion

---

## Epic 4: Sharing & Access Control

Album owners can share albums with users via email, revoke access, view who has access with avatars and names; shared viewers see ownership information and
clear UI distinction between owner and viewer capabilities.

### Story 4.1: Album Sharing UI Components

As a developer,
I want reusable UI components for album sharing,
So that the sharing experience is consistent and maintainable.

**UI Components to Create:**

1. **AlbumSharingDialog** - Main dialog for managing album sharing
   - Props: album name, shared users list, loading state, error message, onGrantAccess, onRevokeAccess, onClose
   - Contains:
     - Album name displayed for context
     - Email address input field with validation
     - "Grant Access" button
     - List of currently shared users with their avatars and names
     - Close button
   - Email input validates format on blur (FR34)
   - Validation error displays inline if email format is invalid
   - Grant Access button disabled if email is invalid or empty
   - Shared users list displays:
     - User avatar image
     - User name
     - User email
     - Remove/revoke button (X icon) next to each user
   - Error message displayed if passed via props (FR44)
   - Recovery option (retry button) shown when error is provided (FR46)
   - Responsive and works on mobile, tablet, desktop

2. **RevokeAccessConfirmation** - Confirmation dialog for revoking access
   - Props: user name, user avatar, album name, loading state, onConfirm, onCancel
   - Contains:
     - User name and avatar being removed
     - Album name for context
     - Warning that user will lose access to the album
     - Cancel and Revoke Access buttons
   - Revoke Access button uses warning color to indicate action
   - Responsive and works on mobile, tablet, desktop

**Acceptance Criteria:**

**Given** the components are implemented
**When** they are used in the album sharing flow
**Then** both components are created and exported
**And** each component has clear prop interfaces
**And** components use MUI styling system
**And** components follow the brand color scheme (#185986)
**And** components have visual regression tests
**And** keyboard navigation works (Tab, Enter, ESC) (NFR5)
**And** focus management is clear and logical

---

### Story 4.2: Integrate Album Sharing UI

As an album owner,
I want to access the album sharing interface from my albums,
So that I can manage who has access to my photos.

**Acceptance Criteria:**

**Given** I am viewing an album I own
**When** I want to share the album
**Then** a "Share" button/icon is visible only on albums I own (FR33)
**And** the share option is not visible on shared albums I don't own
**And** clicking share opens the AlbumSharingDialog component
**And** the dialog displays the album name
**And** the dialog shows a list of currently shared users (FR31)
**And** I can see all users who currently have access to the album (FR31)
**And** focus is set to email input when dialog opens
**And** clicking close button closes the dialog

**Given** I am in the sharing dialog
**When** I interact with the revoke button for a user
**Then** the RevokeAccessConfirmation dialog opens
**And** the dialog shows the user name and avatar
**And** clicking Cancel closes the confirmation dialog
**And** focus defaults to Cancel button (safer default)

**Given** I am a shared viewer of an album
**When** I view the album
**Then** I can see who owns the album (name and avatar) (FR32)
**And** I cannot see the Share button (not an owner) (FR33)
**And** the UI clearly distinguishes that I am viewing, not owning

---

### Story 4.3: Persist Sharing and Revoke Operations

As an album owner,
I want my sharing and revoke actions to be saved,
So that changes to album access are persistent.

**Acceptance Criteria:**

**Given** I have entered a valid email and clicked Grant Access
**When** the grant access action is triggered
**Then** a loading state is shown during API call
**And** user profile information (name, picture) is fetched from API (FR35)
**And** the user is added to the shared users list with avatar and name (FR29)
**And** a success message is displayed
**And** email input field is cleared for next entry
**And** shared user avatars appear on the album card after sharing (FR11, FR5)
**And** error message is displayed if sharing fails (FR44)
**And** retry button is shown on error (FR46)

**Given** I have confirmed revoke access in the confirmation dialog
**When** the revoke action is triggered
**Then** a loading state is shown during API call
**And** the user is removed from the shared users list on success (FR30)
**And** a success message is displayed
**And** album card updates to remove user's avatar if no longer shared
**And** the confirmation dialog closes
**And** the sharing dialog remains open to allow additional changes
**And** error message is displayed if revoke fails (FR44)
**And** retry button is shown on error (FR46)
**And** the revoked user no longer sees the album in their album list

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
**And** the mock implementation uses the random photos API endpoint `GET /api/v1/owners/[ownerId]/randomPhotos?size=8` if available
**And** if API is not available, selects random photos from the album cards displayed below (using sample photos)
**And** each photo is displayed
**And** photos are the small version of the image, size optimised (NFR2)
**And** the layout is responsive:
- Mobile (xs): Horizontal scroll with 2-3 photos visible
- Tablet (sm): 4-5 photos visible
- Desktop (md): 6-8 photos visible
**And** each photo is clickable and links to its source album
**And** clicking a photo navigates to `/albums/[ownerId]/[albumId]`
**And** the section has a clear heading: "Your Memories" or "Highlights"
**And** photos are displayed with subtle spacing and styling
**And** the section uses brand color #185986 for accent elements
**And** loading skeleton is shown while photos are being selected
**And** the section works on mobile, tablet, and desktop devices
**And** the random selection changes on page reload (different subset)
**And** the feature is demoable and validates the UX flow
**And** if no albums exist, the highlights section is hidden
**And** the mock implementation is clearly commented as temporary

---

### Story 5.2: Backend API for Random Photo Discovery

As a frontend developer,
I want to get a random list of photos from any album from a user's accessible albums,
So that the frontend can display diverse memory highlights without relying on mock data.

**Acceptance Criteria:**

**Given** the backend API infrastructure exists
**When** implementing the random photos endpoint
**Then** a new REST API endpoint is created: `GET /api/v1/owners/[ownerId]/randomPhotos`
**And** the endpoint accepts an optional query parameter `size` (default: 8, max: 20)
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
**Then** the highlights section fetches random photos from `GET /api/v1/owners/[ownerId]/randomPhotos?size=8`
**And** the mock implementation from Story 5.1 is completely removed
**And** photos are displayed using the same UI layout from Story 5.1
**And** photo URLs are constructed using mediaId from API response
**And** clicking a photo navigates to `/albums/[ownerId]/[albumId]` using IDs from API
**And** the highlights refresh on each page reload (new API call fetches different randoms)
**And** loading state is shown while API call is in progress
**And** error handling displays graceful fallback if API fails (hide section or show friendly message)
**And** if API returns empty array (no photos), the section is hidden
**And** the feature works across mobile, tablet, and desktop
**And** photos display variety across different albums (leveraging backend randomization)
**And** users can successfully rediscover forgotten memories through this feature

---

## Ideas and Future Exploration

This section contains ideas and features that are not yet prioritized or fully defined. These may be explored and converted into epics in future iterations.

---

### Random Photo Discovery

**What:** Display random photo highlights on the home page to help users rediscover forgotten memories.

**Why:** Users accumulate large photo collections over time, and many beautiful or meaningful photos get forgotten. Random highlights surface these memories and
create moments of delight when browsing the home page.

**Key Features:**

- 5-8 random photos displayed at top of home page
- Photos drawn from all accessible albums (owned + shared)
- Horizontal layout with responsive design
- Click photo to navigate to source album
- Refreshes on each page reload with new random selection
- API endpoint: `GET /api/v1/owners/[ownerId]/randomPhotos?size=8`

**Requirements:**

- Backend API development (out of scope for web-nextjs epic)
- Frontend UI implementation
- Randomization algorithm that ensures variety across albums

**Status:** Defined as Epic 5 with 3 stories, awaiting prioritization

---

### Album Preview Before Navigation

**What:** Allow users to preview album contents before fully navigating to the album page.

**Why:** Users want to quickly identify albums of interest without committing to full page navigation. A preview mechanism lets them "peek" at album contents
and make informed decisions about which albums to explore.

**Exploration Needed:**

- **Interaction pattern:** Hover (desktop), first tap shows preview / second tap navigates (mobile), or dedicated preview icon?
- **Preview content:** How many sample photos? (3-6 thumbnails suggested)
- **Performance:** Must be lightweight and fast
- **Responsive design:** Must work well on mobile, tablet, and desktop

**Possible Approaches:**

1. **Hover preview (desktop only):** Show preview on mouse hover
2. **Two-tap pattern (mobile):** First tap shows preview overlay, second tap navigates
3. **Dedicated icon:** Small preview icon on album card that triggers preview
4. **Long-press (mobile):** Long press on album card shows preview

**Status:** Requires design exploration and UX validation before implementation

---

### Album Timeline Navigation

**What:** Add a timeline component to the AlbumPageHeader showing nearby albums with sample photos.

**Why:** When viewing photos from a specific time period, users may want to explore adjacent albums from similar timeframes. A timeline provides context and
enables seamless temporal navigation.

**Key Features:**

- Visual timeline integrated into AlbumPageHeader
- Shows current album and 2-3 nearby albums (before/after)
- Each album preview shows random sample photos
- Click to navigate to adjacent album
- Helps users maintain temporal context while browsing

**Exploration Needed:**

- Visual design of timeline component
- Number of adjacent albums to show
- How to handle albums with large time gaps
- Mobile responsiveness (timeline may need to collapse or scroll)
- Performance impact of loading preview photos

**Status:** Early concept, requires design and technical exploration
