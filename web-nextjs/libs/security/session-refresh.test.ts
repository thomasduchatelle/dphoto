// @vitest-environment node

import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {refreshSessionIfNecessary, SessionCookies} from './session-refresh';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('refreshSessionIfNecessary', () => {
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

    it('should return active status when access token is valid and not expired', async () => {
        const now = Math.floor(Date.now() / 1000);
        const validAccessToken = createCognitoAccessToken({exp: now + 3600});

        const cookies: SessionCookies = {
            accessToken: validAccessToken,
            refreshToken: 'SOME_REFRESH_TOKEN',
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('active');
        expect(result.newAccessToken).toBeUndefined();
        expect(result.newRefreshToken).toBeUndefined();
    });

    it('should return none status when no tokens are provided', async () => {
        const cookies: SessionCookies = {};

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('none');
    });

    it('should return none status when only expired access token is provided without refresh token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100});

        const cookies: SessionCookies = {
            accessToken: expiredAccessToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('none');
    });

    it('should refresh and return active status when access token is expired but refresh token is valid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100});
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            expires_in: 3600,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const cookies: SessionCookies = {
            accessToken: expiredAccessToken,
            refreshToken: refreshToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('active');
        expect(result.newAccessToken).toBeDefined();
        expect(result.newAccessToken?.token).toBe(newTokenResponse.access_token);
        expect(result.newAccessToken?.expiresIn).toBe(3600);
        expect(result.newRefreshToken).toBe('NEW_REFRESH_TOKEN');
    });

    it('should refresh and return active status when no access token is provided but refresh token is valid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            expires_in: 3600,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const cookies: SessionCookies = {
            refreshToken: refreshToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('active');
        expect(result.newAccessToken).toBeDefined();
        expect(result.newAccessToken?.token).toBe(newTokenResponse.access_token);
        expect(result.newAccessToken?.expiresIn).toBe(3600);
        expect(result.newRefreshToken).toBe('NEW_REFRESH_TOKEN');
    });

    it('should return expired status when refresh token is invalid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100});
        const refreshToken = 'INVALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return an error
        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        const cookies: SessionCookies = {
            accessToken: expiredAccessToken,
            refreshToken: refreshToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('expired');
        expect(result.newAccessToken).toBeUndefined();
        expect(result.newRefreshToken).toBeUndefined();
    });

    it('should handle token response without refresh token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens without refresh token
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            expires_in: 3600,
        });
        delete newTokenResponse.refresh_token;
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const cookies: SessionCookies = {
            refreshToken: refreshToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('active');
        expect(result.newAccessToken).toBeDefined();
        expect(result.newRefreshToken).toBeUndefined();
    });

    it('should handle token response without expires_in', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';

        // Setup fake OIDC server to return new tokens without expires_in
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
        });
        delete newTokenResponse.expires_in;
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const cookies: SessionCookies = {
            refreshToken: refreshToken,
        };

        const result = await refreshSessionIfNecessary(cookies);

        expect(result.status).toBe('active');
        expect(result.newAccessToken).toBeDefined();
        expect(result.newAccessToken?.expiresIn).toBe(0);
    });
});
