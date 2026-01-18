// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {proxy, skipProxyForPageMatching} from './proxy';
import {COOKIE_AUTH_CODE_VERIFIER, COOKIE_AUTH_STATE, COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN} from '@/libs/security/constants';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {fakeNextHeaders} from "@/__tests__/helpers/fake-next-headers";

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

const fakeHeaders = fakeNextHeaders()

vi.mock('next/headers', () => {
    return {
        cookies: vi.fn(() => fakeHeaders.mock().cookies()),
        headers: vi.fn(() => fakeHeaders.mock().headers()),
    };
});

describe('authentication middleware/proxy', () => {
    let fakeOIDCServer: FakeOIDCServer;

    beforeAll(() => {
        fakeOIDCServer = new FakeOIDCServer(TEST_ISSUER_URL, TEST_CLIENT_ID, TEST_CLIENT_SECRET);
        fakeOIDCServer.start();
    });

    beforeEach(() => {
        vi.clearAllMocks();
        fakeHeaders.reset()
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
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.client_id).toBe(TEST_CLIENT_ID);
        expect(redirection.params.redirect_uri).toBe('https://example.com/nextjs/auth/callback');
        expect(redirection.params.scope).toBe('openid profile email');
        expect(redirection.params.code_challenge_method).toBe('S256');
        expect(redirection.params.state).toBeDefined();
        expect(redirection.params.code_challenge).toBeDefined();

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toMatchObject({
            maxAge: 300,
            sameSite: 'lax',
            path: '/',
            value: redirection.params.state,
        });

        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toMatchObject({
            maxAge: 300,
            sameSite: 'lax',
            path: '/',
        });
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER].value).toBeTruthy();
    });

    it('should redirect to Cognito authorization using Forwarded header when behind API Gateway', async () => {
        const request = new NextRequest('https://internal-gateway.my-domain.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
            },
        });
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://my-domain.com/nextjs/auth/callback');

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]?.value).toBeTruthy();
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]?.value).toBeTruthy();
    });

    it('should allow authenticated request to proceed with valid non-expired token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const accessToken = createCognitoAccessToken({exp: now + 3600}); // expires in 1 hour
        const idToken = 'ID_TOKEN_VALUE';

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${accessToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });
        fakeHeaders.withRequest(testRequest)
        fakeHeaders.setCookie(COOKIE_SESSION_REFRESH_TOKEN, 'SOME_REFRESH_TOKEN');

        const response = await proxy(testRequest);

        expect(response.status).toBe(200);
    });

    it('should use refresh token to get new access token when only refresh token is provided', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        // After successful token refresh, the request should be allowed through
        expect(response.status).toBe(200);
    });

    it('should refresh expired access token with valid refresh token and allow request to proceed', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100}); // expired 100 seconds ago
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${expiredAccessToken}; ${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}; ${COOKIE_SESSION_USER_INFO}=${idToken}`,
            },
        });
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        // After successful token refresh, the request should be allowed through
        expect(response.status).toBe(200);
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
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://example.com/nextjs/auth/callback');
    });

    it('should redirect to Cognito authorization when access token is expired and no refresh token is available', async () => {
        const testRequest = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });
        fakeHeaders.withRequest(testRequest)

        const response = await proxy(testRequest);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://example.com/nextjs/auth/callback');
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
