# DPhoto - Copilot Instructions

## Repository Overview

Photo backup and sharing application deployed as serverless AWS app. CLI for backup, React web app for viewing/sharing.

**Stack**: Go 1.23, TypeScript/React 19, Node 20, AWS Lambda, S3, DynamoDB  
**Architecture**: Hexagonal (domain-driven design)

## Project Structure

**Key Directories:**
- `pkg/` - Core Go domain logic (acl, archive, backup, catalog domains with *adapters/ for DynamoDB/S3)
- `cmd/dphoto/` - CLI source
- `api/lambdas/` - 20+ Lambda handlers
- `web/` - React 19 + Waku framework (Vite, Vitest, Playwright; npm only, NOT yarn)
- `deployments/cdk/` - AWS CDK infrastructure (TypeScript/Jest)
- `internal/` - Mocks, utilities
- `Makefile` - **ALWAYS use for builds/tests**

## Build & Test Commands

### CRITICAL: LocalStack Required for Go Tests
```bash
docker compose up -d  # ALWAYS run before Go tests, wait 5-10s
```
Provides S3, DynamoDB, SNS, SQS, SSM, ACM on ports 4563-4599.

### Go Commands
```bash
make build-go          # Build CLI (~6s) -> ./dphoto
make build-api         # Build Lambdas (~50s) -> bin/*.zip
make build-cli         # Cross-compile for Linux/Darwin (~70s) -> bin-cli/dphoto-*
make test-pkg          # Test pkg/ (~60s) - expects jwks test failure (internet required)
make test-api          # Test lambdas (~13s)
make test-go           # Both test-pkg + test-api
```
**Note**: Use `AWS_PROFILE=""` when running `go test` directly.

### Web Commands
```bash
cd web
npm ci                 # Setup (~36s)
npx playwright install # One-time browser install
npm run dev            # Dev server on :3000 (proxies to :8080)
npm run test:unit      # Vitest (~17s)
npm run test:visual    # Playwright visual regression
npm run build:lambda   # Production build for Lambda
npm run ladle          # Component viewer on :61000

# OR use Makefile
make setup-web         # npm ci + playwright install
make test-web-ci       # CI mode with chromium only
make build-web         # Calls build:lambda
make update-snapshots  # Update visual regression snapshots
```

### CDK Commands
```bash
cd deployments/cdk
npm install            # Setup (~13s)
npm test               # Jest (~23s)
cdk deploy --context environment=next --all

# OR
make setup-cdk test-cdk deploy-cdk
```

### Master Commands
```bash
make all               # clean + test + build (all components)
make setup             # setup-cdk + setup-web + setup-go
make clean             # Clean artifacts (web-waku error expected)
make dcup              # Start LocalStack + create bucket/topic
make dcdown            # Full LocalStack reset
```

## GitHub Actions CI/CD

**Commit Convention**: `<domain>[/<area>] [+tags] - <message>`
- Domains: catalog, archive, backup, acl, api, web, proj, infra, ui, test, docs
- Tags: `+update-snapshots`, `+patch/minor/major`, `+pr`, `+next`

**Main Workflow** (push to main): test-go → test-ts → test-cdk → build → deploy-next → deploy-live → release
**Feature Workflow**: Creates PR if `+pr` in commit, runs tests/builds, comments CDK diff, deploys to "next" if `+next`

**Key Jobs**:
- `job-test-go.yml` - 10min timeout, starts localstack, runs `make test-go`
- `job-test-ts.yml` - Node 20, caches node_modules/playwright, runs `make test-web-ci`, can update snapshots
- `job-test-cdk.yml` - Runs `make test-cdk`
- `job-build-go.yml` - Cross-compiles CLI + Lambda zips, uploads artifacts
- `job-build-ts.yml` - Runs `make build-web`, uploads for deployment

## Known Issues & Workarounds

**Go**: 
- `pkg/acl/jwks` test fails without internet (accounts.google.com) - EXPECTED, ignore it
- ALWAYS run `docker compose up -d` before tests, wait 5-10s
- Use `AWS_PROFILE=""` for direct `go test` commands

**Web**: 
- Use `npm` NOT `yarn`
- Playwright snapshots: `-local/` (local) vs `-snapshots/` (CI)
- Ports: 3000 (dev), 8080 (wiremock), 61000 (ladle)

**Makefile**: 
- `make clean` errors on non-existent `web-waku/` - EXPECTED, ignore
- Use `make clean-web clean-api` for targeted clean

**CDK**: Always use `--context environment=<env>` for deploys

## Architecture

**Hexagonal**: Domain in `pkg/`, adapters in `pkg/*adapters/`, ports = interfaces  
**DynamoDB**: Single table (see DATA_MODEL.md), keys like `{OWNER}#ALBUM`, `USER#{EMAIL}`, 4 GSIs  
**Frontend**: React 19 + Waku, Redux-like reducers (`web/src/core/`), React Router 7, Axios, OAuth2/Cognito

## Validation Checklist

**Go**: `docker compose up -d` → `make test-go` (jwks fail OK) → `make build-go` (~6s) → `make build-api` (~50s)  
**Web**: `cd web && npm ci` → `npm run test:unit` (~17s) → `npm run test:visual` (if UI changed) → `npm run build`  
**CDK**: `cd deployments/cdk && npm install` → `npm test` (~23s)  
**Full**: `make clean && make all`

## Quick Reference

**Ports**: 3000 (dev), 4563-4599 (LocalStack), 8055 (LS dash), 8080 (wiremock), 61000 (ladle)  
**Env Vars**: `AWS_PROFILE=""`, `AWS_*=localstack`, `CI=true`  
**Tests**: `*_test.go` (Go), `web/src/**/*.test.{ts,tsx}` (unit), `web/playwright/*.spec.ts` (visual), `deployments/cdk/lib/**/*.test.ts` (CDK)  
**Ignored**: `bin/`, `bin-cli/`, `.build/`, `dist/`, `dphoto`, `.aider*`

---

**Trust these instructions** - validated against repository. Search only if incomplete/incorrect.
