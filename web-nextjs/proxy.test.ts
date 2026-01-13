// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {proxy, skipProxyForPageMatching} from './proxy';
import {COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN} from '@/libs/security/constants';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {redirectionOf, setCookiesOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('authentication middleware', () => {
    let fakeOIDCServer: FakeOIDCServer;

    beforeAll(() => {
        fakeOIDCServer = new FakeOIDCServer(TEST_ISSUER_URL, TEST_CLIENT_ID, TEST_CLIENT_SECRET);
        fakeOIDCServer.start();
    });

    afterEach(() => {
        fakeOIDCServer.reset();
    });

    afterAll(() => {
        fakeOIDCServer.stop();
    });

    it('should redirect to login page when requesting home page without access token', async () => {
        const request = new NextRequest('https://example.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await proxy(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/nextjs/auth/login');
    });

    it('should redirect to login page using Forwarded header when behind API Gateway', async () => {
        const request = new NextRequest('https://internal-gateway.my-domain.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
            },
        });

        const response = await proxy(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://my-domain.com/nextjs/auth/login');
    });

    it('should allow authenticated request to proceed with valid non-expired token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const accessToken = createCognitoAccessToken({exp: now + 3600}); // expires in 1 hour

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${accessToken}`,
            },
        });

        const response = await proxy(request);

        expect(response.status).toBe(200);
    });

    it('should use refresh token to get new access token when only refresh token is provided', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}`,
            },
        });

        const response = await proxy(request);

        expect(redirectionOf(response).url).toBe('https://example.com/nextjs/albums');

        // Check that new tokens were set in cookies
        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toBeDefined();
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toMatchObject({
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });

        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toBeDefined();
        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN].value).toBe('NEW_REFRESH_TOKEN');
    });

    it('should refresh expired access token with valid refresh token and allow request to proceed', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100}); // expired 100 seconds ago
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_ACCESS_TOKEN}=${expiredAccessToken}; ${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}`,
            },
        });

        const response = await proxy(request);

        expect(redirectionOf(response).url).toBe('https://example.com/nextjs/albums');

        // Check that new tokens were set in cookies
        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toBeDefined();
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN].value).not.toBe(expiredAccessToken);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toMatchObject({
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });

        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toBeDefined();
        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN].value).toBe('NEW_REFRESH_TOKEN');
    });

    it('should redirect to login when refresh token fails', async () => {
        const refreshToken = 'INVALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return an error
        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${COOKIE_SESSION_REFRESH_TOKEN}=${refreshToken}`,
            },
        });

        const response = await proxy(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/nextjs/auth/login');
    });

    it('should redirect to login when access token is expired and no refresh token is available', async () => {
        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await proxy(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/nextjs/auth/login');
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
