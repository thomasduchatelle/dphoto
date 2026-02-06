# Story 1.2: State Management Migration

**Status**: ready-for-dev

---

<!-- Part to be completed by the Scrum Master -->

## Story

As a **developer**,  
I want **to migrate the catalog state management from the existing web app**,  
so that **I can reuse battle-tested state logic with 90+ tests**.

## Acceptance Criteria

1. **State Management Migration**: State management code is copied to `web-nextjs/domains/catalog/` maintaining folder structure:
    - `language/` (state types - CatalogViewerState)
    - `actions.ts` (reducer and actions)
    - `album-create/`, `album-edit-*/`, `album-delete/` (thunk folders)

2. **Fetch Adapter**: A new fetch adapter is created at `domains/catalog/adapters/fetch-adapter.ts` replacing axios:
    - Works in both Server Components and Client Components
    - Implements the same interface as the original CatalogAPIAdapter

3. **Custom Image Loader**: Custom image loader is created at `libs/image-loader.ts` mapping:
    - width ≤40 → quality=blur
    - width ≤500 → quality=medium
    - width ≤1200 → quality=high
    - width >1200 → quality=full
    - Image loader is configured in `next.config.ts`

4. **Error Boundaries**: Error boundaries are created:
    - `app/error.tsx` (root catch-all)
    - `app/(authenticated)/error.tsx` (authenticated routes)

5. **Not-Found Pages**: Not-found pages are created:
    - `app/not-found.tsx`

6. **Tests Pass**: All migrated tests continue to pass (run with `npm run test`)

## Scope & Context

### Depends on

* **Story 1.1 (Project Foundation Setup)**: Material UI must be installed and configured with the dark theme before state management migration can proceed. The
  MUI theme and components will be used by the state-driven UI components.

### Expected outcomes

**For Story 1.3 (Basic Album List Display):**

- Complete catalog state management ready to use (reducer, actions, thunks)
- Fetch adapter that can be used in Server and Client Components
- Image loader configured for progressive loading
- Error handling infrastructure in place
- DEV agent can focus on building the UI components that consume the state

### Out of scope

* DO NOT implement any UI components - this story is purely about state management infrastructure
* DO NOT create the AlbumListClient or any page components - those belong to Story 1.3
* DO NOT integrate the state management with actual pages yet - only verify it works through tests
* DO NOT modify or add backend API endpoints - use existing REST API as-is

---

## Implementation Guidance

<!-- This part must be completed by the senior-dev -->

### Domain & Layer

- **Domain**: Catalog
- **Layer**: web-nextjs (domain logic + infrastructure utilities)
- **Coding standard, MUST READ BEFORE DEVELOPMENT**: `.github/instructions/nextjs.instructions.md`
- **Scope**: ✓ Single domain, single layer

### Architecture Decisions

This story implements **Architecture Decision #2: State Management Migration** from `specs/designs/architecture.md` (lines 102-186):

1. **Lift-and-Shift Strategy**: Copy battle-tested state management (90+ tests) from `web/src/core/catalog/` → `web-nextjs/domains/catalog/`
2. **Server-Side Initialization**: Server Components compute initial state, pass to Client Components as props
3. **Pure UI Pattern**: Components receive `state` and `handlers` as props - NO internal state management
4. **Adapter Replacement**: Replace axios with fetch adapter compatible with both Server and Client Components
5. **Two States Pattern**: User State (authentication, permissions) + Catalog State (albums, medias, filters, dialogs)
6. **Thunk Pattern**: Stateless, framework-agnostic business logic functions
7. **Action-Based Mutations**: Actions are the sole mechanism for state mutation
8. **Progressive Image Loading**: Custom Next.js image loader mapping width to backend quality levels

**Key Architectural Constraints:**

- ❌ DO NOT create new state management patterns
- ❌ DO NOT reimplement reducers, actions, or thunks
- ❌ DO NOT put state management logic inside UI components
- ❌ DO NOT create multiple separate contexts
- ✅ DO maintain exact same interface contracts (Ports)
- ✅ DO keep all tests passing with minimal changes
- ✅ DO make the state management server/client compatible

### Source Code Analysis

The existing state management in `web/src/core/catalog/` has this structure:

