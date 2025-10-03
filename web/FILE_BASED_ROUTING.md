# File-Based Routing Structure Documentation

## Overview

This document describes the file-based routing structure implemented in the `web/` project to prepare for migration to Waku. The routing structure has been reorganized to follow Waku's file-based routing conventions while maintaining compatibility with React Router.

## Routing Structure

### Current Routes and File Mapping

The application has been restructured to use Waku-compatible file-based routing conventions:

| Route Pattern | File Location | Description |
|---------------|---------------|-------------|
| `/` (login) | `pages/index.tsx` | Login page (unauthenticated users) |
| `/albums` | `pages/albums/index.tsx` | Albums list page |
| `/albums/:owner/:album` | `pages/albums/[owner]/[album]/index.tsx` | Specific album view |
| `/albums/:owner/:album/:encodedId/:filename` | `pages/albums/[owner]/[album]/[encodedId]/[filename].tsx` | Media viewer page |

### File-Based Routing Conventions

Following Waku's conventions:

- **`index.tsx`**: Represents the root route of a directory
  - `pages/index.tsx` → `/`
  - `pages/albums/index.tsx` → `/albums`
  - `pages/albums/[owner]/[album]/index.tsx` → `/albums/:owner/:album`

- **`[param].tsx`**: Represents a dynamic route parameter
  - `pages/albums/[owner]/...` → captures `:owner` parameter
  - `pages/albums/[owner]/[album]/...` → captures `:album` parameter
  - `pages/albums/[owner]/[album]/[encodedId]/[filename].tsx` → captures `:encodedId` and `:filename` parameters

### Directory Structure

```
web/src/pages/
├── index.tsx                                          # Login page (/)
├── ErrorPage/                                         # Error page component (conditional)
│   └── index.tsx
├── Login/                                             # Login domain logic and components
│   ├── GoogleLoginButton/
│   └── domain/
├── albums/                                            # Albums feature
│   ├── index.tsx                                     # Albums list (/albums)
│   └── [owner]/                                      # Dynamic owner parameter
│       └── [album]/                                  # Dynamic album parameter
│           ├── index.tsx                            # Album view (/albums/:owner/:album)
│           └── [encodedId]/                         # Dynamic media ID parameter
│               └── [filename].tsx                   # Media viewer (/albums/:owner/:album/:encodedId/:filename)
├── authenticated/                                     # Authenticated route components (legacy structure)
│   ├── AuthenticatedRouter.tsx                       # React Router configuration for authenticated routes
│   ├── albums/                                       # Album-related components (shared)
│   │   ├── AlbumsList/
│   │   ├── AlbumsListActions/
│   │   ├── MediasList/
│   │   ├── MediasPage/
│   │   ├── MobileNavigation/
│   │   ├── ShareDialog/
│   │   ├── CreateAlbumDialog/
│   │   ├── DeleteAlbumDialog/
│   │   ├── EditDatesDialog/
│   │   ├── EditNameDialog/
│   │   ├── FolderNameInput/
│   │   ├── DateRangePicker/
│   │   ├── CatalogViewerRoot.tsx                    # Catalog context provider wrapper
│   │   └── CatalogViewerPage.tsx                    # Original catalog viewer (kept for reference)
│   └── media/                                        # Media-related components (shared)
│       ├── FullHeightLink/
│       ├── MediaNavBar/
│       ├── logic/
│       ├── useNativeControl/
│       └── index.tsx                                 # Original media page (kept for reference)
└── GeneralRouter.tsx                                  # Main router handling auth state
```

## Implementation Details

### React Router Configuration

The application continues to use React Router with the new file locations referenced in `AuthenticatedRouter.tsx`:

```tsx
<Routes>
    <Route path='/albums' element={<CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>}/>
    <Route path='/albums/:owner/:album' element={<CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>}/>
    <Route path='/albums/:owner/:album/:encodedId/:filename' element={<CatalogViewerRoot><MediaPage/></CatalogViewerRoot>}/>
    <Route path='*' element={<RedirectToDefaultOrPrevious/>}/>
</Routes>
```

Where:
- `CatalogViewerPage` is imported from `pages/albums/[owner]/[album]/index.tsx`
- `MediaPage` is imported from `pages/albums/[owner]/[album]/[encodedId]/[filename].tsx`
- `CatalogViewerRoot` remains in `pages/authenticated/albums/` as a shared wrapper component

### Component Organization

1. **Page Components**: Located in file-based routing structure (`pages/albums/...`)
   - These are the main route handlers that will map directly to Waku routes
   - They use React Router hooks (`useParams`, `useNavigate`, etc.) for now

2. **Shared Components**: Remain in `pages/authenticated/albums/` and `pages/authenticated/media/`
   - These are reusable components used by page components
   - They include dialogs, lists, navigation, and other UI elements
   - No changes needed during Waku migration

3. **Domain Logic**: Kept in their original locations
   - Login domain logic in `pages/Login/domain/`
   - Media page logic in `pages/authenticated/media/logic/`

### Key Changes Made

1. **Created new page files** following Waku conventions:
   - `pages/index.tsx` - Login page (moved from `pages/Login/index.tsx`)
   - `pages/albums/index.tsx` - Albums list page
   - `pages/albums/[owner]/[album]/index.tsx` - Album viewer page
   - `pages/albums/[owner]/[album]/[encodedId]/[filename].tsx` - Media viewer page

2. **Updated imports** in router files:
   - `GeneralRouter.tsx` now imports Login from `pages/index.tsx`
   - `AuthenticatedRouter.tsx` imports page components from new file-based locations

3. **Preserved original files** in `pages/authenticated/` directory:
   - These serve as shared components and reference implementations
   - No breaking changes to existing component structure

## Migration Path to Waku

When migrating to Waku, the following changes will be needed:

1. **Remove React Router**: Replace `<Routes>` and `<Route>` with Waku's file-based routing
2. **Update hooks**: Replace React Router hooks (`useParams`, `useNavigate`, `useLocation`) with Waku equivalents
3. **Add `'use client'` directives**: Since the app uses client-side state and authentication
4. **Create layout files**: Add `_layout.tsx` files as needed for shared layouts
5. **Update navigation**: Replace `<Navigate>` and programmatic navigation with Waku's navigation

## Testing

To verify the file-based routing structure:

1. **TypeScript compilation**: `npx tsc --noEmit`
2. **Build**: `npm run build`
3. **Run application**: `npm start`
4. **Test routes**:
   - Login page at `/`
   - Albums list at `/albums`
   - Specific album at `/albums/:owner/:album`
   - Media viewer at `/albums/:owner/:album/:encodedId/:filename`

## Notes

- The routing behavior remains identical to the previous implementation
- All imports have been updated to reference the new file locations
- The `authenticated/` directory structure is preserved for shared components
- This structure is compatible with both React Router (current) and Waku (future)
- No functionality changes - only file organization and routing structure

## References

- Waku documentation: https://waku.gg
- React Router documentation: https://reactrouter.com
- Migration decision log: `specs/2025-09_Waku_Migration.md`
