# DPhoto Web Application - Waku Migration Specification

## Overview

Migration strategy for DPhoto web application from Create React App (CRA) with React Router to Waku framework, maintaining current functionality while improving performance through server-side rendering and file-based routing.

## Application Context

### Current Architecture
- **Multi-tenant application**: Each user manages their own albums and media
- **Sharing capabilities**: Read-only access to other users' albums/media
- **Authentication required**: All routes except login page require authentication
- **No open subscription**: Closed user system
- **Responsive design**: Different default landing pages for mobile vs desktop

### Current Tech Stack
- Create React App (CRA)
- React Router for client-side routing
- TypeScript
- Material-UI components
- Redux-like state management (custom daction/dthunks)
- Go backend with JWT authentication

## Migration Strategy

### 1. Deployment Architecture

**Target: Static Site Generation (SSG) with CDN**

    ┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
    │   CloudFront    │───▶│  S3 Static Site  │───▶│  Go Backend     │
    │   (CDN)         │    │  (Waku Build)    │    │  (API/Auth)     │
    └─────────────────┘    └──────────────────┘    └─────────────────┘

**Benefits:**
- **Performance**: Static files served from CDN edge locations
- **Cost-effective**: No server running costs, pay-per-request
- **Scalability**: Automatic scaling through CDN
- **SEO**: Server-side rendered HTML for better search indexing
- **Reliability**: High availability through CDN infrastructure

**Build Process:**
1. Waku builds static HTML files for all routes
2. Client-side hydration for interactivity
3. API calls to existing Go backend remain unchanged

### 2. Authentication Strategy

**Hybrid Approach: localStorage + Server-Side Validation**

#### Token Management
```typescript
// Storage Strategy
const authState = {
  accessToken: string,      // Memory only - short lived (15 min)
  user: AuthenticatedUser,  // Memory only - will refetch on page load
  refreshToken: string      // localStorage - long lived (30 days)
}
```

#### Authentication Flow
1. **Login**: Store refresh token in localStorage, access token in memory
2. **Page Load**: Attempt session restoration using refresh token
3. **Token Refresh**: Automatic renewal using existing logic
4. **Server Components**: Can validate authentication server-side during build/SSR
5. **Client Hydration**: Restores session from localStorage refresh token

#### Security Benefits
- Access tokens never persisted (XSS protection)
- Refresh tokens in localStorage for UX (industry standard)
- Server-side authentication validation capability
- Existing security logic remains unchanged

### 3. Routing Migration

**From React Router to File-Based Routing**

#### Current Routing Structure
    /login                    → Login page
    /albums                   → Albums list (mobile) / First album + sidebar (desktop)
    /albums/:owner/:album     → Specific album view
    /albums/:owner/:album/:encodedId/:filename → Individual media view

#### Target Waku File Structure
    src/
    ├── pages/
    │   ├── _layout.tsx                    # Root auth guard + AppNav
    │   ├── login.tsx                      # Public login page
    │   ├── index.tsx                      # Redirect to /albums
    │   └── albums/
    │       ├── _layout.tsx                # CatalogViewerProvider wrapper
    │       ├── index.tsx                  # Responsive landing page
    │       └── [owner]/
    │           └── [album]/
    │               ├── index.tsx          # Album media view
    │               └── [encodedId]/
    │                   └── [filename].tsx # Individual media page
    ├── components/
    │   ├── auth/                          # Authentication components
    │   ├── catalog/                       # Business logic components
    │   ├── AppNav/                        # Navigation components
    │   └── user.menu/                     # User interface components
    └── lib/
        ├── auth/                          # Existing security logic
        ├── catalog/                       # Server-side data fetching
        └── contexts/                      # Client-side contexts

#### Responsive Default Behavior
- **Mobile**: Lands on albums list (`/albums`)
- **Desktop**: Shows first album's media with albums sidebar
- **Server-side detection**: Device type detection during SSR
- **Fallback**: Redirect to `/albums` for unknown routes
- **Deep linking**: Maintain API Gateway redirect support (`?path=` parameter)

### 4. Component Architecture

**Server vs Client Component Strategy**

#### Server Components (Performance + SEO)
- Albums list rendering
- Media grid rendering
- Album metadata display
- Navigation structure (albums sidebar)
- Individual media page content
- Static layouts and wrappers

#### Client Components (Interactivity)
- All dialogs (Create/Edit/Delete/Share)
- Filtering dropdown and search
- Mobile navigation toggle
- User menu dropdown
- Authentication logic and forms
- State management components

