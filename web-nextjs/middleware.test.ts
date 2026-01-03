// @vitest-environment node

import { describe, expect, it, beforeAll, afterEach, afterAll, vi } from 'vitest';
import { NextRequest } from 'next/server';
import { middleware } from './middleware';
import {
    ACCESS_TOKEN_COOKIE,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    REFRESH_TOKEN_COOKIE,
} from './lib/security/constants';
import { FakeOIDCServer } from './__tests__/helpers/fake-oidc-server';
import {
    createTokenResponse,
    TEST_CLIENT_ID,
    TEST_CLIENT_SECRET,
    TEST_ISSUER_URL,
} from './__tests__/helpers/test-helper-oidc';

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

        expect(response.status).toBe(307); // NextJS uses 307 for temporary redirects
        const location = response.headers.get('Location');
        expect(location).toContain('https://cognito.example.com/oauth2/authorize');
        expect(location).toContain('client_id=test-client-id');
        expect(location).toContain(`redirect_uri=${encodeURIComponent('https://example.com/auth/callback')}`);
        expect(location).toContain('scope=openid+profile+email');
        expect(location).toContain('code_challenge_method=S256');
        expect(location).toContain('state=');
        expect(location).toContain('code_challenge=');

        const setCookieHeaders = response.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders.find((c) => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toBeDefined();
        expect(stateCookie).toContain('Max-Age=300');
        expect(stateCookie).toMatch(/SameSite=(Lax|lax)/i);

        const codeVerifierCookie = setCookieHeaders.find((c) =>
            c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`)
        );
        expect(codeVerifierCookie).toBeDefined();
        expect(codeVerifierCookie).toContain('Max-Age=300');
        expect(codeVerifierCookie).toMatch(/SameSite=(Lax|lax)/i);
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
        const location = response.headers.get('Location');
        expect(location).toContain('https://cognito.example.com/oauth2/authorize');

        const setCookieHeaders = response.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders.find((c) => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toBeDefined();

        const codeVerifierCookie = setCookieHeaders.find((c) =>
            c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`)
        );
        expect(codeVerifierCookie).toBeDefined();
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

        const setCookieHeaders = response.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();

        const accessTokenCookie = setCookieHeaders.find((c) => c.startsWith(`${ACCESS_TOKEN_COOKIE}=`));
        expect(accessTokenCookie).toContain('ACCESS_TOKEN_VALUE');
        expect(accessTokenCookie).toContain('HttpOnly');
        expect(accessTokenCookie).toContain('Secure');
        expect(accessTokenCookie).toMatch(/SameSite=(Strict|strict)/i);

        const refreshTokenCookie = setCookieHeaders.find((c) =>
            c.startsWith(`${REFRESH_TOKEN_COOKIE}=`)
        );
        expect(refreshTokenCookie).toContain('REFRESH_TOKEN_VALUE');
        expect(refreshTokenCookie).toContain('HttpOnly');
        expect(refreshTokenCookie).toContain('Secure');
        expect(refreshTokenCookie).toMatch(/SameSite=(Strict|strict)/i);

        const stateClearedCookie = setCookieHeaders.find((c) => c.startsWith(`${OAUTH_STATE_COOKIE}=;`));
        expect(stateClearedCookie).toContain('Max-Age=0');

        const codeVerifierClearedCookie = setCookieHeaders.find((c) =>
            c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=;`)
        );
        expect(codeVerifierClearedCookie).toContain('Max-Age=0');
    });

    it('should allow authenticated request to proceed with backendSession', async () => {
        const accessToken =
            'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9hdmF0YXIuanBnIiwiU2NvcGVzIjoib3duZXI6dGVzdHVzZXIiLCJleHAiOjk5OTk5OTk5OTl9.signature';

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${ACCESS_TOKEN_COOKIE}=${accessToken}`,
            },
        });

        const response = await middleware(request);

        // NextJS middleware returns a NextResponse object when continuing
        expect(response.status).toBe(200);

        // Check that backendSession was added to headers
        const backendSessionHeader = response.headers.get('x-backend-session');
        expect(backendSessionHeader).toBeDefined();

        const backendSession = JSON.parse(backendSessionHeader!);
        expect(backendSession).toBeDefined();
        expect(backendSession.type).toBe('authenticated');
        expect(backendSession.accessToken.accessToken).toBe(accessToken);
        expect(new Date(backendSession.accessToken.expiresAt)).toBeInstanceOf(Date);
        expect(backendSession.refreshToken).toBe('');
        expect(backendSession.authenticatedUser).toEqual({
            name: 'John Doe',
            email: 'john@example.com',
            picture: 'https://example.com/avatar.jpg',
            isOwner: true,
        });
    });
});
