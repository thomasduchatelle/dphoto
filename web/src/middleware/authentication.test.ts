// @vitest-environment node

import {describe, expect, it, vi} from 'vitest';
import type {HandlerContext} from 'waku/dist/lib/middleware/types';
import cookieMiddleware from './authentication';
import {ACCESS_TOKEN_COOKIE, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, REFRESH_TOKEN_COOKIE} from '../core/security';
import {FakeOIDCServer} from './__tests__/fake-oidc-server';
import {createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from './__tests__/test-helper-oidc';

vi.mock('waku', () => ({
    getEnv: (key: string) => {
        const env: Record<string, string> = {
            COGNITO_ISSUER: TEST_ISSUER_URL,
            COGNITO_CLIENT_ID: TEST_CLIENT_ID,
            COGNITO_CLIENT_SECRET: TEST_CLIENT_SECRET,
        };
        return env[key] || '';
    }
}));

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
        const middleware = cookieMiddleware();

        const ctx: HandlerContext = {
            req: new Request('https://example.com/', {
                method: 'GET',
                headers: {
                    'Accept': 'text/html',
                },
            }),
            res: new Response(null),
            data: {},
        } as HandlerContext;

        await middleware(ctx, async () => {});

        expect(ctx.res?.status).toBe(302);
        const location = ctx.res?.headers.get('Location');
        expect(location).toContain('https://cognito.example.com/oauth2/authorize');
        expect(location).toContain('client_id=test-client-id');
        expect(location).toContain(`redirect_uri=${encodeURIComponent('https://example.com/auth/callback')}`);
        expect(location).toContain('scope=openid+profile+email');
        expect(location).toContain('code_challenge_method=S256');
        expect(location).toContain('state=');
        expect(location).toContain('code_challenge=');

        const setCookieHeaders = ctx.res?.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders?.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toBeDefined();
        expect(stateCookie).toContain('Max-Age=300');
        expect(stateCookie).toContain('SameSite=Lax');

        const codeVerifierCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`));
        expect(codeVerifierCookie).toBeDefined();
        expect(codeVerifierCookie).toContain('Max-Age=300');
        expect(codeVerifierCookie).toContain('SameSite=Lax');
    });

    it('should redirect to authorization authority when explicitly requesting /auth/login', async () => {
        const middleware = cookieMiddleware();

        const ctx: HandlerContext = {
            req: new Request('https://example.com/auth/login', {
                method: 'GET',
                headers: {
                    'Accept': 'text/html',
                    'Cookie': `${ACCESS_TOKEN_COOKIE}=VALID_ACCESS_TOKEN`,
                },
            }),
            res: new Response(null),
            data: {},
        } as HandlerContext;

        await middleware(ctx, async () => {});

        expect(ctx.res?.status).toBe(302);
        const location = ctx.res?.headers.get('Location');
        expect(location).toContain('https://cognito.example.com/oauth2/authorize');

        const setCookieHeaders = ctx.res?.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders?.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toBeDefined();

        const codeVerifierCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`));
        expect(codeVerifierCookie).toBeDefined();
    });

    it('should handle OAuth callback with valid authorization code', async () => {
        const authCode = 'AUTH_CODE_123';
        const tokenResponse = createTokenResponse();
        fakeOIDCServer.setupSuccessfulTokenExchange(authCode, tokenResponse);

        const middleware = cookieMiddleware();

        const ctx: HandlerContext = {
            req: new Request(`https://example.com/auth/callback?code=${authCode}&state=EXPECTED_STATE`, {
                method: 'GET',
                headers: {
                    'Accept': 'text/html',
                    'Cookie': `${OAUTH_STATE_COOKIE}=EXPECTED_STATE; ${OAUTH_CODE_VERIFIER_COOKIE}=CODE_VERIFIER_123`,
                },
            }),
            res: new Response(null),
            data: {},
        } as HandlerContext;

        await middleware(ctx, async () => {});

        expect(ctx.res?.status).toBe(302);
        expect(ctx.res?.headers.get('Location')).toBe('https://example.com/');

        const setCookieHeaders = ctx.res?.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();

        const accessTokenCookie = setCookieHeaders?.find(c => c.startsWith(`${ACCESS_TOKEN_COOKIE}=`));
        expect(accessTokenCookie).toContain('ACCESS_TOKEN_VALUE');
        expect(accessTokenCookie).toContain('HttpOnly');
        expect(accessTokenCookie).toContain('Secure');
        expect(accessTokenCookie).toContain('SameSite=Strict');

        const refreshTokenCookie = setCookieHeaders?.find(c => c.startsWith(`${REFRESH_TOKEN_COOKIE}=`));
        expect(refreshTokenCookie).toContain('REFRESH_TOKEN_VALUE');
        expect(refreshTokenCookie).toContain('HttpOnly');
        expect(refreshTokenCookie).toContain('Secure');
        expect(refreshTokenCookie).toContain('SameSite=Strict');

        const stateClearedCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_STATE_COOKIE}=;`));
        expect(stateClearedCookie).toContain('Max-Age=0');

        const codeVerifierClearedCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=;`));
        expect(codeVerifierClearedCookie).toContain('Max-Age=0');
    });

    it('should allow authenticated request to proceed with backendSession', async () => {
        const middleware = cookieMiddleware();

        const accessToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9hdmF0YXIuanBnIiwiU2NvcGVzIjoib3duZXI6dGVzdHVzZXIiLCJleHAiOjk5OTk5OTk5OTl9.signature';

        const ctx: HandlerContext = {
            req: new Request('https://example.com/albums', {
                method: 'GET',
                headers: {
                    'Accept': 'text/html',
                    'Cookie': `${ACCESS_TOKEN_COOKIE}=${accessToken}`,
                },
            }),
            res: new Response(null),
            data: {},
        } as HandlerContext;

        let nextCalled = false;
        await middleware(ctx, async () => {
            nextCalled = true;
        });

        expect(nextCalled).toBe(true);
        expect(ctx.data.backendSession).toBeDefined();
        expect(ctx.data.backendSession.type).toBe('authenticated');
        expect(ctx.data.backendSession.accessToken.accessToken).toBe(accessToken);
        expect(ctx.data.backendSession.accessToken.expiresAt).toBeInstanceOf(Date);
        expect(ctx.data.backendSession.refreshToken).toBe('');
        expect(ctx.data.backendSession.authenticatedUser).toEqual({
            name: 'John Doe',
            email: 'john@example.com',
            picture: 'https://example.com/avatar.jpg',
            isOwner: true,
        });
    });
});
