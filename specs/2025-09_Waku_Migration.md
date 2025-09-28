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

## Migration Plan

### Phase 1: Anticipation

**Goal**: Prepare for migration with minimal risk to current system

**Steps**:
1. **Launch empty Waku project** alongside current application under `/waku` path
2. **Audit current dependencies** for React 18+ compatibility and identify breaking changes
3. **Set up parallel build pipeline** for Waku in CDK without affecting current deployment
4. **Create dependency compatibility matrix** documenting upgrade paths for key libraries
5. **Establish testing baseline** by running current test suite and documenting any React 18 compatibility issues

### Phase 2: Swap and Stabilise

**Goal**: Gradually migrate functionality while maintaining system stability

**Steps**:
- TBD - To be planned after Phase 1 completion

### Phase 3: Completion and Cleanup

**Goal**: Complete migration and remove legacy systems

**Steps**:
- TBD - To be planned after Phase 2 completion
