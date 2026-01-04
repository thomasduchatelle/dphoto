// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {ACCESS_TOKEN_COOKIE, OAUTH_CODE_VERIFIER_COOKIE, OAUTH_STATE_COOKIE, REFRESH_TOKEN_COOKIE} from '@/libs/security';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '../../../__tests__/helpers/test-helper-oidc';
import {setCookiesOf} from '@/__tests__/helpers/test-assertions';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

const deletedCookie = {
    value: '',
    maxAge: 0,
    path: '/',
};

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

        const response = await GET(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/nextjs/');

        const cookies = setCookiesOf(response);

        expect(cookies[ACCESS_TOKEN_COOKIE]).toMatchObject({
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });
        expect(cookies[ACCESS_TOKEN_COOKIE].value).toBeTruthy();

        expect(cookies[REFRESH_TOKEN_COOKIE]).toMatchObject({
            value: 'REFRESH_TOKEN_VALUE',
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });

        expect(cookies[OAUTH_STATE_COOKIE]).toMatchObject(deletedCookie);
        expect(cookies[OAUTH_CODE_VERIFIER_COOKIE]).toMatchObject(deletedCookie);
    });

    it('should redirect to original domain when using Forwarded header', async () => {
        const authCode = 'AUTH_CODE_456';
        const tokenResponse = createTokenResponse();
        fakeOIDCServer.setupSuccessfulTokenExchange(authCode, tokenResponse);

        const request = new NextRequest(
            `https://internal-gateway.my-domain.com/auth/callback?code=${authCode}&state=EXPECTED_STATE`,
            {
                method: 'GET',
                headers: {
                    Accept: 'text/html',
                    Cookie: `${OAUTH_STATE_COOKIE}=EXPECTED_STATE; ${OAUTH_CODE_VERIFIER_COOKIE}=CODE_VERIFIER_123`,
                    'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
                },
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://my-domain.com/nextjs/');
    });
});
