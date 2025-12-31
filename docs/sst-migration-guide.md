# SST Migration Guide

## Overview

This guide documents the migration from AWS CDK to SST (Serverless Stack) for DPhoto infrastructure. The migration is being done in phases to ensure zero downtime and maintain data integrity.

## Current Status

**Step 1: Initial Setup** ✅ (CURRENT)
- SST v3 installed and configured
- Project structure established
- No resources deployed yet
- CDK remains active

**Step 2: Storage Layer** 🔜 (NEXT)
- Will recreate DynamoDB, S3, Cognito in SST config
- Configuration only, no deployment

**Step 3: Compute Layer** 🔜 (FUTURE)
- Will recreate Lambda, API Gateway, SNS/SQS in SST config
- Configuration only, no deployment

## Prerequisites

- Node.js >= 18.0.0
- AWS CLI configured
- Existing AWS CDK infrastructure running
- Access to AWS account

## Installation

### 1. Install Dependencies

From the repository root:

```bash
npm install
```

This installs SST v3 and its dependencies defined in `package.json`.

### 2. Verify Installation

```bash
npm run sst -- version
```

You should see SST version 3.3.44 or higher.

## SST Commands

### Development

```bash
# Start SST development mode (when ready to deploy)
npm run dev

# Install SST providers
npm run install:sst

# Deploy to AWS (only after Steps 2-3 are complete)
npm run deploy

# See what changes will be made
npm run diff

# Remove SST infrastructure (use with caution!)
npm run remove
```

### Validation

```bash
# Check SST version
npm run sst -- version

# See available SST commands
npm run sst
```

## Project Structure

```
dphoto/
├── sst.config.ts              # Main SST configuration
├── package.json               # SST dependencies and scripts
├── .gitignore                 # Updated with SST artifacts
│
├── deployments/
│   └── cdk/                   # Existing CDK (still active)
│       ├── lib/
│       ├── bin/
│       └── package.json
│
├── api/lambdas/               # Lambda function code (unchanged)
├── cmd/dphoto/                # CLI code (unchanged)
├── pkg/                       # Core business logic (unchanged)
├── web/                       # Web application (unchanged)
│
├── specs/
│   └── 2025-12-SST-integration.md  # Complete migration plan
│
└── docs/
    └── sst-migration-guide.md      # This file
```

## Configuration Files

### `sst.config.ts`

Main SST configuration file that defines infrastructure resources.

**Current State:** Minimal setup with no resources defined yet.

**Future State:** Will contain all DynamoDB, S3, Lambda, API Gateway definitions.

### `package.json`

Root-level package file for SST dependencies.

**Scripts:**
- `npm run sst` - Run any SST command
- `npm run dev` - Start live development
- `npm run deploy` - Deploy to AWS
- `npm run remove` - Remove infrastructure
- `npm run diff` - See what changes will be made
- `npm run install:sst` - Install SST providers

### `.gitignore`

Updated to exclude SST build artifacts:
- `.sst/` - SST internal files
- `.open-next/` - Next.js SST adapter files  
- `sst-env.d.ts` - Auto-generated TypeScript types
- `node_modules/` - Dependencies
- `package-lock.json` - Lock file
- `.env`, `.env.local` - Environment files

## Migration Phases

### Phase 1: Initial Setup (COMPLETED)

**Goal:** Set up SST alongside CDK without affecting production.

**What was done:**
1. ✅ Created `package.json` with SST v3 dependency
2. ✅ Created `sst.config.ts` with minimal configuration
3. ✅ Updated `.gitignore` for SST artifacts
4. ✅ Created this documentation
5. ✅ Created migration specification (`specs/2025-12-SST-integration.md`)

**What was NOT done:**
- ❌ No resources deployed
- ❌ No changes to existing CDK
- ❌ No changes to running services

**Validation:**
```bash
# Verify SST is installed
npm run sst -- version

# See available commands
npm run sst
```

