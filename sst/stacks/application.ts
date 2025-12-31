/**
 * Application Stack for DPhoto SST
 * 
 * This stack contains all STATELESS resources that can be freely redeployed.
 * These resources process requests, serve the UI, and implement business logic.
 * 
 * Resources include:
 * - API Gateway (HTTP API v2)
 * - Lambda functions (API endpoints, workers)
 * - Cognito User Pool and configuration
 * - Web application (Waku SSR)
 * - Lambda authorizers
 * 
 * This mirrors the CDK ApplicationStack from deployments/cdk/lib/stacks/application-stack.ts
 */

import { EnvironmentConfig } from "../config/environments";
import { InfrastructureStack } from "./infrastructure";

export interface ApplicationStackProps {
  /** Environment-specific configuration */
  config: EnvironmentConfig;
  /** References to infrastructure resources */
  infrastructure: InfrastructureStack;
}

/**
 * Create the Application Stack
 * 
 * @param props Configuration and infrastructure references
 */
export function createApplicationStack(props: ApplicationStackProps): void {
  const { config, infrastructure } = props;
  const stage = $app.stage;

  // TODO: Implement API Gateway
  // This will include:
  // - HTTP API v2 with custom domain
  // - CORS configuration
  // - Access logs
  // - Integration with Lambda authorizer
  console.log(`Application Stack: Setting up API Gateway for ${stage}`);

  // TODO: Implement Cognito Stack
  // This will include:
  // - User Pool with Google SSO integration
  // - User Pool Client
  // - User Pool Domain (custom domain with certificate)
  // - User groups (admins, owners, visitors)
  console.log(`Application Stack: Setting up Cognito for ${stage}`);

  // TODO: Implement Lambda Authorizer
  // This will include:
  // - JWT validation
  // - User context enrichment
  // - Integration with Cognito
  console.log(`Application Stack: Setting up Lambda Authorizer for ${stage}`);

  // TODO: Implement Catalog Endpoints
  // This will include Lambda functions for:
  // - List albums
  // - Get album details
  // - Create/update/delete albums
  // - List medias
  // - Album sharing
  console.log(`Application Stack: Setting up Catalog Endpoints for ${stage}`);

  // TODO: Implement Archive Endpoints
  // This will include Lambda functions for:
  // - Get media
  // - Upload media
  // - Process media (thumbnails, compression)
  console.log(`Application Stack: Setting up Archive Endpoints for ${stage}`);

  // TODO: Implement User Endpoints
  // This will include Lambda functions for:
  // - Get current user
  // - User preferences
  // - User permissions
  console.log(`Application Stack: Setting up User Endpoints for ${stage}`);

  // TODO: Implement Authentication Endpoints
  // Note: May be deprecated after Cognito migration completes
  console.log(`Application Stack: Setting up Authentication Endpoints for ${stage}`);

  // TODO: Implement Version Endpoint
  // Simple endpoint to return API version information
  console.log(`Application Stack: Setting up Version Endpoint for ${stage}`);

  // TODO: Implement Waku Web UI
  // This will include:
  // - Lambda function for SSR
  // - S3 bucket for static assets
  // - CloudFront distribution (or API Gateway integration)
  // - Environment variables for OAuth2 config
  console.log(`Application Stack: Setting up Waku Web UI for ${stage}`);

  // Stack dependencies are implicit in SST through the use of resource references
  // No explicit dependency declaration needed like in CDK
  console.log(`Application Stack: Configuration complete for ${stage}`);
  console.log(`  - Archive Bucket: ${infrastructure.archiveBucket.name}`);
  console.log(`  - Cache Bucket: ${infrastructure.cacheBucket.name}`);
  console.log(`  - Catalog Table: ${infrastructure.catalogTable.name}`);
}
