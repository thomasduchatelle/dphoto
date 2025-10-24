import * as cookie from 'cookie';
import {SerializeOptions} from 'cookie';

import type {Middleware} from 'waku/config';
import {Handler, HandlerContext} from "waku/dist/lib/middleware/types";
import {Session} from "../components/AuthProvider";

// Cookie names for authentication tokens
export const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
export const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';

// Cookie for OAuth state during authentication flow
export const OAUTH_STATE_COOKIE = 'dphoto-oauth-state';

// XXX we would probably like to extend config.
interface Cookies {
    accessToken: string | undefined
    refreshToken: string | undefined
    // state in only used during the authentication flow
    state: string | undefined
}

const COOKIE_OPTS: SerializeOptions = {
    maxAge: 60,
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
            ctx.res = new Response("<div>You are not logged in ! <a href='/'>Click here to see the website !</a> </div>", {
                status: 200,
                statusText: "OK",
                headers,
            });

            return
        }

        const clientSession: Session = {
            accessToken: {
                accessToken: sessions.accessToken,
                expiryTime: 3600,
            },
            user: {
                name: "Security Middleware",
                email: "security@middleware.com",
                isOwner: true,
            }
        }
        console.log(`[middleware] ${ctx.req.url} authenticated with ${JSON.stringify(clientSession)}`)
        ctx.data.session = clientSession
        await next();
    };
};

export default cookieMiddleware;