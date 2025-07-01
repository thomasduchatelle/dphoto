# CDK Migration Requirements Document

## Executive Summary

**Migration Goal:** Replace Terraform + Serverless Framework with AWS CDK for infrastructure management, reducing maintenance overhead, improving developer
experience, and creating a more consistent codebase.

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
- Deploy to isolated test environment (`next`)
- Validate resource creation and SSM parameter compatibility
- **Value:** Proof of concept complete, risk validation done
- **GitHub Actions Provisioning:**
    - **Current:** Terraform (dev/live) → Serverless (dev/live)
  - **Added:** CDK deployment to isolated `next` environment
  - **Workflow:** Terraform (dev/live) + CDK(next) → Serverless (dev/live/next)

### Phase 2: Infrastructure Migration (Milestone 3-4)

**Goal:** Replace Terraform with CDK while maintaining SLS compatibility

**Milestone 3: Terraform State Import**

- Use CDK import functionality to adopt existing Terraform-managed resources
- Migrate `dev` environment first
- Validate SLS deployment still works with CDK-managed infrastructure
- **Value:** `dev` environment fully migrated, reduced Terraform maintenance
- **GitHub Actions Provisioning:**
  - **Current:** Terraform (dev/live) + CDK(next) → Serverless (dev/live/next)
  - **Changed:** CDK (dev) → Serverless (dev)
  - **Workflow:** Terraform (live) + CDK(next/dev) → Serverless (dev/live/next)

**Milestone 4: Production Migration**

- Migrate `live` environment using proven process from Milestone 3
- Remove Terraform configuration and state
- Update CI/CD to use CDK only for infrastructure
- **Value:** Complete infrastructure migration, single tool for infrastructure
- **GitHub Actions Provisioning:**
  - **Target:** CDK(dev/live/next) → Serverless (dev/live/next)
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
