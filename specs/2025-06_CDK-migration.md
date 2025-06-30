# CDK Migration Requirements Document

## Executive Summary

**Migration Goal:** Replace Terraform + Serverless Framework with AWS CDK for infrastructure management, reducing maintenance overhead, improving developer experience, and creating a more consistent codebase.

**Technologies Being Retired:**
- Terraform (for infrastructure management)
- Serverless Framework (for Lambda deployment)

**Technology Being Introduced:**

- AWS CDK (TypeScript) for unified infrastructure and application deployment

## Current State Analysis

### Current Architecture
- **Terraform:** Manages core AWS infrastructure (DynamoDB, S3, SNS/SQS, IAM, SSM parameters)
- **Serverless Framework:** Manages Lambda functions, API Gateway, and some additional resources
- **Coupling:** SLS depends on SSM parameters created by Terraform
- **Environments:** `dev` (next) and `live` with separate Terraform workspaces and SLS stages
- **CI/CD:** GitHub Actions orchestrating Terraform → Serverless deployment sequence

### Pain Points Addressed
- Terraform verbosity compared to CDK constructs
- Serverless Framework licensing costs for newer versions
- Maintenance overhead of two separate IaC tools
- Version upgrade complexity across multiple tools

## Target State Vision

### Post-Migration Architecture

- **Single CDK Application:** Unified TypeScript codebase managing all infrastructure
- **CDK Stacks:** Organized by logical boundaries (networking, storage, compute, monitoring)
- **Environment Management:** CDK context-based environment configuration
- **CI/CD Integration:** Simplified GitHub Actions workflow using `cdk deploy`
- **Developer Experience:** Single tool for infrastructure changes, better IDE support, type safety

### Benefits Delivered
- **Reduced Complexity:** One tool instead of two
- **Better Maintainability:** Higher-level constructs, less boilerplate
- **Improved DX:** Type safety, better IDE support, unified workflow
- **Cost Optimization:** No Serverless Framework licensing fees
- **Consistency:** Single language/framework for all infrastructure

## Migration Strategy: Iterative Approach

### Phase 1: Foundation & Preparation (Milestone 1-2)
**Goal:** Set up CDK infrastructure without disrupting current deployments

**Milestone 1: CDK Project Setup**
- Create CDK project structure in `/deployments/cdk`
- Set up TypeScript configuration and dependencies
- Create basic CDK app with environment configuration
- Add CDK deployment to CI/CD pipeline (parallel to existing)
- **Value:** Team can start learning CDK, CI/CD pipeline ready for future steps
- **GitHub Actions Provisioning:**
    - **Current:** Terraform (dev/live) → Serverless (dev/live)
    - **Added:** CDK validation job (syntax check, `cdk synth`, unit tests)
    - **No deployment yet** - validation only to ensure CDK setup is correct

**Milestone 2: CDK Infrastructure Parity**
- Implement all Terraform resources in CDK (DynamoDB, S3, SNS/SQS, IAM)
- Create identical SSM parameters for SLS compatibility
- Deploy to isolated test environment
- Validate resource creation and SSM parameter compatibility
- **Value:** Proof of concept complete, risk validation done
- **GitHub Actions Provisioning:**
    - **Current:** Terraform (dev/live) → Serverless (dev/live)
    - **Added:** CDK deployment to isolated `test` environment
    - **Workflow:** Terraform (dev/live) → Serverless (dev/live) + CDK (test)

### Phase 2: Infrastructure Migration (Milestone 3-4)
**Goal:** Replace Terraform with CDK while maintaining SLS compatibility

**Milestone 3: Terraform State Import**
- Use CDK import functionality to adopt existing Terraform-managed resources
- Migrate `dev` environment first
- Validate SLS deployment still works with CDK-managed infrastructure
- **Value:** `dev` environment fully migrated, reduced Terraform maintenance
- **GitHub Actions Provisioning:**
    - **Current:** Terraform (live) → Serverless (live) + CDK (test)
    - **Changed:** CDK (dev) → Serverless (dev) + Terraform (live) → Serverless (live)
    - **Workflow:** CDK (dev) → Serverless (dev) + Terraform (live) → Serverless (live)

**Milestone 4: Production Migration**
- Migrate `live` environment using proven process from Milestone 3
- Remove Terraform configuration and state
- Update CI/CD to use CDK only for infrastructure
- **Value:** Complete infrastructure migration, single tool for infrastructure
- **GitHub Actions Provisioning:**
    - **Target:** CDK (dev/live) → Serverless (dev/live)
    - **Workflow:** CDK infrastructure deployment followed by Serverless application deployment
    - **Terraform removed** from CI/CD pipeline entirely

### Phase 3: Application Migration (Milestone 5-7)
**Goal:** Migrate Lambda functions and API Gateway from SLS to CDK

**Milestone 5: Lambda Migration Framework**
- Create CDK constructs for Lambda functions
- Implement build process integration for Go binaries
- Migrate 2-3 simple Lambda functions (e.g., `version`, `not-found`)
- Deploy alongside existing SLS functions
- **Value:** Lambda migration pattern established, some functions migrated
- **GitHub Actions Provisioning:**
    - **Current:** CDK (dev/live) → Serverless (dev/live)
    - **Enhanced:** CDK (dev/live) with Lambda builds → Serverless (dev/live) with reduced functions
    - **Workflow:** Go binary builds → CDK deployment (infra + some lambdas) → SLS deployment (remaining lambdas)

