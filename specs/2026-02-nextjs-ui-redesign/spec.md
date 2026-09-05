---
stepsCompleted: ['step-01-init', 'step-02-discovery', 'step-03-success', 'step-04-journeys', 'step-05-domain', 'step-06-innovation', 'step-07-project-type', 'step-08-scoping', 'step-09-functional', 'step-10-nonfunctional', 'step-11-polish']
inputDocuments:
  - '/home/dush/dev/git/dphoto/AGENTS.md'
  - '/home/dush/dev/git/dphoto/README.md'
  - '/home/dush/dev/git/dphoto/specs/2026-01-ux-functionnal.md'
  - '/home/dush/dev/git/dphoto/specs/2026-01-ux-reference.md'
workflowType: 'prd'
briefCount: 0
researchCount: 0
brainstormingCount: 0
projectDocsCount: 4
classification:
  projectType: 'web_app'
  domain: 'general'
  complexity: 'low'
  projectContext: 'brownfield'
---

# Product Requirements Document - dphoto

**Author:** Arch
**Date:** Wed Jan 28 2026

## User Journeys

### Journey 1: Thomas - The Photographer & Organizer

**Persona: Thomas, 39, Photography Enthusiast**

Thomas captures family moments via phone and DSLR. Photos are backed up to DPhoto via CLI. He needs to organize 300+ vacation photos into albums and share with family.

**Context:** After week-long vacation, photos sit unorganized in DPhoto.

**Journey:**

1. **Opens redesigned DPhoto** - Dark interface displays timeline of albums. Random photos from across his collection surface at the top, sparking curiosity about forgotten memories.

2. **Discovers natural groupings** - Timeline view shows photo density clusters around specific days: "3 days in Santa Cruz, July 2026"

3. **Creates album** - Smooth dialog (not form table) accepts name and date range for "Santa Cruz Beach Trip 2026". Polished animation.

4. **Browses album** - Photos grouped by day. Clicks beach sunset photo - smooth zoom to full view. Arrow keys navigate through sequence. Fluid transitions.

5. **Shares with Claire** - Enters Claire's email. Her avatar appears in shared-with list.

**Outcome:** Claire texts: "These beach photos are gorgeous! Which ones should we print?" Thomas feels proud - photos feel curated, not dumped. He wants to organize more trips.

---

### Journey 2: Claire - The Browser & Curator

**Persona: Claire, 45, Family Memory Keeper**

Claire browses albums, selects favorites to print and frame, shares with grandparents. Wants photo book experience, not file system navigation.

**Context:** Receives notification Thomas shared "Santa Cruz Beach Trip 2026".

**Journey:**

1. **Opens shared album on phone** - Smooth mobile experience. Dark theme highlights photos. Random photo highlights spark curiosity.

2. **Scrolls timeline** - Photos grouped by day with date headers. July 15th - Beach Day, July 16th - Amusement Park. Narrative presentation.

3. **Discovers favorite moment** - Taps candid photo of kids building sandcastles. Smooth zoom. Swipes through series.

4. **Explores other albums** - Random highlights from other albums (sunset from last year, family dinner) lead her to albums she hadn't revisited in months.

5. **Uses keyboard on desktop** - Arrow keys browse photos efficiently. Marks mental notes of prints.

6. **Downloads selections** - 8 photos for family room wall. Feels like curating, not downloading files.

**Outcome:** Creates photo wall showing trip story. Random highlights keep drawing her back to rediscover forgotten memories.

---

### Journey 3: Marie - Extended Family Viewer

**Persona: Marie, 68, Claire's Mother**

Marie isn't tech-savvy but loves seeing grandchildren photos. Needs intuitive, forgiving interface.

**Context:** Email notification: "Thomas shared an album with you: Santa Cruz Beach Trip 2026"

**Journey:**

1. **Logs in with Google** - Gmail account authentication familiar and works.

2. **Lands on album view** - Clean interface. Random photo highlights show grandchildren's faces - instant emotional connection.

3. **Browses on iPad** - Touch interactions natural. Swipes through photos, pinches to zoom. Responds smoothly on older iPad. Random highlights between day sections maintain engagement.

4. **Recovers from error** - Accidentally closes photo. Clear album structure helps navigate back.

5. **Shares excitement** - Calls Claire: "I saw the beach photos! The kids looked so happy!"

