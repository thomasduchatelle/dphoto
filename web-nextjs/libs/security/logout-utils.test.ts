// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi, beforeEach} from 'vitest';
import {clearAuthCookies, getLogoutUrl} from '@/libs/security';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

vi.mock('next/headers', () => {
    const mockCookies = {
        set: vi.fn(),
    };
    const mockHeaders = {
        get: vi.fn((key: string) => {
            if (key === 'host') return 'example.com';
            if (key === 'x-forwarded-proto') return 'https';
            return null;
        }),
    };

    return {
        cookies: vi.fn(() => Promise.resolve(mockCookies)),
        headers: vi.fn(() => Promise.resolve(mockHeaders)),
    };
});

describe('logout-utils', () => {
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

    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('getLogoutUrl', () => {
        it('should generate Cognito logout URL with correct logout_uri parameter', async () => {
            const logoutUrl = await getLogoutUrl();

            const url = new URL(logoutUrl);
            expect(url.origin + url.pathname).toBe(`${TEST_ISSUER_URL}/logout`);
            expect(url.searchParams.get('client_id')).toBe(TEST_CLIENT_ID);
            expect(url.searchParams.get('logout_uri')).toBe('https://example.com/nextjs/auth/logout');
        });
    });

    describe('clearAuthCookies', () => {
        it('should clear all authentication cookies', async () => {
            const {cookies} = await import('next/headers');
            const mockCookieStore = await cookies();

            await clearAuthCookies();

            const cookieOptions = {
                maxAge: 0,
                path: '/',
            };

            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-access-token', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-refresh-token', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-state', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-code-verifier', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-nonce', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-redirect-after-login', '', cookieOptions);
            expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-user-info', '', cookieOptions);
        });
    });
});