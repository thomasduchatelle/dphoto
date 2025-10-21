// Server-side token management utilities
// These functions handle tokens in cookies during SSR

import { parse, serialize } from 'cookie';
import { ACCESS_TOKEN_COOKIE, ACCESS_TOKEN_COOKIE_CLIENT, REFRESH_TOKEN_COOKIE, OAUTH_STATE_COOKIE } from './consts';

// Cookie store helper that works with headers
let currentCookieHeader: string | undefined;
let cookiesToSet: string[] = [];

export function initCookieStore(cookieHeader: string | undefined) {
  currentCookieHeader = cookieHeader;
  cookiesToSet = [];
}

export function getCookiesToSet(): string[] {
  return cookiesToSet;
}

function getCookie(name: string): string | undefined {
  if (!currentCookieHeader) return undefined;
  const cookies = parse(currentCookieHeader);
  return cookies[name];
}

function setCookie(name: string, value: string, options: any) {
  cookiesToSet.push(serialize(name, value, options));
}

function deleteCookie(name: string) {
  cookiesToSet.push(serialize(name, '', { maxAge: 0, path: '/' }));
}

export interface AccessToken {
  value: string;
  expiresAt: number;
}

export interface RefreshToken {
  value: string;
}

export interface ServerSession {
  accessToken: AccessToken;
  refreshToken: RefreshToken;
}

export interface OAuthState {
  codeVerifier: string;
  nonce: string;
  state: string;
  returnUrl: string;
}

/**
 * Load the current session from cookies (server-side only)
 * Returns null if no valid session exists
 */
export function loadServerSession(): ServerSession | null {
  const accessTokenValue = getCookie(ACCESS_TOKEN_COOKIE);
  const refreshTokenValue = getCookie(REFRESH_TOKEN_COOKIE);
  
  if (!accessTokenValue || !refreshTokenValue) {
    return null;
  }

  // Decode access token to get expiration
  try {
    const payload = JSON.parse(
      Buffer.from(accessTokenValue.split('.')[1], 'base64').toString()
    );
    
    return {
      accessToken: {
        value: accessTokenValue,
        expiresAt: payload.exp * 1000, // Convert to milliseconds
      },
      refreshToken: {
        value: refreshTokenValue,
      },
    };
  } catch {
    return null;
  }
}

/**
 * Store a new session in cookies (server-side only)
 */
export function storeServerSession(session: ServerSession): void {
  const isProduction = process.env.NODE_ENV === 'production';
  
  const secureCookieOptions = {
    httpOnly: true,
    secure: isProduction,
    sameSite: 'strict' as const,
    path: '/',
  };

  const clientCookieOptions = {
    httpOnly: false, // Allow JavaScript access for Authorization header
    secure: isProduction,
    sameSite: 'strict' as const,
    path: '/',
  };

  // Set HttpOnly cookies for server-side use
  setCookie(ACCESS_TOKEN_COOKIE, session.accessToken.value, {
    ...secureCookieOptions,
    maxAge: 60 * 60, // 1 hour
  });

  setCookie(REFRESH_TOKEN_COOKIE, session.refreshToken.value, {
    ...secureCookieOptions,
    maxAge: 60 * 60 * 24 * 30, // 30 days
  });

  // Set client-accessible cookie for Authorization header
  setCookie(ACCESS_TOKEN_COOKIE_CLIENT, session.accessToken.value, {
    ...clientCookieOptions,
    maxAge: 60 * 60, // 1 hour
  });
}

/**
 * Clear the session cookies (server-side only)
 */
export function clearSession(): void {
  deleteCookie(ACCESS_TOKEN_COOKIE);
  deleteCookie(ACCESS_TOKEN_COOKIE_CLIENT);
  deleteCookie(REFRESH_TOKEN_COOKIE);
  deleteCookie(OAUTH_STATE_COOKIE);
}

/**
 * Store OAuth state in a cookie during authentication flow
 */
export function storeOAuthState(state: OAuthState): void {
  const isProduction = process.env.NODE_ENV === 'production';
  
  setCookie(OAUTH_STATE_COOKIE, JSON.stringify(state), {
    httpOnly: true,
    secure: isProduction,
    sameSite: 'strict' as const,
    path: '/',
    maxAge: 10 * 60, // 10 minutes
  });
}

/**
 * Load OAuth state from cookie
 */
export function loadOAuthState(): OAuthState | null {
  const stateValue = getCookie(OAUTH_STATE_COOKIE);
  
  if (!stateValue) {
    return null;
  }

  try {
    return JSON.parse(stateValue);
  } catch {
    return null;
  }
}

/**
 * Clear OAuth state cookie
 */
export function clearOAuthState(): void {
  deleteCookie(OAUTH_STATE_COOKIE);
}

/**
 * Check if access token is expired or expiring soon
 */
export function isTokenExpiringSoon(token: AccessToken, thresholdMinutes: number = 5): boolean {
  const thresholdTime = Date.now() + (thresholdMinutes * 60 * 1000);
  return token.expiresAt < thresholdTime;
}