### Phase 2: Storage Layer (NEXT STEP)

**Goal:** Define stateful resources in SST configuration (DynamoDB, S3, Cognito).

**Plan:**
1. Add DynamoDB table definition matching `DATA_MODEL.md`
2. Add S3 bucket definitions (archive + cache)
3. Add Cognito User Pool definition
4. Validate configuration with `sst build`
5. **DO NOT DEPLOY** - configuration only

**Validation:**
- `sst version` shows SST 3.x
- Configuration is valid TypeScript
- Documentation updated with resource mappings

### Phase 3: Compute Layer (FUTURE)

**Goal:** Define Lambda functions and API Gateway in SST configuration.

**Plan:**
1. Add API Gateway HTTP API definition
2. Add all Lambda function definitions
3. Add SNS/SQS definitions
4. Configure IAM permissions
5. **DO NOT DEPLOY** - configuration only

**Validation:**
- All Lambda builds succeed
- API routes match existing endpoints
- IAM permissions equivalent to CDK

### Phase 4-6: Deployment and Cleanup (FUTURE)

See `specs/2025-12-SST-integration.md` for detailed plans on:
- Deployment strategy and testing
- Production migration
- CDK cleanup

## Comparison: CDK vs SST

### CDK (Current)

**Location:** `deployments/cdk/`

**Commands:**
```bash
cd deployments/cdk
npm install
npm test
npm run synth
cdk deploy --context environment=next
```

**Pros:**
- ✅ Currently working and deployed
- ✅ Familiar AWS construct library
- ✅ Well-tested in production

**Cons:**
- ❌ Complex configuration
- ❌ Slow deployment cycles
- ❌ Limited local development
- ❌ Manual environment management

### SST (Target)

**Location:** Repository root

**Commands:**
```bash
npm install
npm run dev          # Live development
npm run deploy       # Production deployment
```

**Pros:**
- ✅ Simpler configuration
- ✅ Live Lambda Development
- ✅ Better TypeScript integration
- ✅ Faster deployment cycles
- ✅ Built-in environment management
- ✅ Better local testing

**Cons:**
- ❌ Newer framework (less battle-tested)
- ❌ Migration effort required
- ❌ Team learning curve

## Environment Variables

### CDK (Current)

Environments managed via CDK context:
```bash
cdk deploy --context environment=next
cdk deploy --context environment=live
```

### SST (Future)

Environments managed via SST stages:
```bash
sst deploy --stage next
sst deploy --stage live
```

SST automatically provides:
- `$app.name` - Application name
- `$app.stage` - Current stage
- `$dev` - True if in dev mode

## Resource Naming

### Current CDK Pattern

```
dphoto-{env}-{resource}
```

Examples:
- `dphoto-next-infrastructure`
- `dphoto-next-catalog-table`
- `dphoto-live-archive-bucket`

### Future SST Pattern

SST automatically namespaces with app and stage:
```
{app}-{stage}-{resource}
```

Examples:
- `dphoto-next-catalog-table`
- `dphoto-live-archive-bucket`

**Note:** Resource names will be compatible with current naming scheme.

## Development Workflow

### Current (CDK)

```bash
# Make infrastructure change
cd deployments/cdk
vim lib/some-construct.ts

# Test locally
npm test

# Deploy to dev
cdk deploy --context environment=next

# Wait 5-10 minutes for deployment

# Test the change
# (requires actual AWS deployment)
```

### Future (SST)

```bash
# Make infrastructure change
vim sst.config.ts

# Start live development
npm run dev

# SST automatically:
# - Deploys to AWS (first time)
# - Watches for changes
# - Hot-reloads Lambda functions
# - Provides instant feedback

# Test changes immediately
# (Live Lambda Development)
```

## Cost Implications

### Current Cost: ~$2-4/month

- DynamoDB (on-demand)
- S3 storage (~$0.50)
- Lambda invocations (minimal)
- Cognito users (~6 users)
- CloudFront (minimal)

