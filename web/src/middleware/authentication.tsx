import * as cookie from 'cookie';
import {SerializeOptions} from 'cookie';
import type {Middleware} from 'waku/config';
import {Handler, HandlerContext} from "waku/dist/lib/middleware/types";
import {ACCESS_TOKEN_COOKIE, OAUTH_STATE_COOKIE, REFRESH_TOKEN_COOKIE} from "../core/security/consts";
import {BackendSession} from "../core/security/security-model";

interface Cookies {
    accessToken?: string
    refreshToken?: string
    /** state in only used during the authentication flow */
    state?: string
}

const COOKIE_OPTS: SerializeOptions = {
    maxAge: 60, // TODO AGENT - use the real expiration time of the JWT token (access and refresh)
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: true,
};

const cookieMiddleware: Middleware = (): Handler => {
    return async (ctx: HandlerContext, next: () => Promise<void>) => {
        const cookies = cookie.parse(ctx.req.headers.get('cookie') || '');

        const sessions: Cookies = {
            accessToken: cookies[ACCESS_TOKEN_COOKIE],
            refreshToken: cookies[REFRESH_TOKEN_COOKIE],
            state: cookies[OAUTH_STATE_COOKIE],
        }

        if (!sessions.accessToken) {
            const headers = new Headers(ctx.res?.headers);
            headers.append(
                'set-cookie',
                cookie.serialize(ACCESS_TOKEN_COOKIE, "jwt-access-token-test", COOKIE_OPTS),
            );
            headers.append(
                'Content-Type',
                'text/html',
            )
            ctx.res = new Response("<html lang='en'><body><div>You are not logged in ! <a href='/'>Click here to see the website !</a> </div></body></html>", {
                status: 200,
                statusText: "OK",
                headers,
            });

            return
        }

        const backendSession: BackendSession = {
            type: "authenticated",
            accessToken: {
                accessToken: sessions.accessToken,
                expiresAt: new Date(),
            },
            refreshToken: sessions.refreshToken ?? "",
            authenticatedUser: {
                name: "Security Middleware",
                email: "security@middleware.com",
                isOwner: true,
            },
        }
        ctx.data.backendSession = backendSession

        // backendSession is read by JotialProvider to hydrate the client session

        await next();
    };
};

export default cookieMiddleware;