**Milestone 6: API Gateway Migration**
- Migrate API Gateway and all remaining Lambda functions to CDK
- Implement blue-green deployment for zero-downtime migration
- Remove SLS configuration for migrated functions
- **Value:** Most application components migrated
- **GitHub Actions Provisioning:**
    - **Current:** CDK (dev/live) with some lambdas → Serverless (dev/live) with remaining functions
    - **Target:** CDK (dev/live) with all lambdas and API Gateway → Serverless (dev/live) with minimal resources
    - **Workflow:** Go binary builds → CDK deployment (infra + all lambdas + API Gateway) → SLS deployment (UI bucket only)

**Milestone 7: Complete SLS Removal**
- Migrate remaining SLS resources (S3 bucket for UI, custom resources)
- Remove Serverless Framework entirely
- Clean up CI/CD pipeline
- **Value:** Complete migration, single tool for all infrastructure
- **GitHub Actions Provisioning:**
    - **Target:** CDK (dev/live) only
    - **Final Workflow:** Go binary builds → Web build → CDK deployment (everything)
    - **Serverless Framework removed** from CI/CD pipeline entirely

## Environment Configuration Strategy

**Approach:** Single configuration file with typed environment definitions.

**File Structure:**

    deployments/cdk/
    ├── bin/
    │   └── dphoto.ts                        # Single CDK app entry point
    ├── lib/
    │   ├── dphoto-infrastructure-stack.ts   # DynamoDB, S3, SNS/SQS, IAM
    │   ├── dphoto-application-stack.ts      # Lambdas, API Gateway, CloudFront, web UI
    │   └── config/
    │       └── environments.ts             # Both dev and live configs
    ├── cdk.json                            # CDK configuration
    ├── package.json                        # Dependencies
    └── tsconfig.json                       # TypeScript configuration

**Configuration Pattern:**

    // lib/config/environments.ts
    export interface EnvironmentConfig {
      domainName?: string;
      enableMonitoring: boolean;
      lambdaMemory: number;
      dynamoDbBillingMode: 'PAY_PER_REQUEST' | 'PROVISIONED';
      // ... other typed properties
    }

    export const environments: Record<string, EnvironmentConfig> = {
      dev: {
        domainName: undefined,
        enableMonitoring: false,
        lambdaMemory: 512,
        dynamoDbBillingMode: 'PAY_PER_REQUEST'
      },
      live: {
        domainName: 'photos.example.com',
        enableMonitoring: true,
        lambdaMemory: 1024,
        dynamoDbBillingMode: 'PAY_PER_REQUEST'
      }
    };

**Deployment Pattern:**

    // bin/dphoto.ts - Single environment per deployment
    const envName = app.node.tryGetContext('environment') || 'dev';
    const config = environments[envName];

    new InfrastructureStack(app, `DPhoto-Infrastructure-${envName}`, config);
    new ApplicationStack(app, `DPhoto-Application-${envName}`, config);

**CLI Usage:**

- Development: `cdk deploy --context environment=dev`
- Production: `cdk deploy --context environment=live`

**Rationale:** Single configuration file with explicit environment selection provides type safety, prevents accidental cross-environment deployments, and
follows CDK best practices while maintaining fast command execution.

## Risk Mitigation Strategies

### High-Risk Areas & Mitigation

**1. State Management During Migration**
- **Risk:** Resource recreation causing data loss
- **Mitigation:** Use CDK import functionality, extensive testing in dev environment first

**2. SSM Parameter Compatibility**
- **Risk:** SLS breaking due to parameter format changes
- **Mitigation:** Exact parameter replication in CDK, validation testing

**3. CI/CD Pipeline Disruption**
- **Risk:** Deployment failures during transition
- **Mitigation:** Parallel deployment pipelines, gradual cutover, rollback procedures

**4. Environment Drift**
- **Risk:** Dev and live environments becoming inconsistent
- **Mitigation:** Migrate dev first, validate thoroughly, use identical process for live

### Rollback Strategy
- Maintain Terraform state and configuration until Phase 2 complete
- Keep SLS configuration until Phase 3 complete
- Feature flags for CDK vs SLS Lambda routing during transition
- Automated rollback procedures in CI/CD

## Success Criteria

### Technical Metrics
- Zero downtime during migration
- All existing functionality preserved
- CI/CD pipeline execution time maintained or improved
- Infrastructure drift eliminated

### Business Metrics
- Reduced deployment complexity (single tool)
- Faster development cycles for infrastructure changes
- Reduced maintenance overhead
- Cost savings from SLS licensing elimination

## Timeline Estimation

- **Phase 1:** 2-3 weeks (setup and validation)
- **Phase 2:** 2-3 weeks (infrastructure migration)
- **Phase 3:** 3-4 weeks (application migration)
- **Total:** 7-10 weeks for complete migration

## Resource Requirements

### Team Skills
- CDK/TypeScript knowledge (can be learned during Phase 1)
- AWS infrastructure understanding (existing)
- CI/CD pipeline management (existing)

### Tools & Dependencies
- AWS CDK CLI
- TypeScript/Node.js development environment
- Updated GitHub Actions workflows

## Conclusion

This migration strategy prioritizes incremental value delivery while minimizing risk. Each milestone delivers tangible benefits and can be safely merged to main, ensuring continuous progress without waiting for complete project completion. The approach leverages existing team knowledge while gradually introducing CDK concepts, ensuring a smooth transition to the target state.