**Outcome:** Checks DPhoto regularly. Random highlights show different memories each visit, keeping experience fresh.

---

### Journey 4: Thomas - Album Management

**Persona: Thomas (returning user)**

Thomas has dozens of albums. Needs to reorganize - adjust date ranges, rename albums.

**Context:** "Summer 2025" album includes early September photos that should be in "Fall Adventures".

**Journey:**

1. **Opens album list** - Timeline order with visual temperature indicators. Random photos remind him of albums needing reorganization.

2. **Filters to owned albums** - "My Albums" view focuses on editable albums. Clear UI distinction between owned and shared.

3. **Edits dates** - Opens "Summer 2025", clicks edit dates. Dialog shows current range, accepts end date adjustment. Warns if changes orphan photos. Saves smoothly.

4. **Renames album** - Changes "Trip Photos 2025" to "Iceland Adventure 2025". Updates everywhere.

5. **Deletes duplicate** - Two albums cover same dates. Deletes one. Clear confirmation.

**Outcome:** Collection tells coherent story. Random photo feature keeps him engaged with entire collection, not just recent albums.

## Success Criteria

### Vision

Transform DPhoto from "functional student project" to "polished, curiosity-driven experience" through craft and intentional design, not feature bloat.

### User Success Metrics

**Emotional Engagement:**
- Users describe experience as "polished" and "thoughtfully designed"
- Browsing feels "engaging" and "discovery-driven"
- Interface invites exploration through curiosity, not static CRUD operations

**Functional Capabilities:**
- Keyboard navigation completes core flows (arrows, esc, enter)
- Responsive design works seamlessly on mobile, tablet, desktop
- Progressive image loading performs well on slow networks
- 100% feature parity with existing system

### Technical Success Metrics

**Performance:**
- 60fps animations with no jank
- Lighthouse Performance score ≥90 on mobile
- Lighthouse Accessibility score ≥95
- Progressive image loading (blur-up low to high quality)
- Fast initial page load and time-to-interactive

**Architecture:**
- NextJS + Material UI foundation
- Type-safe TypeScript implementation
- Reusable components with clear separation of concerns
- Keyboard navigation, semantic HTML, ARIA labels
- Same backend REST API (no backend changes)

## Product Scope

### Scope Definition

Frontend redesign of existing functionality. MVP equals complete product - all current capabilities redesigned with modern UX before launch.

**Rationale:** Brownfield redesign requires feature parity to replace existing UI. Small surface area (6 core capabilities) makes complete launch feasible.

### Core Capabilities (All Included in Launch)

1. **Browse Albums** - Timeline view, modern card layouts, visual activity indicators, owner info, sharing status
2. **View Photos** - Smooth zoom transitions, keyboard navigation, progressive loading, day-grouped presentation, full-screen viewing
3. **Album Management** - Create/edit/delete albums, date range selection, modern dialog patterns
4. **Sharing Management** - Grant/revoke access via email, view shared users and avatars
5. **Filtering & Navigation** - Filter by owner, "My Albums" vs "All Albums" views
6. **Authentication** - Google OAuth integration, user profile display, permission handling

### Design Principles

- **Dark-first theme** - modern aesthetic, photos as primary focus
- **Timeline/chronological navigation** - narrative flow through content
- **Random photo surfacing** - spark curiosity, rediscovery of forgotten memories
- **Purposeful transitions** - smooth, performant animations (60fps target)
- **Mobile-responsive** - works seamlessly across all device sizes
- **Keyboard shortcuts** - power user efficiency for core flows

### Technical Foundation

- NextJS (App Router), Material UI, TypeScript
- Same REST API backend (no API changes)
- Modern evergreen browsers only (latest 2 versions)
- Responsive design: mobile, tablet, desktop

### Out of Scope

**Post-Launch Enhancements:**
- Light theme support
- Advanced animations beyond core flows
- Additional keyboard shortcuts
- Additional timeline visualization options
- Enhanced mobile-specific interactions

**Future Vision (Not Planned):**
- New capabilities beyond current feature set
- AI-powered features (tagging, facial recognition, search)
- Photo editing or manipulation
- Social features beyond sharing
- Integration with other photo services
- Offline-first PWA capabilities

## Innovation & Novel Patterns

