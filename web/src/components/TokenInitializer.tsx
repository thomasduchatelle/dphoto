'use client';

import { useEffect } from 'react';
import { getClientAccessToken } from '../libs/auth/token-context';

/**
 * TokenInitializer ensures tokens are loaded from cookies on mount.
 * The token context will automatically read from cookies when getClientAccessToken is called,
 * but this component triggers an initial check on mount.
 */
export function TokenInitializer() {
  useEffect(() => {
    // Trigger initial token load from cookies
    getClientAccessToken();
  }, []);

  return null;
}
