// Custom Lambda handler wrapper for authentication
// This file wraps the Waku handler and adds authentication endpoints

import * as client from 'openid-client';
import { serialize, parse } from 'cookie';

// Configuration from environment variables
function getConfig() {
  return {
    userPoolId: process.env.COGNITO_USER_POOL_ID || '',
    clientId: process.env.COGNITO_CLIENT_ID || '',
    clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
    domain: process.env.COGNITO_DOMAIN || '',
    issuer: process.env.COGNITO_ISSUER || '',
  };
}

function getAppUrl(event) {
  // Try to get from environment first
  if (process.env.APP_URL) {
    return process.env.APP_URL;
  }
  
  // Construct from request headers
  const headers = event.headers || {};
  const host = headers.host || headers['x-forwarded-host'] || 'localhost:3000';
  const proto = headers['x-forwarded-proto'] || 'https';
  
  return `${proto}://${host}`;
}

// Cookie names
const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';

// In-memory session store (simple implementation for MVP)
const sessions = new Map();
const SESSION_TTL = 10 * 60 * 1000; // 10 minutes

// Cached config
let cachedConfig = null;

async function getCognitoConfiguration() {
  if (cachedConfig) {
    return cachedConfig;
  }

  const config = getConfig();
  const issuerUrl = new URL(config.issuer);
  
  cachedConfig = await client.discovery(
    issuerUrl,
    config.clientId,
    { client_secret: config.clientSecret },
    client.ClientSecretPost(config.clientSecret),
  );

  return cachedConfig;
}

function setTokenCookies(accessToken, refreshToken) {
  const isProduction = process.env.NODE_ENV === 'production';
  
  const secureCookieOptions = {
    httpOnly: true,
    secure: isProduction,
    sameSite: 'strict',
    path: '/',
  };

  const clientCookieOptions = {
    httpOnly: false, // Allow JavaScript access
    secure: isProduction,
    sameSite: 'strict',
    path: '/',
  };

  return [
    // HttpOnly cookie for server-side use (most secure)
    serialize(ACCESS_TOKEN_COOKIE, accessToken, {
      ...secureCookieOptions,
      maxAge: 60 * 60, // 1 hour
    }),
    // Non-HttpOnly cookie for client-side JavaScript (needed for Authorization header)
    serialize(`${ACCESS_TOKEN_COOKIE}-client`, accessToken, {
      ...clientCookieOptions,
      maxAge: 60 * 60, // 1 hour
    }),
    serialize(REFRESH_TOKEN_COOKIE, refreshToken, {
      ...secureCookieOptions,
      maxAge: 60 * 60 * 24 * 30, // 30 days
    }),
  ];
}

function clearTokenCookies() {
  const secureCookieOptions = {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    path: '/',
    maxAge: 0,
  };

  const clientCookieOptions = {
    httpOnly: false,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    path: '/',
    maxAge: 0,
  };

  return [
    serialize(ACCESS_TOKEN_COOKIE, '', secureCookieOptions),
    serialize(`${ACCESS_TOKEN_COOKIE}-client`, '', clientCookieOptions),
    serialize(REFRESH_TOKEN_COOKIE, '', secureCookieOptions),
  ];
}