### Core Innovation: Curiosity-Driven Discovery for Intimate Collections

DPhoto targets private family collections (2-10 members) rather than massive scale. Innovation is interaction design, not technology.

**Key Differentiators:**

1. **Random Photo Surfacing**
   - Serendipitous rediscovery of forgotten memories
   - Treats photos as memories to rediscover, not files to organize
   - Sparks curiosity and emotional engagement

2. **Timeline as Narrative Device**
   - Chronological flow creates story-driven browsing
   - Visual density indicators show activity patterns
   - Moves beyond table-based CRUD to narrative presentation

3. **Designed for Intimacy**
   - Every photo has personal meaning
   - Focus on craft and polish over feature breadth
   - Small surface area, exceptional execution

### Market Context

**Existing Solutions:**
- Google Photos / iCloud: Scale-focused, algorithmic organization
- Amazon Photos: Storage-focused, basic organization
- Instagram Memories: Random surfacing in public/social context
- Apple Photos: Timeline view with file-system metaphor

**DPhoto's Differentiation:**
Timeline + random discovery + intimate scale + curiosity-driven interactions + dark-first aesthetic + private family focus.

### Validation Approach

Thomas (builder and primary user) validates through direct usage:
- Frequency of visits to browse vs upload
- Click-through on random photo highlights
- Time spent browsing vs managing
- Claire and family feedback
- Personal satisfaction with craft

### Risk Mitigation

**Risk:** Random photos annoying rather than engaging
- **Mitigation:** Conservative randomization, easy to adjust algorithm, can disable if ineffective

**Risk:** Timeline + transitions feel gimmicky
- **Mitigation:** Focus on purposeful animations, maintain 60fps, dial back if excessive

**Risk:** Dark theme doesn't work for all photos
- **Mitigation:** Iterate based on actual content, add light theme post-launch if needed

**Overall Risk:** Low - personal project with direct feedback loop, immediate adjustments possible.

## Web App Technical Requirements

### Architecture Overview

**Stack:**
- NextJS (App Router) - server and client rendering
- Material UI - component library and design system
- TypeScript - type safety throughout
- React Server Components where beneficial

**Approach:**
- Multi-page application (MPA) structure aligns with album/photo navigation
- Component-based architecture, clear separation between presentation and logic
- Architecture prioritizes design flexibility over technical constraints

### Browser Support

**Supported:** Latest 2 versions of Chrome, Edge, Firefox, Safari (including iOS Safari)
**Not Supported:** IE11, older browser versions
**Allowed:** CSS Grid, Flexbox, ES2020+, Web APIs without polyfills

### Authentication & Security

- Existing Google OAuth backend (already implemented)
- No public-facing pages - all routes require authentication
- Session management handled by existing backend
- Frontend integrates with existing authentication mechanism

### Performance Targets

- Initial page load optimized through NextJS server rendering
- Progressive image loading using existing API quality parameters
- 60fps animations and transitions
- Lighthouse Performance score ≥90 (mobile)
- Efficient bundle size through code splitting

### Data Flow & State Management

- RESTful API communication with existing backend
- Manual refresh model (no real-time updates required)
- State management TBD during implementation (React Context, Zustand, or similar)
- Optimistic UI updates for perceived performance
- Error handling and retry logic for network failures

### Responsive Design

- Mobile-first approach
- Breakpoints: Mobile (<600px), Tablet (600-960px), Desktop (>960px)
- Touch-friendly interactions on mobile/tablet
- Progressive enhancement from mobile to desktop

### Simplifications (No Complex Requirements)

**No SEO:** All pages behind authentication - no meta tags, structured data, or sitemap needed
**No Real-Time:** Traditional request/response model, manual refresh (F5 or pull-to-refresh)
**Accessibility Baseline:** Keyboard navigation, semantic HTML, focus management, alt text - no formal WCAG compliance required
**Design First:** Architecture should not constrain design decisions

## Implementation Strategy

### Launch Strategy: Complete Redesign

Not traditional MVP - complete feature parity required before cutover from old UI to new UI.

**Rationale:**
- Brownfield redesign requires all 6 capabilities functional
- Small surface area makes phased launch unnecessary
- Private family system allows thorough testing before cutover
- Complete experience ensures design cohesion

