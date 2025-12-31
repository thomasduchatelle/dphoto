/**
 * Infrastructure Stack for DPhoto SST
 * 
 * This stack contains all STATEFUL resources that must be retained and protected.
 * These resources store critical data and should never be accidentally deleted.
 * 
 * Resources include:
 * - DynamoDB tables (catalog data, user permissions, album information)
 * - S3 buckets (original media storage, cached miniatures)
 * - SNS topics (asynchronous job notifications)
 * - SQS queues (background job processing)
 * 
 * This mirrors the CDK InfrastructureStack from deployments/cdk/lib/stacks/infrastructure-stack.ts
 */

import { EnvironmentConfig } from "../config/environments";

/**
 * Infrastructure stack return type
 * Contains references to all stateful resources that other stacks need to access
 */
export interface InfrastructureStack {
  /** S3 bucket for storing original media files */
  archiveBucket: {
    name: string;
    arn: string;
  };
  /** S3 bucket for caching processed media (miniatures, thumbnails) */
  cacheBucket: {
    name: string;
    arn: string;
  };
  /** DynamoDB table for catalog data (single-table design) */
  catalogTable: {
    name: string;
    arn: string;
  };
  /** SNS topic for archive job notifications */
  archiveTopic: {
    name: string;
    arn: string;
  };
  /** SQS queue for processing archive jobs */
  archiveQueue: {
    name: string;
    arn: string;
    url: string;
  };
}

/**
 * Create the Infrastructure Stack
 * 
 * @param config Environment-specific configuration
 * @returns Infrastructure stack resources for cross-stack references
 */
export function createInfrastructureStack(
  config: EnvironmentConfig
): InfrastructureStack {
  const stage = $app.stage;

  // TODO: Implement Archive Store - S3 buckets for media storage
  // This will include:
  // - Main storage bucket with lifecycle policies
  // - Cache bucket for processed media
  // - Appropriate bucket policies and CORS configuration
  const archiveBucket = {
    name: `dphoto-${stage}-archive-placeholder`,
    arn: `arn:aws:s3:::dphoto-${stage}-archive-placeholder`,
  };

  const cacheBucket = {
    name: `dphoto-${stage}-cache-placeholder`,
    arn: `arn:aws:s3:::dphoto-${stage}-cache-placeholder`,
  };

  // TODO: Implement Catalog Store - DynamoDB table
  // This will include:
  // - Single-table design with appropriate indexes
  // - Point-in-time recovery enabled (for production)
  // - Appropriate capacity settings based on environment
  const catalogTable = {
    name: `dphoto-${stage}-catalog-placeholder`,
    arn: `arn:aws:dynamodb:eu-west-1:123456789012:table/dphoto-${stage}-catalog-placeholder`,
  };

  // TODO: Implement Archivist - SNS/SQS for async processing
  // This will include:
  // - SNS topic for job fanout
  // - SQS queue with dead-letter queue
  // - Appropriate message retention and visibility timeout
  const archiveTopic = {
    name: `dphoto-${stage}-archive-topic-placeholder`,
    arn: `arn:aws:sns:eu-west-1:123456789012:dphoto-${stage}-archive-topic-placeholder`,
  };

  const archiveQueue = {
    name: `dphoto-${stage}-archive-queue-placeholder`,
    arn: `arn:aws:sqs:eu-west-1:123456789012:dphoto-${stage}-archive-queue-placeholder`,
    url: `https://sqs.eu-west-1.amazonaws.com/123456789012/dphoto-${stage}-archive-queue-placeholder`,
  };

  // TODO: Implement CLI User Access - IAM users for CLI access
  // This will be migrated from cli-user-access-construct.ts

  // TODO: Implement Serverless Integration constructs
  // These provide SSM parameters for backward compatibility with any remaining
  // Serverless Framework deployments

  // Add stack outputs
  // Note: In SST, outputs are automatically exposed via $app.outputs
  
  return {
    archiveBucket,
    cacheBucket,
    catalogTable,
    archiveTopic,
    archiveQueue,
  };
}
