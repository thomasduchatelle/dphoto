# DPhoto - Copilot Instructions

## Repository Overview

DPhoto is a photo backup and sharing application deployed as a serverless application on AWS. It provides a CLI tool for backing up photos/videos and a web application for viewing and sharing albums.

**Repository Size**: Medium (~25K lines of code across Go, TypeScript/React)  
**Primary Languages**: Go 1.23, TypeScript/React 19, Node.js 20  
**Target Runtime**: AWS Lambda (Go & Node), S3, DynamoDB  
**Architecture**: Hexagonal Architecture with domain-driven design

## Project Structure

### Root Directory Files
- `Makefile` - Main build orchestration (ALWAYS use this for builds/tests)
- `go.mod`, `go.sum` - Go module definitions (Go 1.23)
- `docker-compose.yml` - LocalStack for local development/testing
- `DATA_MODEL.md` - DynamoDB single-table schema documentation
- `TODO.md` - Project roadmap and pending features
- `README.md` - Setup and contribution guidelines

### Key Directories

**Backend (Go)**
- `pkg/` - Core domain logic and business rules (hexagonal architecture)
  - Main domains: `acl/`, `archive/`, `backup/`, `catalog/`, `catalogviews/`, `dns/`
  - Adapters: `*adapters/` directories (DynamoDB, S3, SNS/SQS implementations)
- `cmd/dphoto/` - CLI application source code
- `api/lambdas/` - AWS Lambda functions for REST API (20+ lambda handlers)
- `internal/` - Private packages (mocks, utilities)

**Frontend (React/Waku)**
- `web/` - React 19 web application using Waku framework
  - Uses Vite, Vitest for unit tests, Playwright for visual regression
  - Runs on port 3000 in development
  - Configured to proxy `/api` and `/oauth` to port 8080 (wiremock)
  - Uses `npm` (NOT yarn despite some legacy references)

**Infrastructure**
- `deployments/cdk/` - AWS CDK project (TypeScript) for infrastructure as code
  - Defines stacks for data stores, APIs, and web hosting
  - Tests via Jest

**Testing/Tools**
- `test/` - Integration test resources
- `tools/` - Utilities (localstack-init, dphotoops, infra-bootstrap)
- `scripts/` - Helper scripts (wiremock, playwright report viewer)

### Configuration Files
- `web/tsconfig.json` - TypeScript configuration
- `web/vitest.config.ts` - Vitest unit test configuration
- `web/playwright.config.ts` - Playwright visual regression test configuration
- `web/waku.config.ts` - Waku framework configuration (port 3000, proxies)
- `deployments/cdk/cdk.json` - CDK configuration
- `deployments/cdk/jest.config.js` - CDK tests configuration

## Build & Test Commands

### CRITICAL: Always Start LocalStack First

**Before running ANY Go tests**, ALWAYS start LocalStack:
```bash
docker compose up -d
# Wait 5-10 seconds for services to be ready
```

LocalStack provides S3, DynamoDB, SNS, SQS, SSM, and ACM services on ports 4563-4599.

### Go (Backend & CLI)

**Build CLI:**
```bash
make build-go          # Builds dphoto CLI (~6 seconds)
# Output: ./dphoto executable
```

**Build API Lambdas:**
```bash
make build-api         # Builds all Lambda functions (~50 seconds)
# Output: bin/*.zip files
```

**Build CLI for Distribution:**
```bash
make build-cli         # Cross-compiles for Linux (amd64), Darwin (amd64, arm64)
# Output: bin-cli/dphoto-*
```

**Test Go Code:**
```bash
# ALWAYS start localstack first!
docker compose up -d

make test-pkg          # Tests pkg/ packages (~60 seconds)
make test-api          # Tests API lambdas (~13 seconds)
make test-go           # Runs both test-pkg and test-api
```

**Known Test Issue**: `pkg/acl/jwks` tests fail without internet access (tries to reach accounts.google.com). This is expected in sandboxed environments - NOT your responsibility to fix.

**Important Go Commands:**
- ALWAYS use `AWS_PROFILE=""` when running go test directly (to prevent AWS credential issues)
- Use `-race` flag for race condition detection
- Tests use DynamoDB and S3 in LocalStack

### TypeScript/React (Web)

**Setup (First Time):**
```bash
cd web
npm ci                          # ~36 seconds, installs all dependencies
npx playwright install          # Installs Playwright browsers
```

