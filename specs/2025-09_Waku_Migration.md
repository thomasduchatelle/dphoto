# Waku Migration Decision Log

## Context

This document records decisions for migrating a React application from Create React App (CRA) to Waku, a modern React framework with SSR capabilities. The current application uses React 18, Material-UI v5, Emotion for styling, React Router for routing, and JWT authentication. The application runs on AWS using CDK for infrastructure.

## Migration Topics Discussed

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

### Visual Testing Strategy

**Decision**: Replace Storybook with Ladle + Playwright, fallback to Storybook 8 + Chromatic if needed

**Approach**:
- Prototype Ladle + Playwright visual testing solution during Phase 1
- Complete migration from Storybook 6.5.16 before introducing Waku
- Fallback to Storybook 8 + Chromatic (free tier) if Ladle prototype doesn't meet requirements
- Maintain same visual regression workflow: before/after/diff images on failures

**Rationale**:
- **Cost optimization**: Ladle + Playwright is completely free vs Storybook licensing concerns
- **Modern tooling**: Replace deprecated addon-storyshots with actively maintained solutions
- **Early validation**: Test visual tooling compatibility before build system migration
- **Risk mitigation**: Validate solution works with current React 18 + MUI setup before Waku

## Migration Plan

### Phase 1: Anticipation

**Goal**: Prepare for migration with minimal risk to current system

**Steps**:
1. **Launch empty Waku project** alongside current application under `/waku` path
2. **Set up parallel build pipeline** for Waku in CDK without affecting current deployment
3. **Prototype Ladle + Playwright visual testing** with current React 18 + CRA setup
4. **Complete visual testing migration** from Storybook to chosen solution (Ladle + Playwright or Storybook 8 + Chromatic)
5. **Reorganize components to file-based structure** and update React Router route definitions to match Waku expectations

### Phase 2: Swap and Stabilise

**Goal**: Migrate from CRA to Waku while maintaining React 18

**Steps**:
- **Switch build system from CRA to Waku** while keeping React 18
- **Configure Jest explicitly** for Waku environment instead of via react-scripts
- **Test routing functionality** and fix any issues discovered during file structure migration
- **Validate all functionality** works with Waku + React 18
- **Verify visual testing integration** works with Waku build system

### Phase 3: Completion and Cleanup

**Goal**: Complete migration and optimize system

**Steps**:
- **Upgrade to React 19** and test compatibility
- **Migrate from Jest to Vitest** for improved test performance
- **Review authentication to use HTTP+cookies** instead of JWT in requests
- **Optimize components to use SSR** by removing unnecessary `'use client'` directives
- **Migrate to server-side data fetching patterns** where appropriate for SSR components
- **Optimize styling performance for SSR** by improving CSS extraction and reducing hydration mismatches (low priority)

## Phase 1 Implementation Prompts

### Task 1: Launch Empty Waku Project

**Context**: A minimal WAKU project has been created in `/web-waku`. You need to have it tested, build and deployed alongside the current root application as part of the build pipelines.

* `/`: goes to the current WEB application hosted in S3 (unchanged)
* `/waku`: goes to the new WAKU project (new)

**Scope**:
* update the project `deployments/cdk` to deploy the new WEB application alongside the existing one
- create a new job `.github/workflows/job-test-waku.yml` to run the waku tests (similar example: `.github/workflows/job-test-ts.yml`)
- create a new job `.github/workflows/job-build-waku.yml` to build and publish the artefact (similar example: `.github/workflows/job-build-ts.yml`)
- update `.github/workflows/workflow-feature-branch.yml` and `.github/workflows/workflow-main.yml` to integrate both test abd build jobs
- update `.github/workflows/job-deploy.yml` to deploy the WAKU app (using jobs updated previously)
- Do NOT attempt to migrate any existing components or logic yet

**Requirements**:
- Ensure the Waku dev server starts without conflicts with the existing CRA dev server
- Document the setup process and any configuration needed

### Task 2: Set Up Parallel Build Pipeline in CDK

**Context**: Add Waku build and deployment configuration to the existing CDK infrastructure without affecting the current CRA deployment. This allows testing the deployment process before the actual migration.

