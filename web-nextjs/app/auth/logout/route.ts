import {NextRequest, NextResponse} from 'next/server';
import {
    ACCESS_TOKEN_COOKIE,
    getOidcConfigFromEnv,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    oidcConfig,
    REDIRECT_AFTER_LOGIN_COOKIE,
    REFRESH_TOKEN_COOKIE
} from '@/libs/security';
import {basePath, getOriginalOrigin} from '@/libs/requests';

const COOKIE_CLEAR_OPTS = {
    maxAge: 0,
    path: '/',
};

export async function GET(request: NextRequest) {
    const requestUrl = getOriginalOrigin(request);
    const oidcEnvConfig = getOidcConfigFromEnv();
    const config = await oidcConfig(oidcEnvConfig);

    const logoutUri = new URL(`${basePath}/auth/logout-success`, requestUrl).toString();
    const cognitoLogoutUrl = `${config.serverMetadata().issuer}/logout?client_id=${oidcEnvConfig.clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;

    const response = NextResponse.redirect(cognitoLogoutUrl, 307);

    response.cookies.set(ACCESS_TOKEN_COOKIE, '', COOKIE_CLEAR_OPTS);
    response.cookies.set(REFRESH_TOKEN_COOKIE, '', COOKIE_CLEAR_OPTS);
    response.cookies.set(OAUTH_STATE_COOKIE, '', COOKIE_CLEAR_OPTS);
    response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, '', COOKIE_CLEAR_OPTS);
    response.cookies.set(REDIRECT_AFTER_LOGIN_COOKIE, '', COOKIE_CLEAR_OPTS);

    return response;
}
