import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import * as client from 'openid-client';
import {ACCESS_TOKEN_COOKIE, REFRESH_TOKEN_COOKIE} from '@/libs/security/constants';
import {basePath, getOriginalOrigin} from './libs/requests';
import {decodeJWTPayload, getOidcConfigFromEnv, oidcConfig} from '@/libs/security';

interface Cookies {
    accessToken?: string;
    refreshToken?: string;
}

function readCookies(request: NextRequest): Cookies {
    const cookieHeader = request.headers.get('cookie') || '';
    const cookies = cookie.parse(cookieHeader);
    return {
        accessToken: cookies[ACCESS_TOKEN_COOKIE],
        refreshToken: cookies[REFRESH_TOKEN_COOKIE],
    };
}

function isTokenExpired(token: string): boolean {
    const payload = decodeJWTPayload(token);
    if (!payload || !payload.exp) {
        return true;
    }
    
    const now = Math.floor(Date.now() / 1000);
    return payload.exp <= now;
}

const COOKIE_OPTS: any = {
    httpOnly: true,
    secure: true,
    sameSite: 'lax',
    path: '/',
};

async function refreshAccessToken(refreshToken: string): Promise<client.TokenEndpointResponse | null> {
    try {
        const config = await oidcConfig(getOidcConfigFromEnv());
        const tokens = await client.refreshTokenGrant(config, refreshToken);
        return tokens;
    } catch (error) {
        console.error('Failed to refresh token:', error);
        return null;
    }
}

// if basepath was not set, this would work: `export const config = { matcher: [`/(${skipProxyForPageMatching}`] }`
export const skipProxyForPageMatching = /^(?!_next\/static|_next\/image|favicon.ico|api|auth|.*\.js$|.*\.png$|.*\.svg$|.*\.jpg$|.*\.gif$).*/i

export async function proxy(request: NextRequest) {
    const requestUrl = getOriginalOrigin(request)

    const isPublicPath = !skipProxyForPageMatching.test(request.nextUrl.pathname.substring(1))

    if (isPublicPath) {
        return NextResponse.next();
    }

    const cookies = readCookies(request);

    if (!cookies.accessToken) {
        return NextResponse.redirect(new URL(`${basePath}/auth/login`, requestUrl));
    }

    // Check if access token is expired
    if (isTokenExpired(cookies.accessToken)) {
        // Try to refresh the token if we have a refresh token
        if (cookies.refreshToken) {
            const newTokens = await refreshAccessToken(cookies.refreshToken);
            
            if (newTokens && newTokens.access_token) {
                // Successfully refreshed, update cookies and continue
                const response = NextResponse.next();
                response.cookies.set(ACCESS_TOKEN_COOKIE, newTokens.access_token, {
                    ...COOKIE_OPTS,
                    maxAge: newTokens.expires_in,
                });
                
                // Update refresh token if a new one was provided
                if (newTokens.refresh_token) {
                    response.cookies.set(REFRESH_TOKEN_COOKIE, newTokens.refresh_token, COOKIE_OPTS);
                }
                
                return response;
            }
        }
        
        // Refresh failed or no refresh token available, redirect to login
        return NextResponse.redirect(new URL(`${basePath}/auth/login`, requestUrl));
    }

    return NextResponse.next();
}