```
web/src/core/catalog/
├── language/                        # State type definitions (ubiquitous language)
│   ├── catalog-state.ts            # CatalogViewerState interface + all domain types
│   ├── initial-catalog-state.ts    # Factory for initial state
│   ├── errors.ts                   # CatalogError class
│   └── index.ts                    # Exports
│
├── actions.ts                       # Reducer factory + exports
│
├── thunks.ts                        # Aggregates all thunks
│
├── navigation/                      # Album list + filtering logic
│   ├── action-albumsLoaded.ts
│   ├── action-albumsFiltered.ts
│   ├── action-mediasLoaded.ts
│   ├── thunk-onPageRefresh.ts
│   ├── thunk-onAlbumFilterChange.ts
│   ├── selector-*.ts
│   └── *.test.ts
│
├── album-create/                    # Album creation flow
│   ├── action-*.ts
│   ├── thunk-*.ts
│   ├── selector-*.ts
│   └── *.test.ts
│
├── album-edit-dates/                # Edit album date range
├── album-edit-name/                 # Edit album name/folder
├── album-delete/                    # Delete album flow
├── sharing/                         # Share album with users
├── base-name-edit/                  # Shared name editing logic
├── date-range/                      # Shared date range logic
├── common/                          # Shared utilities
│
├── adapters/
│   └── api/
│       └── CatalogAPIAdapter.ts    # Axios-based API adapter (259 lines)
│
└── tests/
    └── test-helper-state.ts        # Test fixtures and constants
```

**Critical Files to Migrate:**

1. **47 test files** (*.test.ts) - MUST all pass after migration
2. **CatalogViewerState** - 210 lines defining all state shapes
3. **CatalogAPIAdapter** - 259 lines - MUST be reimplemented with fetch
4. **Generic reducer + actions** - Uses custom `daction` library
5. **Thunks** - Uses custom `dthunks` library for dependency injection
6. **Test helpers** - Used by all 47 tests for consistency

**Dependencies to Migrate:**

- `src/libs/daction` - Action factory and generic reducer (73 lines total)
- `src/libs/dthunks` - Thunk declaration patterns (40 lines)

### Migration Strategy

**Phase 1: Infrastructure Libraries (do this FIRST)**

1. Copy `web/src/libs/daction/` → `web-nextjs/libs/daction/`
    - `action-factory.ts` (71 lines) - NO changes needed
    - `reducer.ts` (11 lines) - NO changes needed
    - `index.ts` (2 lines) - NO changes needed
    - Tests: `action-factory.test.ts` - verify it passes

2. Copy `web/src/libs/dthunks/` → `web-nextjs/libs/dthunks/`
    - `api/index.ts` (40 lines) - ThunkDeclaration interface
    - `react/useThunks.ts` - NOT needed yet (Story 1.3)
    - `index.ts` (5 lines) - Export only the `api` subfolder
    - NO tests to run (framework code)

**Phase 2: Core Domain State**

3. Copy `web/src/core/catalog/language/` → `web-nextjs/domains/catalog/language/`
    - **NO code changes required** - pure TypeScript interfaces
    - Files to copy:
        - `catalog-state.ts` (210 lines) - all domain types
        - `initial-catalog-state.ts` - state factory
        - `errors.ts` (30 lines) - CatalogError class
        - `utils-albumIdEquals.ts` - helper
        - `selector-displayedAlbum.ts` - selector
        - `index.ts` - exports
    - **Import path updates**: Change `src/libs/daction` → `@/libs/daction`

4. Copy `web/src/core/catalog/tests/` → `web-nextjs/domains/catalog/tests/`
    - `test-helper-state.ts` (219 lines)
    - **Import path updates**: Change `src/` → `@/`
    - Remove reference to `AlbumsListActions` component (line 18) - replace with inline type definition

**Phase 3: Actions and Reducer**

5. Copy `web/src/core/catalog/actions.ts` → `web-nextjs/domains/catalog/actions.ts`
    - **NO code changes** - just import path updates

