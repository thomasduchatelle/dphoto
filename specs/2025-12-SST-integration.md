# SST (Serverless Stack) Integration Migration

## Context

This document records the plan for migrating DPhoto's infrastructure from AWS CDK to SST (Serverless Stack). SST is a modern framework built on top of CDK that provides better developer experience, type-safe infrastructure, and seamless integration with modern web frameworks like Waku.

## Rationale for SST Migration

1. **Better Developer Experience**: SST provides `sst dev` for local development with AWS resources
2. **Type-Safe Infrastructure**: Infrastructure types are automatically available in the application code
3. **Modern Framework Integration**: Native support for frameworks like Waku, making the migration complementary
4. **Simplified Configuration**: Less boilerplate compared to pure CDK
5. **Live Lambda Development**: Test Lambda functions locally with real AWS resources
6. **Cost Management**: Better visibility and control over AWS resource costs

## Migration Strategy

### Phase 1: Preparation and Setup

**Step 1: SST Installation and Initial Setup**
- Install SST CLI and dependencies
- Create initial SST configuration file (`sst.config.ts`)
- Set up SST project structure alongside existing CDK
- Configure environments (next, live) in SST
- Verify SST can deploy a minimal stack without affecting existing infrastructure

**Step 2: Create SST Stack Scaffolding**
- Create placeholder SST stacks that mirror the current CDK stacks
- Set up Infrastructure stack (stateful resources)
- Set up Application stack (stateless resources)
- Configure stack dependencies and cross-stack references
- Add basic documentation for SST stack structure

**Step 3: Local Development Setup**
- Configure `sst dev` for local development
- Set up environment variables and secrets management
- Configure AWS profile and region settings
- Document local development workflow
- Test local development setup with a simple function

### Phase 2: Incremental Migration

**Step 4: Migrate Storage Infrastructure**
- Migrate DynamoDB tables to SST constructs
- Migrate S3 buckets to SST constructs
- Verify resource naming and logical IDs match
- Test data persistence and access patterns

**Step 5: Migrate Lambda Functions**
- Migrate API Lambda functions to SST Function constructs
- Migrate worker Lambda functions
- Configure function permissions and environment variables
- Test function execution and integrations

**Step 6: Migrate API Gateway**
- Migrate API Gateway HTTP v2 to SST Api construct
- Configure routes and authorizers
- Migrate custom domain configuration
- Test API endpoints and authentication

**Step 7: Migrate Web Application**
- Integrate Waku application with SST
- Configure SSR Lambda deployment
- Set up static asset deployment to S3
- Test web application routing and rendering

### Phase 3: Completion and Cleanup

**Step 8: Migration Verification**
- Run full test suite with SST infrastructure
- Verify all endpoints and functionality
- Performance testing and optimization
- Security audit of new infrastructure

**Step 9: Documentation and Training**
- Update deployment documentation
- Create troubleshooting guides
- Document differences from CDK approach
- Train team on SST workflows

**Step 10: CDK Deprecation**
- Archive CDK deployment code
- Update CI/CD pipelines to use SST exclusively
- Remove CDK dependencies
- Final verification and sign-off

## Current Status

- [ ] Step 1: SST Installation and Initial Setup
- [x] Step 2: Create SST Stack Scaffolding ✅ **COMPLETED**
- [ ] Step 3: Local Development Setup
- [ ] Step 4: Migrate Storage Infrastructure
- [ ] Step 5: Migrate Lambda Functions
- [ ] Step 6: Migrate API Gateway
- [ ] Step 7: Migrate Web Application
- [ ] Step 8: Migration Verification
- [ ] Step 9: Documentation and Training
- [ ] Step 10: CDK Deprecation

## Technical Decisions

### Stack Organization

**Decision**: Maintain two-stack pattern (Infrastructure + Application)

**Rationale**: 
- Consistent with current CDK architecture
- Separates stateful and stateless resources
- Allows for independent deployment schedules
- Protects critical data during deployments

### Resource Naming Strategy

**Decision**: Preserve existing logical IDs for stateful resources

**Rationale**:
- Prevents data loss during migration
- Maintains existing resource references
- Simplifies rollback if needed
- Reduces migration risk

### Environment Configuration

**Decision**: Use SST's built-in stage management for environments

**Rationale**:
- Native SST feature for environment management
- Better than CDK context-based configuration
- Simplified CI/CD pipeline configuration
- Improved local development experience

## Next Steps

Current focus: **Step 2 - Create SST Stack Scaffolding**

This step involves creating the basic SST stack structure that will eventually replace the CDK stacks, without yet migrating actual resources.
