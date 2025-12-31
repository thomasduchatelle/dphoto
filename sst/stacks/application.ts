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
 * 
 * TODO: Implement the following components in future steps:
 * 
 * 1. API Gateway - HTTP API v2 with custom domain, CORS, access logs
 * 2. Cognito Stack - User Pool with Google SSO, User Pool Client, custom domain, user groups
 * 3. Lambda Authorizer - JWT validation, user context enrichment
 * 4. Catalog Endpoints - List/create/update/delete albums, list medias, album sharing
 * 5. Archive Endpoints - Get/upload/process media (thumbnails, compression)
 * 6. User Endpoints - Current user, preferences, permissions
 * 7. Authentication Endpoints - Legacy endpoints (to be deprecated after Cognito migration)
 * 8. Version Endpoint - API version information
 * 9. Waku Web UI - SSR Lambda, S3 static assets, OAuth2 config
 */
export function createApplicationStack(props: ApplicationStackProps): void {
  const { config, infrastructure } = props;
  const stage = $app.stage;

  // Log scaffolding summary
  console.log(`\n=== Application Stack Scaffolding for ${stage} ===`);
  console.log(`Components to be implemented in future steps:`);
  console.log(`  • API Gateway (HTTP v2 with custom domain)`);
  console.log(`  • Cognito User Pool (with Google SSO)`);
  console.log(`  • Lambda Authorizer (JWT validation)`);
  console.log(`  • Catalog Endpoints (albums, medias, sharing)`);
  console.log(`  • Archive Endpoints (media get/upload/process)`);
  console.log(`  • User Endpoints (current user, preferences)`);
  console.log(`  • Authentication Endpoints (to be deprecated)`);
  console.log(`  • Version Endpoint (API version info)`);
  console.log(`  • Waku Web UI (SSR with static assets)`);
  console.log(`\nInfrastructure references:`);
  console.log(`  • Archive Bucket: ${infrastructure.archiveBucket.name}`);
  console.log(`  • Cache Bucket: ${infrastructure.cacheBucket.name}`);
  console.log(`  • Catalog Table: ${infrastructure.catalogTable.name}`);
  console.log(`===================================================\n`);
}
