import {NextRequest, NextResponse} from 'next/server';
import * as client from 'openid-client';
import {OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE} from '../../../lib/security/constants';
import {ResponseCookie} from "next/dist/compiled/@edge-runtime/cookies";

type OpenIdConfig = {
    issuer: string;
    clientId: string;
    clientSecret: string;
};

export const basePath = '/nextjs'

async function oidcConfig({ issuer, clientId, clientSecret }: OpenIdConfig): Promise<client.Configuration> {
    return client.discovery(new URL(issuer), clientId, clientSecret);
}

function getOidcConfigFromEnv(): OpenIdConfig {
    return {
        issuer: process.env.OAUTH_ISSUER_URL || '',
        clientId: process.env.OAUTH_CLIENT_ID || '',
        clientSecret: process.env.OAUTH_CLIENT_SECRET || '',
    };
}

const AUTH_COOKIE_OPTS: Partial<ResponseCookie> = {
    maxAge: 5 * 60,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'lax',
};

export async function GET(request: NextRequest) {
    const config = await oidcConfig(getOidcConfigFromEnv());
    const {origin} = request.nextUrl;

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