**Scope**:
- Extend existing CDK stack to include Waku build and deployment alongside current CRA setup
- Configure API Gateway + Lambda architecture for Waku SSR as decided
- Set up S3 bucket for Waku static assets separate from current assets
- Use separate paths/domains to avoid conflicts (e.g., `/waku` prefix)
- Do NOT modify or affect the existing CRA deployment pipeline

**Requirements**:
- Deploy to same AWS account but isolated resources where possible
- Follow the decided architecture: API Gateway + Lambda (SSR) + S3 (static assets)
- No CloudFront as per the architecture decision
- Maintain existing CRA deployment unchanged
- Test deployment pipeline with the basic Waku app from Task 1

**Deliverables**:
- Updated CDK code with parallel Waku deployment
- Successful deployment of basic Waku app to AWS
- Separate build pipeline that doesn't interfere with existing processes
- Documentation of the deployment process and architecture

### Task 3: Prototype Ladle + Playwright Visual Testing

**Context**: Test Ladle + Playwright as a replacement for the current Storybook 6.5.16 + addon-storyshots setup. The current workflow uses Storybook for component development and addon-storyshots for visual regression testing with before/after/diff images on failures.

**Scope**:
- Set up Ladle for component development with existing MUI components
- Configure Playwright for visual regression testing
- Test with 2-3 existing components from the current application
- Replicate the current workflow: develop components, take screenshots, detect changes, generate diff images
- Compare development experience with current Storybook workflow
- Do NOT migrate the entire component library yet

**Requirements**:
- Use current React 18 + CRA setup for testing
- Test with existing Material-UI components and DPhotoTheme
- Verify screenshot consistency and diff generation works reliably
- Document any limitations or issues compared to current Storybook setup
- Test build integration and CI/CD compatibility

**Deliverables**:
- Working Ladle setup with sample components
- Playwright visual testing configuration
- Test results comparing Ladle vs current Storybook workflow
- Recommendation report on whether to proceed with Ladle or fallback to Storybook 8 + Chromatic
- Setup documentation and configuration files

### Task 4: Complete Visual Testing Migration

**Context**: Based on the prototype results from Task 3, complete the migration from Storybook 6.5.16 to the chosen solution (Ladle + Playwright or Storybook 8 + Chromatic). This must be completed before introducing Waku to reduce the number of simultaneous changes.

**Scope**:
- Migrate all existing Storybook stories to the chosen solution
- Replace addon-storyshots with the new visual testing setup
- Update CI/CD pipeline to use the new visual testing approach
- Ensure all team members can use the new development workflow
- Remove old Storybook 6.5.16 dependencies and configuration

**Requirements**:
- Maintain exact same visual regression testing workflow as before
- All existing component stories must work in the new setup
- CI/CD integration must work with build failures on visual changes
- Team training/documentation for the new workflow
- Verify all visual tests pass with the new setup

**Deliverables**:
- Complete migration to chosen visual testing solution
- All component stories migrated and working
- Updated CI/CD pipeline with new visual testing
- Team documentation and training materials
- Removal of old Storybook dependencies

### Task 5: Reorganize Components to File-Based Structure

You're going to implement the changes required to achieve the last step of the anticipation phase of `specs/2025-09_Waku_Migration.md`.

I'd need you to do the following:

1. refactor the code to create the **file-based routing structure compatible with Waku routing**
   * each index file must be of the type `export default async function NameOfThePage({pathArg1, ...})`
   * each `NameOfThePage` must be the full page: fetching data, layout, ... They will be client components, keep using `useEffect` as it is currently done. This code is certainly spread through several files at the moment: you need to regroup it into 1 page-component per page.
   * the page `src/pages/index.tsx` is expected to look like: `export async function GET(request: Request) { return Response.redirect('/albums/album-1') }`. Make a signature that can be used in a CRA context but is functionally similar.
   * extract the layout of the pages under `/albums/...` into `/albums/_layout.tsx` 
   * remove code becoming duplicated or dead because of the refactoring 
2. create a file `src/pages/_cra-router.tsx` which is using the React router, and the page components created on the previous step.
   * the resulting should be an extremely simple routing: `<Routes><Route path='/albums' element={<NameOfThePage />}/>...</Routes>`
3. move all the others UI components in `src/components/` in `src/pages/` ; any component which are not a page, a layout, or the cra router.

The refactoring should result into the same URLs navigation. The code must compile. And make sure you don't leave duplicated code or dead code behind your refactoring.
