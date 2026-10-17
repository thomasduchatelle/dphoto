// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';

vi.mock('server-only', () => ({}));

const {POST} = await import('./route');
import {COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN, COOKIE_SESSION_USER_INFO} from '@/libs/security/constants';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {setCookiesOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('POST /auth/refresh', () => {
    let fakeOIDCServer: FakeOIDCServer;

    beforeAll(() => {
        fakeOIDCServer = new FakeOIDCServer(TEST_ISSUER_URL, TEST_CLIENT_ID, TEST_CLIENT_SECRET);
        fakeOIDCServer.start();
    });

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        fakeOIDCServer.reset();
    });

    afterAll(() => {
        fakeOIDCServer.stop();
    });

    it('should return 401 and clear session cookies (access token cleared as non-httpOnly to match how it is stored)', async () => {
        const request = new NextRequest('https://example.com/auth/refresh', {method: 'POST'});

        const response = await POST(request);

        expect(response.status).toBe(401);
        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toMatchObject({value: '', maxAge: 0});
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN].httpOnly).toBeUndefined();
        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toMatchObject({value: '', maxAge: 0, httpOnly: true});
        expect(cookies[COOKIE_SESSION_USER_INFO]).toMatchObject({value: '', maxAge: 0, httpOnly: true});
    });

    it('should return 401 and clear session cookies when the refresh fails upstream', async () => {
        const refreshToken = 'INVALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';
        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        const request = new NextRequest('https://example.com/auth/refresh', {
            method: 'POST',
            headers: {Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`},
        });

        const response = await POST(request);

        expect(response.status).toBe(401);
        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toMatchObject({value: '', maxAge: 0});
        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toMatchObject({value: '', maxAge: 0});
    });

    it('should still return 401 when the OIDC server returns tokens but openid-client rejects them (fake JWTs)', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            id_token: idToken,
            expires_in: 3600,
        }));

        const request = new NextRequest('https://example.com/auth/refresh', {
            method: 'POST',
            headers: {Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`},
        });

        const response = await POST(request);

        expect(response.status).toBe(401);
    });
});
