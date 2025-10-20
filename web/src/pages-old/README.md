# pages-old Directory

This directory contains the original page structure from the CRA (Create React App) version of the application.

With the migration to Waku, we're using a client-side only approach where:
- The common layout with Providers is defined in `src/pages/_layout.tsx` (Waku's layout convention)
- The layout renders the `GeneralRouter` component which uses React Router (BrowserRouter) to handle all routing
- All the original page components and routing logic remain in this `pages-old` directory
- The individual Waku pages in `src/pages/` are minimal and return null since all routing is handled by GeneralRouter

This approach maintains backward compatibility and minimal changes during the Phase 2 migration to Waku.
The navigation and routing work exactly as they did before, just using Waku as the build system instead of CRA.