6. Copy ALL action folders maintaining structure:
    - `navigation/` (11 files) → keep folder structure
    - `album-create/` (9 files)
    - `album-edit-dates/` (11 files)
    - `album-edit-name/` (11 files)
    - `album-delete/` (9 files)
    - `sharing/` (9 files)
    - `base-name-edit/` (5 files)
    - `date-range/` (5 files)
    - `common/` (3 files - shared utilities)
    - **Import path updates only**: `src/libs/` → `@/libs/`, `src/core/` → `@/domains/`

7. Copy integration tests:
    - `catalog-acceptance.test.ts`
    - `catalog-state-behaviour.test.ts`
    - `catalog-factories.ts`

**Phase 4: Thunks**

8. Copy `web/src/core/catalog/thunks.ts` → `web-nextjs/domains/catalog/thunks.ts`
    - Aggregates all thunks from subdomains
    - **Import path updates only**

**Phase 5: API Adapter Replacement**

9. Create `web-nextjs/domains/catalog/adapters/fetch-adapter.ts`
    - **DO NOT copy** the axios adapter
    - Implement from scratch using fetch (see detailed implementation section below)
    - Must implement ALL interfaces from the original adapter:
        - `CreateAlbumPort`
        - `GrantAlbumAccessAPI`
        - `RevokeAlbumAccessAPI`
        - `DeleteAlbumPort`
        - `UpdateAlbumDatesPort`
        - `SaveAlbumNamePort`
    - Preserve exact same method signatures
    - Preserve exact same error handling (throw CatalogError)

10. Create index file `web-nextjs/domains/catalog/index.ts`
    - Export all public APIs from language, actions, thunks, adapters

**Phase 6: Run Tests and Fix Import Paths**

11. Run `npm run test` from `web-nextjs/`
    - Fix any remaining import path issues
    - ALL 47+ tests MUST pass

### Fetch Adapter Implementation

**Location**: `web-nextjs/domains/catalog/adapters/fetch-adapter.ts`

**Design Goals:**

1. Drop-in replacement for `CatalogAPIAdapter` (same interface)
2. Works in both Server Components and Client Components
3. Identical error handling (throw `CatalogError`)
4. Identical data transformations

**Key Differences from Axios Adapter:**

| Aspect               | Axios (old)                         | Fetch (new)                                     |
|----------------------|-------------------------------------|-------------------------------------------------|
| Dependency injection | `authenticatedAxios: AxiosInstance` | `baseUrl: string, getAccessToken: () => string` |
| Error handling       | `catch((err: AxiosError) => ...)`   | `if (!response.ok) throw ...`                   |
| Request config       | `axios.get(url, {params})`          | `fetch(url + '?' + params.toString())`          |
| Response parsing     | `resp.data`                         | `await resp.json()`                             |
| Server/Client        | Client only                         | Server + Client compatible                      |

**Implementation Pattern:**

```typescript
// web-nextjs/domains/catalog/adapters/fetch-adapter.ts
import {CatalogError} from "../language/errors";
import type {Album, AlbumId, Media,

...
}
from
"../language";
import type {CreateAlbumPort, DeleteAlbumPort,

...
}
from
"../thunks";

export class CatalogFetchAdapter implements CreateAlbumPort,
    DeleteAlbumPort,
    UpdateAlbumDatesPort,
    SaveAlbumNamePort,
    GrantAlbumAccessAPI,
    RevokeAlbumAccessAPI {

    constructor(
        private readonly baseUrl: string,
        private readonly getAccessToken: () => string
    ) {
    }

    // Helper: make authenticated request
    private async request<T>(
        method: string,
        path: string,
        options?: { body?: any, params?: Record<string, string> }
    ): Promise<T> {
        const token = this.getAccessToken();
        const url = new URL(path, this.baseUrl);

        if (options?.params) {
            Object.entries(options.params).forEach(([key, value]) => {
                url.searchParams.append(key, value);
            });
        }

        const response = await fetch(url.toString(), {
            method,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: options?.body ? JSON.stringify(options.body) : undefined,
            cache: 'no-store', // Always fetch fresh data
        });

        if (!response.ok) {
            await this.handleError(method, path, response);
        }

        if (response.status === 204) {
            return undefined as T;
        }

        return await response.json();
    }

    private async handleError(method: string, path: string, response: Response): Promise<never> {
        const defaultMessage = `'${method.toUpperCase()} ${path}' failed with status ${response.status} ${response.statusText}`;

        try {
            const errorData = await response.json();

            if (
                response.status >= 400 &&
                response.status < 500 &&
                errorData.code &&
                errorData.message
            ) {
                throw new CatalogError(errorData.code, errorData.message ?? defaultMessage);
            }
        } catch (e) {
            if (e instanceof CatalogError) throw e;
        }

        throw new CatalogError("", defaultMessage);
    }

    // Copy method implementations from CatalogAPIAdapter.ts
    // Replace axios calls with this.request<T>(method, path, options)
    // Keep all data transformations IDENTICAL

    public async fetchAlbums(): Promise<Album[]> {
        const albums = await this.request<RestAlbum[]>('GET', '/api/v1/albums');

        // ... same logic as original adapter (lines 66-122)
        // Data transformation MUST be identical
    }

    public async deleteAlbum(albumId: AlbumId): Promise<void> {
        await this.request<void>(
            'DELETE',
            `/api/v1/owners/${albumId.owner}/albums/${albumId.folderName}`
        );
    }

    // ... implement all other methods with same pattern
}
```

