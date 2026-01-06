# DPhoto - Agent Development Guide

This guide provides essential information for AI coding agents working on the DPhoto codebase.

## Quick Reference

### Build Commands

```bash
# Root level - build everything
make all                    # clean, test, build all projects
make test                   # run all tests
make build                  # build all projects

# Go - pkg/ and cmd/dphoto/
make setup-go              # REQUIRED before testing (starts Docker services)
go test ./...              # run all tests
go test ./pkg/catalog/...  # test specific package
go build ./cmd/dphoto      # build CLI

# Go - api/lambdas/
cd api/lambdas && go test ./...    # test API lambdas
make build-api                      # build all lambda handlers

# TypeScript - web-nextjs/ (current web app)
cd web-nextjs
npm install                # REQUIRED before other commands
npm run test               # run unit tests (~5s)
npm run test:watch         # watch mode
npm run build              # build for production
npm run dev                # start dev server (port 3000)

# TypeScript - web/ (legacy, being replaced)
cd web
npm install                # REQUIRED before other commands
npx vitest run            # run unit tests (~17s)
npx playwright test       # run visual regression tests
npm run build             # build for deployment
npm run ladle             # component viewer (port 61000)

# TypeScript - deployments/cdk/
cd deployments/cdk
npm install               # REQUIRED before other commands
npm test                  # run unit tests
npm run synth:test        # verify CDK can build
```

### Running Single Tests

```bash
# Go - run specific test function
go test -run TestCreateAlbumStateless_Create ./pkg/catalog/

# Go - run with verbose output
go test -v ./pkg/backup/

# TypeScript (web-nextjs) - run specific test file
npm run test -- access-token-service.test.ts

# TypeScript (web) - run specific test file
npx vitest run src/core/catalog/navigation/thunk-onPageRefresh.test.ts

# Playwright - run specific test
npx playwright test --grep "timeline view"
```

## Project Architecture

### Directory Structure

```
pkg/                  - Go core business logic (Hexagonal Architecture)
  ├── acl/           - Access control and permissions
  ├── archive/       - Long-term storage, compression, miniatures
  ├── catalog/       - Media organization into albums
  └── backup/        - Analyze and load medias
cmd/dphoto/          - Go CLI (Cobra framework)
api/lambdas/         - Go REST API handlers (AWS Lambda)
web-nextjs/          - TypeScript/React 19/NextJS (current)
web/                 - TypeScript/React 19/Waku (DEPRECATED)
deployments/cdk/     - AWS CDK infrastructure (TypeScript)
internal/            - Go mocks and test utilities
DATA_MODEL.md        - DynamoDB single-table structure
```

### Key Principles

1. **No data loss** - medias are irreplaceable, never lose data
2. **Architecture integrity** - follow established patterns strictly
3. **Simplicity** - prefer refactoring over complexity
4. **Cost efficiency** - minimize AWS operating costs
5. **Security** - prevent data leaks with reasonable practices

## Go Code Style

### Import Organization

```go
import (
    // 1. Standard library
    "context"
    "fmt"
    
    // 2. Third-party packages
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    
    // 3. Internal packages
    "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)
```

### Error Handling

```go
// Sentinel errors with Err suffix
var (
    AlbumNotFoundErr = errors.New("album hasn't been found")
    EmptyOwnerError = errors.New("owner is mandatory and must be not empty")
)

// Always wrap errors with context
return nil, errors.Wrapf(err, "CreateNewAlbum(%s) failed", request)

// Check errors immediately
if err != nil {
    return nil, err
}
```

### Naming Conventions

```go
// Interfaces - Port suffix for hexagonal architecture
type FindAlbumsByOwnerPort interface {
    FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error)
}

// Interfaces - Observer suffix for event handlers
type CreateAlbumObserver interface {
    ObserveCreateAlbum(ctx context.Context, album Album) error
}

// Constructors - New prefix, return pointer
func NewAlbumCreate(
    FindAlbumsByOwnerPort FindAlbumsByOwnerPort,
    InsertAlbumPort InsertAlbumPort,
) *CreateAlbum

// Value objects - type aliases with validation methods
type Owner string
func (o Owner) IsValid() error { /* ... */ }
```

### Type Usage

```go
// Pointer receivers for mutable structs
func (c *CreateAlbum) Create(ctx context.Context, req CreateAlbumRequest) (*AlbumId, error)

// Value receivers for immutable types
func (a Album) IsEqual(other *Album) bool
func (a AlbumId) String() string

// Small, focused interfaces (Interface Segregation Principle)
type InsertAlbumPort interface {
    InsertAlbum(ctx context.Context, album Album) error
}
```

### Testing Patterns

```go
// Table-driven tests
func TestCreateAlbum_Create(t *testing.T) {
    tests := []struct {
        name    string
        args    args
        wantErr assert.ErrorAssertionFunc
    }{
        {
            name: "it should create album successfully",
            // ...
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

// Test doubles with Fake suffix
type FindAlbumsByOwnerFake map[Owner][]*Album

// Helper functions with stub/expect prefix
func stubFindAlbumsWith(albums ...*Album) FindAlbumsByOwnerPort
func expectAlbumInserted(album Album) InsertAlbumPort
```

### Documentation

```go
// Package documentation at top
// Package catalog provides commands to organize medias into albums.
package catalog

// Type documentation (godoc style)
// Album is a logical grouping of medias stored together.
type Album struct {
    Name  string    // Name for display, not unique
    Start time.Time // Start datetime inclusive
}

// Function documentation starts with function name
// NewAlbumCreate creates the service to create a new album.
func NewAlbumCreate(...) *CreateAlbum
```

