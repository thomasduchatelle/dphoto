import * as client from 'openid-client';

export const basePath = '/nextjs';

export type OpenIdConfig = {
    issuer: string;
    clientId: string;
    clientSecret: string;
};

export function getOidcConfigFromEnv(): OpenIdConfig {
    return {
        issuer: process.env.OAUTH_ISSUER_URL || '',
        clientId: process.env.OAUTH_CLIENT_ID || '',
        clientSecret: process.env.OAUTH_CLIENT_SECRET || '',
    };
}

function isValidIssuer(issuer: string): boolean {
    return !!issuer && /^https?:\/\//.test(issuer);
}

export async function oidcConfig({issuer, clientId, clientSecret}: OpenIdConfig): Promise<client.Configuration> {
    if (!isValidIssuer(issuer)) {
        throw new Error('OIDC issuer is missing or invalid.');
    }
    return client.discovery(new URL(issuer), clientId, clientSecret);
}

/**
 * Locally, wiremock stubs the backend REST API without requiring authentication, so `npm run dev` can run without a
 * real OIDC issuer configured. This is only ever true for `next dev` (NODE_ENV=development) and stops applying as
 * soon as a real OAUTH_ISSUER_URL is set, e.g. to test against a real Cognito pool.
 */
export function isLocalAuthBypassed(): boolean {
    return process.env.NODE_ENV === 'development' && !isValidIssuer(getOidcConfigFromEnv().issuer);
}
