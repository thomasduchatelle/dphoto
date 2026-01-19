// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {getLogoutUrl} from './logout-utils';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {clearFullSession} from "@/libs/security/backend-store";
import {fakeNextHeaders} from "@/__tests__/helpers/fake-next-headers";
import {NextRequest} from "next/server";
import {
    COOKIE_AUTH_CODE_VERIFIER, COOKIE_AUTH_NONCE, COOKIE_AUTH_REDIRECT_AFTER_LOGIN,
    COOKIE_AUTH_STATE,
    COOKIE_SESSION_ACCESS_TOKEN,
    COOKIE_SESSION_REFRESH_TOKEN,
    COOKIE_SESSION_USER_INFO
} from "@/libs/security/constants";

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

    describe('clearFullSession', () => {
        it('should return cookie values to clear all authentication cookies', () => {
            const cookies = clearFullSession();

            const deleteCookieValue = {value: '', maxAge: 0};

            expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_SESSION_USER_INFO]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_AUTH_STATE]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_AUTH_NONCE]).toStrictEqual(deleteCookieValue);
            expect(cookies[COOKIE_AUTH_REDIRECT_AFTER_LOGIN]).toStrictEqual(deleteCookieValue);
        });
    });
});