async function handleAuthLogin(event) {
  try {
    const appUrl = getAppUrl(event);
    const config = await getCognitoConfiguration();
    
    // Generate security parameters
    const state = client.randomState();
    const nonce = client.randomNonce();
    const codeVerifier = client.randomPKCECodeVerifier();
    const codeChallenge = await client.calculatePKCECodeChallenge(codeVerifier);
    
    // Get original URL from query parameter or default to home
    const originalUrl = event.queryStringParameters?.returnUrl || '/';
    
    // Store session data
    sessions.set(state, {
      originalUrl,
      nonce,
      codeVerifier,
      createdAt: Date.now(),
    });
    
    // Auto-cleanup after TTL
    setTimeout(() => sessions.delete(state), SESSION_TTL);
    
    // Build authorization URL
    const redirectUri = `${appUrl}/auth/callback`;
    const authUrl = new URL(config.serverMetadata().authorization_endpoint);
    
    authUrl.searchParams.set('client_id', config.clientMetadata().client_id);
    authUrl.searchParams.set('redirect_uri', redirectUri);
    authUrl.searchParams.set('response_type', 'code');
    authUrl.searchParams.set('scope', 'openid email profile');
    authUrl.searchParams.set('state', state);
    authUrl.searchParams.set('nonce', nonce);
    authUrl.searchParams.set('code_challenge', codeChallenge);
    authUrl.searchParams.set('code_challenge_method', 'S256');
    
    return {
      statusCode: 302,
      headers: {
        'Location': authUrl.toString(),
      },
      body: '',
    };
  } catch (error) {
    console.error('Login error:', error);
    return {
      statusCode: 500,
      headers: {
        'Content-Type': 'text/plain',
      },
      body: 'Authentication error occurred',
    };
  }
}

async function handleAuthCallback(event) {
  try {
    const { code, state } = event.queryStringParameters || {};
    
    if (!code || !state) {
      return {
        statusCode: 400,
        headers: {
          'Content-Type': 'text/plain',
        },
        body: 'Missing code or state parameter',
      };
    }
    
    const session = sessions.get(state);
    
    if (!session) {
      return {
        statusCode: 400,
        headers: {
          'Content-Type': 'text/plain',
        },
        body: 'Invalid or expired session',
      };
    }
    
    // Check if session has expired
    if (Date.now() - session.createdAt > SESSION_TTL) {
      sessions.delete(state);
      return {
        statusCode: 400,
        headers: {
          'Content-Type': 'text/plain',
        },
        body: 'Session expired',
      };
    }
    
    // Exchange code for tokens
    const appUrl = getAppUrl(event);
    const config = await getCognitoConfiguration();
    const callbackUrl = `${appUrl}/auth/callback`;
    
    // Build current URL for token exchange
    const currentUrl = new URL(callbackUrl);
    currentUrl.searchParams.set('code', code);
    currentUrl.searchParams.set('state', state);
    
    const tokens = await client.authorizationCodeGrant(
      config,
      currentUrl,
      {
        pkceCodeVerifier: session.codeVerifier,
        expectedNonce: session.nonce,
      }
    );
    
    // Clean up session
    sessions.delete(state);
    
    // Set cookies
    const cookies = setTokenCookies(
      tokens.access_token,
      tokens.refresh_token
    );
    
    return {
      statusCode: 302,
      headers: {
        'Location': session.originalUrl,
        'Set-Cookie': cookies.join(', '),
      },
      body: '',
    };
  } catch (error) {
    console.error('Callback error:', error);
    return {
      statusCode: 500,
      headers: {
        'Content-Type': 'text/plain',
      },
      body: 'Authentication callback error occurred',
    };
  }
}

async function handleAuthLogout(event) {
  try {
    const cookies = clearTokenCookies();
    
    return {
      statusCode: 302,
      headers: {
        'Location': '/auth/login',
        'Set-Cookie': cookies.join(', '),
      },
      body: '',
    };
  } catch (error) {
    console.error('Logout error:', error);
    return {
      statusCode: 500,
      headers: {
        'Content-Type': 'text/plain',
      },
      body: 'Logout error occurred',
    };
  }
}

// Import the Waku handler
import { handler as wakuHandler } from './server/index.js';

export async function handler(event, context) {
  const path = event.rawPath || event.requestContext?.http?.path || '/';
  
  // Handle authentication routes
  if (path === '/auth/login') {
    return handleAuthLogin(event);
  }
  
  if (path === '/auth/callback') {
    return handleAuthCallback(event);
  }
  
  if (path === '/auth/logout') {
    return handleAuthLogout(event);
  }
  
  // Pass all other requests to Waku handler
  return wakuHandler(event, context);
}
