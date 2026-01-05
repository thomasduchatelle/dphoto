import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import {ACCESS_TOKEN_COOKIE} from './lib/security/constants';
import {getOriginalOrigin} from './lib/request-utils';

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

// if basepath was not set, this would work: `export const config = { matcher: [`/(${skipProxyForPageMatching}`] }`
export const skipProxyForPageMatching = /^(?!_next\/static|_next\/image|favicon.ico|api|auth|.*\.js$|.*\.png$|.*\.svg$|.*\.jpg$|.*\.gif$).*/i

export async function middleware(request: NextRequest) {
    const {pathname, basePath} = request.nextUrl;

    const isPublicPath = !skipProxyForPageMatching.test(pathname.substring(1))

    if (isPublicPath) {
        return NextResponse.next();
    }

    const cookies = readCookies(request);

    if (!cookies.accessToken) {
        const origin = getOriginalOrigin(request);
        return NextResponse.redirect(new URL(`${basePath}/auth/login`, origin));
    }

    return NextResponse.next();
}
