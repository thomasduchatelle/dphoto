import {cookies, headers} from 'next/headers';
import {
    ACCESS_TOKEN_COOKIE,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_NONCE_COOKIE,
    OAUTH_STATE_COOKIE,
    REDIRECT_AFTER_LOGIN_COOKIE,
    REFRESH_TOKEN_COOKIE,
    USER_INFO_COOKIE
} from './constants';
import {basePath, getOidcConfigFromEnv, oidcConfig} from './oidc-config';

/**
 * Gets the origin URL from headers (for use in server components and routes)
 */
async function getOriginFromHeaders(): Promise<string> {
    const headersList = await headers();
    const host = headersList.get('host') || 'localhost:3000';
    const proto = headersList.get('x-forwarded-proto') || 'http';
    return `${proto}://${host}`;
}

/**
 * Generates the Cognito logout URL with the logout_uri parameter
 */
export async function getLogoutUrl(): Promise<string> {
    const oidcEnvConfig = getOidcConfigFromEnv();
    const config = await oidcConfig(oidcEnvConfig);
    const origin = await getOriginFromHeaders();

    const logoutUri = new URL(`${basePath}/auth/logout`, origin).toString();
    return `${config.serverMetadata().issuer}/logout?client_id=${oidcEnvConfig.clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;
}

/**
 * Clears all authentication-related cookies
 */
export async function clearAuthCookies(): Promise<void> {
    const cookieStore = await cookies();
    const cookieOptions = {
        maxAge: 0,
        path: '/',
    };

    cookieStore.set(ACCESS_TOKEN_COOKIE, '', cookieOptions);
    cookieStore.set(REFRESH_TOKEN_COOKIE, '', cookieOptions);
    cookieStore.set(OAUTH_STATE_COOKIE, '', cookieOptions);
    cookieStore.set(OAUTH_CODE_VERIFIER_COOKIE, '', cookieOptions);
    cookieStore.set(OAUTH_NONCE_COOKIE, '', cookieOptions);
    cookieStore.set(REDIRECT_AFTER_LOGIN_COOKIE, '', cookieOptions);
    cookieStore.set(USER_INFO_COOKIE, '', cookieOptions);
}

