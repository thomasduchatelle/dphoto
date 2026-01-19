import {decodeJWTPayload} from "./jwt-utils";
import {clearFullSession, loadSession, storeSession} from "./backend-store";
import {getOidcConfigFromEnv, oidcConfig} from "@/libs/security/oidc-config";
import * as client from 'openid-client';
import {ReadCookieStore, SetCookies} from "@/libs/nextjs-cookies";

interface AccessTokenClaims {
    expiresAt: Date
    isOwner: boolean
    aboutToExpire: boolean
}

export type AccessToken = AccessTokenClaims & {
    accessToken: string;
}

function expiresInMoreThanFiveMinutes(expiresAt: Date): boolean {
    return expiresAt.getTime() - Date.now() > 5 * 60 * 1000;
}

export async function parseCurrentAccessToken(accessToken?: string): Promise<AccessToken | null> {
    const claims = readAccessTokenClaims(accessToken);
    if (!accessToken || !claims) {
        return null;
    }

    return {
        ...claims,
        accessToken,
        aboutToExpire: !expiresInMoreThanFiveMinutes(claims.expiresAt),
    }
}

export async function refreshSession(cookiesStore: ReadCookieStore): Promise<{
    cookies: SetCookies,
    success: boolean,
}> {
    const session = loadSession(cookiesStore)
    if (!session.refreshToken || !session.idToken) {
        // Unexpected path: refreshing without refresh token shows an error in the flow.
        return {success: false, cookies: clearFullSession()}
    }

    const refreshedTokens = await refreshAccessToken(session.refreshToken);
    if (!refreshedTokens) {
        // Unexpected path: refreshing without refresh token shows an error in the flow.
        return {success: false, cookies: clearFullSession()}
    }

    return {
        success: true,
        cookies: storeSession(refreshedTokens),
    }
}


function readAccessTokenClaims(token?: string): AccessTokenClaims | null {
    if (!token) {
        return null;
    }

    const payload = decodeJWTPayload(token);
    if (!payload || !payload.exp) {
        return null;
    }

    const scopes = payload.Scopes?.split(' ') ?? []
    let expiresAt = payload?.exp ? new Date(payload.exp * 1000) : new Date();
    return {
        expiresAt,
        isOwner: scopes.some((scope: string) => scope.startsWith('owner:')),
        aboutToExpire: !expiresInMoreThanFiveMinutes(expiresAt),
    }
}

async function refreshAccessToken(refreshToken ?: string) {
    if (!refreshToken) {
        return null;
    }

    const newTokens = await oidcRefresh(refreshToken);
    if (!newTokens || !newTokens.access_token) {
        return null;
    }
    return {
        accessToken: newTokens.access_token,
        accessTokenExpiresIn: newTokens.expires_in || 3600,
        refreshToken: newTokens.refresh_token || refreshToken,
        idToken: newTokens.id_token || '',
    };
}

async function oidcRefresh(refreshToken: string): Promise<client.TokenEndpointResponse | null> {
    try {
        const config = await oidcConfig(getOidcConfigFromEnv());
        return await client.refreshTokenGrant(config, refreshToken);

    } catch (error) {
        console.error('Failed to refresh token:', error instanceof Error ? error.message : 'Unknown error');
        return null;
    }
}