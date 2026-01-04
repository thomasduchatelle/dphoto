import {NextRequest, NextResponse} from 'next/server';
import * as client from 'openid-client';
import {getOidcConfigFromEnv, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_NONCE_COOKIE, OAUTH_STATE_COOKIE, oidcConfig} from '@/libs/security';
import {basePath, getOriginalOrigin} from "@/libs/requests";

const AUTH_COOKIE_OPTS: any = {
    maxAge: 5 * 60,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'lax', // lax is required when the Referer is a different site (which happens during OAuth flow when user is not already authenticated on Cognito: user comes from the Social login page)
};

export async function GET(request: NextRequest) {
    console.log("GET /auth/login called");
    const requestUrl = getOriginalOrigin(request);
    try {
        const config = await oidcConfig(getOidcConfigFromEnv());

        const codeVerifier: string = client.randomPKCECodeVerifier();
        const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

        const parameters: Record<string, string> = {
            redirect_uri: new URL(`${basePath}/auth/callback`, requestUrl).toString(),
            scope: 'openid profile email',
            code_challenge,
            code_challenge_method: 'S256',
            state: client.randomState(),
            nonce: client.randomNonce(),
        };

        const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);

        const response = NextResponse.redirect(redirectTo);

        response.cookies.set(OAUTH_STATE_COOKIE, parameters.state, AUTH_COOKIE_OPTS);
        response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, codeVerifier, AUTH_COOKIE_OPTS);
        response.cookies.set(OAUTH_NONCE_COOKIE, parameters.nonce, AUTH_COOKIE_OPTS);

        return response;
    } catch (e) {
        console.error('Error during OAuth login initiation:', e);
        return NextResponse.redirect(new URL(`${basePath}/auth/error`, requestUrl));
    }
}
