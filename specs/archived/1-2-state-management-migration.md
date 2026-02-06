# Story 1.2: State Management Migration

Status: ready-for-dev

## Story

As a developer,
I want to migrate the catalog state management from the existing web app,
So that I can reuse battle-tested state logic with 221+ passing tests.

## Acceptance Criteria

**Given** the existing state management in `web/src/core/catalog/` has 221+ passing tests
**When** I migrate the state management to web-nextjs
**Then** state management code is copied to `web-nextjs/domains/catalog/` maintaining folder structure:

- `language/` (state types - CatalogViewerState, Album, Media, etc.)
- `actions.ts` (reducer export)
- `album-create/`, `album-delete/`, `album-edit-dates/`, `album-edit-name/`, `sharing/` (thunk folders with actions, selectors, thunks)
- `date-range/`, `navigation/`, `common/`, `base-name-edit/` (supporting modules)
- `tests/` (test helpers like test-helper-state.ts)

**And** a new fetch adapter is created at `domains/catalog/adapters/fetch-adapter.ts` replacing axios
**And** the fetch adapter works in both Server Components and Client Components
**And** the fetch adapter implements the same interface as the original CatalogAPIAdapter (all Port interfaces)
**And** the fetch adapter converts REST responses to domain types (Album, Media, UserDetails, OwnerDetails)
**And** the fetch adapter includes authentication headers using existing session cookies (no explicit token management)
**And** the fetch adapter uses try-catch for error handling with CatalogError casting
**And** custom image loader is created at `libs/image-loader.ts` mapping:

- width ≤360 → width=360 (backend MiniatureCachedWidth)
- width >360 → width=2400 (backend MediumQualityCachedWidth, maximum cacheable resolution)

**And** image loader is configured in `next.config.ts` with custom loader function
**And** if the migrated code depends on a library (e.g., daction), that library must also be migrated to web-nextjs
**And** imports are fixed to make the application start and pass tests
**And** error boundaries are created:

- `app/error.tsx` (root catch-all)
- `app/(authenticated)/error.tsx` (authenticated routes)

**And** not-found pages are created:

- `app/not-found.tsx` (root 404)

**And** all migrated tests continue to pass when run with `npm run test` in web-nextjs
**And** the application builds successfully with `npm run build`

## Tasks / Subtasks

- [ ] Copy state management structure from web to web-nextjs (AC: folder structure)
    - [ ] Copy `web/src/core/catalog/language/` → `web-nextjs/domains/catalog/language/`
    - [ ] Copy `web/src/core/catalog/actions.ts` → `web-nextjs/domains/catalog/actions.ts`
    - [ ] Copy `web/src/core/catalog/album-create/` → `web-nextjs/domains/catalog/album-create/`
    - [ ] Copy `web/src/core/catalog/album-delete/` → `web-nextjs/domains/catalog/album-delete/`
    - [ ] Copy `web/src/core/catalog/album-edit-dates/` → `web-nextjs/domains/catalog/album-edit-dates/`
    - [ ] Copy `web/src/core/catalog/album-edit-name/` → `web-nextjs/domains/catalog/album-edit-name/`
    - [ ] Copy `web/src/core/catalog/sharing/` → `web-nextjs/domains/catalog/sharing/`
    - [ ] Copy supporting modules: `date-range/`, `navigation/`, `common/`, `base-name-edit/`, `tests/`
    - [ ] Copy `web/src/core/catalog/index.ts` (exports)
    - [ ] Copy `web/src/core/catalog/catalog-factories.ts` (if needed)

- [ ] Create fetch adapter replacing axios (AC: fetch adapter)
    - [ ] Create `domains/catalog/adapters/fetch-adapter.ts`
    - [ ] Implement fetchAlbums() - GET /api/v1/albums with owner/user details
    - [ ] Implement fetchMedias(albumId) - GET /api/v1/owners/{owner}/albums/{folderName}/medias
    - [ ] Implement createAlbum(request) - POST /api/v1/albums
    - [ ] Implement deleteAlbum(albumId) - DELETE /api/v1/owners/{owner}/albums/{folderName}
    - [ ] Implement renameAlbum(albumId, name, folderName?) - PUT /api/v1/owners/{owner}/albums/{folderName}/name
    - [ ] Implement updateAlbumDates(albumId, start, end) - PUT /api/v1/owners/{owner}/albums/{folderName}/dates
    - [ ] Implement grantAccessToAlbum(albumId, email) - PUT /api/v1/owners/{owner}/albums/{folderName}/shares/{email}
    - [ ] Implement revokeSharingAlbum(albumId, email) - DELETE /api/v1/owners/{owner}/albums/{folderName}/shares/{email}
    - [ ] Implement loadUserDetails(email) - GET /api/v1/users?emails={email}
    - [ ] Add error handling with CatalogError casting for 4xx/5xx responses
    - [ ] Calculate temperature and relativeTemperature for albums (numberOfDays helper)
    - [ ] Sort albums by start date descending

- [ ] Migrate dependencies (AC: libraries)
    - [ ] Copy daction library from `web/src/libs/daction/` → `web-nextjs/libs/daction/`
    - [ ] Copy any other library dependencies required by catalog domain
    - [ ] Fix import paths to reference migrated libraries

