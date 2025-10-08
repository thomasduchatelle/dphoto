# SPA Navigation Implementation

## Summary

This implementation updates the navigation system to provide SPA-like behavior during the Waku migration phase. Previously, clicking links would trigger full page reloads. Now, navigation happens client-side without reloading the page.

## Changes Made

### 1. ClientRouter (`web/src/components/ClientRouter.tsx`)

**Before:** Used Waku's `useRouter()` which triggered server-side navigation and page reloads.

**After:** Uses browser's History API for client-side navigation:
- `window.history.pushState()` for navigation
- `window.history.replaceState()` for URL replacement
- State tracking for `currentPath` and `currentSearch` to trigger re-renders
- `popstate` event listener for browser back/forward buttons

### 2. AppNav Component (`web/src/components/AppNav/index.tsx`)

**Before:** Logo links used plain `<a href="/">` which caused page reloads.

**After:** Logo links now use `useClientRouter()` with click handlers that prevent default and call `navigate()`.

### 3. MobileNavigation Component (`web/src/components/albums/MobileNavigation/index.tsx`)

**Before:** Albums breadcrumb link used MUI's `Link` with plain href which caused page reloads.

**After:** Albums link now uses `useClientRouter()` with click handler for SPA navigation.

### 4. Tests (`web/src/components/ClientRouter.test.tsx`)

Added comprehensive tests covering:
- Path navigation without reload
- URL replacement without reload
- Path parameter parsing for albums and media
- Query parameter handling
- `ClientLink` component behavior
- Browser back/forward navigation

## How It Works

### User Journey (As Required)

1. **User loads or refreshes a page from URL**
   - Waku file-based navigation loads the page
   - All pages render the same `<App />` component

2. **User clicks on a link**
   - Click handler prevents default browser navigation
   - `navigate()` or `replace()` is called
   - URL updates using `pushState`/`replaceState` (no reload)
   - State change triggers re-render
   - `AuthenticatedRouter` reads new path and renders appropriate component

3. **Browser back/forward buttons**
   - `popstate` event fires
   - State updates with new path
   - Component re-renders with correct content

### Technical Flow

```
Click → preventDefault() → pushState(newURL) → setState(newPath) → 
Component Re-render → AuthenticatedRouter reads path → Render new page
```

## Benefits

- ✅ No page flashes or reload delays
- ✅ No redirect to login page flash
- ✅ Smooth, instant page transitions
- ✅ Browser back/forward works correctly
- ✅ URL in address bar always reflects current page
- ✅ All existing navigation patterns preserved

## Migration Path

This implementation is designed for the migration phase where:
- Waku handles initial page loads (SSR)
- Client-side router handles all subsequent navigation
- Future: Can be replaced with full Waku routing when migration is complete
