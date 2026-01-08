import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import {ACCESS_TOKEN_COOKIE, REFRESH_TOKEN_COOKIE, refreshSessionIfNecessary, SessionCookies, SessionRefresh} from '@/libs/security';
import {basePath, getOriginalOrigin} from './libs/requests';

function readCookies(request: NextRequest): SessionCookies {
    const cookieHeader = request.headers.get('cookie') || '';
    const cookies = cookie.parse(cookieHeader);
    return {
        accessToken: cookies[ACCESS_TOKEN_COOKIE],
        refreshToken: cookies[REFRESH_TOKEN_COOKIE],
    };
}

interface CookieOptions {
    httpOnly: boolean;
    secure: boolean;
    sameSite: 'lax' | 'strict' | 'none';
    path: string;
}

const COOKIE_OPTS: CookieOptions = {
    httpOnly: true,
    secure: true,
    sameSite: 'lax',
    path: '/',
};

// if basepath was not set, this would work: `export const config = { matcher: [`/(${skipProxyForPageMatching}`] }`
export const skipProxyForPageMatching = /^(?!_next\/static|_next\/image|favicon.ico|api|auth|.*\.js$|.*\.png$|.*\.svg$|.*\.jpg$|.*\.gif$).*/i

function reloadWithTheCookies(requestUrl: URL, sessionRefresh: SessionRefresh) {
    const response = NextResponse.redirect(requestUrl);

    if (sessionRefresh.newAccessToken) {
        if (sessionRefresh.newAccessToken.expiresIn > 0) {
            response.cookies.set(ACCESS_TOKEN_COOKIE, sessionRefresh.newAccessToken.token, {
                ...COOKIE_OPTS,
                maxAge: sessionRefresh.newAccessToken.expiresIn,
            });
        } else {
            response.cookies.set(ACCESS_TOKEN_COOKIE, sessionRefresh.newAccessToken.token, COOKIE_OPTS);
        }
    }

    if (sessionRefresh.newRefreshToken) {
        response.cookies.set(REFRESH_TOKEN_COOKIE, sessionRefresh.newRefreshToken, COOKIE_OPTS);
    }

    return response
}

export async function proxy(request: NextRequest) {
    const requestUrl = getOriginalOrigin(request)

    const isPublicPath = !skipProxyForPageMatching.test(request.nextUrl.pathname.substring(1))

    if (isPublicPath) {
        return NextResponse.next();
    }

    const cookies = readCookies(request);

    // Check session and refresh if necessary
    const sessionRefresh = await refreshSessionIfNecessary(cookies);

    if (sessionRefresh.status === 'anonymous' || sessionRefresh.status === 'expired') {
        return NextResponse.redirect(new URL(`${basePath}/auth/login`, requestUrl));
    }

    // Update cookies if tokens were refreshed
    if (sessionRefresh.newRefreshToken || sessionRefresh.newAccessToken) {
        return reloadWithTheCookies(requestUrl, sessionRefresh)
    }

    return NextResponse.next();
}

