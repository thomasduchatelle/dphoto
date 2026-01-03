import { NextRequest, NextResponse } from 'next/server';
import * as cookie from 'cookie';
import * as client from 'openid-client';
import {
    ACCESS_TOKEN_COOKIE,
    REFRESH_TOKEN_COOKIE,
    OAUTH_STATE_COOKIE,
    OAUTH_CODE_VERIFIER_COOKIE,
    BackendSession,
} from './lib/security/constants';
import { decodeJWTPayload, isOwnerFromJWT } from './lib/security/jwt-utils';

const USER_INFO_COOKIE = 'dphoto-user-info';

interface IDTokenPayload {
    name?: string;
    email?: string;
    picture?: string;
    exp?: number;
    [key: string]: any;
}

interface AccessTokenWithUserInfo {
    name?: string;
    email?: string;
    picture?: string;
    exp?: number;
    Scopes?: string;
    [key: string]: any;
}

interface UserInfo {
    name: string;
    email: string;
    picture?: string;
}

interface Cookies {
    accessToken?: string;
    refreshToken?: string;
    userInfo?: string;
    state?: string;
    codeVerifier?: string;
}

const COOKIE_OPTS: cookie.SerializeOptions = {
    maxAge: 3600,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'strict',
};

function readCookies(request: NextRequest): Cookies {
    const cookieHeader = request.headers.get('cookie') || '';
    const cookies = cookie.parse(cookieHeader);
    return {
        accessToken: cookies[ACCESS_TOKEN_COOKIE],
        refreshToken: cookies[REFRESH_TOKEN_COOKIE],
        userInfo: cookies[USER_INFO_COOKIE],
        state: cookies[OAUTH_STATE_COOKIE],
        codeVerifier: cookies[OAUTH_CODE_VERIFIER_COOKIE],
    };
}

type OpenIdConfig = {
    issuer: string;
    clientId: string;
    clientSecret: string;
};

async function oidcConfig({ issuer, clientId, clientSecret }: OpenIdConfig): Promise<client.Configuration> {
    // TODO: Cache the configuration
    return client.discovery(new URL(issuer), clientId, clientSecret);
}

function getOidcConfigFromEnv(): OpenIdConfig {
    return {
        issuer: process.env.COGNITO_ISSUER || '',
        clientId: process.env.COGNITO_CLIENT_ID || '',
        clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
    };
}

export async function middleware(request: NextRequest) {
    const { pathname, origin } = request.nextUrl;
    const cookies = readCookies(request);

    // Handle OAuth callback
    if (pathname === '/auth/callback') {
        const config = await oidcConfig(getOidcConfigFromEnv());
        const url = new URL(request.url);

        try {
            const tokens: client.TokenEndpointResponse = await client.authorizationCodeGrant(
                config,
                url,
                {
                    pkceCodeVerifier: cookies.codeVerifier,
                    expectedState: cookies.state,
                }
            );

            let userInfo: UserInfo = {
                name: '',
                email: '',
            };

            if (tokens.id_token) {
                const idTokenPayload = decodeJWTPayload(tokens.id_token) as IDTokenPayload | null;
                if (idTokenPayload) {
                    userInfo = {
                        name: idTokenPayload.name || '',
                        email: idTokenPayload.email || '',
                        picture: idTokenPayload.picture,
                    };
                }
            }

            const response = NextResponse.redirect(new URL('/', request.url));

            response.cookies.set(ACCESS_TOKEN_COOKIE, tokens.access_token ?? '', {
                ...COOKIE_OPTS,
                maxAge: tokens.expires_in,
            });
            response.cookies.set(REFRESH_TOKEN_COOKIE, tokens.refresh_token ?? '', COOKIE_OPTS);
            response.cookies.set(USER_INFO_COOKIE, JSON.stringify(userInfo), COOKIE_OPTS);
            response.cookies.set(OAUTH_STATE_COOKIE, '', { maxAge: 0, path: '/' });
            response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, '', { maxAge: 0, path: '/' });

            return response;
        } catch (error) {
            console.error('OAuth callback error:', error);
            // TODO: Handle errors properly by showing error page
            return NextResponse.redirect(new URL('/', request.url));
        }
    }

    // Handle login redirect (no access token or explicit /auth/login)
    if (!cookies.accessToken || pathname === '/auth/login') {
        const config = await oidcConfig(getOidcConfigFromEnv());

        const codeVerifier: string = client.randomPKCECodeVerifier();
        const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

        const parameters: Record<string, string> = {
            redirect_uri: `${origin}/auth/callback`,
            scope: 'openid profile email',
            code_challenge,
            code_challenge_method: 'S256',
            state: client.randomState(),
        };

        const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);

        const response = NextResponse.redirect(redirectTo);

        const authCookiesOptions: cookie.SerializeOptions = {
            ...COOKIE_OPTS,
            maxAge: 5 * 60, // 5 minutes
            sameSite: 'lax',
        };

        response.cookies.set(OAUTH_STATE_COOKIE, parameters.state, authCookiesOptions);
        response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, codeVerifier, authCookiesOptions);

        return response;
    }

    // For authenticated requests, attach backendSession to headers
    const accessTokenPayload = cookies.accessToken
        ? (decodeJWTPayload(cookies.accessToken) as AccessTokenWithUserInfo | null)
        : null;
    const expiresAt = accessTokenPayload?.exp ? new Date(accessTokenPayload.exp * 1000) : new Date();

    let userInfo: UserInfo | null = null;

    // Try to get user info from the access token itself
    if (accessTokenPayload && (accessTokenPayload.name || accessTokenPayload.email)) {
        userInfo = {
            name: accessTokenPayload.name || '',
            email: accessTokenPayload.email || '',
            picture: accessTokenPayload.picture,
        };
    }

    // Fall back to the user info cookie
    if (!userInfo && cookies.userInfo) {
        try {
            userInfo = JSON.parse(cookies.userInfo);
        } catch (e) {
            console.error('Failed to parse user info cookie:', e);
        }
    }

    const backendSession: BackendSession = {
        type: 'authenticated',
        accessToken: {
            accessToken: cookies.accessToken,
            expiresAt: expiresAt,
        },
        refreshToken: cookies.refreshToken ?? '',
        authenticatedUser: {
            name: userInfo?.name || '',
            email: userInfo?.email || '',
            picture: userInfo?.picture,
            isOwner: cookies.accessToken ? isOwnerFromJWT(cookies.accessToken) : false,
        },
    };

    // In NextJS, we can't add custom data to the request directly like Waku
    // Instead, we'll add it as a custom header that can be read by the app
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set('x-backend-session', JSON.stringify(backendSession));

    const response = NextResponse.next({
        request: {
            headers: requestHeaders,
        },
    });

    // Also set on response headers for testing purposes
    response.headers.set('x-backend-session', JSON.stringify(backendSession));

    return response;
}

export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         */
        '/((?!api|_next/static|_next/image|favicon.ico).*)',
    ],
};
