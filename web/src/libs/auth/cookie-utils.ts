import { serialize, parse } from 'cookie';

export interface TokenCookies {
  accessToken?: string;
  refreshToken?: string;
}

const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';

export function setTokenCookies(accessToken: string, refreshToken: string): string[] {
  const isProduction = process.env.NODE_ENV === 'production';
  
  const cookieOptions = {
    httpOnly: true,
    secure: isProduction,
    sameSite: 'strict' as const,
    path: '/',
    maxAge: 60 * 60 * 24 * 30, // 30 days
  };

  return [
    serialize(ACCESS_TOKEN_COOKIE, accessToken, {
      ...cookieOptions,
      maxAge: 60 * 60, // 1 hour for access token
    }),
    serialize(REFRESH_TOKEN_COOKIE, refreshToken, cookieOptions),
  ];
}

export function clearTokenCookies(): string[] {
  const cookieOptions = {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict' as const,
    path: '/',
    maxAge: 0,
  };

  return [
    serialize(ACCESS_TOKEN_COOKIE, '', cookieOptions),
    serialize(REFRESH_TOKEN_COOKIE, '', cookieOptions),
  ];
}

export function parseTokenCookies(cookieHeader?: string): TokenCookies {
  if (!cookieHeader) {
    return {};
  }

  const cookies = parse(cookieHeader);
  
  return {
    accessToken: cookies[ACCESS_TOKEN_COOKIE],
    refreshToken: cookies[REFRESH_TOKEN_COOKIE],
  };
}

export function getAccessTokenCookieName(): string {
  return ACCESS_TOKEN_COOKIE;
}

export function getRefreshTokenCookieName(): string {
  return REFRESH_TOKEN_COOKIE;
}
