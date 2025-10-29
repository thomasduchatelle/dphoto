import * as cookie from 'cookie';
import type {Middleware} from 'waku/config';
import {Handler, HandlerContext} from "waku/dist/lib/middleware/types";
import {ACCESS_TOKEN_COOKIE, BackendSession, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, REFRESH_TOKEN_COOKIE} from "../core/security";
import * as client from 'openid-client'
import {getEnv} from "waku";

interface Cookies {
    accessToken?: string
    refreshToken?: string
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
        state: cookies[OAUTH_STATE_COOKIE],
        codeVerifier: cookies[OAUTH_CODE_VERIFIER_COOKIE],
    }
}

async function oidcConfig() {
    // TODO AGENT - fetch the configuration once, and then cache it forever.
    let config: client.Configuration = await client.discovery(
        new URL(getEnv("COGNITO_ISSUER")),
        getEnv("COGNITO_CLIENT_ID"),
        getEnv("COGNITO_CLIENT_SECRET"),
    )
    return config;
}

function parse(url: string): { scheme: string; host: string; path: string } {
    const u = new URL(url);
    return {
        scheme: u.protocol.replace(':', ''),
        host: u.host,
        path: u.pathname,
    };
}

const cookieMiddleware: Middleware = (): Handler => {

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

            console.log("Cookies on callback:", cookies);
            const config = await oidcConfig();

            const tokens: client.TokenEndpointResponse = await client.authorizationCodeGrant(
                config,
                new URL(ctx.req.url),
                {
                    pkceCodeVerifier: cookies.codeVerifier,
                    expectedState: cookies.state,
                },
            )

            // TODO AGENT - capture the identifier token which is required for the full name and picture of the user. The details must be available in the context, and be stored in the dynamodb table (see pkg/acl/aclcore/authenticate_sso.go)
            // TODO AGENT - use the real expiration time of the JWT token (access and refresh) for the expiration of the cookies.
            // TODO AGENT - redirect to the original URL that was requested before login (store it in a cookie before redirecting to /auth/login)
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
                cookie.serialize(OAUTH_STATE_COOKIE, '', {maxAge: 0, path: '/'}),
            );
            headers.append(
                'set-cookie',
                cookie.serialize(OAUTH_CODE_VERIFIER_COOKIE, '', {maxAge: 0, path: '/'}),
            );
            headers.append(
                'Location',
                `${scheme}://${host}/`,
            )
            ctx.res = new Response(null, {
                status: 302,
                statusText: "Found",
                headers,
            });
            return
        } else if (!cookies.accessToken || path === '/auth/login') {
            const config = await oidcConfig();

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
            headers.append(
                'set-cookie',
                cookie.serialize(OAUTH_STATE_COOKIE, parameters.state, authCookiesOptions),
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
        const backendSession: BackendSession = {
            type: "authenticated",
            accessToken: {
                accessToken: cookies.accessToken,
                expiresAt: new Date(),
            },
            refreshToken: cookies.refreshToken ?? "",
            authenticatedUser: {
                name: "Security Middleware",
                email: "security@middleware.com",
                isOwner: true,
            },
        }
        ctx.data.backendSession = backendSession


        await next();
    };
};

export default cookieMiddleware;