**Launch Criteria:**
- All 6 core capabilities redesigned (see Product Scope section)
- Performance targets met (60fps animations, Lighthouse 90+)
- Thomas and Claire validate as primary interface for 2-4 weeks
- Design vision fully realized (timeline, random photos, transitions)

### Implementation Sequence (Suggested)

1. Core browsing (album list, photo viewing, basic navigation)
2. Random photo surfacing and timeline visualization
3. Transitions and animations (zoom, smooth interactions)
4. Album management (create, edit, delete)
5. Sharing management (grant, revoke)
6. Polish and refinement (loading states, error handling, edge cases)

### Cut-Over Strategy

- Build new frontend alongside existing one
- Test thoroughly with Thomas and Claire
- Switch traffic when complete and validated
- Deprecate old frontend after validation period

### Resource Requirements

- Frontend developer (TypeScript, React, NextJS, Material UI)
- Design implementation and animation work
- Integration with existing REST API
- Testing across mobile, tablet, desktop
- Iterative refinement based on usage

### Risk Mitigation

**Animation performance on older mobile devices**
- Mitigation: 60fps target, performance monitoring, graceful degradation
- Validation: Test on range of devices

**Progressive image loading might not feel smooth**
- Mitigation: Leverage existing API quality parameters, optimize blur-up
- Validation: Test on slow networks

**Timeline + random photos interaction complexity**
- Mitigation: Start simple, iterate based on usage
- Validation: Direct usage feedback from Thomas and Claire

**Family might prefer old interface (familiarity)**
- Mitigation: Thorough testing before cutover
- Validation: 2-4 weeks as primary interface

**Implementation takes longer than expected**
- Mitigation: Small surface area, well-defined requirements, existing backend
- Fallback: Launch with reduced polish, iterate after cutover

**Overall Risk Posture:** Low to Medium
- ✅ Stable backend (no API changes)
- ✅ Well-defined requirements (existing functionality)
- ✅ Direct user feedback (Thomas is builder and user)
- ✅ No external dependencies or compliance
- ⚠️ Design/animation complexity requires careful execution
- ⚠️ Complete feature parity required before launch

## Functional Requirements

### Album Discovery & Browsing

- **FR1:** Album owners can view all albums they own
- **FR2:** Shared viewers can view all albums shared with them
- **FR3:** Users can filter albums by owner (my albums, all albums, specific owner)
- **FR4:** Users can view album metadata (name, date range, media count, owner information)
- **FR5:** Users can see which users an album is shared with
- **FR6:** Users can view albums in chronological order
- **FR7:** Users can see visual indicators of album activity density (temperature)
- **FR8:** Users can discover random photos from across their accessible collection
- **FR9:** Users can navigate from random photo highlights to the source album

### Photo Viewing & Navigation

- **FR10:** Users can view photos grouped by capture date within an album
- **FR11:** Users can open individual photos in full-screen view
- **FR12:** Users can navigate between photos using keyboard controls (arrow keys, esc, enter)
- **FR13:** Users can navigate between photos using touch gestures (swipe on mobile/tablet)
- **FR14:** Users can view photos with progressive quality loading (low to high resolution)
- **FR15:** Users can zoom into photo details
- **FR16:** Users can navigate back to album list from photo view
- **FR17:** Users can see date headers separating photos by day

### Album Management

- **FR18:** Album owners can create new albums by specifying name and date range
- **FR19:** Album owners can specify custom folder names for albums (optional)
- **FR20:** Album owners can edit album names
- **FR21:** Album owners can edit album date ranges
- **FR22:** Album owners can delete albums they own
- **FR23:** System validates that date edits don't orphan media
- **FR24:** System re-indexes photos when album date ranges change
- **FR25:** System provides feedback when album operations succeed or fail

### Sharing & Access Control

- **FR26:** Album owners can share albums with other users via email address
- **FR27:** Album owners can revoke access from users they previously shared with
- **FR28:** Album owners can view who has access to their albums
- **FR29:** Shared viewers can see who owns albums shared with them
- **FR30:** System distinguishes between owner capabilities and viewer capabilities in UI
- **FR31:** System validates email addresses when granting access
- **FR32:** System loads user profile information (name, picture) for display

### Authentication & User Management

