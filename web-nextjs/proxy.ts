import {NextRequest, NextResponse} from 'next/server';
import {completeLogout, getCurrentAuthenticationStatus, initiateAuthenticationFlow, refreshSession} from "@/libs/security";
import {appliesCookies, buildRedirectResponse, newReadCookieStore} from "@/libs/nextjs-cookies";

// if basepath was not set, this would work: `export const config = { matcher: [`/(${skipProxyForPageMatching}`] }`
export const skipProxyForPageMatching = /^(?!_next\/static|_next\/image|favicon.ico|api|auth|.*\.js$|.*\.png$|.*\.svg$|.*\.jpg$|.*\.gif$).*/i

export async function proxy(request: NextRequest) {
    const isPublicPath = !skipProxyForPageMatching.test(request.nextUrl.pathname.substring(1))

    if (request.nextUrl.pathname.replaceAll("^/nextjs", "") === '/auth/logout') {
        // Cookies cannot be set from a server component
        // https://nextjs.org/docs/app/api-reference/functions/cookies#understanding-cookie-behavior-in-server-components
        return appliesCookies(NextResponse.next(), completeLogout());
    }

    if (isPublicPath) {
        return NextResponse.next();
    }

    const authentication = await getCurrentAuthenticationStatus(request)

    if (!authentication.authenticated) {
        const redirection = await initiateAuthenticationFlow(request);
        return buildRedirectResponse(redirection);
    }

    if (authentication.aboutToExpire) {
        const refresh = await refreshSession(newReadCookieStore(request))
        if (!refresh.success) {
            const redirection = await initiateAuthenticationFlow(request);
            return buildRedirectResponse(redirection);
        }

        // the page redirect to itself to get the new cookies in the request, available to the server components.
        return buildRedirectResponse({
            ...refresh,
            redirectTo: request.nextUrl,
        });
    }

    return NextResponse.next();
}

