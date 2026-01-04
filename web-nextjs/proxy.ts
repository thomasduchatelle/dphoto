import {NextRequest, NextResponse} from 'next/server';
import * as cookie from 'cookie';
import {ACCESS_TOKEN_COOKIE} from '@/libs/security/constants';
import {basePath, getOriginalOrigin} from './libs/requests';

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

    return NextResponse.next();
}
