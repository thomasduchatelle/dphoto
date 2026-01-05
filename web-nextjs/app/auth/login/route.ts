import {NextRequest, NextResponse} from 'next/server';
import * as client from 'openid-client';
import {basePath, getOidcConfigFromEnv, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, oidcConfig} from '@/lib/security';
import {getOriginalOrigin} from '@/lib/request-utils';

const AUTH_COOKIE_OPTS: any = {
    maxAge: 5 * 60,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'lax',
};

export async function GET(request: NextRequest) {
    const config = await oidcConfig(getOidcConfigFromEnv());
    const origin = getOriginalOrigin(request);
    
    // Debug: Only print headers if 'forwarded' header is not found
    if (!request.headers.get('forwarded')) {
        console.log('=== Login Route - Request Headers (no forwarded header found) ===');
        request.headers.forEach((value, key) => {
            console.log(`${key}: ${value}`);
        });
        console.log('request.url:', request.url);
        console.log('request.nextUrl.origin:', request.nextUrl.origin);
        console.log('=====================================');
    }
    
    console.log('Computed origin:', origin);

    const codeVerifier: string = client.randomPKCECodeVerifier();
    const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

    const parameters: Record<string, string> = {
        redirect_uri: `${origin}${basePath}/auth/callback`,
        scope: 'openid profile email',
        code_challenge,
        code_challenge_method: 'S256',
        state: client.randomState(),
    };

    const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);

    const response = NextResponse.redirect(redirectTo);

    response.cookies.set(OAUTH_STATE_COOKIE, parameters.state, AUTH_COOKIE_OPTS);
    response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, codeVerifier, AUTH_COOKIE_OPTS);

    return response;
}
