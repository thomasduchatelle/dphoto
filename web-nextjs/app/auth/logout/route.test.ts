// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {
    ACCESS_TOKEN_COOKIE,
    OAUTH_CODE_VERIFIER_COOKIE,
    OAUTH_STATE_COOKIE,
    REDIRECT_AFTER_LOGIN_COOKIE,
    REFRESH_TOKEN_COOKIE
} from '@/libs/security';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '../../../__tests__/helpers/test-helper-oidc';
import {redirectionOf, setCookiesOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

const deletedCookie = {
    value: '',
    maxAge: 0,
    path: '/',
};

describe('authentication middleware - logout', () => {
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

    it('should clear all cookies and redirect to OIDC logout URL', async () => {
        const request = new NextRequest('https://example.com/auth/logout', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${ACCESS_TOKEN_COOKIE}=VALID_TOKEN; ${REFRESH_TOKEN_COOKIE}=REFRESH_TOKEN`,
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/logout`);
        expect(redirection.params).toMatchObject({
            client_id: TEST_CLIENT_ID,
            logout_uri: 'https://example.com/nextjs/auth/logout-success',
        });

        const cookies = setCookiesOf(response);
        expect(cookies[ACCESS_TOKEN_COOKIE]).toMatchObject(deletedCookie);
        expect(cookies[REFRESH_TOKEN_COOKIE]).toMatchObject(deletedCookie);
        expect(cookies[OAUTH_STATE_COOKIE]).toMatchObject(deletedCookie);
        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]).toMatchObject(deletedCookie);
        expect(cookies[REDIRECT_AFTER_LOGIN_COOKIE]).toMatchObject(deletedCookie);
    });

    it('should redirect to cognito logout URL even if there is no cookies', async () => {
        const request = new NextRequest('https://example.com/auth/logout', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/logout`);
        expect(redirection.params).toMatchObject({
            client_id: TEST_CLIENT_ID,
            logout_uri: 'https://example.com/nextjs/auth/logout-success',
        });
    });
});
