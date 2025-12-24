import {describe, it, expect, vi, beforeEach} from 'vitest';
import type {HandlerContext} from 'waku/dist/lib/middleware/types';
import type * as client from 'openid-client';
import cookieMiddleware from './authentication';
import {
    ACCESS_TOKEN_COOKIE,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    REFRESH_TOKEN_COOKIE
} from '../core/security';

vi.mock('waku', () => ({
    getEnv: (key: string) => {
        const env: Record<string, string> = {
            COGNITO_ISSUER: 'https://cognito.example.com',
            COGNITO_CLIENT_ID: 'test-client-id',
            COGNITO_CLIENT_SECRET: 'test-client-secret',
        };
        return env[key] || '';
    }
}));

const mockConfig: client.Configuration = {
    authorization_endpoint: 'https://cognito.example.com/oauth2/authorize',
    token_endpoint: 'https://cognito.example.com/oauth2/token',
    issuer: 'https://cognito.example.com',
} as client.Configuration;

let mockRandomPKCECodeVerifier: () => string;
let mockRandomState: () => string;
let mockCalculatePKCECodeChallenge: (verifier: string) => Promise<string>;
let mockBuildAuthorizationUrl: (config: client.Configuration, params: Record<string, string>) => URL;
let mockAuthorizationCodeGrant: (config: client.Configuration, url: URL, options: any) => Promise<client.TokenEndpointResponse>;
let mockDiscovery: (issuer: URL, clientId: string, clientSecret: string) => Promise<client.Configuration>;

vi.mock('openid-client', () => {
    return {
        discovery: (...args: any[]) => mockDiscovery(...args),
        randomPKCECodeVerifier: () => mockRandomPKCECodeVerifier(),
        randomState: () => mockRandomState(),
        calculatePKCECodeChallenge: (verifier: string) => mockCalculatePKCECodeChallenge(verifier),
        buildAuthorizationUrl: (config: any, params: any) => mockBuildAuthorizationUrl(config, params),
        authorizationCodeGrant: (config: any, url: any, options: any) => mockAuthorizationCodeGrant(config, url, options),
    };
});

describe('authentication middleware', () => {
    beforeEach(() => {
        vi.resetAllMocks();

        mockDiscovery = vi.fn().mockResolvedValue(mockConfig);
        mockRandomPKCECodeVerifier = vi.fn().mockReturnValue('CODE_VERIFIER_123');
        mockRandomState = vi.fn().mockReturnValue('GENERATED_STATE_123');
        mockCalculatePKCECodeChallenge = vi.fn().mockResolvedValue('CODE_CHALLENGE_HASH');
        mockBuildAuthorizationUrl = vi.fn().mockReturnValue(
            new URL('https://cognito.example.com/oauth2/authorize?client_id=test-client-id&redirect_uri=https://example.com/auth/callback&scope=openid+profile+email&code_challenge=CODE_CHALLENGE_HASH&code_challenge_method=S256&state=GENERATED_STATE_123')
        );
        mockAuthorizationCodeGrant = vi.fn().mockResolvedValue({
            access_token: 'ACCESS_TOKEN_VALUE',
            refresh_token: 'REFRESH_TOKEN_VALUE',
            expires_in: 3600,
            id_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9hdmF0YXIuanBnIiwiU2NvcGVzIjoib3duZXI6dGVzdHVzZXIifQ.signature',
            token_type: 'Bearer',
        } as client.TokenEndpointResponse);
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
        expect(ctx.res?.headers.get('Location')).toContain('https://cognito.example.com/oauth2/authorize');
        expect(ctx.res?.headers.get('Location')).toContain('client_id=test-client-id');
        expect(ctx.res?.headers.get('Location')).toContain('redirect_uri=https://example.com/auth/callback');
        expect(ctx.res?.headers.get('Location')).toContain('scope=openid+profile+email');
        expect(ctx.res?.headers.get('Location')).toContain('code_challenge=CODE_CHALLENGE_HASH');
        expect(ctx.res?.headers.get('Location')).toContain('code_challenge_method=S256');
        expect(ctx.res?.headers.get('Location')).toContain('state=GENERATED_STATE_123');

        const setCookieHeaders = ctx.res?.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders?.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toContain('GENERATED_STATE_123');
        expect(stateCookie).toContain('Max-Age=300');
        expect(stateCookie).toContain('SameSite=Lax');

        const codeVerifierCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`));
        expect(codeVerifierCookie).toContain('CODE_VERIFIER_123');
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
        expect(ctx.res?.headers.get('Location')).toContain('https://cognito.example.com/oauth2/authorize');

        const setCookieHeaders = ctx.res?.headers.getSetCookie();
        expect(setCookieHeaders).toBeDefined();
        expect(setCookieHeaders?.length).toBeGreaterThanOrEqual(2);

        const stateCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_STATE_COOKIE}=`));
        expect(stateCookie).toBeDefined();

        const codeVerifierCookie = setCookieHeaders?.find(c => c.startsWith(`${OAUTH_CODE_VERIFIER_COOKIE}=`));
        expect(codeVerifierCookie).toBeDefined();
    });

    it('should handle OAuth callback with valid authorization code', async () => {
        const middleware = cookieMiddleware();

        const ctx: HandlerContext = {
            req: new Request('https://example.com/auth/callback?code=AUTH_CODE_123&state=EXPECTED_STATE', {
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
