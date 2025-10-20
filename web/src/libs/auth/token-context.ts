// Token context for managing access tokens on the client
// Tokens are read from cookies and cached in memory

import { getAccessTokenFromCookie, getTokenExpirationFromCookie } from './client-cookie-utils';

export interface TokenInfo {
  accessToken: string;
  expiresAt: number;
}

let clientAccessToken: TokenInfo | null = null;
let lastCookieCheck = 0;
const COOKIE_CHECK_INTERVAL = 5000; // Check cookies every 5 seconds

function loadTokenFromCookie(): TokenInfo | null {
  const accessToken = getAccessTokenFromCookie();
  const expiresAt = getTokenExpirationFromCookie();
  
  if (accessToken && expiresAt) {
    return { accessToken, expiresAt };
  }
  
  return null;
}

export function setClientAccessToken(token: TokenInfo): void {
  clientAccessToken = token;
}

export function getClientAccessToken(): string | null {
  // Periodically check if the cookie has been updated
  const now = Date.now();
  if (now - lastCookieCheck > COOKIE_CHECK_INTERVAL) {
    lastCookieCheck = now;
    const tokenFromCookie = loadTokenFromCookie();
    if (tokenFromCookie) {
      clientAccessToken = tokenFromCookie;
    }
  }

  if (!clientAccessToken) {
    // Try to load from cookie if not in memory
    clientAccessToken = loadTokenFromCookie();
  }

  if (!clientAccessToken) {
    return null;
  }

  // Check if token is expired
  if (Date.now() >= clientAccessToken.expiresAt) {
    clientAccessToken = null;
    return null;
  }

  return clientAccessToken.accessToken;
}

export function clearClientAccessToken(): void {
  clientAccessToken = null;
}