**Development Server:**
```bash
cd web
npm run dev                     # Starts Waku dev server on port 3000
# OR
make start                      # Also starts wiremock on port 8080
```

**Run Tests:**
```bash
cd web
npm run test:unit              # Vitest unit tests (~17 seconds)
npm run test:visual            # Playwright visual regression tests
npm test                       # Runs both unit and visual tests
```

**CI Test Mode:**
```bash
make test-web-ci               # Uses chromium only, proper reporting
```

**Build for Production:**
```bash
cd web
npm run build                  # Standard build
npm run build:lambda           # Build for AWS Lambda deployment
# OR
make build-web                 # Uses build:lambda
```

**Update Visual Regression Snapshots:**
```bash
make update-snapshots          # Updates local snapshots
make update-snapshots-ci       # Updates CI snapshots (use only on CI)
```

**Ladle (Component Development):**
```bash
cd web
npm run ladle                  # Starts component viewer on port 61000
```

### CDK (Infrastructure)

**Setup:**
```bash
cd deployments/cdk
npm install                    # ~13 seconds
```

**Test:**
```bash
cd deployments/cdk
npm test                       # Jest tests (~23 seconds)
# OR
make test-cdk
```

**Deploy:**
```bash
cd deployments/cdk
cdk deploy --context environment=next --all
# OR
make deploy-cdk
```

### Master Commands

**Complete Build & Test:**
```bash
make all                       # Runs: clean, test, build (all components)
# Components: test-go, test-web, test-cdk, build-go, build-web
```

**Setup Everything:**
```bash
make setup                     # Runs: setup-cdk, setup-web, setup-go
```

**Clean:**
```bash
make clean                     # Cleans all build artifacts and test cache
```

### LocalStack Utilities

**Initialize LocalStack:**
```bash
make dcup                      # Starts LocalStack and creates S3 bucket + SNS topic
```

**Clear LocalStack Data:**
```bash
make clearlocal                # Removes all S3 objects and DynamoDB tables
```

**Full Reset:**
```bash
make dcdown                    # Stops containers, removes volumes, clears data
```

## GitHub Actions CI/CD

### Main Workflows

**workflow-main.yml** (on push to `main`):
1. test-go, test-ts, test-cdk, build-go, build-ts
2. deploy-next (to "next" environment)
3. deploy-live (to production)
4. Creates GitHub release with CLI binaries

**workflow-feature-branch.yml** (on push to feature branches):
1. Creates/updates Pull Request (if `+pr` or `+next` in commit message)
2. Runs test-go, test-ts, test-cdk
3. Builds: build-go (snapshot=true), build-ts
4. Comments CDK diff on PR
5. Optionally deploys to "next" (if `+next` in commit)

**workflow-comment-update-snapshots.yml**:
- Triggered by PR comment `/update-snapshots`
- Updates visual regression snapshots

### Commit Message Convention

Format: `<domain>[/<area>] [+tags] - <message>`

**Domains**: catalog, archive, backup, acl, api, web, proj, infra, ui, test, docs  
**Tags**:
- `+update-snapshots` - Updates visual regression snapshots
- `+patch`, `+minor`, `+major` - Controls version bump
- `+pr` - Creates Pull Request automatically
- `+next` - Deploys to "next" environment

### Test Jobs Details

**job-test-go.yml**: 
- Timeout: 10 minutes
- Starts localstack via `docker compose up -d`
- Runs `make test-go`

**job-test-ts.yml**:
- Installs Node 20, caches node_modules and playwright
- Runs `make test-web-ci`
- Can update snapshots if `+update-snapshots` in commit message
- Uploads Playwright reports on failure

**job-test-cdk.yml**:
- Runs `make test-cdk` (Jest tests)

**job-build-go.yml**:
- Builds API lambdas and CLI for 3 platforms
- Uploads artifacts: `dist-api/` and `bin-cli/`

**job-build-ts.yml**:
- Builds web app with `make build-web`
- Uploads artifact for deployment

## Known Issues & Workarounds

### Go Tests
1. **Network-dependent test failure**: `pkg/acl/jwks/jwks_test.go` fails without internet (tries to reach accounts.google.com). This is expected - ignore this specific test failure.
2. **LocalStack requirement**: ALWAYS run `docker compose up -d` before Go tests. Wait 5-10 seconds for services to initialize.
3. **AWS_PROFILE**: Set `AWS_PROFILE=""` when running `go test` directly to avoid AWS credential issues.

