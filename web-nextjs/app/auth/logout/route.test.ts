// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {redirectionOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('logout route', () => {
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

    it('should redirect to Cognito logout endpoint with logout_uri parameter', async () => {
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
        expect(redirection.params.client_id).toBe(TEST_CLIENT_ID);
        expect(redirection.params.logout_uri).toBe('https://example.com/nextjs/auth/logout-callback');
    });

    it('should use Forwarded header for logout_uri when behind API Gateway', async () => {
        const request = new NextRequest('https://internal-gateway.my-domain.com/auth/logout', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
            },
        });

        const response = await GET(request);

        expect(response.status).toBe(307);

        const redirection = redirectionOf(response);
        expect(redirection.url).toBe(`${TEST_ISSUER_URL}/logout`);
        expect(redirection.params.logout_uri).toBe('https://my-domain.com/nextjs/auth/logout-callback');
    });
});
