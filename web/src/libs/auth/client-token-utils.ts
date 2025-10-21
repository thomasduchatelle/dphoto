// Client-side token utilities
// These functions read tokens from cookies on the client

import Cookies from 'js-cookie';
import { ACCESS_TOKEN_COOKIE_CLIENT } from './consts';

export interface AccessToken {
  value: string;
  expiresAt: number;
}

export interface ClientSession {
  accessToken: AccessToken;
}

/**
 * Load the access token from client-accessible cookie
 * Returns null if no valid token exists
 */
export function loadClientSession(): ClientSession | null {
  const tokenValue = Cookies.get(ACCESS_TOKEN_COOKIE_CLIENT);
  
  if (!tokenValue) {
    return null;
  }

  try {
    // Decode token to get expiration
    const payload = JSON.parse(atob(tokenValue.split('.')[1]));
    const expiresAt = payload.exp * 1000; // Convert to milliseconds
    
    // Check if expired
    if (expiresAt < Date.now()) {
      return null;
    }

    return {
      accessToken: {
        value: tokenValue,
        expiresAt,
      },
    };
  } catch {
    return null;
  }
}
