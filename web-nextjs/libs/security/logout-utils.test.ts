// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {getLogoutUrl} from './logout-utils';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {clearFullSession} from "@/libs/security/backend-store";
import {fakeNextHeaders} from "@/__tests__/helpers/fake-next-headers";
import {NextRequest} from "next/server";

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
            fakeHeaders.withRequest(new NextRequest('http://cloudfront.example.com/auth/logout', {
                method: 'GET',
                headers: {
                    Accept: 'text/html',
                    'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com;proto=https',
                },
            }))
            const logoutUrl = await getLogoutUrl();

            const url = new URL(logoutUrl);
            expect(url.origin + url.pathname).toBe(`${TEST_ISSUER_URL}/logout`);
            expect(url.searchParams.get('client_id')).toBe(TEST_CLIENT_ID);
            expect(url.searchParams.get('logout_uri')).toBe('https://example.com/nextjs/auth/logout');
        });
    });

    describe('clearAuthCookies', () => {
        it('should clear all authentication cookies', async () => {
            await clearFullSession();

            const deleteCookie = {
                maxAge: 0,
                path: '/',
            };

            expect(fakeHeaders.getSetCookie('dphoto-access-token')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-refresh-token')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-user-info')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-oauth-state')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-oauth-code-verifier')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-oauth-nonce')).toStrictEqual({value: '', options: deleteCookie});
            expect(fakeHeaders.getSetCookie('dphoto-redirect-after-login')).toStrictEqual({value: '', options: deleteCookie});
        });
    });
});