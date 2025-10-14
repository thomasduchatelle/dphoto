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
    │   ├── dphoto.ts                           # CDK app entry point
    │   └── dphoto.test.ts                      # Integration tests
    ├── lib/
    │   ├── stacks/
    │   │   ├── infrastructure-stack.ts         # Aggregate of store constructs (stateful resources that must never be lost)
    │   │   ├── infrastructure-stack.test.ts    # Tests of "InfrastructureStack" follwing the Testing Philosophy  
    │   │   ├── application-stack.ts            # Aggregate of other constructs (stateless, frequently deployed resources)
    │   │   └── application-stack.test.ts       # Tests of "ApplicationStack" follwing the Testing Philosophy
    │   ├── archive/                            # "archive" is the name of one of the domain
    │   │   ├── archive-store-construct.ts
    │   │   ├── archive-workers-construct.ts
    │   │   ├── archive-workers-construct.test.ts       # Unit tests can be added only if required
    │   │   └── ... # constructs abstract complex resource combinations from stacks
    │   ├── <other domains>/
    │   │   └── ...
    │   └── utils/
    │       ├── golang-lambda-function.ts               # Lambda construct for GoLang deployment
    │       └── ...
    ├── cdk.json                                # CDK configuration
    ├── package.json                            # Dependencies
    ├── tsconfig.json                           # TypeScript configuration
    └── jest.config.js                          # Test configuration

## Stack Architecture

### Two-Stack Pattern

**Infrastructure Stack** (`infrastructure-stack.ts`):

* Constructs with stateful resources **that must never be lost** (examples: DynamoDB table, S3 buckets)
* Constructs with resources exposed to other components of the system
  * SSM parameters used by legacy deployments (Serverless Framework)
  * SNS topics and SQS queues used a communication between two domains
* Underlying AWS resource must have their Logical ID pinned, and tested, to prevent resources to be re-created when they are moved between constructs

**Application Stack** (`application-stack.ts`):

* Constructs with **stateless, or frequently deployed resources**
  * examples: Lambda functions, API Gateway, CloudFront, IAM roles

### Clean Code

#### Stack Design Principles

- **Descriptive & Functional:** Stacks read like high-level blueprints, not implementation details
- **Single Responsibility:** Each construct handles one logical component
- **Abstraction Layers:** Complex resource combinations moved to constructs, complex logic to utils

_Implementation Pattern:_

```typescript
// Stack: High-level, declarative, readable
export class InfrastructureStack extends Stack {
  public readonly archiveStore: ArchiveStoreConstruct; // exposed for cross-stacks communication

  constructor(scope: Construct, id: string, props: StackProps) {
        super(scope, id, props);

        // Clean, descriptive resource creation
        this.archiveStore = new ArchiveStoreConstruct(this, 'ArchiveStore', config);
        // ... other "store" constructs
  }
}

export class ApplicationStack extends Stack {
  constructor(scope: Construct, id: string, props: StackProps & {
    archiveStore: ArchiveStoreConstruct,
  }) {
        super(scope, id, props);

        new ArchiveWorkersConstruct(this, 'ArchiveWorkers', {...config, archiveStore});
        // ... other "worker" or "endpoints" constructs
  }
}
```

_Construct Extraction Rules:_

1. **Readability** No AWS resource can be declared directly in a Stack: all must be abstracted through constructs
2. **Type** Stateful and persistent resources must be declared in a different construct from the volatile and stateless resources
3. **Complexity:** a construct must always be about a main resource with its supporting resources (example: IAM role and Log group are supporting a Lambda) ;
   when construct responsibility becomes unclear, it should be broke down into several sub-constructs
4. **Affinity** IAM roles and SQS queues specific for a consumer must be declared in the same construct as the consumer

### Cross-Stack Communication

Each Stack exposes the constructs with capabilities required by other stacks as public read-only properties.

Each Construct exposes as `public` functions its capabilities: it will patch the resources provided as arguments. This arguments must be interfaces to avoid
tight coupling and cycle dependencies (examples: `iam.IGrantable`, or custom interfaces, or combination of both).

The reference of the dependency construct is passed to the construct requiring its capabilities.

