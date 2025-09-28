# DPhoto Web Application - Waku Migration Plan

## Target Architecture

**Deployment**: Static Site Generation (SSG) with CloudFront CDN + S3
**Authentication**: localStorage refresh tokens + memory access tokens
**Routing**: File-based routing replacing React Router
**Components**: Server components for data display, client components for interactivity

## File Structure

    src/
    ├── pages/
    │   ├── _layout.tsx                    # Auth guard + AppNav
    │   ├── login.tsx                      # Public login
    │   ├── index.tsx                      # Redirect to /albums
    │   └── albums/
    │       ├── _layout.tsx                # CatalogViewerProvider
    │       ├── index.tsx                  # Responsive landing (mobile: list, desktop: first album)
    │       └── [owner]/
    │           └── [album]/
    │               ├── index.tsx          # Album media view
    │               └── [encodedId]/
    │                   └── [filename].tsx # Individual media
    ├── components/
    │   ├── auth/                          # Authentication (client)
    │   ├── catalog/                       # Business components
    │   ├── AppNav/                        # Navigation (server)
    │   └── user.menu/                     # User interface (client)
    └── lib/
        ├── auth/                          # Security logic (preserved)
        ├── catalog/                       # Data fetching + business logic
        └── contexts/                      # Client contexts

## Component Strategy

**Server Components**:
- Albums list rendering
- Media grid rendering
- Navigation structure
- Static layouts

**Client Components** (`"use client"`):
- All dialogs (Create/Edit/Delete/Share)
- Filtering and search
- Mobile navigation toggle
- User menu dropdown
- Authentication forms

## Migration Mapping

    Current                           →    Waku Target
    ─────────────────────────         →    ───────────────────
    src/core/security/               →    src/lib/auth/
    src/core/catalog/                →    src/lib/catalog/
    src/components/AppNav/           →    src/components/AppNav/ (server)
    src/components/user.menu/        →    src/components/user.menu/ (client)
    src/components/catalog-react/    →    src/lib/contexts/ + server fetching
    src/pages/authenticated/         →    src/pages/albums/
    src/libs/daction/dthunks/        →    Preserved for client state

## Authentication Flow

1. **Login**: Store refresh token in localStorage, access token in memory
2. **Page Load**: Restore session from refresh token
3. **Server Components**: Optional server-side auth validation
4. **Client Hydration**: Restore from localStorage

## Migration Phases

1. **Foundation**: Waku setup + auth system + root layout
2. **Core Features**: Albums/media components + navigation + dialogs
3. **Optimization**: Performance tuning + error handling
4. **Production**: CDN setup + deployment pipeline
