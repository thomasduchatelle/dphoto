# Waku Migration Decision Log

## Migration Topics to Discuss

- Build System & Configuration
- Routing Architecture  
- Server-Side Rendering (SSR)
- React Version Compatibility
- State Management & Data Fetching
- Styling & CSS
- Development Experience
- Testing & Quality Assurance
- Compatibility of libraries with the new version of React
- Visual testing with Storybook or an alternative

## Decisions Made

### Target Runtime Architecture

**Decision**: Use API Gateway + Lambda architecture without CloudFront for Waku SSR deployment

**Architecture**:

User → API Gateway → Lambda (Waku SSR)
                   ↘ S3 (static assets)

**Alternative options considered**:
- Lambda@Edge + CloudFront (higher per-request costs)
- Static Generation + Incremental rendering (not suitable for user-specific content)
- API Gateway + Lambda + CloudFront (adds unnecessary caching layer and costs)

**Rationale**:
- **Cost optimization**: Eliminates CloudFront minimum monthly costs (~$0.60/month)
- **Consistency with existing system**: Matches the current Golang API architecture using API Gateway + Lambda
- **Poor cacheability**: Application content is user-specific, making CloudFront caching benefits minimal
- **Simplified deployment**: Single deployment pattern for all services in the CDK stack

### Authentication and SSR Strategy

**Decision**: Use client components with JWT authentication for initial migration

**Approach**:
- All components will use `'use client'` directive to maintain current behavior
- Continue using JWT tokens sent in each API request
- Migrate from React Router to file-based routing (mandatory for Waku)

**Rationale**:
- **Minimize migration risk**: Keeps existing authentication flow unchanged
- **Faster migration**: Avoids complex SSR authentication patterns during initial cutover
- **Gradual optimization**: Allows post-migration improvements without blocking delivery

### Routing Migration Strategy

**Decision**: Use gradual manual migration approach

**Approach**:
- Reorganize components to match Waku's expected file structure before migration
- Manually update React Router route definitions to reference new file locations
- Validate routing works correctly with current system before switching to Waku

**Rationale**:
- **Simpler implementation**: No need to build automated file-scanning helpers
- **Manual control**: Each route change can be tested individually
- **Reduced migration surface area**: File structure changes are completed before Waku switch
- **Early validation**: Routing conflicts identified while still using React Router

### React Version Migration Strategy

**Decision**: Stay on React 18 for CRA to Waku migration, upgrade to React 19 post-migration

**Approach**:
- Phase 2: CRA + React 18 → Waku + React 18 (build system migration only)
- Phase 3: Waku + React 18 → Waku + React 19 (React version upgrade only)

**Rationale**:
- **CRA incompatibility**: React Scripts 5.0.1 does not support React 19
- **Risk reduction**: Isolate build system changes from React version changes
- **Dependency compatibility**: Current libraries are tested with React 18
- **Sequential validation**: Test Waku stability before introducing React 19 breaking changes

### State Management and Data Fetching Strategy

**Decision**: Continue using React state patterns, migrate to server-side data fetching with SSR adoption

**Approach**:
- Keep current React state management (useState, useEffect, useContext) during initial migration
- Maintain client-side data fetching patterns with JWT authentication
- Migrate to server component data fetching patterns when adopting SSR optimizations

**Rationale**:
- **Consistency with client-first approach**: Aligns with keeping all components client-side initially
- **Minimize migration complexity**: Avoids changing state management and build system simultaneously
- **Natural progression**: Server-side data fetching becomes relevant when adopting SSR patterns

### Styling and CSS Strategy

**Decision**: Continue using Material-UI with Emotion for CSS-in-JS

**Approach**:
- Keep current MUI v5 components and theme system
- Maintain Emotion for styled components
- Continue using DPhotoTheme wrapper component
- No changes to styling approach during initial migration

**Rationale**:
- **Zero migration effort**: Current styling stack is compatible with Waku
- **Proven compatibility**: MUI v5 + Emotion work with SSR when properly configured
- **Risk reduction**: Avoid changing styling and build system simultaneously
- **Future optimization**: SSR styling improvements can be addressed post-migration

### Unit Testing Strategy

**Decision**: Continue using Jest with React Testing Library, migrate to Vitest post-migration

**Approach**:
- Keep Jest as test runner during CRA to Waku migration
- Maintain current React Testing Library patterns and setup
- Configure Jest explicitly for Waku instead of via react-scripts
- Migrate to Vitest after Waku migration is stable

**Rationale**:
- **CRA compatibility**: Vitest cannot be easily used with Create React App
- **Risk reduction**: Avoid changing testing framework and build system simultaneously
- **Proven patterns**: Current Jest + RTL setup works well with Waku
- **Performance optimization**: Vitest migration can provide faster test execution post-migration

## Migration Plan

### Phase 1: Anticipation

**Goal**: Prepare for migration with minimal risk to current system

**Steps**:
1. **Launch empty Waku project** alongside current application under `/waku` path
2. **Set up parallel build pipeline** for Waku in CDK without affecting current deployment
3. **Reorganize components to file-based structure** and update React Router route definitions to match Waku expectations

### Phase 2: Swap and Stabilise

**Goal**: Migrate from CRA to Waku while maintaining React 18

**Steps**:
- **Switch build system from CRA to Waku** while keeping React 18
- **Configure Jest explicitly** for Waku environment instead of via react-scripts
- **Test routing functionality** and fix any issues discovered during file structure migration
- **Validate all functionality** works with Waku + React 18

### Phase 3: Completion and Cleanup

**Goal**: Complete migration and optimize system

**Steps**:
- **Upgrade to React 19** and test compatibility
- **Migrate from Jest to Vitest** for improved test performance
- **Review authentication to use HTTP+cookies** instead of JWT in requests
- **Optimize components to use SSR** by removing unnecessary `'use client'` directives
- **Migrate to server-side data fetching patterns** where appropriate for SSR components
- **Optimize styling performance for SSR** by migrating to CSS modules (CSS extraction and reducing hydration mismatches) - low priority
