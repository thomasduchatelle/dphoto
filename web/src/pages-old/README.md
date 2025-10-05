# pages-old Directory

This directory contains the original page structure from the CRA (Create React App) version of the application.

With the migration to Waku, we're using a client-side only approach where:
- The main Waku page is in `src/pages/index.tsx`
- It simply imports and renders the `App` component from `src/App.tsx`
- The `App` component uses React Router (BrowserRouter) to handle all routing
- All the original page components and routing logic remain in this `pages-old` directory

This approach maintains backward compatibility and minimal changes during the Phase 2 migration to Waku.
The navigation and routing work exactly as they did before, just using Waku as the build system instead of CRA.
