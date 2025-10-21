// Server component that checks authentication and redirects if needed
import * as client from 'openid-client';
import { initCookieStore, loadServerSession, storeOAuthState, isTokenExpiringSoon, getCookiesToSet } from '../libs/auth/server-token-utils';

async function getCognitoConfig(): Promise<client.Configuration> {
  const issuerUrl = new URL(process.env.COGNITO_ISSUER || '');
  const clientId = process.env.COGNITO_CLIENT_ID || '';
  const clientSecret = process.env.COGNITO_CLIENT_SECRET || '';

  return await client.discovery(
    issuerUrl,
    clientId,
    { client_secret: clientSecret },
    client.ClientSecretPost(clientSecret)
  );
}

function getAppUrl(): string {
  return process.env.APP_URL || 'http://localhost:3000';
}

async function buildAuthUrl(returnUrl: string): Promise<string> {
  const config = await getCognitoConfig();
  const appUrl = getAppUrl();
  
  // Generate OAuth security parameters
  const state = client.randomState();
  const nonce = client.randomNonce();
  const codeVerifier = client.randomPKCECodeVerifier();
  const codeChallenge = await client.calculatePKCECodeChallenge(codeVerifier);
  
  // Store OAuth state in cookie
  storeOAuthState({
    codeVerifier,
    nonce,
    state,
    returnUrl,
  });
  
  // Build authorization URL
  const redirectUri = `${appUrl}/auth/callback`;
  const authUrl = new URL(config.serverMetadata().authorization_endpoint!);
  
  authUrl.searchParams.set('client_id', config.clientMetadata().client_id);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', 'openid email profile');
  authUrl.searchParams.set('state', state);
  authUrl.searchParams.set('nonce', nonce);
  authUrl.searchParams.set('code_challenge', codeChallenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  
  return authUrl.toString();
}

interface AuthProviderProps {
  children: React.ReactNode;
  cookieHeader?: string;
}

export async function AuthProvider({ children, cookieHeader }: AuthProviderProps) {
  // Initialize cookie store with current cookies
  initCookieStore(cookieHeader);
  
  // Check if we have a valid session
  const session = loadServerSession();
  
  if (!session) {
    // No session, need to redirect to Cognito
    // Return a client component that will handle the redirect
    const authUrl = await buildAuthUrl('/');
    return (
      <html>
        <head>
          <meta httpEquiv="refresh" content={`0;url=${authUrl}`} />
        </head>
        <body>
          <p>Redirecting to login...</p>
        </body>
      </html>
    );
  }
  
  // Check if token is expiring soon (would need refresh)
  if (isTokenExpiringSoon(session.accessToken)) {
    // TODO: Implement token refresh
    // For now, initiate re-authentication
    const authUrl = await buildAuthUrl('/');
    return (
      <html>
        <head>
          <meta httpEquiv="refresh" content={`0;url=${authUrl}`} />
        </head>
        <body>
          <p>Session expired. Redirecting to login...</p>
        </body>
      </html>
    );
  }
  
  // Valid session exists, render children
  return <>{children}</>;
}
