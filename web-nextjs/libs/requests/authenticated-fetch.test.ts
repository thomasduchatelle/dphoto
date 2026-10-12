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

describe('authenticatedFetch', () => {
    let fetchMock: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        vi.resetModules();
        clearCookies();
        fetchMock = vi.fn();
        vi.stubGlobal('fetch', fetchMock);
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        clearCookies();
    });

    it('should pass through a successful response without any refresh when the token is fresh', async () => {
        setExpiresAtCookie(600);
        fetchMock.mockResolvedValue(new Response('{"ok":true}', {status: 200}));
        const {authenticatedFetch} = await import('./authenticated-fetch');

        const response = await authenticatedFetch('/api/v1/albums');

        expect(response.status).toBe(200);
        expect(fetchMock).toHaveBeenCalledTimes(1);
        expect(fetchMock).toHaveBeenCalledWith('/api/v1/albums', expect.objectContaining({credentials: 'include'}));
    });

    it('should proactively refresh the session when the token is close to expiry, then make the request', async () => {
        setExpiresAtCookie(30);
        fetchMock
            .mockResolvedValueOnce(new Response(null, {status: 204}))
            .mockResolvedValueOnce(new Response('{"ok":true}', {status: 200}));
        const {authenticatedFetch} = await import('./authenticated-fetch');

        const response = await authenticatedFetch('/api/v1/albums');

        expect(response.status).toBe(200);
        expect(fetchMock).toHaveBeenCalledTimes(2);
        expect(fetchMock.mock.calls[0][0]).toContain('/auth/refresh');
        expect(fetchMock.mock.calls[1][0]).toBe('/api/v1/albums');
    });

    it('should refresh and retry once on a 401 response, returning the retried response', async () => {
        setExpiresAtCookie(600);
        fetchMock
            .mockResolvedValueOnce(new Response(null, {status: 401}))
            .mockResolvedValueOnce(new Response(null, {status: 204}))
            .mockResolvedValueOnce(new Response('{"ok":true}', {status: 200}));
        const {authenticatedFetch} = await import('./authenticated-fetch');

        const response = await authenticatedFetch('/api/v1/albums');

        expect(response.status).toBe(200);
        expect(fetchMock).toHaveBeenCalledTimes(3);
        expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/albums');
        expect(fetchMock.mock.calls[1][0]).toContain('/auth/refresh');
        expect(fetchMock.mock.calls[2][0]).toBe('/api/v1/albums');
    });

    it('should return the retried 401 unchanged when the forced refresh also fails to recover', async () => {
        setExpiresAtCookie(600);
        fetchMock
            .mockResolvedValueOnce(new Response(null, {status: 401}))
            .mockResolvedValueOnce(new Response(null, {status: 401}))
            .mockResolvedValueOnce(new Response(null, {status: 401}));
        const {authenticatedFetch} = await import('./authenticated-fetch');

        const response = await authenticatedFetch('/api/v1/albums');

        expect(response.status).toBe(401);
        expect(fetchMock).toHaveBeenCalledTimes(3);
    });

    it('should not retry on non-401 errors', async () => {
        setExpiresAtCookie(600);
        fetchMock.mockResolvedValueOnce(new Response(null, {status: 500}));
        const {authenticatedFetch} = await import('./authenticated-fetch');

        const response = await authenticatedFetch('/api/v1/albums');

        expect(response.status).toBe(500);
        expect(fetchMock).toHaveBeenCalledTimes(1);
    });
});