**Critical Implementation Rules:**

1. **Preserve error handling exactly**: Cast HTTP 4xx errors with code+message to `CatalogError` (lines 218-234 of original)
2. **Preserve data transformations exactly**: Date parsing, sorting, temperature calculation (lines 66-122)
3. **Handle 404 specially**: `fetchMedias` returns empty array on 404 (lines 158-161)
4. **Access token in URLs**: Media content paths must include `?access_token=${token}` (line 170)
5. **Parallel requests**: Use `Promise.allSettled` for owners + users (lines 71-92)

**Testing Strategy:**

- Run existing adapter tests against new implementation
- Tests use fake implementations - they should pass without changes
- If tests fail, the adapter interface contract is broken

### Image Loader Implementation

**Location**: `web-nextjs/libs/image-loader/index.ts`

**Purpose**: Map Next.js Image component width requests to backend API quality parameters.

**Backend API Contract** (from architecture doc):

```
GET /api/v1/media/{mediaId}/image?quality={quality}&width={width}

Quality Levels:
- blur: 20-40px, <2KB (instant placeholder)
- medium: 50-150KB, <500ms on 3G (grid display)
- high: near-original (full-screen viewer)
- full: original (zoom interactions)

Cache Headers: Cache-Control: max-age=31536000, immutable
```

**Implementation:**

```typescript
// web-nextjs/libs/image-loader/index.ts

export type ImageQuality = 'blur' | 'medium' | 'high' | 'full';

export interface ImageLoaderProps {
    src: string;
    width: number;
    quality?: number;
}

export function imageLoader({src, width, quality}: ImageLoaderProps): string {
    const backendQuality = mapWidthToQuality(width);

    // Parse existing URL to preserve query params
    const url = new URL(src, 'http://placeholder.local');
    url.searchParams.set('quality', backendQuality);
    url.searchParams.set('width', width.toString());

    // Return path + search params (Next.js will add the domain)
    return url.pathname + url.search;
}

function mapWidthToQuality(width: number): ImageQuality {
    if (width <= 40) return 'blur';
    if (width <= 500) return 'medium';
    if (width <= 1200) return 'high';
    return 'full';
}
```

**Configuration in `next.config.ts`:**

```typescript
// web-nextjs/next.config.ts
import type {NextConfig} from "next";

const nextConfig: NextConfig = {
    output: "standalone",
    images: {
        loader: 'custom',
        loaderFile: './libs/image-loader/index.ts',
        remotePatterns: [
            {
                protocol: 'https',
                hostname: '**',
            }
        ]
    },
    basePath: '/nextjs',
};

export default nextConfig;
```

**Usage Example (for Story 1.3+):**

```tsx
import Image from 'next/image';

<Image
    src="/api/v1/owners/myself/medias/media-123/photo.jpg"
    width={500}
    height={300}
    alt="Album photo"
    placeholder="blur"
    blurDataURL="/api/v1/owners/myself/medias/media-123/photo.jpg?quality=blur&width=40"
/>
```

**Testing Strategy:**

