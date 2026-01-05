// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE} from '@/lib/security/constants';
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
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/oauth2/authorize`);

        const cookies = setCookiesOf(response);
        expect(cookies[OAUTH_STATE_COOKIE]?.value).toBeTruthy();
        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]?.value).toBeTruthy();
    });

    it('should use DPHOTO_DOMAIN_NAME for redirect_uri when environment variable is set', async () => {
        vi.stubEnv('DPHOTO_DOMAIN_NAME', 'dphoto.example.com');

        const request = new NextRequest('https://cloudfront-distribution.cloudfront.net/auth/login', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.params.redirect_uri).toBe('https://dphoto.example.com/nextjs/auth/callback');

        vi.unstubAllEnvs();
    });
});