#### Component Boundaries
```typescript
// Server Component Example
// src/components/catalog/AlbumsList.tsx
export default function AlbumsList({ albums }: { albums: Album[] }) {
  return (
    <div>
      {albums.map(album => (
        <AlbumCard key={album.id} album={album} />
      ))}
      <AlbumsListActions /> {/* Client component for filtering */}
    </div>
  )
}

// Client Component Example
// src/components/catalog/AlbumsListActions.tsx
"use client"
export default function AlbumsListActions() {
  // Interactive filtering logic
}
```

## Migration Benefits

### Performance Improvements
- **Faster initial load**: Server-rendered HTML
- **Better Core Web Vitals**: Reduced JavaScript bundle size
- **CDN delivery**: Static assets served from edge locations
- **Progressive enhancement**: Works without JavaScript

### Developer Experience
- **File-based routing**: No router configuration needed
- **TypeScript support**: Built-in TypeScript support
- **Component co-location**: Routes and components in same directory
- **Simplified deployment**: Static build output

### SEO & Accessibility
- **Server-side rendering**: Search engines can index content
- **Direct URLs**: Deep links work without client-side routing
- **Progressive enhancement**: Accessible without JavaScript
- **Meta tags**: Server-rendered meta information

## Migration Phases

### Phase 1: Foundation
1. Set up Waku project structure
2. Migrate authentication system
3. Create root layout with auth guard
4. Implement responsive albums landing page

### Phase 2: Core Features
1. Migrate albums list and media grid (server components)
2. Migrate navigation and routing
3. Implement client-side dialogs and interactions
4. Test authentication flow

### Phase 3: Polish & Optimization
1. Optimize server-side data fetching
2. Implement error boundaries and loading states
3. Performance testing and optimization
4. Deployment pipeline setup

### Phase 4: Production
1. CDN configuration
2. Domain setup and SSL
3. Monitoring and analytics
4. Rollback strategy

## Risks & Mitigations

### Technical Risks
- **SSR complexity**: Server-side data fetching changes
  - *Mitigation*: Gradual migration, maintain API compatibility
- **Client-side state**: Managing hydration mismatches
  - *Mitigation*: Careful server/client component boundaries
- **Bundle size**: Potential JavaScript bundle increase
  - *Mitigation*: Code splitting and tree shaking optimization

### Business Risks
- **Migration downtime**: User experience disruption
  - *Mitigation*: Feature flagging and gradual rollout
- **Authentication issues**: Login/session problems
  - *Mitigation*: Extensive testing and rollback plan
- **Performance regression**: Slower than current CRA
  - *Mitigation*: Performance benchmarking throughout migration

## Success Criteria

- **Performance**: Lighthouse score > 90 for all metrics
- **Functionality**: All existing features work identically
- **Authentication**: Seamless user login and session management
- **Responsive**: Mobile and desktop experiences maintained
- **SEO**: All public pages indexed and discoverable
- **Deployment**: Successful CDN deployment with < 1 second global load times

## Current Architecture Analysis

### Existing Project Structure
Based on the current `src/` directory structure:

#### Core Business Logic (`src/core/`)
- **Application layer**: Context, hooks, state management (daction/dthunks)
- **Security**: Authentication cases and state management
- **Catalog**: Business logic for albums and media management
- **Utils**: Shared utilities and helpers

#### Component Architecture (`src/components/`)
- **AppNav**: Main navigation component
- **catalog-react**: Catalog context and hooks
- **user.menu**: User interface components
- **DPhotoTheme**: Material-UI theme configuration

#### Current Routing (`src/pages/`)
- **GeneralRouter**: Main router configuration
- **Login**: Authentication pages
- **authenticated/**: Protected routes with nested album/media routing

#### State Management
- **daction**: Custom action factory system
- **dthunks**: Custom async action handling
- **Redux-like patterns**: Reducers and state management

### Migration Mapping

#### Core Logic Preservation
The existing `src/core/` business logic can be largely preserved:
- Security authentication flows remain unchanged
- Catalog business logic stays intact
- State management patterns adapt to Waku's server/client model

#### Component Migration Strategy
    Current Structure              →    Waku Structure
    ─────────────────────────      →    ─────────────────
    src/components/AppNav/         →    src/components/AppNav/ (server component)
    src/components/user.menu/      →    src/components/user.menu/ ("use client")
    src/components/catalog-react/  →    src/lib/contexts/ + server data fetching
    src/pages/authenticated/       →    src/pages/albums/ (file-based routing)
    src/core/security/            →    src/lib/auth/ (preserved)
    src/core/catalog/             →    src/lib/catalog/ (preserved)

#### State Management Evolution
- **daction/dthunks**: Preserved for client-side state
- **Server components**: Direct data fetching without state management
- **Context providers**: Migrate to client components for interactive features
- **Authentication state**: Hybrid server/client validation

This analysis ensures the migration preserves all existing business logic while optimizing for Waku's architecture patterns.