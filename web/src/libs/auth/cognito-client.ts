import { Issuer, Client, generators } from 'openid-client';

interface CognitoConfig {
  userPoolId: string;
  clientId: string;
  clientSecret: string;
  domain: string;
  issuer: string;
  callbackUrl: string;
}

let cachedClient: Client | null = null;

export async function getCognitoClient(): Promise<Client> {
  if (cachedClient) {
    return cachedClient;
  }

  const config: CognitoConfig = {
    userPoolId: process.env.COGNITO_USER_POOL_ID || '',
    clientId: process.env.COGNITO_CLIENT_ID || '',
    clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
    domain: process.env.COGNITO_DOMAIN || '',
    issuer: process.env.COGNITO_ISSUER || '',
    callbackUrl: `${process.env.APP_URL || ''}/auth/callback`,
  };

  const issuer = await Issuer.discover(config.issuer);
  
  cachedClient = new issuer.Client({
    client_id: config.clientId,
    client_secret: config.clientSecret,
    redirect_uris: [config.callbackUrl],
    response_types: ['code'],
  });

  return cachedClient;
}

export function generateAuthorizationUrl(client: Client, state: string, nonce: string, codeVerifier: string): string {
  const codeChallenge = generators.codeChallenge(codeVerifier);
  
  return client.authorizationUrl({
    scope: 'openid email profile',
    state,
    nonce,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
  });
}

export { generators };
