import { getCognitoClient } from './cognito-client';
import type { TokenSet } from 'openid-client';

export interface DecodedToken {
  sub: string;
  email: string;
  'cognito:groups'?: string[];
  exp: number;
  iat: number;
}

export async function validateAccessToken(accessToken: string): Promise<DecodedToken | null> {
  try {
    const client = await getCognitoClient();
    const tokenSet = new client.TokenSet({ access_token: accessToken });
    
    // Validate token signature and expiration using openid-client
    const userinfo = await client.userinfo(tokenSet);
    
    // Decode the token to get claims
    const decoded = JSON.parse(
      Buffer.from(accessToken.split('.')[1], 'base64').toString()
    ) as DecodedToken;
    
    // Check if token is expired
    if (decoded.exp * 1000 < Date.now()) {
      return null;
    }

    return decoded;
  } catch (error) {
    console.error('Token validation failed:', error);
    return null;
  }
}

export function isTokenExpiringSoon(accessToken: string, thresholdMinutes: number = 5): boolean {
  try {
    const decoded = JSON.parse(
      Buffer.from(accessToken.split('.')[1], 'base64').toString()
    ) as DecodedToken;
    
    const expirationTime = decoded.exp * 1000;
    const thresholdTime = Date.now() + (thresholdMinutes * 60 * 1000);
    
    return expirationTime < thresholdTime;
  } catch {
    return true;
  }
}

export async function exchangeCodeForTokens(
  code: string,
  codeVerifier: string,
  redirectUri: string
): Promise<TokenSet> {
  const client = await getCognitoClient();
  
  return await client.callback(redirectUri, { code }, { code_verifier: codeVerifier });
}
