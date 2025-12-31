# SST Integration Plan

## Overview

This document outlines the migration plan from AWS CDK to SST (Serverless Stack) for the DPhoto application infrastructure. SST provides a better developer experience with features like Live Lambda Development, type-safe infrastructure, and improved local testing capabilities.

## Motivation

**Current State: AWS CDK**
- Complex TypeScript configuration
- Slow deployment cycles
- Limited local development experience
- Manual environment management

**Target State: SST v3**
- Modern, simpler infrastructure as code
- Live Lambda Development for instant feedback
- Better TypeScript integration
- Improved cost visibility
- Enhanced developer experience

## Migration Strategy

The migration will be done in three phases to ensure zero downtime and maintain data integrity:

### Step 1: Initial Setup and Preparation

**Objective:** Set up SST alongside existing CDK infrastructure without affecting production.

**Tasks:**
1. Install SST v3 CLI and dependencies
2. Initialize SST configuration (`sst.config.ts`)
3. Create basic SST project structure
4. Document the new project layout
5. Add SST to CI/CD pipeline preparation
6. Create migration guide documentation

**Deliverables:**
- `sst.config.ts` configuration file
- Updated `package.json` with SST dependencies
- Documentation in `docs/sst-migration-guide.md`
- `.gitignore` updates for SST artifacts

**Success Criteria:**
- SST CLI can run successfully
- Basic configuration is valid
- Documentation is complete
- No changes to existing CDK deployment

### Step 2: Infrastructure Parity - Storage Layer

**Objective:** Recreate stateful resources (DynamoDB, S3, Cognito) in SST configuration while keeping CDK as the active deployment.

**Tasks:**
1. Define DynamoDB table in SST (`sst.config.ts`)
   - Single-table design matching `DATA_MODEL.md`
   - Same partition/sort keys and GSIs
   - Point-in-time recovery enabled
   - Retention policies configured
2. Define S3 buckets in SST
   - Archive bucket with lifecycle rules
   - Cache bucket configuration
   - CORS settings
3. Define Cognito User Pool in SST
   - User groups (admins, owners, visitors)
   - Google OAuth integration
   - Custom domain configuration
4. **DO NOT DEPLOY** - only validate configuration

**Deliverables:**
- SST resource definitions matching CDK
- Validation scripts to compare configurations
- Documentation of differences (if any)

**Success Criteria:**
- `sst build` succeeds
- Configuration matches existing CDK resources
- No actual deployment occurs
- Side-by-side comparison documentation

### Step 3: Infrastructure Parity - Compute Layer

**Objective:** Recreate Lambda functions and API Gateway in SST configuration.

**Tasks:**
1. Define API Gateway in SST
   - HTTP API (v2) configuration
   - Routes matching existing endpoints
   - Authorizer configuration
2. Define all Lambda functions in SST
   - Go Lambda build configuration
   - Environment variables
   - IAM permissions
   - Timeout and memory settings
3. Define SNS/SQS resources for async processing
4. Configure CloudFront (if applicable)
5. **DO NOT DEPLOY** - only validate configuration

**Deliverables:**
- Complete SST configuration for all compute resources
- Lambda build validation
- API endpoint mapping documentation

**Success Criteria:**
- All Lambda functions build successfully
- API Gateway configuration matches CDK
- IAM permissions are equivalent
- No actual deployment occurs

### Step 4: Deployment Strategy and Testing (Future)

**Note:** Steps 4-6 are outlined for completeness but are NOT part of the current implementation scope.

**Objective:** Plan the actual switchover from CDK to SST.

**Tasks:**
1. Create deployment runbook
2. Plan data migration (if needed)
3. Set up rollback procedures
4. Test deployment in `next` environment
5. Performance testing and validation

### Step 5: Production Migration (Future)

**Objective:** Switch production to SST-managed infrastructure.

**Tasks:**
1. Deploy SST infrastructure to new stack
2. Migrate data (if necessary)
3. Update DNS/CloudFront to point to new resources
4. Monitor and validate
5. Maintain CDK as backup for 30 days

### Step 6: Cleanup (Future)

**Objective:** Remove CDK infrastructure after successful migration.

**Tasks:**
1. Verify SST infrastructure stability (30 days)
2. Remove CDK code
3. Update all documentation
4. Remove CDK from CI/CD pipelines

## Key Considerations

### Data Safety
- **CRITICAL:** DynamoDB and S3 data must never be lost
- Use AWS resource import where possible
- Maintain backups throughout migration
- Test restore procedures

### Cost Implications
- Running dual infrastructure briefly acceptable (< $10/month overhead)
- Monitor costs during migration
- Plan immediate cleanup after switchover

### Zero Downtime
- Use Blue/Green deployment strategy
- API Gateway can route to either CDK or SST Lambdas
- CloudFront can serve from either stack
- Cognito User Pool can be shared initially

### Testing Requirements
- Validate in `next` environment first
- Test all critical user journeys
- Performance benchmarking
- Security audit

## SST Configuration Structure

```
/
├── sst.config.ts              # Main SST configuration
├── sst-env.d.ts               # Auto-generated types
├── .sst/                      # SST build artifacts (gitignored)
├── deployments/
│   ├── cdk/                   # Existing CDK (to be deprecated)
│   └── sst/                   # SST-specific modules (if needed)
├── api/lambdas/               # Existing Lambda code (no changes)
├── web/                       # Existing web code (no changes)
└── docs/
    └── sst-migration-guide.md # Migration documentation
```

## Timeline Estimate

- **Step 1:** 2-4 hours (initial setup)
- **Step 2:** 4-6 hours (storage layer parity)
- **Step 3:** 6-8 hours (compute layer parity)
- **Step 4-6:** TBD (actual migration and cleanup)

**Total for Steps 1-3:** 12-18 hours
**Risk Buffer:** +4-6 hours for unexpected issues

## Rollback Plan

At any point before Step 5:
1. Continue using CDK infrastructure
2. Remove SST configuration files if desired
3. No impact to running services

After Step 5:
1. Update CloudFront to point back to CDK stack
2. Restore from backups if data issues occur
3. Keep SST stack for 30 days before cleanup

## Success Metrics

- All existing functionality works identically
- Deployment time reduced by 50%+
- Local development experience improved
- Cost remains under $5/month
- Zero data loss
- Zero downtime during migration

## References

- [SST Documentation](https://docs.sst.dev/)
- [SST v3 Migration Guide](https://docs.sst.dev/upgrade-guide)
- [AWS CDK to SST Migration](https://docs.sst.dev/migrating/cdk)
- [DPhoto DATA_MODEL.md](/DATA_MODEL.md)
- [DPhoto README.md](/README.md)

## Notes

- This migration is opportunistic - improve developer experience while maintaining stability
- Current CDK infrastructure works well - migration is enhancement, not fix
- Take time to validate each step thoroughly
- Document learnings for other projects
