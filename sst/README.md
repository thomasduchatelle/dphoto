# SST Infrastructure

This directory contains the SST (Serverless Stack) infrastructure configuration for DPhoto, which is being migrated from AWS CDK.

## Overview

SST is a modern framework built on top of AWS CDK that provides:
- Better developer experience with `sst dev`
- Type-safe infrastructure as code
- Simplified Lambda function deployment
- Built-in support for modern web frameworks like Waku

## Project Structure

```
sst/
├── config/
│   └── environments.ts    # Environment-specific configuration
├── stacks/
│   ├── infrastructure.ts  # Stateful resources (S3, DynamoDB, SQS, SNS)
│   └── application.ts     # Stateless resources (Lambda, API Gateway, Cognito)
└── constructs/            # Reusable SST constructs (to be added)
```

## Configuration

The application supports three environments (stages):
- **next**: Development/staging environment
- **live**: Production environment  
- **test**: Testing environment for CI/CD

Configuration is defined in `sst/config/environments.ts` and includes:
- Domain names
- OAuth2 settings
- Environment-specific flags

## Stack Architecture

### Infrastructure Stack (`infrastructure.ts`)

Contains **stateful resources** that must be protected:
- S3 buckets for media storage (archive and cache)
- DynamoDB table for catalog data
- SNS topics and SQS queues for async processing

These resources have retention policies and should never be accidentally deleted.

### Application Stack (`application.ts`)

Contains **stateless resources** that can be freely redeployed:
- API Gateway (HTTP API v2)
- Lambda functions (endpoints and workers)
- Cognito User Pool
- Waku web application
- Lambda authorizers

## Getting Started

### Prerequisites

- Node.js >= 18.0.0
- AWS CLI configured with appropriate credentials
- Access to the DPhoto AWS account

### Installation

Install dependencies from the project root:

```bash
npm install
```

### Development

Start local development with live Lambda reloading:

```bash
npm run sst:dev
```

This will:
- Deploy a personal development stack to AWS
- Enable live Lambda development
- Watch for code changes

### Deployment

Deploy to a specific stage:

```bash
# Deploy to next (development)
npm run sst:deploy -- --stage next

# Deploy to live (production)
npm run sst:deploy -- --stage live
```

### Preview Changes

See what will change before deploying:

```bash
npm run sst:diff -- --stage next
```

### Remove Stack

Remove a deployed stack:

```bash
npm run sst:remove -- --stage next
```

**Warning**: Be careful with the `live` stage - it has `removal: "retain"` policy to protect production data.

## Migration Status

This is currently a **placeholder implementation** for Step 2 of the SST migration plan. The stacks are scaffolded but not yet implementing actual resources.

### Current Status

- [x] SST configuration file created
- [x] Environment configuration ported from CDK
- [x] Infrastructure stack scaffolded
- [x] Application stack scaffolded
- [ ] Resource implementations (S3, DynamoDB, Lambda, etc.)
- [ ] Testing and validation
- [ ] Documentation completion

### Next Steps

See `specs/2025-12-SST-integration.md` for the complete migration plan.

Step 3 (next): Configure local development setup with `sst dev`

## Relationship with CDK

The SST stacks mirror the existing CDK stacks:

| SST Stack | CDK Stack | Purpose |
|-----------|-----------|---------|
| `infrastructure.ts` | `infrastructure-stack.ts` | Stateful resources |
| `application.ts` | `application-stack.ts` | Stateless resources |

During the migration, both systems will coexist. Once SST is fully operational, the CDK infrastructure will be deprecated.

## Resources

- [SST Documentation](https://sst.dev/docs)
- [SST v3 Guide](https://sst.dev/docs/start/aws/api)
- [DPhoto CDK Implementation](../deployments/cdk/)
- [Migration Spec](../specs/2025-12-SST-integration.md)
