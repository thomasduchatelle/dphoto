// @vitest-environment jsdom

import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {COOKIE_SESSION_ACCESS_TOKEN} from "@/libs/security/constants";

interface RefreshFake {
    fetch: typeof fetch;
    getRotatedToken(): string | null;
}

interface RefreshFakeConfig {
    nextStatus: number;
    rotateCookie: boolean;
}

function makeRefreshFake(config: Partial<RefreshFakeConfig> = {}): RefreshFake {
    const settings: RefreshFakeConfig = {
        nextStatus: 204,
        rotateCookie: true,
        ...config,
    };
    let rotatedToken: string | null = null;

    return {
        fetch: (async (input: RequestInfo | URL): Promise<Response> => {
            const url = typeof input === "string" ? input : input.toString();
            if (!url.includes("/auth/refresh")) {
                throw new Error(`Refresh fake received an unexpected call: ${url}`);
            }
            if (settings.rotateCookie && settings.nextStatus >= 200 && settings.nextStatus < 300) {
                rotatedToken = setAccessTokenCookie(3600, "rotated");
            }
            return new Response(null, {status: settings.nextStatus});
        }) as typeof fetch,
        getRotatedToken: () => rotatedToken,
    };
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

describe('getAccessToken', () => {
    beforeEach(() => {
        vi.resetModules();
        clearCookies();
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        clearCookies();
    });

    it('should return the current cookie value when the token is still fresh', async () => {
        const token = setAccessTokenCookie(600, "fresh");
        vi.stubGlobal('fetch', makeRefreshFake().fetch);
        const {getAccessToken} = await import('./get-access-token');

        await expect(getAccessToken()).resolves.toBe(token);
    });

    it('should return the rotated token after a successful refresh when no cookie was present', async () => {
        const fake = makeRefreshFake();
        vi.stubGlobal('fetch', fake.fetch);
        const {getAccessToken} = await import('./get-access-token');

        const result = await getAccessToken();

        expect(result).toBe(fake.getRotatedToken());
    });

    it('should return the rotated token after a successful refresh when the cookie was close to expiry', async () => {
        setAccessTokenCookie(10, "stale");
        const fake = makeRefreshFake();
        vi.stubGlobal('fetch', fake.fetch);
        const {getAccessToken} = await import('./get-access-token');

        const result = await getAccessToken();

        expect(result).toBe(fake.getRotatedToken());
    });

    it('should throw SessionExpiredException when the refresh endpoint rejects the request', async () => {
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 401, rotateCookie: false}).fetch);
        const {getAccessToken, SessionExpiredException} = await import('./index');

        await expect(getAccessToken()).rejects.toBeInstanceOf(SessionExpiredException);
    });

    it('should throw SessionExpiredException when the refresh succeeds but no cookie ends up written', async () => {
        vi.stubGlobal('fetch', makeRefreshFake({nextStatus: 204, rotateCookie: false}).fetch);
        const {getAccessToken, SessionExpiredException} = await import('./index');

        await expect(getAccessToken()).rejects.toBeInstanceOf(SessionExpiredException);
    });
});