- **FR33:** Users can authenticate using Google OAuth
- **FR34:** System maintains user session across page navigation
- **FR35:** Users can view their profile information (name, picture)
- **FR36:** System identifies if authenticated user is an owner (can create albums)
- **FR37:** System restricts all pages to authenticated users only

### Visual Presentation & UI State

- **FR38:** System displays loading states while fetching albums or photos
- **FR39:** System displays empty states when no albums exist
- **FR40:** System displays error messages when operations fail
- **FR41:** System indicates selected albums and active filters
- **FR42:** System provides visual transitions when navigating between views
- **FR43:** System displays album sharing status with user avatars
- **FR44:** System provides responsive layouts for mobile, tablet, and desktop devices

### Error Handling & Validation

- **FR45:** System validates album date ranges (end date after start date)
- **FR46:** System validates album names are not empty
- **FR47:** System handles and displays errors for failed album operations
- **FR48:** System handles and displays errors for failed sharing operations
- **FR49:** System handles and displays errors when albums are not found
- **FR50:** System provides recovery options when operations fail

## Non-Functional Requirements

### Performance

**Animation & Interaction Performance:**
- **NFR1:** Visual transitions and animations must maintain 60fps on modern mobile and desktop devices
- **NFR2:** User interactions (clicks, taps, swipes) must provide immediate visual feedback (<100ms)
- **NFR3:** Keyboard navigation must respond without perceptible lag
- **NFR4:** Page transitions and photo zoom animations must feel smooth and purposeful

**Image Loading & Display:**
- **NFR5:** System must display thumbnail-quality images immediately while full-resolution images load in background
- **NFR6:** System must request minimum image size appropriate for current screen dimensions (mobile, tablet, desktop)
- **NFR7:** Progressive image loading (blur-up from low to high quality) must complete within 3 seconds on slow network conditions
- **NFR8:** System must use existing API quality parameters to optimize image delivery

**Page Load Performance:**
- **NFR9:** Initial page load must achieve Lighthouse Performance score ≥90 on mobile devices
- **NFR10:** Time to interactive must be optimized through efficient code splitting and lazy loading
- **NFR11:** Subsequent page navigations must feel instantaneous through appropriate caching strategies

### Integration

**Backend API Integration:**
- **NFR12:** Frontend must consume existing REST API endpoints without requiring backend modifications
- **NFR13:** System must handle API response times gracefully with appropriate loading states
- **NFR14:** System must handle API errors with clear user feedback and retry mechanisms
- **NFR15:** System must maintain data contract compatibility with existing API response formats
- **NFR16:** System must leverage existing API image quality parameters for progressive loading

### Security

**Authentication & Authorization:**
- **NFR17:** System relies on existing backend Google OAuth authentication mechanism
- **NFR18:** System must maintain user session security as provided by backend
- **NFR19:** System must restrict all pages to authenticated users only (enforced by backend)
- **NFR20:** Frontend requires no additional security measures beyond consuming authenticated API endpoints

### Usability

**Keyboard Navigation:**
- **NFR21:** Core photo browsing flows must be fully keyboard-accessible (arrow keys, esc, enter)
- **NFR22:** Focus management must be clear and logical during navigation and in dialogs
- **NFR23:** Keyboard shortcuts must not conflict with browser defaults

**Mobile & Responsive Design:**
- **NFR24:** Touch interactions must feel responsive with appropriate visual feedback
- **NFR25:** Layouts must adapt appropriately across mobile (<600px), tablet (600-960px), and desktop (>960px) breakpoints
- **NFR26:** Mobile gestures (swipe, pinch-to-zoom) must feel natural and performant
- **NFR27:** Mobile performance must not degrade below acceptable animation frame rates

**Browser Compatibility:**
- **NFR28:** System must support latest 2 versions of Chrome, Firefox, Safari, and Edge (evergreen browsers)
- **NFR29:** System may use modern web features (CSS Grid, Flexbox, ES2020+) without polyfills for legacy browsers
- **NFR30:** No support required for Internet Explorer or older browser versions

### Reliability

**Availability & Error Recovery:**
- **NFR31:** System operates on best-effort availability basis (no uptime SLA required)
- **NFR32:** Manual page refresh is acceptable recovery mechanism for transient errors
- **NFR33:** System must provide clear error messages when operations fail with guidance on recovery
- **NFR34:** Network failures must not cause data loss for in-progress operations
