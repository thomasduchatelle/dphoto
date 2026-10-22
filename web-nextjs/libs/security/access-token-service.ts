import {decodeJWTPayload} from "./jwt-utils";
import {clearFullSession, loadSession, storeSession} from "./backend-store";
import {getOidcConfigFromEnv, isLocalAuthBypassed, oidcConfig} from "@/libs/security/oidc-config";
import * as client from 'openid-client';
import {ReadCookieStore, SetCookies} from "@/libs/nextjs-cookies";

const LOCAL_DEV_SESSION_TTL_SECONDS = 30 * 24 * 3600;

function base64UrlEncode(value: string): string {
    return Buffer.from(value, 'utf-8').toString('base64')
        .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

/**
 * Builds an unsigned JWT impersonating the local dev owner, matching the LOCAL_DEV_USER used server-side when
 * isLocalAuthBypassed() is true. Wiremock does not verify the signature, and the long expiry keeps the client
 * from re-hitting /auth/refresh on every call.
 */
function localDevStubSession(): Parameters<typeof storeSession>[0] {
    const header = base64UrlEncode(JSON.stringify({alg: 'none', typ: 'JWT'}));
    const payload = base64UrlEncode(JSON.stringify({
        sub: 'tony@stark.com',
        Scopes: 'owner:tony@stark.com',
        exp: Math.floor(Date.now() / 1000) + LOCAL_DEV_SESSION_TTL_SECONDS,
    }));
    const stubToken = `${header}.${payload}.`;
    return {
        accessToken: stubToken,
        accessTokenExpiresIn: LOCAL_DEV_SESSION_TTL_SECONDS,
        refreshToken: 'local-dev',
        idToken: stubToken,
    };
}

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
    if (!accessToken) {
        return null;
    }
    const claims = readAccessTokenClaims(accessToken);
    if(!claims) {
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
    console.log("AccessTokenService - refreshSession")

    if (isLocalAuthBypassed()) {
        console.log("AccessTokenService - local auth bypass: minting a stub session")
        return {success: true, cookies: storeSession(localDevStubSession())}
    }

    const session = loadSession(cookiesStore)
    if (!session.refreshToken || !session.idToken) {
        console.log("AccessTokenService - cannot refresh: no refreshToken, and/or no identity token")
        // Unexpected path: refreshing without refresh token shows an error in the flow.
        return {success: false, cookies: clearFullSession()}
    }

    const refreshedTokens = await refreshAccessToken(session.refreshToken);
    if (!refreshedTokens) {
        console.log("AccessTokenService - refresh failed, no access token returned")
        // Unexpected path: refreshing without refresh token shows an error in the flow.
        return {success: false, cookies: clearFullSession()}
    }

    return {
        success: true,
        cookies: storeSession(refreshedTokens),
    }
}


function readAccessTokenClaims(token: string): AccessTokenClaims | null {
    const payload = decodeJWTPayload(token);
    if (!payload || !payload.exp) {
        return null;
    }

    const scopes = payload.Scopes?.split(' ') ?? []
    let expiresAt = payload?.exp ? new Date(payload.exp * 1000) : new Date();
    return {
        expiresAt,
        isOwner: true, // TODO Original is `scopes.some((scope: string) => scope.startsWith('owner:')),` but the scope is not loaded.
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