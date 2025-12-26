import * as cookie from 'cookie';
import type {Middleware} from 'waku/config';
import {Handler, HandlerContext} from "waku/dist/lib/middleware/types";
import {
    ACCESS_TOKEN_COOKIE,
    BackendSession,
    decodeJWTPayload,
    isOwnerFromJWT,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    REFRESH_TOKEN_COOKIE
} from "../core/security";
import * as client from 'openid-client'
import {getEnv} from "waku";

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
    accessToken?: string
    refreshToken?: string
    userInfo?: string
    /** state in only used during the authentication flow */
    state?: string
    codeVerifier?: string
}

const COOKIE_OPTS: cookie.SerializeOptions = {
    maxAge: 3600,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: true,
};

function readCookies(ctx: HandlerContext): Cookies {
    const cookies = cookie.parse(ctx.req.headers.get('cookie') || '');
    return {
        accessToken: cookies[ACCESS_TOKEN_COOKIE],
        refreshToken: cookies[REFRESH_TOKEN_COOKIE],
        userInfo: cookies[USER_INFO_COOKIE],
        state: cookies[OAUTH_STATE_COOKIE],
        codeVerifier: cookies[OAUTH_CODE_VERIFIER_COOKIE],
    }
}

type OpenIdConfig = {
    issuer: string,
    clientId: string,
    clientSecret: string,
};

function oidcConfig({issuer, clientId, clientSecret}: OpenIdConfig): Promise<client.Configuration> {
    // TODO AGENT - fetch the configuration once, and then cache it forever.
    return client.discovery(
        new URL(issuer),
        clientId,
        clientSecret,
    )
}

function parse(url: string): { scheme: string; host: string; path: string } {
    const u = new URL(url);
    return {
        scheme: u.protocol.replace(':', ''),
        host: u.host,
        path: u.pathname,
    };
}

