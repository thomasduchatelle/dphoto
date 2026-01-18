import {decodeJWTPayload} from "./jwt-utils";
import {clearFullSession, loadSession, storeSession} from "./backend-store";
import {getOidcConfigFromEnv, oidcConfig} from "@/libs/security/oidc-config";
import * as client from 'openid-client';

interface AccessTokenClaims {
    expiresAt: Date
    isOwner: boolean
}

export type AccessToken = AccessTokenClaims & {
    accessToken: string;
}

function expiresInMoreThanFiveMinutes(accessToken: AccessToken): boolean {
    return accessToken.expiresAt.getTime() - Date.now() > 5 * 60 * 1000;
}

async function parseCurrentAccessToken(accessToken?: string): Promise<AccessToken | null> {
    const claims = readAccessTokenClaims(accessToken);
    if (!accessToken || !claims) {
        return null;
    }

    return {
        ...claims,
        accessToken,
    }
}

export async function getValidAccessToken(): Promise<{ accessToken: AccessToken, idToken: string } | null> {
    const session = await loadSession()
    if (!session.refreshToken || !session.idToken) {
        // await clearFullSession() // clear any partial session
        return null
    }

    const accessToken = await parseCurrentAccessToken(session.accessToken);
    if (accessToken && expiresInMoreThanFiveMinutes(accessToken)) {
        return {accessToken, idToken: session.idToken};
    }

    const refreshedToken = await refreshAccessToken(session.refreshToken);
    if (!refreshedToken) {
        await clearFullSession() // disconnected
        return null
    }

    return {accessToken: refreshedToken, idToken: session.idToken}
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
    return {
        expiresAt: payload?.exp ? new Date(payload.exp * 1000) : new Date(),
        isOwner: scopes.some((scope: string) => scope.startsWith('owner:')),
    }
}

async function refreshAccessToken(refreshToken ?: string): Promise<AccessToken | null> {
    if (!refreshToken) {
        return null;
    }

    const newTokens = await oidcRefresh(refreshToken);
    if (newTokens && newTokens.access_token) {
        await storeSession({
            accessToken: newTokens.access_token,
            accessTokenExpiresIn: newTokens.expires_in || 3600,
            refreshToken: newTokens.refresh_token || refreshToken,
            idToken: newTokens.id_token || '',
        })

        const claims = readAccessTokenClaims(newTokens.access_token)
        if (claims) {
            return {
                ...claims,
                accessToken: newTokens.access_token,
            }
        }
    }

    return null;
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