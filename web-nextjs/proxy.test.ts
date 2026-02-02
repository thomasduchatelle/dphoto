// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {proxy, skipProxyForPageMatching} from './proxy';
import {
    COOKIE_AUTH_CODE_VERIFIER,
    COOKIE_AUTH_STATE,
    COOKIE_SESSION_ACCESS_TOKEN,
    COOKIE_SESSION_REFRESH_TOKEN,
    COOKIE_SESSION_USER_INFO
} from '@/libs/security/constants';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {redirectionOf, setCookiesOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('authentication middleware/proxy', () => {
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

    it('should redirect to Cognito authorization when requesting home page without access token', async () => {
        const testRequest = new NextRequest('https://example.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://example.com/nextjs/auth/callback');

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toBeDefined();
        expect(cookies[COOKIE_AUTH_STATE].maxAge).toBe(600);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toBeDefined();
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER].maxAge).toBe(600);
    });

    it('should redirect to Cognito authorization using Forwarded header when behind API Gateway', async () => {
        const testRequest = new NextRequest('https://internal-gateway.my-domain.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
            },
        });

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://my-domain.com/nextjs/auth/callback');

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toBeDefined();
        expect(cookies[COOKIE_AUTH_STATE].maxAge).toBe(600);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toBeDefined();
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER].maxAge).toBe(600);
    });

    it('should allow authenticated request to proceed with valid non-expired token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const accessToken = createCognitoAccessToken({exp: now + 3600}); // expires in 1 hour
        const idToken = 'ID_TOKEN_VALUE';

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${accessToken}; ${COOKIE_SESSION_REFRESH_TOKEN}=SOME_REFRESH_TOKEN; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });

        const response = await proxy(testRequest);

        expect(response.status).toBe(200);
    });

    it('should redirect to Cognito authorization when only refresh token is provided (refresh always fails with fake tokens)', async () => {
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return new tokens (but openid-client will reject fake tokens)
        const now = Math.floor(Date.now() / 1000);
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            id_token: idToken,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });

        const response = await proxy(testRequest);

        // Since openid-client rejects fake tokens, it should redirect to auth
        expect(response.status).toBe(307);
        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
    });

    it('should redirect to Cognito authorization when access token is expired (refresh always fails with fake tokens)', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100}); // expired 100 seconds ago
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return new tokens (but openid-client will reject fake tokens)
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            id_token: idToken,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${expiredAccessToken}; ${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });

        const response = await proxy(testRequest);

        // Since openid-client rejects fake tokens, it should redirect to auth
        expect(response.status).toBe(307);
        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
    });

    it('should redirect to Cognito authorization when refresh token fails', async () => {
        const refreshToken = 'INVALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return an error
        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
    });

    it('should redirect to Cognito authorization when access token is expired and no refresh token is available', async () => {
        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
    });
});

describe('skipProxyForPageMatching regex', () => {
    it('should NOT match static image paths', () => {
        expect(skipProxyForPageMatching.test('images/photo.png')).toBe(false);
        expect(skipProxyForPageMatching.test('images/photo.jpg')).toBe(false);
        expect(skipProxyForPageMatching.test('images/photo.gif')).toBe(false);
        expect(skipProxyForPageMatching.test('static/image.svg')).toBe(false);
        expect(skipProxyForPageMatching.test('photo.png')).toBe(false);
    });

    it('should NOT match /auth/* paths', () => {
        expect(skipProxyForPageMatching.test('auth')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/login')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/callback')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/logout')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/refresh')).toBe(false);
    });

    it('should NOT match /api/* paths', () => {
        expect(skipProxyForPageMatching.test('api/albums')).toBe(false);
        expect(skipProxyForPageMatching.test('api/photos/123')).toBe(false);
    });

    it('should NOT match favicon.ico', () => {
        expect(skipProxyForPageMatching.test('favicon.ico')).toBe(false);
    });

    it('should NOT match Next.js internal paths', () => {
        expect(skipProxyForPageMatching.test('_next/static/chunks/main.js')).toBe(false);
        expect(skipProxyForPageMatching.test('_next/image')).toBe(false);
    });

    it('should NOT match JavaScript files', () => {
        expect(skipProxyForPageMatching.test('scripts/main.js')).toBe(false);
        expect(skipProxyForPageMatching.test('bundle.js')).toBe(false);
    });

    it('should match application pages that require authentication', () => {
        expect(skipProxyForPageMatching.test('')).toBe(true);
        expect(skipProxyForPageMatching.test('albums')).toBe(true);
        expect(skipProxyForPageMatching.test('albums/123')).toBe(true);
        expect(skipProxyForPageMatching.test('photos')).toBe(true);
        expect(skipProxyForPageMatching.test('settings')).toBe(true);
    });
});
