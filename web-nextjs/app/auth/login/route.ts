import {NextRequest, NextResponse} from 'next/server';
import * as client from 'openid-client';
import {basePath, getOidcConfigFromEnv, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, oidcConfig} from '@/lib/security';

const AUTH_COOKIE_OPTS: any = {
    maxAge: 5 * 60,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'lax',
};

export async function GET(request: NextRequest) {
    const oidcConfigFromEnv = getOidcConfigFromEnv();
    const config = await oidcConfig(oidcConfigFromEnv);
    const {origin} = request.nextUrl;

    // Use the configured domain name if available, otherwise fall back to request origin
    const baseUrl = oidcConfigFromEnv.domainName 
        ? `https://${oidcConfigFromEnv.domainName}`
        : origin;

    const codeVerifier: string = client.randomPKCECodeVerifier();
    const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

    const parameters: Record<string, string> = {
        redirect_uri: `${baseUrl}${basePath}/auth/callback`,
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