### During Migration: ~$4-8/month (temporary)

- Dual infrastructure if deploying SST alongside CDK
- Expected duration: < 1 week
- Additional cost: < $10 total

### Post-Migration: ~$2-4/month

- Same as current
- SST doesn't add infrastructure costs
- Possible savings from better resource optimization

## Troubleshooting

### SST Installation Issues

**Problem:** `npm install` fails

**Solution:**
```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules if exists
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

### SST Version Issues

**Problem:** SST version mismatch

**Solution:**
```bash
# Check installed version
npm run sst -- version

# Update to latest
npm update sst
```

### Configuration Errors

**Problem:** SST configuration issues

**Solution:**
```bash
# Check SST version
npm run sst -- version

# Try SST diff to validate config
npm run diff

# Check for detailed error messages
npm run sst -- help
```

### AWS Credentials

**Problem:** SST can't access AWS

**Solution:**
```bash
# Verify AWS CLI is configured
aws sts get-caller-identity

# Ensure credentials are set
export AWS_PROFILE=your-profile
# or
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
```

## Testing Strategy

### Step 1 (Current)

**What to test:**
- ✅ SST CLI is installed
- ✅ Configuration file is valid
- ✅ TypeScript compiles
- ✅ Documentation is complete

**How to test:**
```bash
npm run sst -- version
# Check sst.config.ts has valid TypeScript syntax
```

### Step 2 (Next)

**What to test:**
- ✅ DynamoDB configuration matches DATA_MODEL.md
- ✅ S3 bucket settings match existing
- ✅ Cognito pool configuration correct
- ✅ SST configuration valid

**How to test:**
```bash
npm run diff
# Compare generated configuration with CDK output
```

### Step 3 (Future)

**What to test:**
- ✅ All Lambda functions build
- ✅ API Gateway routes correct
- ✅ IAM permissions equivalent
- ✅ Environment variables set

## Rollback Plan

### Current Status (Step 1)

**If issues arise:**
1. Delete SST files:
   ```bash
   rm -rf package.json sst.config.ts .sst/
   git checkout .gitignore
   ```
2. Continue using CDK as before
3. Zero impact on running services

### Future Steps (2-3)

**If configuration issues arise:**
1. Don't deploy SST
2. Continue using CDK
3. Fix SST configuration offline
4. Zero impact on running services

### Post-Deployment (Steps 4-6)

See `specs/2025-12-SST-integration.md` for detailed rollback procedures.

## Next Steps

1. ✅ **Step 1 Complete:** Initial setup done
2. ⏭️ **Step 2 Next:** Define storage resources in `sst.config.ts`
   - Review `DATA_MODEL.md`
   - Review existing CDK DynamoDB configuration
   - Add DynamoDB table definition to SST
   - Add S3 bucket definitions
   - Add Cognito User Pool definition
   - Validate with `sst build`
3. ⏭️ **Step 3 Future:** Define compute resources
4. ⏭️ **Step 4-6 Future:** Deploy and migrate

## References

- [SST Documentation](https://docs.sst.dev/)
- [SST v3 Guide](https://docs.sst.dev/docs)
- [AWS CDK to SST Migration](https://docs.sst.dev/migrating/cdk)
- [SST Resource Types](https://docs.sst.dev/docs/component/)
- [DPhoto Migration Spec](../specs/2025-12-SST-integration.md)
- [DPhoto Data Model](../DATA_MODEL.md)

## Support

For questions or issues:
1. Check this documentation
2. Review `specs/2025-12-SST-integration.md`
3. Consult [SST Discord](https://discord.gg/sst)
4. Check [SST GitHub Issues](https://github.com/sst/sst/issues)

## Version History

- **2025-12-31:** Initial setup (Step 1 completed)
  - SST v3.3.44 installed
  - Basic configuration created
  - Documentation completed
  - No resources deployed
