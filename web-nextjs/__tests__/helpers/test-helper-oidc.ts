export const TEST_ISSUER_URL = 'https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_7CivTjR7R';
export const TEST_CLIENT_ID = '7k53mt7hv23fffi7dqe9sfi1b2';
export const TEST_CLIENT_SECRET = 'test-client-secret';

/**
 * Creates a Cognito-style access token (used during OAuth flow)
 */
export const createCognitoAccessToken = (claims: {
    sub?: string;
    username?: string;
    scope?: string;
    exp?: number;
    iat?: number;
}): string => {
    const now = Math.floor(Date.now() / 1000);
    const payload = {
        sub: claims.sub || '52d5a4b4-4021-70f5-51c2-6bbdc346fc76',
        'cognito:groups': ['eu-west-1_7CivTjR7R_Google'],
        iss: TEST_ISSUER_URL,
        version: 2,
        client_id: TEST_CLIENT_ID,
        origin_jti: 'd66449ad-6296-47a5-a920-81c7f18a4edd',
        token_use: 'access',
        scope: claims.scope || 'openid profile email',
        auth_time: claims.iat || now,
        exp: claims.exp || now + 3600,
        iat: claims.iat || now,
        jti: '0e3c8081-1807-437b-a269-cf1b8757cd73',
        username: claims.username || 'Google_109832440005038882419',
    };

    const header = Buffer.from(JSON.stringify({ alg: 'RS256', typ: 'JWT' })).toString('base64url');
    const body = Buffer.from(JSON.stringify(payload)).toString('base64url');
    return `${header}.${body}.fake-signature`;
};

/**
 * Creates a Cognito-style ID token with user information
 */
export const createValidIDToken = (claims: {
    sub?: string;
    email?: string;
    given_name?: string;
    family_name?: string;
    picture?: string;
    exp?: number;
    iat?: number;
}): string => {
    const now = Math.floor(Date.now() / 1000);
    const payload = {
        at_hash: 'qBwIhsgyC7EgWaD3_SgUXA',
        sub: claims.sub || '52d5a4b4-4021-70f5-51c2-6bbdc346fc76',
        'cognito:groups': ['eu-west-1_7CivTjR7R_Google'],
        email_verified: false,
        iss: TEST_ISSUER_URL,
        'cognito:username': 'Google_109832440005038882419',
        given_name: claims.given_name || 'Thomas',
        picture: claims.picture || 'https://lh3.googleusercontent.com/a/ACg8ocKBKtsO86UaxMwMaQpnykZv5Qb38FLYJlMzQi3FrriBcDaxAUxP=s96-c',
        origin_jti: 'd66449ad-6296-47a5-a920-81c7f18a4edd',
        aud: TEST_CLIENT_ID,
        identities: [
            {
                dateCreated: '1767445640097',
                userId: '109832440005038882419',
                providerName: 'Google',
                providerType: 'Google',
                issuer: null,
                primary: 'true',
            },
        ],
        token_use: 'id',
        auth_time: claims.iat || now,
        exp: claims.exp || now + 3600,
        iat: claims.iat || now,
        family_name: claims.family_name || 'Duchatelle',
        jti: '0eace6ef-3d8c-4c8d-9733-4c00b7fb15f0',
        email: claims.email || 'tomdush@gmail.com',
    };

    const header = Buffer.from(JSON.stringify({ alg: 'RS256', typ: 'JWT' })).toString('base64url');
    const body = Buffer.from(JSON.stringify(payload)).toString('base64url');
    return `${header}.${body}.fake-signature`;
};

/**
 * Creates a DPhoto backend access token with Scopes for authorization
 * (This is what the backend returns after exchanging the Cognito token)
 */
export const createBackendAccessToken = (claims: {
    sub?: string;
    email?: string;
    Scopes?: string;
    exp?: number;
    iat?: number;
}): string => {
    const now = Math.floor(Date.now() / 1000);
    const payload = {
        sub: claims.sub || claims.email || 'tomdush@gmail.com',
        iss: 'dphoto',
        aud: ['dphoto'],
        exp: claims.exp || now + 3600,
        iat: claims.iat || now,
        Scopes: claims.Scopes || 'owner:tomdush@gmail.com',
    };

    const header = Buffer.from(JSON.stringify({ alg: 'HS512', typ: 'JWT' })).toString('base64url');
    const body = Buffer.from(JSON.stringify(payload)).toString('base64url');
    return `${header}.${body}.fake-signature`;
};

export const createTokenResponse = (overrides?: {
    access_token?: string;
    refresh_token?: string;
    id_token?: string;
    expires_in?: number;
}) => ({
    access_token: overrides?.access_token || createCognitoAccessToken({}),
    refresh_token: overrides?.refresh_token || 'REFRESH_TOKEN_VALUE',
    id_token: overrides?.id_token || createValidIDToken({}),
    expires_in: overrides?.expires_in || 3600,
    token_type: 'Bearer',
});
