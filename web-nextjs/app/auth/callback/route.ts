import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import * as client from 'openid-client';
import {
    ACCESS_TOKEN_COOKIE,
    basePath,
    decodeJWTPayload,
    getOidcConfigFromEnv,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    oidcConfig,
    REFRESH_TOKEN_COOKIE
} from '@/lib/security';
import {getOriginalOrigin} from '@/lib/request-utils';

const USER_INFO_COOKIE = 'dphoto-user-info';

interface IDTokenPayload {
    given_name?: string;
    family_name?: string;
    email?: string;
    picture?: string;
    exp?: number;

    [key: string]: any;
}

interface UserInfo {
    name: string;
    email: string;
    picture?: string;
}

interface Cookies {
    state?: string;
    codeVerifier?: string;
}

function readCookies(request: NextRequest): Cookies {
    const cookieHeader = request.headers.get('cookie') || '';
    const cookies = cookie.parse(cookieHeader);
    return {
        state: cookies[OAUTH_STATE_COOKIE],
        codeVerifier: cookies[OAUTH_CODE_VERIFIER_COOKIE],
    };
}

const COOKIE_OPTS: any = {
    httpOnly: true,
    secure: true,
    sameSite: 'strict',
    path: '/',
};

export async function GET(request: NextRequest) {
    const config = await oidcConfig(getOidcConfigFromEnv());
    const origin = getOriginalOrigin(request);
    const url = new URL(request.url);
    const cookies = readCookies(request);

    const targetUrl = new URL(`${basePath ?? '/'}`, `${origin}${url.pathname}`);

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
                const firstName = idTokenPayload.given_name || '';
                const lastName = idTokenPayload.family_name || '';
                const fullName = [firstName, lastName].filter(Boolean).join(' ');

                userInfo = {
                    name: fullName,
                    email: idTokenPayload.email || '',
                    picture: idTokenPayload.picture,
                };
            }
        }

        const response = NextResponse.redirect(targetUrl);

        response.cookies.set(ACCESS_TOKEN_COOKIE, tokens.access_token ?? '', {
            ...COOKIE_OPTS,
            maxAge: tokens.expires_in,
        });
        response.cookies.set(REFRESH_TOKEN_COOKIE, tokens.refresh_token ?? '', COOKIE_OPTS);
        response.cookies.set(USER_INFO_COOKIE, JSON.stringify(userInfo), COOKIE_OPTS);
        response.cookies.set(OAUTH_STATE_COOKIE, '', {maxAge: 0, path: '/'});
        response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, '', {maxAge: 0, path: '/'});

        return response;
    } catch (error) {
        console.error('OAuth callback error:', error);
        return NextResponse.redirect(targetUrl);
    }
}
