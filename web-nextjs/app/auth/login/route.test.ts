// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {COOKIE_AUTH_CODE_VERIFIER, COOKIE_AUTH_STATE} from '@/libs/security/constants';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
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

    it('should redirect to authorization authority when requesting home page without access token', async () => {
        const request = new NextRequest('https://example.com/auth/login', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await GET(request);

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

    it('should redirect to authorization authority when explicitly requesting /auth/login', async () => {
        const request = new NextRequest('https://example.com/auth/login', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]?.value).toBeTruthy();
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]?.value).toBeTruthy();
    });

    it('should use Forwarded header for redirect_uri when behind API Gateway', async () => {
        const request = new NextRequest('https://internal-gateway.my-domain.com/auth/login', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);
        expect(redirection.params.redirect_uri).toBe('https://my-domain.com/nextjs/auth/callback');
    });
});
