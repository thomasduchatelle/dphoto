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

## Migration Plan

### Phase 1: Anticipation

**Goal**: Prepare for migration with minimal risk to current system

**Steps**:
1. **Launch empty Waku project** alongside current application under `/waku` path
2. **Set up parallel build pipeline** for Waku in CDK without affecting current deployment
3. **Reorganize components to file-based structure** and update React Router route definitions to match Waku expectations

### Phase 2: Swap and Stabilise

**Goal**: Gradually migrate functionality while maintaining system stability

**Steps**:
- **Test routing functionality** and fix any issues discovered during file structure migration

### Phase 3: Completion and Cleanup

**Goal**: Complete migration and remove legacy systems

**Steps**:
- **Review authentication to use HTTP+cookies** instead of JWT in requests
- **Optimize components to use SSR** by removing unnecessary `'use client'` directives
