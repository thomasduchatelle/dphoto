import * as cookie from 'cookie';
import {SerializeOptions} from 'cookie';

import type {Middleware} from 'waku/config';
import {Handler, HandlerContext} from "waku/dist/lib/middleware/types";

// XXX we would probably like to extend config.
const COOKIE_OPTS: SerializeOptions = {
    maxAge: 24 * 30 * 3600,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: true,
};

const cookieMiddleware: Middleware = (): Handler => {
    return async (ctx: HandlerContext, next: () => Promise<void>) => {
        const cookies = cookie.parse(ctx.req.headers.get('cookie') || '');
        ctx.data.count = (Number(cookies.count) || 0);
        console.log(`[middleware] req=${JSON.stringify(ctx.req.url)} ...`)
        console.log(`[middleware] before count=${ctx.data.count} ...`)
        await next();
        console.log(`[middleware] res=${JSON.stringify(ctx.res?.headers)} ...`)
        if (ctx.res) {
            const newCount = (Number(cookies.count) || 0) + 1
            const headers = new Headers(ctx.res.headers);
            console.log(`[middleware] after count=${newCount} ...`)
            headers.append(
                'set-cookie',
                cookie.serialize('count', String(newCount), COOKIE_OPTS),
            );
            ctx.res = new Response(ctx.res.body, {
                status: ctx.res.status,
                statusText: ctx.res.statusText,
                headers,
            });
        }
    };
};

export default cookieMiddleware;