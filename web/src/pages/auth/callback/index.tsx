// OAuth callback page - completes the authentication flow
import * as client from 'openid-client';
import { initCookieStore, loadOAuthState, clearOAuthState, storeServerSession, getCookiesToSet } from '../../../libs/auth/server-token-utils';

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

interface CallbackPageProps {
  searchParams?: {
    code?: string;
    state?: string;
    error?: string;
    error_description?: string;
  };
  cookieHeader?: string;
}

export default async function AuthCallbackPage({ searchParams, cookieHeader }: CallbackPageProps) {
  // Initialize cookie store
  initCookieStore(cookieHeader);
  // Check for errors from Cognito
  if (searchParams?.error) {
    console.error('OAuth error:', searchParams.error, searchParams.error_description);
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <h1>Authentication Error</h1>
        <p>{searchParams.error_description || searchParams.error}</p>
        <a href="/">Return to Home</a>
      </div>
    );
  }

  const code = searchParams?.code;
  const state = searchParams?.state;

  if (!code || !state) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <h1>Invalid Request</h1>
        <p>Missing authorization code or state parameter.</p>
        <a href="/">Return to Home</a>
      </div>
    );
  }

  // Load OAuth state from cookie
  const oauthState = loadOAuthState();

  if (!oauthState || oauthState.state !== state) {
    clearOAuthState();
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <h1>Invalid State</h1>
        <p>OAuth state mismatch or expired session. Please try again.</p>
        <a href="/">Return to Home</a>
      </div>
    );
  }

  try {
    // Exchange authorization code for tokens
    const config = await getCognitoConfig();
    const appUrl = getAppUrl();
    const callbackUrl = `${appUrl}/auth/callback`;

    // Build current URL for token exchange
    const currentUrl = new URL(callbackUrl);
    currentUrl.searchParams.set('code', code);
    currentUrl.searchParams.set('state', state);

    const tokens = await client.authorizationCodeGrant(
      config,
      currentUrl,
      {
        pkceCodeVerifier: oauthState.codeVerifier,
        expectedNonce: oauthState.nonce,
      }
    );

    if (!tokens.access_token || !tokens.refresh_token) {
      throw new Error('Missing tokens in response');
    }

    // Decode access token to get expiration
    const payload = JSON.parse(
      Buffer.from(tokens.access_token.split('.')[1], 'base64').toString()
    );

    // Store session in cookies
    storeServerSession({
      accessToken: {
        value: tokens.access_token,
        expiresAt: payload.exp * 1000,
      },
      refreshToken: {
        value: tokens.refresh_token,
      },
    });

    // Clear OAuth state
    clearOAuthState();

    // Redirect to original URL with cookies set
    const redirectUrl = oauthState.returnUrl || '/';
    return (
      <html>
        <head>
          <meta httpEquiv="refresh" content={`0;url=${redirectUrl}`} />
          {getCookiesToSet().map((cookie, i) => (
            <meta key={i} httpEquiv="Set-Cookie" content={cookie} />
          ))}
        </head>
        <body>
          <p>Authentication successful. Redirecting...</p>
        </body>
      </html>
    );
  } catch (error) {
    console.error('Token exchange failed:', error);
    clearOAuthState();
    
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <h1>Authentication Failed</h1>
        <p>Failed to complete authentication. Please try again.</p>
        <a href="/">Return to Home</a>
      </div>
    );
  }
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
