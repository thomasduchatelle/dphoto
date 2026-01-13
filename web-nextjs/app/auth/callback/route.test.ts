// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {NextRequest} from 'next/server';
import {GET} from './route';
import {COOKIE_AUTH_CODE_VERIFIER, COOKIE_AUTH_NONCE, COOKIE_AUTH_STATE, COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN} from '@/libs/security';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '../../../__tests__/helpers/test-helper-oidc';
import {redirectionOf, setCookiesOf} from '@/__tests__/helpers/test-assertions';

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
                    Cookie: `${COOKIE_AUTH_STATE}=EXPECTED_STATE; ${COOKIE_AUTH_CODE_VERIFIER}=CODE_VERIFIER_123; ${COOKIE_AUTH_NONCE}=NONCE_VALUE`,
                },
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/nextjs/');

        const cookies = setCookiesOf(response);

        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toMatchObject({
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN].value).toBeTruthy();

        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toMatchObject({
            value: 'REFRESH_TOKEN_VALUE',
            httpOnly: true,
            secure: true,
            sameSite: 'lax',
            path: '/',
        });

        expect(cookies[COOKIE_AUTH_STATE]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toMatchObject(deletedCookie);
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
                    Cookie: `${COOKIE_AUTH_STATE}=EXPECTED_STATE; ${COOKIE_AUTH_CODE_VERIFIER}=CODE_VERIFIER_123; ${COOKIE_AUTH_NONCE}=NONCE_VALUE`,
                    'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=my-domain.com;proto=https',
                },
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://my-domain.com/nextjs/');
    });

    it('should redirect to the error page with the same parameters when an error occurred', async () => {
        const request = new NextRequest(
            'https://example.com/auth/callback?error=invalid_request&error_description=user.email%3A+Attribute+cannot+be+updated.',
            {
                method: 'GET',
                headers: {
                    Accept: 'text/html',
                    Cookie: `${COOKIE_AUTH_STATE}=EXPECTED_STATE; ${COOKIE_AUTH_CODE_VERIFIER}=CODE_VERIFIER; ${COOKIE_AUTH_NONCE}=NONCE_VALUE`,
                },
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        const redirection = redirectionOf(response);
        expect(redirection.url).toBe('https://example.com/nextjs/auth/error');
        expect(redirection.params).toEqual({
            error: 'invalid_request',
            error_description: 'user.email: Attribute cannot be updated.',
        });

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_NONCE]).toMatchObject(deletedCookie);
    });

    it('should redirect to error page when the state mismatch', async () => {
        const request = new NextRequest(
            'https://example.com/auth/callback?code=AUTH_CODE&state=WRONG_STATE',
            {
                method: 'GET',
                headers: {
                    Cookie: `${COOKIE_AUTH_STATE}=EXPECTED_STATE; ${COOKIE_AUTH_CODE_VERIFIER}=CODE_VERIFIER; ${COOKIE_AUTH_NONCE}=NONCE_VALUE`,
                },
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        const redirection = redirectionOf(response);
        expect(redirection.url).toBe('https://example.com/nextjs/auth/error');
        expect(redirection.params).toEqual({
            error: 'state-mismatch',
        });

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_NONCE]).toMatchObject(deletedCookie);
    });

    it('should redirect when authentication cookies are not present', async () => {
        const request = new NextRequest(
            'https://example.com/auth/callback?code=AUTH_CODE&state=SOME_STATE',
            {
                method: 'GET',
            }
        );

        const response = await GET(request);

        expect(response.status).toBe(307);
        const redirection = redirectionOf(response);
        expect(redirection.url).toBe('https://example.com/nextjs/auth/error');
        expect(redirection.params).toEqual({
            error: 'missing-authentication-cookies',
        });

        const cookies = setCookiesOf(response);
        expect(cookies[COOKIE_AUTH_STATE]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_CODE_VERIFIER]).toMatchObject(deletedCookie);
        expect(cookies[COOKIE_AUTH_NONCE]).toMatchObject(deletedCookie);
    });
});
