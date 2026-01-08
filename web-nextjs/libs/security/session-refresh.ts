import * as client from 'openid-client';
import {decodeJWTPayload, getOidcConfigFromEnv, oidcConfig} from '@/libs/security';

export interface SessionCookies {
    accessToken?: string;
    refreshToken?: string;
}

export interface SessionRefresh {
    status: 'active' | 'expired' | 'anonymous';
    newAccessToken?: {
        token: string;
        expiresIn: number;
    };
    newRefreshToken?: string;
}

function isTokenExpired(token: string): boolean {
    const payload = decodeJWTPayload(token);
    if (!payload || !payload.exp) {
        return true;
    }
    
    const now = Math.floor(Date.now() / 1000);
    return payload.exp <= now;
}

async function refreshAccessToken(refreshToken: string): Promise<client.TokenEndpointResponse | null> {
    try {
        const config = await oidcConfig(getOidcConfigFromEnv());
        const tokens = await client.refreshTokenGrant(config, refreshToken);
        return tokens;
    } catch (error) {
        console.error('Failed to refresh token:', error instanceof Error ? error.message : 'Unknown error');
        return null;
    }
}

export async function refreshSessionIfNecessary(cookies: SessionCookies): Promise<SessionRefresh> {
    // If no refresh token is available, check if access token is valid
    if (!cookies.refreshToken) {
        if (cookies.accessToken && !isTokenExpired(cookies.accessToken)) {
            return { status: 'active' };
        }
        return { status: 'anonymous' };
    }

    // If we have a valid access token, no refresh needed
    if (cookies.accessToken && !isTokenExpired(cookies.accessToken)) {
        return { status: 'active' };
    }

    // Access token is missing or expired, try to refresh
    const newTokens = await refreshAccessToken(cookies.refreshToken);
    
    if (newTokens && newTokens.access_token) {
        return {
            status: 'active',
            newAccessToken: {
                token: newTokens.access_token,
                expiresIn: newTokens.expires_in || 0,
            },
            newRefreshToken: newTokens.refresh_token,
        };
    }

    // Refresh failed
    return { status: 'expired' };
}
