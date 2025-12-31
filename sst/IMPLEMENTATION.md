# SST Integration - Step 2 Implementation

This document describes the implementation of Step 2 from the SST Integration migration plan.

## What Was Implemented

Step 2 focuses on creating the SST stack scaffolding that mirrors the current CDK infrastructure.

### Files Created

1. **`package.json`** (root level)
   - Added SST dependencies
   - Defined npm scripts for SST commands
   - Set Node.js engine requirement

2. **`sst.config.ts`**
   - Main SST configuration file
   - Defines the app configuration
   - Orchestrates stack creation
   - Configures AWS provider

3. **`sst/config/environments.ts`**
   - Environment-specific configuration
   - Ported from CDK's `deployments/cdk/lib/config/environments.ts`
   - Supports three stages: next, live, test

4. **`sst/stacks/infrastructure.ts`**
   - Placeholder for stateful resources
   - Mirrors CDK's `InfrastructureStack`
   - Defines resource interfaces for cross-stack references

5. **`sst/stacks/application.ts`**
   - Placeholder for stateless resources
   - Mirrors CDK's `ApplicationStack`
   - Logs intended functionality

6. **`sst/README.md`**
   - Comprehensive documentation
   - Usage instructions
   - Migration status
   - Relationship with CDK

7. **`tsconfig.json`** (root level)
   - TypeScript configuration for SST
   - Module resolution for ESM

8. **`.gitignore`** (updated)
   - Added SST artifacts exclusions
   - Node modules
   - Environment files

## Stack Architecture

The implementation creates two main stacks:

### Infrastructure Stack
Contains placeholder definitions for:
- S3 buckets (archive and cache)
- DynamoDB catalog table
- SNS archive topic
- SQS archive queue

All resources currently return placeholder values with proper TypeScript interfaces.

### Application Stack
Contains placeholder definitions for:
- API Gateway
- Lambda functions (endpoints and workers)
- Cognito configuration
- Web application (Waku)
- Lambda authorizers

All components log their intended setup without creating actual resources.

## Design Decisions

### TypeScript Configuration

The SST configuration uses TypeScript with ESM modules. The `/// <reference path="./.sst/platform/config.d.ts" />` directive at the top of `sst.config.ts` is intentional - SST generates this file when initialized.

TypeScript errors are expected until SST is initialized with `sst dev` or `sst deploy`, which creates the `.sst` directory with type definitions.

### Resource Naming

Resource names follow the pattern: `dphoto-{stage}-{resource}-placeholder`

This makes it clear these are scaffolds and prevents accidental conflicts with existing CDK resources.

### Configuration Porting

The environment configuration was directly ported from the CDK configuration to ensure consistency. This includes:
- Domain names
- OAuth2 settings
- Production flags
- CLI access keys

### Stack Dependencies

Unlike CDK where dependencies are explicitly declared with `addDependency()`, SST handles dependencies implicitly through resource references. The Application Stack receives infrastructure resources as parameters, establishing the dependency relationship.

## Testing the Implementation

### Prerequisites

```bash
npm install
```

### Validate Configuration

The configuration can be validated by checking that:
1. All files are syntactically correct
2. Imports resolve properly
3. Environment configuration loads correctly

### Note on TypeScript Errors

TypeScript compilation will show errors for SST globals (`$config`, `$app`) until SST is initialized. This is expected and normal. These globals are injected by SST at runtime.

To initialize SST and generate type definitions:
```bash
# This will fail without AWS credentials but will generate types
npm run sst:dev 2>&1 | grep -i "error"
```

## What's NOT Included

This is a **scaffolding only** implementation. The following are explicitly NOT included:

- ❌ Actual AWS resource creation
- ❌ Lambda function implementations
- ❌ S3 bucket configurations
- ❌ DynamoDB table definitions
- ❌ API Gateway routes
- ❌ Cognito User Pool setup
- ❌ Testing infrastructure
- ❌ CI/CD integration

These will be implemented in subsequent steps of the migration plan.

## Relationship to CDK

Both systems can coexist:

| Component | CDK Location | SST Location | Status |
|-----------|--------------|--------------|--------|
| Infrastructure | `deployments/cdk/lib/stacks/infrastructure-stack.ts` | `sst/stacks/infrastructure.ts` | Scaffolded |
| Application | `deployments/cdk/lib/stacks/application-stack.ts` | `sst/stacks/application.ts` | Scaffolded |
| Config | `deployments/cdk/lib/config/environments.ts` | `sst/config/environments.ts` | Ported |

The CDK infrastructure remains the active deployment. The SST scaffolding is preparation for future migration.

## Next Steps

According to the migration plan:

**Step 3: Local Development Setup**
- Configure `sst dev` for local development
- Set up environment variables and secrets management
- Configure AWS profile and region settings
- Document local development workflow
- Test local development setup with a simple function

## Validation Checklist

- [x] SST dependencies installed
- [x] Configuration file created with proper structure
- [x] Environment configuration ported from CDK
- [x] Infrastructure stack scaffolded with proper interfaces
- [x] Application stack scaffolded with placeholders
- [x] Documentation created (README.md)
- [x] .gitignore updated for SST artifacts
- [x] TypeScript configuration created
- [x] Stack dependencies configured
- [x] Resource naming conventions established

## Known Limitations

1. **TypeScript Errors**: Expected until SST initialization
2. **No AWS Resources**: Placeholder values only
3. **No Testing**: Test infrastructure not yet created
4. **No CI/CD**: GitHub Actions not yet configured for SST
5. **CDK Active**: CDK remains the active deployment system

## References

- [SST Documentation](https://sst.dev/docs)
- [Migration Spec](../specs/2025-12-SST-integration.md)
- [CDK Implementation](../deployments/cdk/)
