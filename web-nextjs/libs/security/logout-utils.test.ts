// @vitest-environment node

import {describe, expect, it, vi, beforeEach} from 'vitest';
import {clearAuthCookies} from '@/libs/security';

vi.mock('next/headers', () => {
    const mockCookies = {
        set: vi.fn(),
    };

    return {
        cookies: vi.fn(() => Promise.resolve(mockCookies)),
    };
});

describe('clearAuthCookies', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should clear all authentication cookies', async () => {
        const {cookies} = await import('next/headers');
        const mockCookieStore = await cookies();

        await clearAuthCookies();

        const cookieOptions = {
            maxAge: 0,
            path: '/',
        };

        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-access-token', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-refresh-token', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-state', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-code-verifier', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-oauth-nonce', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-redirect-after-login', '', cookieOptions);
        expect(mockCookieStore.set).toHaveBeenCalledWith('dphoto-user-info', '', cookieOptions);
    });
});
