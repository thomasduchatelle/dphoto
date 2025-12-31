/// <reference path="./.sst/platform/config.d.ts" />

/**
 * SST Configuration for DPhoto
 * 
 * This configuration defines the infrastructure for DPhoto using SST v3.
 * It mirrors the existing CDK stacks (Infrastructure and Application) but uses
 * SST's simplified constructs and improved developer experience.
 * 
 * Environments:
 * - next: Development/staging environment
 * - live: Production environment
 * - test: Testing environment for CI/CD
 * 
 * Note: The $config, $app globals are provided by SST at runtime.
 * TypeScript errors are expected until SST initializes (.sst directory created).
 */

export default $config({
  app(input) {
    return {
      name: "dphoto",
      removal: input?.stage === "live" ? "retain" : "remove",
      home: "aws",
      providers: {
        aws: {
          region: "eu-west-1",
        },
      },
    };
  },
  async run() {
    // Get the current stage (environment)
    const stage = $app.stage;
    
    // Import environment configuration
    const { getEnvironmentConfig } = await import("./sst/config/environments");
    const config = getEnvironmentConfig(stage);

    // Infrastructure Stack - Stateful resources that must never be lost
    // Includes: DynamoDB tables, S3 buckets, SNS topics, SQS queues
    const infrastructure = await import("./sst/stacks/infrastructure");
    const infra = infrastructure.createInfrastructureStack(config);

    // Application Stack - Stateless resources that can be freely redeployed
    // Includes: Lambda functions, API Gateway, web application
    const application = await import("./sst/stacks/application");
    application.createApplicationStack({
      config,
      infrastructure: infra,
    });

    // Outputs
    return {
      infrastructure: {
        archiveBucket: infra.archiveBucket.name,
        cacheBucket: infra.cacheBucket.name,
        catalogTable: infra.catalogTable.name,
      },
    };
  },
});