const cookieMiddleware: Middleware = (openIdConfig: OpenIdConfig | undefined = {
    issuer: getEnv("COGNITO_ISSUER"),
    clientId: getEnv("COGNITO_CLIENT_ID"),
    clientSecret: getEnv("COGNITO_CLIENT_SECRET"),
}): Handler => {

    return async (ctx: HandlerContext, next: () => Promise<void>) => {
        console.log("Middleware/authentication: processing request ", ctx.req.url);

        const {scheme, host, path} = parse(ctx.req.url);
        const cookies = readCookies(ctx);

        if (path === '/auth/callback') {
            // TODO AGENT - handle errors callback like "http://localhost:3000/auth/callback?error_description=user.email%3A+Attribute+cannot+be+updated.%0A+&state=NJPSQ2B59ghT6oggnSls2SXCx45CJ4Z1XtZIU9oBknU&error=invalid_request". It needs to show an html error with the actual error in the callback. (Just process next() in the middleware, and create a new waku page that show the error.)
            // TODO AGENT - handle errors when authorizationCodeGrant fails with a 400 error and a body content like:
            //   cause: { error: 'invalid_grant' },
            //   code: 'OAUTH_RESPONSE_BODY_ERROR',
            //   error: 'invalid_grant',
            //   status: 400,
            //   error_description: undefined
            // }
            // for that, place the error in ctx.data.oidcError and then create a waku page that shows the error to the user.

            const config = await oidcConfig(openIdConfig);

            const tokens: client.TokenEndpointResponse = await client.authorizationCodeGrant(
                config,
                new URL(ctx.req.url),
                {
                    pkceCodeVerifier: cookies.codeVerifier,
                    expectedState: cookies.state,
                },
            )

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

            const headers = new Headers(ctx.res?.headers);
            headers.append(
                'set-cookie',
                cookie.serialize(ACCESS_TOKEN_COOKIE, tokens.access_token ?? "", {
                    ...COOKIE_OPTS,
                    maxAge: tokens.expires_in,
                }),
            );
            headers.append(
                'set-cookie',
                cookie.serialize(REFRESH_TOKEN_COOKIE, tokens.refresh_token ?? "", COOKIE_OPTS),
            );
            headers.append(
                'set-cookie',
                cookie.serialize(USER_INFO_COOKIE, JSON.stringify(userInfo), COOKIE_OPTS),
            );
            headers.append(
                'set-cookie',
                cookie.serialize(OAUTH_STATE_COOKIE, '', {maxAge: 0, path: '/'}),
            );
            headers.append(
                'set-cookie',
                cookie.serialize(OAUTH_CODE_VERIFIER_COOKIE, '', {maxAge: 0, path: '/'}),
            );
            headers.append(
                'Location',
                `${scheme}://${host}/`, // TODO AGENT Redirect to the original URL requested before login
            )
            ctx.res = new Response(null, {
                status: 302,
                statusText: "Found",
                headers,
            });
            return
        } else if (!cookies.accessToken || path === '/auth/login') {
            const config = await oidcConfig(openIdConfig);

            const codeVerifier: string = client.randomPKCECodeVerifier()
            const code_challenge: string =
                await client.calculatePKCECodeChallenge(codeVerifier)

            let parameters: Record<string, string> = {
                redirect_uri: `${scheme}://${host}/auth/callback`,
                scope: "openid profile email",
                code_challenge,
                code_challenge_method: 'S256',
                state: client.randomState(),
            }

            let redirectTo: URL = client.buildAuthorizationUrl(config, parameters)

            const headers = new Headers(ctx.res?.headers);
            const authCookiesOptions: cookie.SerializeOptions = {
                ...COOKIE_OPTS,
                maxAge: 5 * 60, // authentication flow in Cognito is set to 3 minutes.
                sameSite: 'lax', // required when the Referer is expected to be a different site (which is the case during OIDC flows)
            };
            let value = cookie.serialize(OAUTH_STATE_COOKIE, parameters.state, authCookiesOptions);
            headers.append(
                'FOO',
                value,
            );
            headers.append(
                'set-cookie',
                cookie.serialize(OAUTH_CODE_VERIFIER_COOKIE, codeVerifier, authCookiesOptions),
            );
            headers.append(
                'Location',
                redirectTo.toString(),
            )
            ctx.res = new Response(null, {
                status: 302,
                statusText: "Found",
                headers,
            });

            return
        }

        // backendSession is read by JotialProvider to hydrate the client session
        const accessTokenPayload = cookies.accessToken ? decodeJWTPayload(cookies.accessToken) as AccessTokenWithUserInfo | null : null;
        const expiresAt = accessTokenPayload?.exp ? new Date(accessTokenPayload.exp * 1000) : new Date();

        let userInfo: UserInfo | null = null;

        // First try to get user info from the access token itself (if present)
        // Note: access token may contain name, email, picture for testing purposes
        if (accessTokenPayload && (accessTokenPayload.name || accessTokenPayload.email)) {
            userInfo = {
                name: accessTokenPayload.name || '',
                email: accessTokenPayload.email || '',
                picture: accessTokenPayload.picture,
            };
        }

        // Otherwise, fall back to the user info cookie
        if (!userInfo && cookies.userInfo) {
            try {
                userInfo = JSON.parse(cookies.userInfo);
            } catch (e) {
                console.error('Failed to parse user info cookie:', e);
            }
        }

        const backendSession: BackendSession = {
            type: "authenticated",
            accessToken: {
                accessToken: cookies.accessToken,
                expiresAt: expiresAt,
            },
            refreshToken: cookies.refreshToken ?? "",
            authenticatedUser: {
                name: userInfo?.name || '',
                email: userInfo?.email || '',
                picture: userInfo?.picture,
                isOwner: cookies.accessToken ? isOwnerFromJWT(cookies.accessToken) : false,
            },
        }
        ctx.data.backendSession = backendSession


        await next();
    };
};

export default cookieMiddleware;