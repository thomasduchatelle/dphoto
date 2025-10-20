'use client';

import { useEffect } from 'react';
import { setClientAccessToken, clearClientAccessToken } from '../libs/auth/token-context';

interface TokenInitializerProps {
  accessToken?: string;
  expiresAt?: number;
}

export function TokenInitializer({ accessToken, expiresAt }: TokenInitializerProps) {
  useEffect(() => {
    if (accessToken && expiresAt) {
      setClientAccessToken({ accessToken, expiresAt });
    } else {
      clearClientAccessToken();
    }
  }, [accessToken, expiresAt]);

  return null;
}
