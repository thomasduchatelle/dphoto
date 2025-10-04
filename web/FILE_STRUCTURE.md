# File-Based Routing Structure Migration

This document describes the new file-based routing structure implemented in preparation for Waku migration.

## Overview

The application has been refactored from a centralized router approach to a file-based routing structure that is compatible with Waku's requirements while maintaining compatibility with React Router (CRA).

## New Structure

```
src/pages/
├── index.tsx                           # Root page - redirects to /albums
├── login.tsx                           # Login page
├── _cra-router.tsx                     # React Router configuration (CRA only)
├── albums/
│   ├── index.tsx                       # Albums list page (GET /albums)
│   ├── _layout.tsx                     # Layout wrapper for all /albums/* routes
│   ├── [owner]/
│   │   └── [album]/
│   │       ├── index.tsx               # Album detail page (GET /albums/:owner/:album)
│   │       └── [encodedId]/
│   │           └── [filename].tsx      # Media viewer page (GET /albums/:owner/:album/:encodedId/:filename)
│   ├── media/                          # Media viewer components
│   ├── AlbumsList/                     # Album list UI components
│   ├── MediasPage/                     # Media page UI components
│   └── ...                             # Other UI components for albums
├── AppNav/                             # Global navigation component
├── catalog-react/                      # Catalog context provider
├── DPhotoTheme/                        # Theme provider
├── user.menu/                          # User menu component
├── error-page/                         # Error page component
├── google-login-button/                # Google SSO button component
└── login-domain/                       # Login domain logic
```

## Page Components

All page components follow the Waku-compatible pattern:

```tsx
export default function PageName() {
  // Page logic here
  return (/* JSX */)
}
```

### Route Mapping

| URL Pattern | File Path | Component |
|------------|-----------|-----------|
| `/` | `pages/index.tsx` | IndexPage (redirects to /albums) |
| `/login` | `pages/login.tsx` | LoginPage |
| `/albums` | `pages/albums/index.tsx` | AlbumsIndexPage |
| `/albums/:owner/:album` | `pages/albums/[owner]/[album]/index.tsx` | AlbumPage |
| `/albums/:owner/:album/:encodedId/:filename` | `pages/albums/[owner]/[album]/[encodedId]/[filename].tsx` | MediaPage |

## Layout Components

### `albums/_layout.tsx`

The `_layout.tsx` file provides a shared layout for all routes under `/albums/*`. It wraps child pages with the `CatalogViewerProvider` context and handles album ID extraction from URL parameters.

## CRA Router (`_cra-router.tsx`)

The `_cra-router.tsx` file is a temporary component that uses React Router to map routes to the new page components. This allows the application to work with Create React App while maintaining a file structure compatible with Waku.

**Important**: This file will be removed when migrating to Waku, as Waku uses automatic file-based routing.

## Migration Notes

### Changes Made

1. **Reorganized files**: Moved all components from `src/components/` and `src/pages/authenticated/` into `src/pages/` with a Waku-compatible structure
2. **Extracted layout**: Created `albums/_layout.tsx` to handle the `CatalogViewerProvider` logic
3. **Created page components**: Each route now has a dedicated page component that follows Waku conventions
4. **Updated imports**: Fixed all import paths to work with the new structure
5. **Removed dead code**: Deleted old router files and duplicate components

### Compatibility

- ✅ Works with current CRA setup via `_cra-router.tsx`
- ✅ File structure matches Waku file-based routing conventions
- ✅ All URLs remain the same
- ✅ All functionality preserved

### Next Steps for Waku Migration

1. Remove `_cra-router.tsx`
2. Update `App.tsx` to use Waku's routing
3. Add `'use client'` directives where needed
4. Migrate to Waku's data fetching patterns

## Component Organization

Components are now organized by feature/domain rather than by type:

- **Page-level components**: In `pages/albums/`, `pages/login.tsx`, etc.
- **Shared UI components**: `AppNav/`, `user.menu/`, etc. in `pages/`
- **Context providers**: `catalog-react/`, `DPhotoTheme/` in `pages/`
- **Domain logic**: `login-domain/` and similar in `pages/`

This structure makes it easier to identify which components belong to which pages and simplifies the future Waku migration.
