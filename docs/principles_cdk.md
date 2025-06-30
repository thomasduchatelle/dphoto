# CDK Development Principles for DPhoto

## Overview

This document provides comprehensive guidance for LLM agents and developers working on DPhoto's AWS CDK infrastructure. It covers architectural decisions,
development patterns, testing requirements, and operational considerations.

## Project Context

DPhoto: Personal photo backup system with Go Lambda backend, React frontend, CDK infrastructure.

- **Environments:** `next` (dev) and `live` (prod)
- **Constraint:** <$5/month total cost

## CDK Project Structure

    deployments/cdk/
    ├── bin/
    │   └── dphoto.ts                        # CDK app entry point
    ├── lib/
    │   ├── stacks/
    │   │   ├── dphoto-infrastructure-stack.ts   # Core AWS resources
    │   │   └── dphoto-application-stack.ts      # Application components
    │   ├── constructs/
    │   │   ├── lambda-function.ts               # Reusable Lambda construct
    │   │   ├── s3-bucket.ts                     # Standardized S3 buckets
    │   │   └── dynamodb-table.ts               # DynamoDB table construct
    │   ├── config/
    │   │   └── environments.ts                 # Environment configurations
    │   └── utils/
    │       ├── ssm-parameters.ts               # SSM parameter helpers
    │       └── naming.ts                       # Resource naming conventions
    ├── test/
    │   ├── unit/
    │   │   ├── infrastructure-stack.test.ts
    │   │   └── application-stack.test.ts
    │   └── integration/
    │       └── deployment.test.ts
    ├── cdk.json                                # CDK configuration
    ├── package.json                           # Dependencies
    ├── tsconfig.json                          # TypeScript configuration
    └── jest.config.js                         # Test configuration

## Stack Architecture

### Two-Stack Pattern

**Infrastructure Stack** (`dphoto-infrastructure-stack.ts`):

- Stateful resources **that must never be lost**
- DynamoDB table, S3 buckets (main/cache), SNS/SQS
- Exports SSM parameters for application stack

**Application Stack** (`dphoto-application-stack.ts`):

- Stateless, frequently deployed resources
- Lambda functions, API Gateway, CloudFront, IAM roles
- Imports infrastructure ARNs via SSM parameters

### Cross-Stack Communication

All references via SSM parameters at `/dphoto/{env}/{service}/{param}`:

    /dphoto/{environment}/
    ├── dynamodb/
    │   ├── table-name
    │   └── table-arn
    ├── s3/
    │   ├── main-bucket-name
    │   ├── cache-bucket-name
    │   └── web-bucket-name
    ├── sns/
    │   └── archive-arn
    ├── sqs/
    │   └── archive-url
    ├── iam/
    │   ├── lambda-execution-role-arn
    │   └── storage-access-role-arn
    └── api/
        ├── gateway-url
        └── cloudfront-domain

## Resource Naming

**Pattern:** `dphoto-{env}-{resource-type}`

**NamingConvention class** in `lib/utils/naming.ts`:

- `dynamoTable()`: `dphoto-${env}-index`
- `s3Bucket(type)`: `dphoto-${env}-${type}`
- `lambdaFunctionName(name)`: `dphoto-${env}-${name}`

## Environment Configuration

    export interface EnvironmentConfig {
      account: string;
      region: string;
      domainName?: string;
      enableMonitoring: boolean;
      lambdaMemory: number;
      dynamoDbBillingMode: 'PAY_PER_REQUEST' | 'PROVISIONED';
      useDefaultVpc: boolean;
      useDefaultEncryption: boolean;
    }

**Context-based selection:**

- `cdk deploy --context environment=next`
- `cdk deploy --context environment=live`

## Lambda Patterns

**GoLambdaFunction construct:**

- Runtime: `PROVIDED_AL2`
- Build: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0`
- Handler: `bootstrap`
- Source: relative to CDK root

## Testing Philosophy

**Test only critical contracts:**

1. SSM parameter exports/imports between stacks
2. Data protection (S3 deletion protection, DynamoDB PITR)
3. Critical API Gateway integrations

**Do NOT test:** CDK properties, resource creation, configuration values

    // test/unit/infrastructure-stack.test.ts
    describe('InfrastructureStack', () => {
      test('exports all required SSM parameters for application stack', () => {
        // Verify SSM parameters exist with correct paths
      });
      
      test('main S3 bucket has deletion protection', () => {
        // Verify bucket cannot be accidentally deleted
      });
    });

## Development Workflow

1. **Deploy order:** Infrastructure → Application
2. **Development:** Use `next` environment
3. **Updates:** Application stack independent, infrastructure requires planning
4. **Migration:** Use CDK import, preserve existing resource names

## Cost Optimization

**Critical Constraint: <$5/month total**

**Infrastructure:**

- Default VPC (no NAT costs)
- Default encryption (no KMS costs)
- Single AZ deployment
- Minimal CloudWatch alarms
- On-demand billing

## Security

**IAM:** Least privilege, SSM parameter references, no hardcoded ARNs
**Encryption:** AWS managed keys only
**Network:** Default VPC security groups