- [ ] Create custom image loader (AC: image loader)
    - [ ] Create `libs/image-loader.ts` with width mapping logic
    - [ ] Map width ≤360 → width=360 (MiniatureCachedWidth)
    - [ ] Map width >360 → width=2400 (MediumQualityCachedWidth)
    - [ ] Build URL: `/api/v1/media/{mediaId}/image?width={width}`
    - [ ] Export loader function compatible with Next.js Image component
    - [ ] Update `next.config.ts` to configure custom loader

- [ ] Create error boundaries and not-found pages (AC: error boundaries)
    - [ ] Create `app/error.tsx` with reset() button and error display
    - [ ] Create `app/(authenticated)/error.tsx` with context-aware error handling
    - [ ] Create `app/not-found.tsx` with navigation back to home
    - [ ] Use MUI components (Paper, Typography, Button) for consistent styling

- [ ] Fix imports (AC: imports fixed)
    - [ ] Fix all import paths to work in web-nextjs structure
    - [ ] Remove axios references from migrated code
    - [ ] Ensure all Port interfaces are correctly imported
    - [ ] Verify TypeScript types resolve correctly

- [ ] Run tests and verify build (AC: all)
    - [ ] Run `npm run test` - verify all tests pass
    - [ ] Run `npm run build` - verify successful build
    - [ ] Verify no TypeScript errors
    - [ ] Verify no missing imports

## Dev Notes

### Migration Context

**This is a lift-and-shift migration, NOT a reimplementation.** Copy `web/src/core/catalog/` → `web-nextjs/domains/catalog/` maintaining 221 passing tests. Only
change: replace axios adapter with fetch adapter.

### Key Migration Requirements

**1. Dependency Libraries:**

- daction library must be migrated: `web/src/libs/daction/` → `web-nextjs/libs/daction/`
- Fix all import paths after migration

**2. Image Loader - Backend Constraints:**

Backend supports only 2 width levels:

```go
MiniatureCachedWidth = 360      // Minimum cached size
MediumQualityCachedWidth = 2400 // Maximum cacheable resolution
```

**Loader Implementation:**

```typescript
export default function imageLoader({src, width}: ImageLoaderProps): string {
    const targetWidth = width <= 360 ? 360 : 2400;
    return `/api/v1/media/${src}/image?width=${targetWidth}`;
}
```

Configure in `next.config.ts`:

```typescript
const nextConfig = {
    images: {
        loader: 'custom',
        loaderFile: './libs/image-loader.ts',
    },
};
```

**3. Fetch Adapter:**

Reference: `web/src/core/catalog/adapters/api/CatalogAPIAdapter.ts`

Replace axios with fetch:

- Use `fetch()` with `credentials: 'include'` for session cookies
- Parse responses with `.json()`
- Error handling: Check `response.ok`, parse error body, cast to CatalogError
- Maintain exact Port interface signatures (CreateAlbumPort, DeleteAlbumPort, SaveAlbumNamePort, UpdateAlbumDatesPort, GrantAlbumAccessAPI,
  RevokeAlbumAccessAPI)
- Include helper functions: `numberOfDays()`, `convertToType()`
- Calculate temperature: `totalCount / numberOfDays(start, end)`

**4. Error Boundaries:**

Create 3 files using MUI components (Paper, Typography, Button):

- `app/error.tsx` - Root catch-all with "Return to Home"
- `app/(authenticated)/error.tsx` - Auth routes with reset() and "Try Again"
- `app/not-found.tsx` - 404 page (Server Component)

### Previous Story Context (1.1)

**Foundation Available:**

- Material UI 6.x with dark theme (#185986 brand color)
- `@/components/Link` wrapper for MUI Button integration
- ThemeProvider configured in `components/theme/`

**Patterns Learned:**

- Server Components use async searchParams
- Client wrappers needed for conditional rendering
- All tests passing (39), build succeeds

### State Management Structure

```
catalog/
├── language/           # State types (CatalogViewerState, Album, Media)
├── actions.ts          # Reducer using daction
├── album-create/       # Feature: action-*.ts, thunk-*.ts, selector-*.ts
├── album-delete/       
├── album-edit-dates/   
├── album-edit-name/    
├── sharing/           
├── date-range/        
├── navigation/        
├── common/            
├── base-name-edit/    
├── tests/             # test-helper-state.ts
└── adapters/api/      # CatalogAPIAdapter.ts (reference)
```

**Key Concepts:**

- **daction**: Custom action/reducer pattern
- **Port/Adapter**: Interfaces with implementations
- **Thunks**: Business logic dispatching actions
- **Domain isolation**: NO React/NextJS in domains/

### Success Validation

1. ✅ All 221 tests pass: `npm run test`
2. ✅ Build succeeds: `npm run build`
3. ✅ No TypeScript errors
4. ✅ daction library integrated
5. ✅ Fetch adapter implements all Ports
6. ✅ Image loader configured

### Anti-Patterns

❌ **DO NOT:**

- Reimplement state management
- Change action/reducer patterns
- Add React to domains/
- Modify Port interfaces

✅ **DO:**

- Exact copy (except adapter)
- Fix imports as needed
- Maintain all tests
- Follow existing naming

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
