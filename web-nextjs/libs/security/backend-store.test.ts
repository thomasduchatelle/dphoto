// @vitest-environment node

import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {clearFullSession, storeSession} from './backend-store';
import {
    COOKIE_SESSION_ACCESS_TOKEN,
    COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT,
    COOKIE_SESSION_REFRESH_TOKEN,
    COOKIE_SESSION_USER_INFO,
} from './constants';

describe('storeSession', () => {
    beforeEach(() => {
        vi.useFakeTimers();
        vi.setSystemTime(new Date('2026-01-01T00:00:00Z'));
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('should write the three httpOnly session cookies and the non-httpOnly expiry hint', () => {
        const cookies = storeSession({
            accessToken: 'ACCESS',
            accessTokenExpiresIn: 3600,
            refreshToken: 'REFRESH',
            idToken: 'ID',
        });

        const nowSeconds = Math.floor(new Date('2026-01-01T00:00:00Z').getTime() / 1000);
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN]).toEqual({value: 'ACCESS', maxAge: 3600});
        expect(cookies[COOKIE_SESSION_REFRESH_TOKEN]).toEqual({value: 'REFRESH', maxAge: 30 * 24 * 3600});
        expect(cookies[COOKIE_SESSION_USER_INFO]).toEqual({value: 'ID', maxAge: 30 * 24 * 3600});
        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT]).toEqual({
            value: (nowSeconds + 3600).toString(),
            maxAge: 3600,
            httpOnly: false,
        });
    });
});

describe('clearFullSession', () => {
    it('should clear the expiry hint cookie as non-httpOnly to match how it was written', () => {
        const cookies = clearFullSession();

        expect(cookies[COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT]).toEqual({value: '', maxAge: 0, httpOnly: false});
    });
});
