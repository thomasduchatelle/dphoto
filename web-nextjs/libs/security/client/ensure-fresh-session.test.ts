// @vitest-environment jsdom

import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT} from "@/libs/security/constants";

function clearCookies(): void {
    for (const entry of document.cookie.split("; ")) {
        const name = entry.split("=")[0];
        if (name) {
            document.cookie = `${name}=; Max-Age=0; path=/`;
        }
    }
}

function setExpiresAtCookie(secondsFromNow: number): void {
    const expiresAt = Math.floor(Date.now() / 1000) + secondsFromNow;
    document.cookie = `${COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT}=${expiresAt}; path=/`;
}

describe('ensureFreshSession', () => {
    let fetchMock: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        vi.resetModules();
        clearCookies();
        fetchMock = vi.fn().mockResolvedValue(new Response(null, {status: 204}));
        vi.stubGlobal('fetch', fetchMock);
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        clearCookies();
    });

    it('should trigger a refresh when the expiry cookie is missing', async () => {
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await ensureFreshSession();

        expect(fetchMock).toHaveBeenCalledTimes(1);
        expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining('/auth/refresh'), expect.objectContaining({method: 'POST'}));
    });

    it('should not trigger a refresh when the token has more than threshold seconds left', async () => {
        setExpiresAtCookie(600);
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await ensureFreshSession();

        expect(fetchMock).not.toHaveBeenCalled();
    });

    it('should trigger a refresh when the token is within the freshness threshold', async () => {
        setExpiresAtCookie(30);
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await ensureFreshSession();

        expect(fetchMock).toHaveBeenCalledTimes(1);
    });

    it('should trigger a single refresh for two concurrent calls', async () => {
        setExpiresAtCookie(10);
        let resolveFetch: (value: Response) => void = () => undefined;
        fetchMock.mockImplementation(() => new Promise<Response>((resolve) => {
            resolveFetch = resolve;
        }));
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        const first = ensureFreshSession();
        const second = ensureFreshSession();

        expect(fetchMock).toHaveBeenCalledTimes(1);
        resolveFetch(new Response(null, {status: 204}));
        await Promise.all([first, second]);
    });

    it('should always refresh when thresholdSeconds is Infinity, even with a fresh cookie', async () => {
        setExpiresAtCookie(3600);
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await ensureFreshSession({thresholdSeconds: Number.POSITIVE_INFINITY});

        expect(fetchMock).toHaveBeenCalledTimes(1);
    });

    it('should resolve silently when the refresh request fails', async () => {
        fetchMock.mockRejectedValue(new Error('network boom'));
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await expect(ensureFreshSession()).resolves.toBeUndefined();
    });
});