- Create `image-loader.test.ts` with vitest
- Test width boundaries: 40, 41, 500, 501, 1200, 1201
- Test URL preservation: ensure existing query params are not lost
- Test integration: verify Next.js can import and use the loader

### Error Boundaries Implementation

**Goal**: Implement NextJS App Router error handling pattern (Architecture Decision #7, lines 325-356).

**Structure:**

```
app/
├── error.tsx                              # Root catch-all
├── not-found.tsx                          # Root 404
└── (authenticated)/
    ├── error.tsx                          # Authenticated routes error
    └── ... (future: album-specific errors in Story 2.x)
```

**File 1: `app/error.tsx`** (Root Error Boundary)

```tsx
'use client';

import {useEffect} from 'react';
import {Box, Button, Typography} from '@mui/material';

export default function Error({
                                  error,
                                  reset,
                              }: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        console.error('Root error boundary:', error);
    }, [error]);

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100vh',
                padding: 3,
                gap: 2,
            }}
        >
            <Typography variant="h4" component="h1">
                Something went wrong
            </Typography>
            <Typography variant="body1" color="text.secondary">
                {error.message || 'An unexpected error occurred'}
            </Typography>
            {error.digest && (
                <Typography variant="caption" color="text.secondary">
                    Error ID: {error.digest}
                </Typography>
            )}
            <Button variant="contained" onClick={reset}>
                Try Again
            </Button>
        </Box>
    );
}
```

**File 2: `app/(authenticated)/error.tsx`** (Authenticated Routes)

```tsx
'use client';

import {useEffect} from 'react';
import {Box, Button, Typography} from '@mui/material';
import Link from '@/components/Link';

export default function AuthenticatedError({
                                               error,
                                               reset,
                                           }: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        console.error('Authenticated route error:', error);
    }, [error]);

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '70vh',
                padding: 3,
                gap: 2,
            }}
        >
            <Typography variant="h5" component="h1">
                Unable to Load Page
            </Typography>
            <Typography variant="body1" color="text.secondary" textAlign="center">
                {error.message || 'An error occurred while loading this page'}
            </Typography>
            {error.digest && (
                <Typography variant="caption" color="text.secondary">
                    Error ID: {error.digest}
                </Typography>
            )}
            <Box sx={{display: 'flex', gap: 2}}>
                <Button variant="contained" onClick={reset}>
                    Try Again
                </Button>
                <Button
                    variant="outlined"
                    component={Link}
                    href="/"
                    prefetch={false}
                >
                    Go Home
                </Button>
            </Box>
        </Box>
    );
}
```

**File 3: `app/not-found.tsx`** (Root 404)

```tsx
import {Box, Button, Typography} from '@mui/material';
import Link from '@/components/Link';

export default function NotFound() {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100vh',
                padding: 3,
                gap: 2,
            }}
        >
            <Typography variant="h4" component="h1">
                Page Not Found
            </Typography>
            <Typography variant="body1" color="text.secondary">
                The page you are looking for does not exist.
            </Typography>
            <Button
                variant="contained"
                component={Link}
                href="/"
                prefetch={false}
            >
                Go Home
            </Button>
        </Box>
    );
}
```

**Critical Rules:**

1. Error boundaries MUST be Client Components (`'use client'`)
2. Not-found pages SHOULD be Server Components (no interactivity needed)
3. Use `@/components/Link` for Material UI Button links (NOT `next/link` directly)
4. Log errors with `useEffect` for debugging
5. Display error digest when available (helps trace errors in production)
6. Provide contextual recovery: "Try Again" + "Go Home" for authenticated routes

### Testing Strategy

**Goal**: Ensure all 47+ migrated tests pass without breaking the battle-tested state management.

**Test Execution:**

```bash
cd web-nextjs
npm install                    # Install dependencies
npm run test                   # Run all vitest tests
```

**Expected Test Count:**

- 47+ action/selector/thunk tests from catalog domain
- 2-3 tests for daction library
- 1+ test for image loader
- **Total: ~50 tests, all must pass**

**Common Migration Issues and Fixes:**

1. **Import path errors**:
    - Change `src/libs/` → `@/libs/`
    - Change `src/core/` → `@/domains/`
    - Change relative imports to absolute when crossing domain boundaries

2. **Missing test helper types**:
    - Remove import of `AlbumsListActions` from `test-helper-state.ts`
    - Define `AlbumListActionsProps` inline in the test helpers:
   ```typescript
   // Replace line 18 in test-helper-state.ts
   export interface AlbumListActionsProps {
       albumFilter: AlbumFilterEntry;
       albumFilterOptions: AlbumFilterEntry[];
       displayedAlbumIdIsOwned: boolean;
       hasAlbumsToDelete: boolean;
       canCreateAlbums: boolean;
   }
   ```

3. **Fetch adapter not tested directly**:
    - Tests use FAKE implementations of Ports (not the actual adapter)
    - Tests validate thunk behavior with any adapter that implements the Port interface
    - This is intentional: tests are decoupled from adapter implementation

**Test Debugging Strategy:**

If tests fail:

1. Check import paths first (90% of migration issues)
2. Verify all files copied correctly (use diff tool)
3. Check for TypeScript errors: `npm run type-check` (if available)
4. Run one test file at a time to isolate issues:
   ```bash
   npm run test -- action-albumsLoaded.test.ts
   ```
5. Compare original vs migrated file side-by-side for unintended changes

**DO NOT**:

- ❌ Modify test logic to make them pass
- ❌ Skip failing tests
- ❌ Change state management behavior
- ❌ Simplify or remove tests

**DO**:

- ✅ Fix import paths
- ✅ Update test helper dependencies
- ✅ Ensure exact same behavior as original

### File Structure

Complete directory structure after migration:

```
web-nextjs/
├── app/
│   ├── error.tsx                                    # NEW: Root error boundary
│   ├── not-found.tsx                                # NEW: Root 404 page
│   └── (authenticated)/
│       └── error.tsx                                # NEW: Authenticated error boundary
│
├── domains/
│   └── catalog/
│       ├── language/                                # MIGRATED: State types
│       │   ├── catalog-state.ts                     # (210 lines) Domain types
│       │   ├── initial-catalog-state.ts             # State factory
│       │   ├── errors.ts                            # (30 lines) CatalogError
│       │   ├── selector-displayedAlbum.ts           # Selector
│       │   ├── utils-albumIdEquals.ts               # Helper
│       │   └── index.ts                             # Exports
│       │
│       ├── actions.ts                               # MIGRATED: Reducer + exports
│       │
│       ├── thunks.ts                                # MIGRATED: Thunk aggregator
│       │
│       ├── navigation/                              # MIGRATED: 11 files
│       │   ├── action-albumsLoaded.ts
│       │   ├── action-albumsLoaded.test.ts
│       │   ├── action-albumsFiltered.ts
│       │   ├── action-albumsFiltered.test.ts
│       │   ├── action-mediasLoaded.ts
│       │   ├── action-mediasLoaded.test.ts
│       │   ├── action-albumsAndMediasLoaded.ts
│       │   ├── action-albumsAndMediasLoaded.test.ts
│       │   ├── action-noAlbumAvailable.ts
│       │   ├── action-noAlbumAvailable.test.ts
│       │   ├── action-mediaLoadFailed.ts
│       │   ├── action-mediaLoadFailed.test.ts
│       │   ├── thunk-onPageRefresh.ts
│       │   ├── thunk-onPageRefresh.test.ts
│       │   ├── thunk-onAlbumFilterChange.ts
│       │   ├── thunk-onAlbumFilterChange.test.ts
│       │   ├── selector-catalog-viewer-page.ts
│       │   ├── selector-albumListActions.ts
│       │   ├── utils-loadAlbumsAndMedias.ts
│       │   ├── group-by-day.ts
│       │   └── index.ts
│       │
│       ├── album-create/                            # MIGRATED: 9 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   ├── selector-*.ts
│       │   └── index.ts
│       │
│       ├── album-edit-dates/                        # MIGRATED: 11 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   ├── selector-*.ts
│       │   └── index.ts
│       │
│       ├── album-edit-name/                         # MIGRATED: 11 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   ├── selector-*.ts
│       │   └── index.ts
│       │
│       ├── album-delete/                            # MIGRATED: 9 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   ├── selector-*.ts
│       │   └── index.ts
│       │
│       ├── sharing/                                 # MIGRATED: 9 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   ├── selector-*.ts
│       │   └── index.ts
│       │
│       ├── base-name-edit/                          # MIGRATED: 5 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   └── index.ts
│       │
│       ├── date-range/                              # MIGRATED: 5 files
│       │   ├── action-*.ts + tests
│       │   ├── thunk-*.ts + tests
│       │   └── index.ts
│       │
│       ├── common/                                  # MIGRATED: 3 files
│       │   ├── utils.ts                             # Shared utilities
│       │   ├── catalog-factory-args.ts
│       │   └── index.ts
│       │
│       ├── adapters/
│       │   ├── fetch-adapter.ts                     # NEW: Fetch-based adapter (~300 lines)
│       │   └── index.ts                             # NEW: Exports
│       │
│       ├── tests/
│       │   └── test-helper-state.ts                 # MIGRATED: Test fixtures (219 lines)
│       │
│       ├── catalog-acceptance.test.ts               # MIGRATED: Integration tests
│       ├── catalog-state-behaviour.test.ts          # MIGRATED: State tests
│       ├── catalog-factories.ts                     # MIGRATED: Test factories
│       └── index.ts                                 # NEW: Public API exports
│
├── libs/
│   ├── daction/                                     # MIGRATED: Action framework
│   │   ├── action-factory.ts                        # (71 lines)
│   │   ├── action-factory.test.ts
│   │   ├── reducer.ts                               # (11 lines)
│   │   └── index.ts
│   │
│   ├── dthunks/                                     # MIGRATED: Thunk framework
│   │   ├── api/
│   │   │   └── index.ts                             # (40 lines) ThunkDeclaration
│   │   └── index.ts
│   │
│   ├── image-loader/                                # NEW: Image optimization
│   │   ├── index.ts                                 # Custom loader (~30 lines)
│   │   └── image-loader.test.ts                     # NEW: Tests
│   │
│   └── ... (existing: nextjs-cookies, requests, security)
│
├── next.config.ts                                   # UPDATED: Add image loader config
└── vitest.config.ts                                 # (existing)
```

**File Count Summary:**

- **Migrated files**: ~70 (language, actions, thunks, tests)
- **New files**: 6 (fetch-adapter, image-loader, 3 error boundaries)
- **Updated files**: 1 (next.config.ts)
- **Test files**: 47+ (all must pass)

### Implementation Notes

**1. Order of Implementation is Critical**

Follow the migration phases strictly:

1. Infrastructure libs (daction, dthunks) FIRST
2. Domain state types (language) SECOND
3. Test helpers THIRD
4. Actions, reducer, thunks FOURTH
5. Fetch adapter FIFTH
6. Image loader + error boundaries LAST

**Why?** Dependencies cascade: actions depend on daction, tests depend on test helpers, thunks depend on language types.

**2. Import Path Strategy**

Use TypeScript path aliases consistently:

- `@/libs/*` - Framework utilities
- `@/domains/*` - Domain logic
- `@/components/*` - React components
- Never use `src/*` - that's the old web app

**3. Adapter Constructor Pattern**

The new adapter must work in both environments:

**Server Components:**

```typescript
const adapter = new CatalogFetchAdapter(
    process.env.API_BASE_URL!,
    () => getServerSideToken()  // From cookies
);
```

**Client Components:**

```typescript
const adapter = new CatalogFetchAdapter(
    '/api',  // Relative URL
    () => getClientSideToken()  // From context/state
);
```

**4. Error Handling Edge Cases**

`CatalogError` must be thrown for 4xx errors with structured response:

```json
{
  "code": "ALBUM_NOT_FOUND",
  "message": "Album 'jan-25' does not exist"
}
```

BUT fallback to generic error for:

- Network errors
- 5xx errors
- Malformed error responses

**5. Test Isolation**

Tests use FAKES, not MOCKS:

```typescript
// ✅ GOOD: Fake implementation
class FakeCatalogAdapter implements DeleteAlbumPort {
    constructor(private albumsToDelete: Set<string>) {
    }

    async deleteAlbum(albumId: AlbumId): Promise<void> {
        if (this.albumsToDelete.has(albumId.folderName)) {
            return;
        }
        throw new Error('Album not found');
    }
}

// ❌ BAD: Mock
const mockAdapter = {
    deleteAlbum: jest.fn().mockResolvedValue(undefined)
};
```

**Why?** Fakes test the Port interface contract. Mocks test method calls.

**6. State Immutability**

ALL reducers return new state objects:

```typescript
// ✅ GOOD
return {
    ...state,
    albums: newAlbums,
    error: undefined,
};

// ❌ BAD
state.albums = newAlbums;
state.error = undefined;
return state;
```

**7. Image Loader Configuration**

The loader file path in `next.config.ts` is relative to the config file:

```typescript
loaderFile: './libs/image-loader/index.ts',  // ✅ Relative to next.config.ts
```

NOT:

```typescript
loaderFile: '@/libs/image-loader/index.ts',  // ❌ Alias won't work
    loaderFile
:
'libs/image-loader/index.ts',    // ❌ Missing ./
```

**8. Type-Safety Checklist**

Before committing:

- ✅ No `any` types (except in daction/dthunks framework code)
- ✅ All Ports explicitly typed
- ✅ All actions have explicit State and Payload types
- ✅ All selectors have explicit return types
- ✅ No TypeScript errors: Run build to verify

**9. What NOT to Change**

Do NOT modify these aspects during migration:

- ❌ State shape (CatalogViewerState fields)
- ❌ Action payloads
- ❌ Thunk signatures
- ❌ Port interfaces
- ❌ Test assertions
- ❌ Business logic
- ❌ Selector return types

ONLY change:

- ✅ Import paths
- ✅ Adapter implementation (axios → fetch)
- ✅ Framework compatibility (add server/client support)

**10. Common Pitfalls**

**Pitfall 1**: Forgetting `cache: 'no-store'` in fetch

```typescript
// ❌ BAD: Will cache stale data
await fetch(url, {headers: {...}});

// ✅ GOOD: Always fetch fresh
await fetch(url, {headers: {...}, cache: 'no-store'});
```

**Pitfall 2**: Not handling empty response bodies

```typescript
// ❌ BAD: Will throw on DELETE (204 No Content)
return await response.json();

// ✅ GOOD: Check status first
if (response.status === 204) return undefined as T;
return await response.json();
```

**Pitfall 3**: Incorrect error propagation

```typescript
// ❌ BAD: Swallows CatalogError
try {
    const data = await response.json();
    throw new CatalogError(data.code, data.message);
} catch (e) {
    throw new Error('Request failed');
}

// ✅ GOOD: Preserves CatalogError
try {
    const data = await response.json();
    throw new CatalogError(data.code, data.message);
} catch (e) {
    if (e instanceof CatalogError) throw e;
    throw new Error('Request failed');
}
```

**11. Definition of Done**

Before marking this story complete:

- [ ] All 47+ tests pass (`npm run test`)
- [ ] No TypeScript errors (`npm run build` or `tsc --noEmit`)
- [ ] Image loader is configured in `next.config.ts`
- [ ] Error boundaries render correctly (manual test: throw error in component)
- [ ] Not-found page renders correctly (manual test: visit `/not-a-page`)
- [ ] Fetch adapter implements ALL 6 Port interfaces
- [ ] Import paths use `@/` aliases consistently
- [ ] No references to `src/` remain in migrated code
- [ ] Test helper types are self-contained (no UI component imports)

**12. Story 1.3 Readiness**

After completing this story, Story 1.3 can:

1. Import catalog state: `import {CatalogViewerState, catalogReducer} from '@/domains/catalog'`
2. Import thunks: `import {catalogThunks} from '@/domains/catalog'`
3. Import adapter: `import {CatalogFetchAdapter} from '@/domains/catalog/adapters'`
4. Use image loader: `import Image from 'next/image'` (auto-configured)
5. Rely on error boundaries: Errors will be caught automatically
6. Use test helpers: `import {twoAlbums, loadedStateWithTwoAlbums} from '@/domains/catalog/tests/test-helper-state'`

---

**END OF IMPLEMENTATION GUIDANCE**

---

## Implementation report

This part must be completed by the DEV agent to summarise the changes made to implement this story:

* What was the problem
* What has been done to solve it
* Results and screenshots when possible