## TypeScript Code Style

### Import Organization (web-nextjs/)

```typescript
// 1. Server/client directive at top
import "server-only"  // or 'use client'

// 2. Next.js/React imports
import type {Metadata} from "next";
import {NextRequest, NextResponse} from 'next/server';

// 3. Internal imports with @ alias
import {getValidAuthentication} from "@/libs/security";
import {UserInfo} from "@/components/UserInfo";

// 4. Third-party libraries
import * as client from 'openid-client';
```

### Type vs Interface Usage

```typescript
// Interfaces for object shapes (preferred)
export interface UserInfoProps {
    name: string;
    email: string;
    picture?: string;
}

export interface AuthenticatedUser {
    name: string;
    email: string;
    isOwner: boolean;
}

// Types for unions, intersections, aliases
export type BackendSession = AuthenticatedSession | AnonymousSession
export type MediaId = string
export type AccessToken = AccessTokenClaims & { accessToken: string }
```

### Naming Conventions

```typescript
// Component props: ComponentNameProps
export interface UserInfoProps { /* ... */ }

// Event handlers: onAction prefix
interface DialogHandlers {
    onClose: () => void;
    onSubmit: () => Promise<void>;
    onNameChange: (name: string) => void;
}

// Boolean props: is/has prefix
interface ComponentProps {
    isLoading: boolean;
    hasError: boolean;
}

// Custom hooks: use prefix, camelCase
export function useCatalogContext() { /* ... */ }
```

### React Component Structure

```typescript
// Functional component with destructured props
export function UserInfo({ name, email, picture }: UserInfoProps) {
    return <div>{/* JSX */}</div>
}

// Async server component (Next.js)
export default async function Page({ searchParams }: PageProps) {
    const params = await searchParams;
    // ...
}

// Hooks usage
const theme = useTheme();
const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

const callback = useCallback((action: Action) => {
    dispatch(action)
}, [dispatch])
```

### Error Handling

```typescript
// Error mapping pattern
function getErrorInfo(error: string): ErrorInfo {
    const errorMap: Record<string, ErrorInfo> = {
        'invalid_request': { title: '...', description: '...' },
        'access_denied': { title: '...', description: '...' },
    };
    return errorMap[error] || { title: 'Error', description: '...' };
}

// Try-catch with null return
export function decodeJWT(token: string): Payload | null {
    try {
        // decode logic
        return decoded;
    } catch (error) {
        console.error('Failed to decode JWT:', error);
        return null;
    }
}

// Async error propagation
try {
    const result = await someOperation();
    return NextResponse.json(result);
} catch (e) {
    console.error('Operation failed:', e);
    return NextResponse.redirect('/error');
}
```

### Testing Patterns

```typescript
// Vitest environment directive
// @vitest-environment node

import {describe, it, expect, beforeAll, afterAll, vi} from 'vitest';

// Descriptive test names
describe('getValidAccessToken', () => {
    it('should return null when no tokens are provided', async () => {
        const result = await getValidAccessToken();
        expect(result).toBeNull();
    });
    
    it('should return valid token when access token is not expired', async () => {
        // Arrange
        const validToken = createToken();
        
        // Act
        const result = await getValidAccessToken();
        
        // Assert
        expect(result).not.toBeNull();
    });
});

// Test fakes
class ActionObserverFake {
    public actions: Action[] = []
    onAction = (action: Action) => this.actions.push(action)
}

// Mock patterns
vi.stubEnv('OAUTH_CLIENT_ID', 'test-client-id');
vi.mock('next/headers', () => ({
    cookies: vi.fn(() => fakeCookies),
}));
```

### File Naming Conventions

```
app/
  layout.tsx                # Root layout
  page.tsx                  # Route page
  (authenticated)/          # Route group (not in URL)
  auth/login/route.ts       # API route handler

components/
  UserInfo/index.tsx        # Component with folder

libs/
  security/
    access-token-service.ts      # kebab-case services
    access-token-service.test.ts # co-located tests

src/core/catalog/
  action-albumRenamed.ts         # action- prefix
  thunk-saveAlbumName.ts         # thunk- prefix
  selector-albumActions.ts       # selector- prefix
```

## Commit Message Convention

```
<domain>[/<area>] [+tags] - <message>

Examples:
catalog/web +minor - add album sorting by date
backup/cli - fix duplicate detection logic
archive/api +update-snapshots - improve thumbnail quality

Domains: catalog, archive, backup
Areas: web, cli, proj, api
Tags: +patch, +minor, +major, +pr, +update-snapshots
```

## Before Raising a Pull Request

1. **Follow coding standards strictly** - architecture and design patterns
2. **Simplify resulting code** - clean code, no paraphrasing comments
3. **Adhere to testing strategy** - maintain test coverage
4. **Ensure code builds** - run relevant build commands
5. **Verify tests pass** - run test suite for modified areas

## Common Gotchas

- **Go**: Always run `make setup-go` before testing (starts Docker/localstack)
- **Go**: Test `pkg/acl/jwks` fails without internet - this is expected
- **Node**: Always run `npm install` in the project directory first
- **web-nextjs**: Use `@/` path alias for imports
- **web**: Use relative paths or `src/` prefix for imports
- **Makefile**: Commands run from repo root, use `cd` for subdirectories
