export const TEST_ISSUER_URL = 'https://cognito.example.com';
export const TEST_CLIENT_ID = 'test-client-id';
export const TEST_CLIENT_SECRET = 'test-client-secret';

export const createValidIDToken = (claims: {
    name?: string;
    email?: string;
    picture?: string;
    exp?: number;
    Scopes?: string;
}): string => {
    const now = Math.floor(Date.now() / 1000);
    const payload = {
        name: claims.name || 'Test User',
        email: claims.email || 'test@example.com',
        picture: claims.picture,
        exp: claims.exp || now + 3600,
        iat: now,
        iss: TEST_ISSUER_URL,
        aud: TEST_CLIENT_ID,
        sub: claims.email || 'test@example.com',
        Scopes: claims.Scopes,
    };

    const header = Buffer.from(JSON.stringify({ alg: 'RS256', typ: 'JWT' })).toString('base64url');
    const body = Buffer.from(JSON.stringify(payload)).toString('base64url');
    return `${header}.${body}.fake-signature`;
};

export const createValidAccessToken = (claims: {
    name?: string;
    email?: string;
    picture?: string;
    exp?: number;
    Scopes?: string;
}): string => {
    return createValidIDToken(claims);
};

export const createTokenResponse = (overrides?: {
    access_token?: string;
    refresh_token?: string;
    id_token?: string;
    expires_in?: number;
}) => ({
    access_token: overrides?.access_token || 'ACCESS_TOKEN_VALUE',
    refresh_token: overrides?.refresh_token || 'REFRESH_TOKEN_VALUE',
    id_token: overrides?.id_token || createValidIDToken({
        name: 'John Doe',
        email: 'john@example.com',
        picture: 'https://example.com/avatar.jpg',
        Scopes: 'owner:testuser',
    }),
    expires_in: overrides?.expires_in || 3600,
    token_type: 'Bearer',
});
