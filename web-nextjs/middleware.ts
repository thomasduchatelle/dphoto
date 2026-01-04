import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import {ACCESS_TOKEN_COOKIE} from './lib/security/constants';

interface Cookies {
    accessToken?: string;
}

function readCookies(request: NextRequest): Cookies {
    const cookieHeader = request.headers.get('cookie') || '';
    const cookies = cookie.parse(cookieHeader);
    return {
        accessToken: cookies[ACCESS_TOKEN_COOKIE],
    };
}

export async function middleware(request: NextRequest) {
    const { pathname } = request.nextUrl;
    
    // Whitelist pages that don't require authentication
    const publicPaths = ['/auth/login', '/auth/callback', '/error'];
    const isPublicPath = publicPaths.some(path => pathname.startsWith(path));
    
    if (isPublicPath) {
        return NextResponse.next();
    }

    const cookies = readCookies(request);

    if (!cookies.accessToken) {
        return NextResponse.redirect(new URL('/auth/login', request.url));
    }

    return NextResponse.next();
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
