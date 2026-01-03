// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {middleware} from './middleware';
import {ACCESS_TOKEN_COOKIE, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, REFRESH_TOKEN_COOKIE,} from './lib/security/constants';
import {FakeOIDCServer} from './__tests__/helpers/fake-oidc-server';
import {createBackendAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL,} from './__tests__/helpers/test-helper-oidc';
import {redirectionOf, setCookiesOf} from './__tests__/helpers/test-assertions';

// Mock environment variables
vi.stubEnv('COGNITO_ISSUER', TEST_ISSUER_URL);
vi.stubEnv('COGNITO_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('COGNITO_CLIENT_SECRET', TEST_CLIENT_SECRET);

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

    it('should redirect to authorization authority when requesting home page without access token', async () => {
        const request = new NextRequest('https://example.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await middleware(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.client_id).toBe(TEST_CLIENT_ID);
        expect(redirection.params.redirect_uri).toBe('https://example.com/auth/callback');
        expect(redirection.params.scope).toBe('openid profile email');
        expect(redirection.params.code_challenge_method).toBe('S256');
        expect(redirection.params.state).toBeDefined();
        expect(redirection.params.code_challenge).toBeDefined();

        const cookies = setCookiesOf(response);
        expect(cookies[OAUTH_STATE_COOKIE]).toMatchObject({
            maxAge: 300,
            sameSite: 'lax',
            path: '/',
            value: redirection.params.state,
        });

        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]).toMatchObject({
            maxAge: 300,
            sameSite: 'lax',
            path: '/',
        });
        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE].value).toBeTruthy();
    });

    it('should redirect to authorization authority when explicitly requesting /auth/login', async () => {
        const request = new NextRequest('https://example.com/auth/login', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${ACCESS_TOKEN_COOKIE}=VALID_ACCESS_TOKEN`,
            },
        });

        const response = await middleware(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);

        const cookies = setCookiesOf(response);
        expect(cookies[OAUTH_STATE_COOKIE]?.value).toBeTruthy();
        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]?.value).toBeTruthy();
    });

    it('should handle OAuth callback with valid authorization code', async () => {
        const authCode = 'AUTH_CODE_123';
        const tokenResponse = createTokenResponse();
        fakeOIDCServer.setupSuccessfulTokenExchange(authCode, tokenResponse);

        const request = new NextRequest(
            `https://example.com/auth/callback?code=${authCode}&state=EXPECTED_STATE`,
            {
                method: 'GET',
                headers: {
                    Accept: 'text/html',
                    Cookie: `${OAUTH_STATE_COOKIE}=EXPECTED_STATE; ${OAUTH_CODE_VERIFIER_COOKIE}=CODE_VERIFIER_123`,
                },
            }
        );

        const response = await middleware(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/');

        const cookies = setCookiesOf(response);

        expect(cookies[ACCESS_TOKEN_COOKIE]).toMatchObject({
            httpOnly: true,
            secure: true,
            sameSite: 'strict',
            path: '/',
        });
        expect(cookies[ACCESS_TOKEN_COOKIE].value).toBeTruthy();

        expect(cookies[REFRESH_TOKEN_COOKIE]).toMatchObject({
            value: 'REFRESH_TOKEN_VALUE',
            httpOnly: true,
            secure: true,
            sameSite: 'strict',
            path: '/',
        });

        expect(cookies[OAUTH_STATE_COOKIE]).toMatchObject({
            value: '',
            maxAge: 0,
            path: '/',
        });

        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]).toMatchObject({
            value: '',
            maxAge: 0,
            path: '/',
        });
    });

    it('should allow authenticated request to proceed with backendSession', async () => {
        const accessToken = createBackendAccessToken({
            email: 'user@example.com',
            Scopes: 'owner:user@example.com',
        });

        const userInfoCookie = JSON.stringify({
            name: 'Test User',
            email: 'user@example.com',
            picture: 'https://example.com/avatar.jpg',
        });

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${ACCESS_TOKEN_COOKIE}=${accessToken}; dphoto-user-info=${encodeURIComponent(userInfoCookie)}`,
            },
        });

        const response = await middleware(request);

        expect(response.status).toBe(200);

        const backendSessionHeader = response.headers.get('x-backend-session');
        expect(backendSessionHeader).toBeDefined();

        const backendSession = JSON.parse(backendSessionHeader!);
        expect(backendSession).toBeDefined();
        expect(backendSession.type).toBe('authenticated');
        expect(backendSession.accessToken.accessToken).toBe(accessToken);
        expect(new Date(backendSession.accessToken.expiresAt)).toBeInstanceOf(Date);
        expect(backendSession.refreshToken).toBe('');
        expect(backendSession.authenticatedUser).toEqual({
            name: 'Test User',
            email: 'user@example.com',
            picture: 'https://example.com/avatar.jpg',
            isOwner: true,
        });
    });
});
