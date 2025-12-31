/// <reference path="./.sst/platform/config.d.ts" />

/**
 * DPhoto SST Configuration
 * 
 * This configuration will gradually replace the existing AWS CDK infrastructure.
 * See specs/2025-12-SST-integration.md for the complete migration plan.
 * 
 * Current Status: Step 1 - Initial Setup (no resources deployed yet)
 */

export default $config({
  app(input) {
    return {
      name: "dphoto",
      removal: input?.stage === "production" ? "retain" : "remove",
      protect: ["production"].includes(input?.stage),
      home: "aws",
    };
  },
  async run() {
    // Step 1: Initial Setup
    // - SST configuration initialized
    // - No resources deployed yet
    // - CDK remains the active deployment mechanism
    
    // Step 2 (Future): Storage Layer
    // - DynamoDB table
    // - S3 buckets
    // - Cognito User Pool
    
    // Step 3 (Future): Compute Layer
    // - API Gateway
    // - Lambda functions
    // - SNS/SQS queues
    
    // Note: All resources will be added in subsequent steps
    // This file currently serves as a placeholder and validation point
    
    return {
      // Outputs will be added as resources are defined
    };
  },
});
