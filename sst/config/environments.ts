/**
 * Environment Configuration for DPhoto SST
 * 
 * This module provides environment-specific configuration for different deployment stages.
 * It mirrors the configuration from deployments/cdk/lib/config/environments.ts
 * but is adapted for SST's stage-based approach.
 */

export interface EnvironmentConfig {
  /** Whether this is a production environment */
  production: boolean;
  /** Pairs of access keys used by the CLI */
  cliAccessKeys?: string[];
  /** Domain registered in Route53 */
  rootDomain: string;
  /** (sub)Domain used to expose the root of the application */
  domainName: string;
  /** (subsub)Domain used for Cognito hosted UI (authentication) */
  cognitoDomainName: string;
  /** Other URLs to allow redirection to after login (Cognito Hosted UI), and logout */
  cognitoExtraRedirectURLs: string[];
  /** Email used for SSL certificate registration automated by let's encrypt */
  certificateEmail: string;
  /** OAuth2 Client ID for Google SSO, used by Cognito */
  googleLoginClientId: string;
}

const environments: Record<string, EnvironmentConfig> = {
  live: {
    production: true,
    cliAccessKeys: ["2025-07"],
    rootDomain: "duchatelle.me",
    domainName: "dphoto.duchatelle.me",
    cognitoDomainName: "login.dphoto.duchatelle.me",
    cognitoExtraRedirectURLs: [],
    certificateEmail: "duchatelle.thomas@gmail.com",
    googleLoginClientId: "841197197570-1o0or8ioo9c4m31405q2h2k8hvdb5enh.apps.googleusercontent.com",
  },
  next: {
    production: false,
    cliAccessKeys: ["2025-07"],
    rootDomain: "duchatelle.me",
    domainName: "next.duchatelle.me",
    cognitoDomainName: "login.next.duchatelle.me",
    cognitoExtraRedirectURLs: ["http://localhost:3000"],
    certificateEmail: "duchatelle.thomas@gmail.com",
    googleLoginClientId: "841197197570-7hlq9e86d6u37eoq8nsd8af4aaisl5gb.apps.googleusercontent.com",
  },
  test: {
    production: false,  // Test environment should not have production-level settings
    cliAccessKeys: ["2024-04"],
    rootDomain: "example.com",
    domainName: "dphoto.example.com",
    cognitoDomainName: "login.dphoto.example.com",
    cognitoExtraRedirectURLs: ["http://localhost:3210"],
    certificateEmail: "dphoto@example.com",
    googleLoginClientId: "test-google-client-id",
  },
};

/**
 * Get configuration for a specific environment/stage
 * @param stage The stage name (next, live, test)
 * @returns The environment configuration
 * @throws Error if the stage is not recognized
 */
export function getEnvironmentConfig(stage: string): EnvironmentConfig {
  const config = environments[stage];
  
  if (!config) {
    throw new Error(
      `Unknown environment: ${stage}. Available: ${Object.keys(environments).join(", ")}`
    );
  }
  
  console.log(`Loading configuration for stage: ${stage}`);
  return config;
}
