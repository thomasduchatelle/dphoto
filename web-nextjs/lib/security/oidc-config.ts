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

export async function oidcConfig({issuer, clientId, clientSecret}: OpenIdConfig): Promise<client.Configuration> {
    if (!issuer || !/^https?:\/\//.test(issuer)) {
        throw new Error('OIDC issuer is missing or invalid.');
    }
    return client.discovery(new URL(issuer), clientId, clientSecret);
}
