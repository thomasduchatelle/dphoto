// @vitest-environment node

import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {parseCurrentAccessToken, refreshSession} from './access-token-service';
import {FakeOIDCServer} from '@/__tests__/helpers/fake-oidc-server';
import {createCognitoAccessToken, createTokenResponse, TEST_CLIENT_ID, TEST_CLIENT_SECRET, TEST_ISSUER_URL} from '@/__tests__/helpers/test-helper-oidc';
import {COOKIE_SESSION_ACCESS_TOKEN, COOKIE_SESSION_REFRESH_TOKEN, COOKIE_SESSION_USER_INFO} from './constants';
import {ReadCookieStore} from "@/libs/nextjs-cookies";

vi.stubEnv('OAUTH_ISSUER_URL', TEST_ISSUER_URL);
vi.stubEnv('OAUTH_CLIENT_ID', TEST_CLIENT_ID);
vi.stubEnv('OAUTH_CLIENT_SECRET', TEST_CLIENT_SECRET);

describe('parseCurrentAccessToken', () => {
    it('should return null when no token is provided', async () => {
        const result = await parseCurrentAccessToken(undefined);

        expect(result).toBeNull();
    });

    it('should return null when token is malformed', async () => {
        const result = await parseCurrentAccessToken('invalid-token');

        expect(result).toBeNull();
    });

    it('should parse valid access token and return parsed token with claims', async () => {
        const now = Math.floor(Date.now() / 1000);
        const validAccessToken = createCognitoAccessToken({exp: now + 3600});

        const result = await parseCurrentAccessToken(validAccessToken);

        expect(result).not.toBeNull();
        expect(result?.accessToken).toBe(validAccessToken);
        expect(result?.expiresAt).toBeInstanceOf(Date);
        expect(result?.expiresAt.getTime()).toBeGreaterThan(Date.now());
        expect(result?.aboutToExpire).toBe(false);
    });

    it('should detect when token is about to expire (less than 5 minutes)', async () => {
        const now = Math.floor(Date.now() / 1000);
        const soonToExpireToken = createCognitoAccessToken({exp: now + 60}); // expires in 1 minute

        const result = await parseCurrentAccessToken(soonToExpireToken);

        expect(result).not.toBeNull();
        expect(result?.aboutToExpire).toBe(true);
    });

    it('should detect expired token', async () => {
        const now = Math.floor(Date.now() / 1000);
        const expiredToken = createCognitoAccessToken({exp: now - 100}); // expired 100 seconds ago

        const result = await parseCurrentAccessToken(expiredToken);

        expect(result).not.toBeNull();
        expect(result?.expiresAt.getTime()).toBeLessThan(Date.now());
        expect(result?.aboutToExpire).toBe(true);
    });
});

describe('refreshSession', () => {
    let fakeOIDCServer: FakeOIDCServer;

    beforeAll(() => {
        fakeOIDCServer = new FakeOIDCServer(TEST_ISSUER_URL, TEST_CLIENT_ID, TEST_CLIENT_SECRET);
        fakeOIDCServer.start();
    });

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        fakeOIDCServer.reset();
    });

    afterAll(() => {
        fakeOIDCServer.stop();
    });

    it('should return failure and clear cookies when no refresh token is provided', async () => {
        const cookieStore: ReadCookieStore = {
            get: (name: string) => undefined
        };

        const result = await refreshSession(cookieStore);

        expect(result.success).toBe(false);
        expect(result.cookies[COOKIE_SESSION_ACCESS_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_REFRESH_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_USER_INFO]).toEqual({value: '', maxAge: 0});
    });

    it('should return failure when openid-client rejects the token response (fake tokens)', async () => {
        const now = Math.floor(Date.now() / 1000);
        const refreshToken = 'VALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        // Setup fake OIDC server to return new tokens (but openid-client will reject them as invalid)
        const newTokenResponse = createTokenResponse({
            access_token: createCognitoAccessToken({exp: now + 3600}),
            refresh_token: 'NEW_REFRESH_TOKEN',
            expires_in: 3600,
            id_token: idToken,
        });
        fakeOIDCServer.setupSuccessfulRefreshTokenExchange(refreshToken, newTokenResponse);

        const cookieStore: ReadCookieStore = {
            get: (name: string) => {
                if (name === COOKIE_SESSION_REFRESH_TOKEN) return refreshToken;
                if (name === COOKIE_SESSION_USER_INFO) return idToken;
                return undefined;
            }
        };

        const result = await refreshSession(cookieStore);

        // openid-client library validates responses and rejects fake JWTs
        expect(result.success).toBe(false);
        expect(result.cookies[COOKIE_SESSION_ACCESS_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_REFRESH_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_USER_INFO]).toEqual({value: '', maxAge: 0});
    });

    it('should return failure and clear cookies when refresh token is invalid', async () => {
        const refreshToken = 'INVALID_REFRESH_TOKEN';
        const idToken = 'ID_TOKEN_VALUE';

        fakeOIDCServer.setupRefreshTokenError(refreshToken, 'invalid_grant', 'Refresh token is invalid or expired');

        const cookieStore: ReadCookieStore = {
            get: (name: string) => {
                if (name === COOKIE_SESSION_REFRESH_TOKEN) return refreshToken;
                if (name === COOKIE_SESSION_USER_INFO) return idToken;
                return undefined;
            }
        };

        const result = await refreshSession(cookieStore);

        expect(result.success).toBe(false);
        expect(result.cookies[COOKIE_SESSION_ACCESS_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_REFRESH_TOKEN]).toEqual({value: '', maxAge: 0});
        expect(result.cookies[COOKIE_SESSION_USER_INFO]).toEqual({value: '', maxAge: 0});
    });
});
