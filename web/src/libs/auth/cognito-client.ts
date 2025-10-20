import * as client from 'openid-client';

export interface CognitoConfig {
  userPoolId: string;
  clientId: string;
  clientSecret: string;
  domain: string;
  issuer: string;
}

let cachedConfig: client.Configuration | null = null;

export async function getCognitoConfig(): Promise<client.Configuration> {
  if (cachedConfig) {
    return cachedConfig;
  }

  const config: CognitoConfig = {
    userPoolId: process.env.COGNITO_USER_POOL_ID || '',
    clientId: process.env.COGNITO_CLIENT_ID || '',
    clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
    domain: process.env.COGNITO_DOMAIN || '',
    issuer: process.env.COGNITO_ISSUER || '',
  };

  const issuerUrl = new URL(config.issuer);
  
  cachedConfig = await client.discovery(
    issuerUrl,
    config.clientId,
    { client_secret: config.clientSecret },
    client.ClientSecretPost(config.clientSecret),
  );

  return cachedConfig;
}

export function generateCodeVerifier(): string {
  return client.randomPKCECodeVerifier();
}

export function generateState(): string {
  return client.randomState();
}

export function generateNonce(): string {
  return client.randomNonce();
}

export async function generateCodeChallenge(codeVerifier: string): Promise<string> {
  return await client.calculatePKCECodeChallenge(codeVerifier);
}

export function buildAuthorizationUrl(
  config: client.Configuration,
  redirectUri: string,
  state: string,
  nonce: string,
  codeChallenge: string
): URL {
  const authUrl = new URL(config.serverMetadata().authorization_endpoint!);
  
  authUrl.searchParams.set('client_id', config.clientMetadata().client_id);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', 'openid email profile');
  authUrl.searchParams.set('state', state);
  authUrl.searchParams.set('nonce', nonce);
  authUrl.searchParams.set('code_challenge', codeChallenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  
  return authUrl;
}