* **grant functions**: grant specific access to the argument. Example:
  ```typescript
  // lib/utils/workload.ts
  interface Workload { // we define Workload interface instead of taking the concrete implementaiton "GoLangLambdaFunction"
    role: iam.IGrantable;
    function: {addEnvironment(key: string, value: string)};
  }

  // lib/catalog/catalog-store-construct.ts
  class CatalogStoreConstruct {
    public grantRWToCatalogTable(lambda: Workload) {
        this.catalogTable.grantReadWriteData(lambda.role);
        lambda.function.addEnvironment("CATALOG_TABLE_NAME", this.catalogTable.tableName)
    }
  }
  ```

## Testing Philosophy

**Test only critical contracts:**

1. **Data protection**: databases and stores retention properties, logical IDs, deletion protection, and Point In Time Recovery (if available)
2. **Security**: public access and default endpoints disabled
3. **Interfaces** with external system: SSM parameters (used by legacy deployments mechanism)
4. **Endpoints** exposed on the API Gateway integrations (used by CLI and UI)

**Do not duplicate in the test the declarative code from the constructs.** Tests must be robust to refactoring, if the tests needs to be changed each time a
minor property change, a resource is added or is removed, they lose their value.

```typescript
// deployments/cdk/libs/stacks/infrastructure-stack.test.ts
describe('InfrastructureStack', () => {
  test('exports all required SSM parameters for Serverless deployments', () => {
        // Verify SSM parameters exist with correct paths
  });

  test('main S3 bucket has deletion protection', () => {
        // Verify bucket cannot be accidentally deleted
  });
});
```

## Naming Conventions

* **Stacks**:
  * **Stacks IDs**: hyphen case prefixed with "dphoto-{environment name}" (example: `dphoto-next-infrastructure`, where `next` is the environment)
  * **Stack Class**: class naming is in PascalCase and suffixed by "Stack" (example: `InfrastructureStack`)
  * **Stack File**: the file has the same name in hyphen case (example: `infrastructure-stack.ts`)
* **Resources**:
  * **Resource IDs**: PascalCase suffixed by "Construct" (example: `CatalogStoreConstruct`)
  * **Resource Class**: same as the resource ID (example: `CatalogStoreConstruct`)
  * **Resource File**: same as the class but in hyphen case (example: `catalog-store-construct.ts`)

## Environment Configuration

    export interface EnvironmentConfig {
      domainName: string;
      // ... ONLY the properties that are different between the environment are set explicitly. Variables from environment must be set explicitely, they never have a default value.
    }

**Context-based selection:**

- `cdk deploy --context environment=next`
- `cdk deploy --context environment=live`

## Utils

The following constructs are already implemented and must be used when relevant.

### Utils: GoLang function

It abstracts the Lambda resource function and its associated role, and gives sensible default for each variable.

```typescript
// deployments/cdk/lib/utils/golang-lambda-function.ts
export interface GoLangLambdaFunctionProps {
  environmentName: string;
  functionName: string;
  artifactPath?: string;
  timeout?: cdk.Duration;
  memorySize?: number;
  environment?: Record<string, string>;
}

export class GoLangLambdaFunction extends Construct {
  public readonly function: lambda.Function;
  public readonly role: iam.Role;

  constructor(scope: Construct, id: string, props: GoLangLambdaFunctionProps) {
    super(scope, id, props)
    // ...
  }
}
```

### Utils: SimpleGoEndpoint

Create a GO Lambda and expose it through the API Gateway

```typescript
// deployments/cdk/lib/utils/simple-go-endpoint.ts
export interface RouteConfig {
  path: string;
  method: apigatewayv2.HttpMethod;
}

export interface SimpleGoEndpointProps extends GoLangLambdaFunctionProps {
  httpApi: apigatewayv2.HttpApi;
  routes: RouteConfig[];
}

export class SimpleGoEndpoint extends Construct {
  public readonly lambda: GoLangLambdaFunction;
  private readonly integration: apigatewayv2_integrations.HttpLambdaIntegration;

  constructor(scope: Construct, id: string, props: SimpleGoEndpointProps) {
    super(scope, id, props);
    // ...
    }
}
```

## Edit Mode for LLM

**Important !** When a file needs to be renamed, moved, or deleted, **do not edit the file**: create the new one, then provide a shell commands to delete the
old file(s).