// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {getValidAccessToken} from './access-token-service';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN, COOKIE_SESSION_USER_INFO} from './constants';

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

let mockCookies: Map<string, string>;

vi.mock('next/headers', () => {
    return {
        cookies: vi.fn(() => Promise.resolve({
            get: vi.fn((key: string) => {
                const value = mockCookies.get(key);
                return value ? {value} : undefined;
            }),
            set: vi.fn((key: string, value: string) => {
                mockCookies.set(key, value);
            }),
        })),
    };
});

describe('getValidAccessToken', () => {
    let fakeOIDCServer: FakeOIDCServer;

    beforeAll(() => {
        fakeOIDCServer = new FakeOIDCServer(TEST_ISSUER_URL, TEST_CLIENT_ID, TEST_CLIENT_SECRET);
        fakeOIDCServer.start();
    });

    beforeEach(() => {
        mockCookies = new Map();
        vi.clearAllMocks();
    });

    afterEach(() => {
        fakeOIDCServer.reset();
    });

    afterAll(() => {
        fakeOIDCServer.stop();
    });

    it('should return null when no tokens are provided', async () => {
        const result = await getValidAccessToken();

        expect(result).toBeNull();
    });

    it('should return null when only access token is provided without refresh token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const validAccessToken = createCognitoAccessToken({exp: now + 3600});

        mockCookies.set(COOKIE_SESSION_ACCESS_TOKEN, validAccessToken);

        const result = await getValidAccessToken();

        expect(result).toBeNull();
    });

    it('should return valid access token when access token is valid and not expired', async () => {
        const now = Math.floor(Date.now() / 1000);
        const validAccessToken = createCognitoAccessToken({exp: now + 3600});
        const idToken = 'ID_TOKEN_VALUE';

        mockCookies.set(COOKIE_SESSION_ACCESS_TOKEN, validAccessToken);
        mockCookies.set(COOKIE_SESSION_REFRESH_TOKEN, 'SOME_REFRESH_TOKEN');
        mockCookies.set(COOKIE_SESSION_USER_INFO, idToken);

        const result = await getValidAccessToken();

        expect(result).not.toBeNull();
        expect(result?.accessToken.accessToken).toBe(validAccessToken);
        expect(result?.idToken).toBe(idToken);
    });

    it('should refresh and return new access token when access token is expired but refresh token is valid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100});
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            expires_in: 3600,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        mockCookies.set(COOKIE_SESSION_ACCESS_TOKEN, expiredAccessToken);
        mockCookies.set(COOKIE_SESSION_REFRESH_TOKEN, refreshToken);
        mockCookies.set(COOKIE_SESSION_USER_INFO, idToken);

        const result = await getValidAccessToken();

        expect(result).not.toBeNull();
        expect(result?.accessToken.accessToken).toBe(newTokenResponse.access_token);
        expect(result?.idToken).toBe(idToken);
    });

    it('should refresh and return new access token when no access token is provided but refresh token is valid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            expires_in: 3600,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        mockCookies.set(COOKIE_SESSION_REFRESH_TOKEN, refreshToken);
        mockCookies.set(COOKIE_SESSION_USER_INFO, idToken);

        const result = await getValidAccessToken();

        expect(result).not.toBeNull();
        expect(result?.accessToken.accessToken).toBe(newTokenResponse.access_token);
        expect(result?.idToken).toBe(idToken);
    });

    it('should return null when refresh token is invalid', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredAccessToken = createCognitoAccessToken({exp: now - 100});
        const refreshToken = 'INVALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        mockCookies.set(COOKIE_SESSION_ACCESS_TOKEN, expiredAccessToken);
        mockCookies.set(COOKIE_SESSION_REFRESH_TOKEN, refreshToken);
        mockCookies.set(COOKIE_SESSION_USER_INFO, idToken);

        const result = await getValidAccessToken();

        expect(result).toBeNull();
    });
});
