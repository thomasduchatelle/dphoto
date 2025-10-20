// Server-side utilities for handling tokens during SSR
// These functions run on the server and extract tokens from request cookies

import { parse } from 'cookie';

const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';

export interface ServerTokens {
  accessToken?: string;
  refreshToken?: string;
}

/**
 * Extract tokens from request cookies (server-side only)
 * @param cookieHeader - The Cookie header from the request
 */
export function extractTokensFromCookies(cookieHeader?: string): ServerTokens {
  if (!cookieHeader) {
    return {};
  }

  const cookies = parse(cookieHeader);
  
  return {
    accessToken: cookies[ACCESS_TOKEN_COOKIE],
    refreshToken: cookies[REFRESH_TOKEN_COOKIE],
  };
}

/**
 * Decode JWT token to get expiration time (server-side only)
 * @param token - The JWT token
 */
export function getTokenExpiration(token: string): number | null {
  try {
    const payload = JSON.parse(
      Buffer.from(token.split('.')[1], 'base64').toString()
    );
    return payload.exp * 1000; // Convert to milliseconds
  } catch {
    return null;
  }
}

/**
 * Check if a token is valid and not expired (server-side only)
 * @param token - The JWT token to validate
 */
export function isTokenValid(token: string): boolean {
  const expiration = getTokenExpiration(token);
  if (!expiration) {
    return false;
  }
  
  return Date.now() < expiration;
}
