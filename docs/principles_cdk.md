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
    │   │   ├── golang-lambda-function.ts               # Lambda construct for GoLang deployment
    │   │   ├── infra_media_storage.ts                  # Media storage construct - complex combination of S3 resource and policies
    │   │   ├── infra_catalog_dynamodb.ts               # DynamoDB table construct - abstracted for stack readability
    │   │   └── ... # constructs abstract complex resource combinations from stacks
    │   ├── config/
    │   │   └── environments.ts                 # Environment configurations
    │   └── utils/
    │       └── ...                             # Pure functions for complex logic extraction (clean code)
    ├── test/
    │   ├── unit/
    │   │   ├── infrastructure-stack.test.ts
    │   │   └── application-stack.test.ts
    │   └── integration/
    │       └── deployment.test.ts
    ├── cdk.json                               # CDK configuration
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

### Clean Code

#### Stack Design Principles

- **Descriptive & Functional:** Stacks read like high-level blueprints, not implementation details
- **Single Responsibility:** Each construct handles one logical component
- **Abstraction Layers:** Complex resource combinations moved to constructs, complex logic to utils

_Implementation Pattern:_

    // Stack: High-level, declarative, readable
    export class DPhotoInfrastructureStack extends Stack {
      constructor(scope: Construct, id: string, props: StackProps) {
        super(scope, id, props);
        
        // Clean, descriptive resource creation
        const mediaStorage = new MediaStorageConstruct(this, 'MediaStorage', config);
        const catalogDb = new CatalogDynamoDbConstruct(this, 'CatalogDb', config);
        
        // Export parameters for application stack
        exportSsmParameters(this, mediaStorage, catalogDb);
      }
    }

_Construct Extraction Rules:_

1. **Complexity:** >3 related resources or complex configuration
2. **Readability:** Improves stack comprehension by abstracting implementation

_Utils Extraction Rules:_

1. **Logic Complexity:** Conditional resource creation, parameter transformation
2. **Testability:** Complex business logic that benefits from unit testing

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

**Pattern:** names are in hyphen case and contains the environment name (example: `dphoto-next-catalog-db`, where `next` is the environment)

If imported resources are not named following the pattern, the name must be set in the configuration.

## Environment Configuration

    export interface EnvironmentConfig {
      domainName: string;
      // ... ONLY the properties that are different between the environment are set explicitly. In this case, they should NOT have a default value.
    }

**Context-based selection:**

- `cdk deploy --context environment=next`
- `cdk deploy --context environment=dev`
- `cdk deploy --context environment=live`

## Adding New Endpoints

### Simple Go Endpoints

For standard Go Lambda endpoints, use the `SimpleGoEndpoint` construct:

**Step 1:** Add the endpoint to an appropriate construct (or create a new one):

```typescript
// In lib/constructs/your-endpoints.ts
new SimpleGoEndpoint(this, 'YourEndpoint', {
    environmentName: props.environmentName,
    functionName: 'your-function-name',  // Will be prefixed with dphoto-{env}- automatically
    httpApi: props.apiGateway.httpApi,
    path: '/api/v1/your/path',
    method: apigatewayv2.HttpMethod.GET,
    // Optional overrides:
    memorySize: 512,                     // Default: 256
    timeout: Duration.seconds(30),       // Default: 1 minute
    artifactPath: '../../bin/custom.zip' // Default: ../../bin/{functionName}.zip
});
```

**Step 2:** Add the construct to the application stack:

```typescript
// In lib/stacks/dphoto-application-stack.ts
const yourEndpoints = new YourEndpoints(this, 'YourEndpoints', {
    environmentName: props.environmentName,
    apiGateway: apiGateway
});
```

### Endpoint Organization

**Group related endpoints** into logical constructs:
- `MetadataEndpoints`: version, health, not-found
- `AuthEndpoints`: login, logout, token refresh
- `AlbumEndpoints`: CRUD operations for albums
- `MediaEndpoints`: media upload, download, metadata

**Naming Convention:**
- Construct: `{Domain}Endpoints` (e.g., `AuthEndpoints`)
- Endpoint ID: Descriptive name (e.g., `CreateAlbum`, `ListMedias`)
- Function name: kebab-case matching binary (e.g., `create-album`) - will be automatically prefixed with `dphoto-{env}-`
- Lambda function name: `dphoto-{env}-{function-name}` (e.g., `dphoto-next-create-album`)

### Complex Endpoints

For endpoints requiring custom configuration (environment variables, IAM permissions, etc.), extend `SimpleGoEndpoint` or create a specialized construct:

```typescript
export class CustomEndpoint extends SimpleGoEndpoint {
    constructor(scope: Construct, id: string, props: CustomEndpointProps) {
        super(scope, id, {
            ...props,
            memorySize: 1024,
            timeout: Duration.minutes(5)
        });
        
        // Add custom environment variables
        this.lambda.addEnvironment('CUSTOM_CONFIG', props.customValue);
        
        // Add custom IAM permissions
        this.lambda.addToRolePolicy(new PolicyStatement({
            actions: ['s3:GetObject'],
            resources: ['arn:aws:s3:::bucket/*']
        }));
    }
}
```

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