### Web/TypeScript
1. **Package manager**: Use `npm` NOT `yarn` (legacy references in some files are outdated)
2. **Playwright snapshots**: Separate snapshot directories for local vs CI (`-local` suffix vs `-snapshots`)
3. **Node version**: Use Node 20.x (specified in GitHub Actions)
4. **Port conflicts**: Dev server uses port 3000, ladle uses 61000, wiremock uses 8080

### CDK
1. **Context required**: Always specify `--context environment=<env>` for deploys
2. **Cache suffix**: Used in CI to avoid cache locking issues

### Docker Compose
1. **Volume persistence**: LocalStack data persists in `.build/localstack/`
2. **Wiremock profile**: Use `docker-compose --profile bg up -d` to include wiremock

### Makefile
1. **web-waku references**: The Makefile contains references to `web-waku` directory (in `clean-waku`, `setup-waku`, `test-waku`, `build-waku` targets), but this directory does not exist. The Waku migration is integrated into the `web/` directory. If `make clean` fails due to web-waku, this is expected - the web part still cleans successfully.
2. **Workaround**: Use specific targets like `make clean-web clean-api` or ignore the web-waku error

## Architecture Notes

### Hexagonal Architecture
- **Domain Logic**: In `pkg/` directories (e.g., `pkg/catalog`, `pkg/backup`)
- **Adapters**: In `pkg/*adapters/` directories (e.g., `pkg/catalogadapters/catalogdynamo`)
- **Ports**: Interfaces defined in domain packages

### Data Model
- **Single DynamoDB Table**: All data in one table with composite keys
- **Schema**: Documented in `DATA_MODEL.md`
- **Key Patterns**: `{OWNER}#ALBUM`, `USER#{EMAIL}`, etc.
- **GSIs**: AlbumIndex, ReverseLocationIndex, ReverseGrantIndex, RefreshTokenExpiration

### Frontend Architecture
- **React 19** with **Waku** framework (progressive migration)
- **State Management**: Redux-like reducers (see `web/src/core/`)
- **Routing**: React Router DOM 7
- **API Integration**: Axios with proxy to backend
- **Authentication**: OAuth2 with Google (managed via AWS Cognito)

## Validation Checklist

Before finalizing changes:

1. **Go Changes**:
   - [ ] Start LocalStack: `docker compose up -d`
   - [ ] Run `make test-go` (expect jwks test failure)
   - [ ] Run `make build-go` (should succeed in ~6s)
   - [ ] Run `make build-api` (should succeed in ~50s)

2. **TypeScript Changes**:
   - [ ] Run `cd web && npm ci` (if dependencies changed)
   - [ ] Run `cd web && npm run test:unit` (~17s)
   - [ ] Run `cd web && npm run test:visual` (if UI changed)
   - [ ] Run `cd web && npm run build`

3. **CDK Changes**:
   - [ ] Run `cd deployments/cdk && npm install`
   - [ ] Run `cd deployments/cdk && npm test`

4. **Full Integration**:
   - [ ] Run `make clean && make all` (complete build & test cycle)

## Quick Reference

### File Locations
- Go tests: `*_test.go` files throughout `pkg/`, `cmd/`, `api/`
- Web unit tests: `web/src/**/*.test.ts`, `web/src/**/*.test.tsx`
- Web visual tests: `web/playwright/visual-regression.spec.ts`
- CDK tests: `deployments/cdk/lib/**/*.test.ts`
- Mocks: `internal/mocks/` (generated via mockery)

### Port Assignments
- 3000: Waku dev server (web app)
- 4563-4599: LocalStack services
- 8055: LocalStack dashboard
- 8080: Wiremock (mock API for frontend dev)
- 61000: Ladle component viewer

### Environment Variables
- `AWS_PROFILE=""` - Disable AWS profile for tests
- `AWS_ACCESS_KEY_ID=localstack` - For LocalStack
- `AWS_SECRET_ACCESS_KEY=localstack` - For LocalStack
- `CI=true` - Enables CI mode for various tools

### Git Ignore Patterns
- `bin/`, `bin-cli/` - Build outputs
- `.build/` - LocalStack data
- `dist/` - Web build output
- `dphoto`, `dphoto-ops` - Built executables
- `.aider*` (except `.aider.conf.yml`) - AI assistant files

---

**Trust These Instructions**: This file has been validated against the actual repository. Only perform additional searches if information is incomplete or found to be incorrect.
