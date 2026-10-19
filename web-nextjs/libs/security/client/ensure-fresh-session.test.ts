// @vitest-environment jsdom

import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {COOKIE_SESSION_ACCESS_TOKEN} from "@/libs/security/constants";

interface RefreshFakeConfig {
    nextStatus: number;
    rotateCookieOnRefresh: boolean;
    throwOnRefresh: boolean;
    inflightGate: Promise<void> | null;
}

function makeRefreshFake(config: Partial<RefreshFakeConfig> = {}): typeof fetch {
    const settings: RefreshFakeConfig = {
        nextStatus: 204,
        rotateCookieOnRefresh: true,
        throwOnRefresh: false,
        inflightGate: null,
        ...config,
    };

    return (async (input: RequestInfo | URL): Promise<Response> => {
        const url = typeof input === "string" ? input : input.toString();
        if (!url.includes("/auth/refresh")) {
            throw new Error(`Refresh fake received an unexpected call: ${url}`);
        }
        if (settings.inflightGate) {
            await settings.inflightGate;
        }
        if (settings.throwOnRefresh) {
            throw new Error("network boom");
        }
        if (settings.rotateCookieOnRefresh && settings.nextStatus >= 200 && settings.nextStatus < 300) {
            setAccessTokenCookie(3600, "rotated");
        }
        return new Response(null, {status: settings.nextStatus});
    }) as typeof fetch;
}

function clearCookies(): void {
    for (const entry of document.cookie.split("; ")) {
        const name = entry.split("=")[0];
        if (name) {
            document.cookie = `${name}=; Max-Age=0; path=/`;
        }
    }
}

function base64url(input: string): string {
    return btoa(input).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function fakeJwt(claims: Record<string, unknown>): string {
    const header = base64url(JSON.stringify({alg: "none", typ: "JWT"}));
    const payload = base64url(JSON.stringify(claims));
    return `${header}.${payload}.fake-signature`;
}

function setAccessTokenCookie(secondsFromNow: number, label = "token"): string {
    const exp = Math.floor(Date.now() / 1000) + secondsFromNow;
    const token = fakeJwt({exp, label});
    document.cookie = `${COOKIE_SESSION_ACCESS_TOKEN}=${token}; path=/`;
    return token;
}

describe('ensureFreshSession', () => {
    beforeEach(() => {
        vi.resetModules();
        clearCookies();
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        clearCookies();
    });

    it('should resolve to undefined when the current token expires further than one minute in the future', async () => {
        setAccessTokenCookie(600);
        vi.stubGlobal('fetch', makeRefreshFake());
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await expect(ensureFreshSession()).resolves.toBeUndefined();
    });

    it('should resolve to undefined after a successful refresh when no access-token cookie is present', async () => {
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 204}));
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await expect(ensureFreshSession()).resolves.toBeUndefined();
    });

    it('should resolve to undefined after a successful refresh when the current token is close to expiry', async () => {
        setAccessTokenCookie(30);
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 204}));
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        await expect(ensureFreshSession()).resolves.toBeUndefined();
    });

    it('should throw SessionExpiredException when the refresh endpoint returns 401', async () => {
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 401, rotateCookieOnRefresh: false}));
        const {ensureFreshSession, SessionExpiredException} = await import('./index');

        await expect(ensureFreshSession()).rejects.toBeInstanceOf(SessionExpiredException);
    });

    it('should throw SessionExpiredException when the refresh request errors', async () => {
        vi.stubGlobal('fetch', makeRefreshFake({throwOnRefresh: true}));
        const {ensureFreshSession, SessionExpiredException} = await import('./index');

        await expect(ensureFreshSession()).rejects.toBeInstanceOf(SessionExpiredException);
    });

    it('should resolve all concurrent callers with a single refresh outcome', async () => {
        let gateResolve: () => void = () => undefined;
        const gate = new Promise<void>((resolve) => {
            gateResolve = resolve;
        });
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 204, inflightGate: gate}));
        const {ensureFreshSession} = await import('./ensure-fresh-session');

        const first = ensureFreshSession();
        const second = ensureFreshSession();
        gateResolve();

        await expect(Promise.all([first, second])).resolves.toEqual([undefined, undefined]);
    });

    it('should reject all concurrent callers with SessionExpiredException when the refresh fails', async () => {
        let gateResolve: () => void = () => undefined;
        const gate = new Promise<void>((resolve) => {
            gateResolve = resolve;
        });
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 401, rotateCookieOnRefresh: false, inflightGate: gate}));
        const {ensureFreshSession, SessionExpiredException} = await import('./index');

        const first = ensureFreshSession();
        const second = ensureFreshSession();
        gateResolve();

        await expect(first).rejects.toBeInstanceOf(SessionExpiredException);
        await expect(second).rejects.toBeInstanceOf(SessionExpiredException);
    });
});
