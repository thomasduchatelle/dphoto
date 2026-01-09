import {NextRequest, NextResponse} from 'next/server';
import {getOidcConfigFromEnv, getLogoutUrl} from '@/libs/security';
import {basePath, getOriginalOrigin} from '@/libs/requests';

export async function GET(request: NextRequest) {
    console.log("GET /auth/logout called");
    const requestUrl = getOriginalOrigin(request);
    
    try {
        const config = getOidcConfigFromEnv();
        const logoutCallbackUrl = new URL(`${basePath}/auth/logout-callback`, requestUrl).toString();
        const cognitoLogoutUrl = getLogoutUrl(config.issuer, config.clientId, logoutCallbackUrl);

        return NextResponse.redirect(cognitoLogoutUrl);
    } catch (e) {
        console.error('Error during logout:', e);
        return NextResponse.redirect(new URL(`${basePath}/auth/error`, requestUrl));
    }
}
