import {NextRequest, NextResponse} from 'next/server';
import {redirectUrl} from './libs/requests';
import {completeLogout, getValidAuthentication, initiateAuthenticationFlow} from "@/libs/security";

// if basepath was not set, this would work: `export const config = { matcher: [`/(${skipProxyForPageMatching}`] }`
export const skipProxyForPageMatching = /^(?!_next\/static|_next\/image|favicon.ico|api|auth|.*\.js$|.*\.png$|.*\.svg$|.*\.jpg$|.*\.gif$).*/i

export async function proxy(request: NextRequest) {
    const isPublicPath = !skipProxyForPageMatching.test(request.nextUrl.pathname.substring(1))

    if (request.nextUrl.pathname.replaceAll("^/nextjs", "") === '/auth/logout') {
        // Cookies cannot be set from a server component
        // https://nextjs.org/docs/app/api-reference/functions/cookies#understanding-cookie-behavior-in-server-components
        await completeLogout()
        return NextResponse.next()
    }

    if (isPublicPath) {
        return NextResponse.next();
    }

    // note: the access token is refreshed if required
    const authentication = await getValidAuthentication()
    if (authentication.status == "anonymous") {
        const redirection = await initiateAuthenticationFlow(request.nextUrl.pathname);
        return NextResponse.redirect(redirection.redirectTo);
    }

    return NextResponse.next();
}

