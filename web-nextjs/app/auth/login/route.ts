import {NextRequest, NextResponse} from 'next/server';
import * as client from 'openid-client';
import {OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE} from '../../../lib/security/constants';
import {ResponseCookie} from "next/dist/compiled/@edge-runtime/cookies";

type OpenIdConfig = {
    issuer: string;
    clientId: string;
    clientSecret: string;
};

async function oidcConfig({ issuer, clientId, clientSecret }: OpenIdConfig): Promise<client.Configuration> {
    return client.discovery(new URL(issuer), clientId, clientSecret);
}

function getOidcConfigFromEnv(): OpenIdConfig {
    return {
        issuer: process.env.COGNITO_ISSUER || '',
        clientId: process.env.COGNITO_CLIENT_ID || '',
        clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
    };
}

const COOKIE_OPTS: Partial<ResponseCookie> = {
    maxAge: 3600,
    httpOnly: true,
    path: '/',
    secure: true,
    sameSite: 'strict',
};

export async function GET(request: NextRequest) {
    const config = await oidcConfig(getOidcConfigFromEnv());
    const { origin } = request.nextUrl;

    const codeVerifier: string = client.randomPKCECodeVerifier();
    const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

    const parameters: Record<string, string> = {
        redirect_uri: `${origin}/auth/callback`,
        scope: 'openid profile email',
        code_challenge,
        code_challenge_method: 'S256',
        state: client.randomState(),
    };

    const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);

    const response = NextResponse.redirect(redirectTo);

    const authCookiesOptions: Partial<ResponseCookie> = {
        ...COOKIE_OPTS,
        maxAge: 5 * 60,
        sameSite: 'lax',
    };

    response.cookies.set(OAUTH_STATE_COOKIE, parameters.state, authCookiesOptions);
    response.cookies.set(OAUTH_CODE_VERIFIER_COOKIE, codeVerifier, authCookiesOptions);

    return response;
}
