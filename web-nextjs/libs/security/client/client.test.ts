// @vitest-environment jsdom

import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {COOKIE_SESSION_ACCESS_TOKEN} from "@/libs/security/constants";

interface FakeApiCall {
    url: string;
    method: string;
    authorization: string | null;
}

class FakeApi {
    public readonly calls: FakeApiCall[] = [];
    private readonly responses: Response[] = [];
    private pendingRefreshes: Array<() => void> = [];
    private holdingRefresh = false;
    private onRefreshCallback: (() => void) | null = null;

    public queue(...responses: Response[]): void {
        this.responses.push(...responses);
    }

    public holdRefreshes(): void {
        this.holdingRefresh = true;
    }

    public releaseRefreshes(): void {
        this.holdingRefresh = false;
        const toRelease = this.pendingRefreshes;
        this.pendingRefreshes = [];
        for (const resolve of toRelease) {
            resolve();
        }
    }

    public onRefresh(callback: () => void): void {
        this.onRefreshCallback = callback;
    }

    public asFetch(): typeof fetch {
        return ((input: RequestInfo | URL, init?: RequestInit) => {
            const url = typeof input === "string" ? input : input.toString();
            const method = init?.method ?? "GET";
            const headers = new Headers(init?.headers);
            this.calls.push({url, method, authorization: headers.get("Authorization")});

            if (url.includes("/auth/refresh")) {
                this.onRefreshCallback?.();
                if (this.holdingRefresh) {
                    return new Promise<Response>((resolve) => {
                        this.pendingRefreshes.push(() => {
                            const response = this.responses.shift() ?? new Response(null, {status: 204});
                            resolve(response);
                        });
                    });
                }
            }

            const response = this.responses.shift();
            if (!response) {
                throw new Error(`Unexpected call to ${method} ${url}: no response queued`);
            }
            return Promise.resolve(response);
        }) as typeof fetch;
    }
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

describe('security/client slice (ensureFreshSession + getAccessToken + authenticatedFetch end to end)', () => {
    let api: FakeApi;

    beforeEach(() => {
        vi.resetModules();
        clearCookies();
        api = new FakeApi();
        vi.stubGlobal('fetch', api.asFetch());
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        clearCookies();
    });

    describe('ensureFreshSession', () => {
        it('should POST /auth/refresh when no access-token cookie is present', async () => {
            api.queue(new Response(null, {status: 204}));
            const {ensureFreshSession} = await import('./ensure-fresh-session');

            await ensureFreshSession();

            expect(api.calls).toEqual([{url: expect.stringContaining("/auth/refresh"), method: "POST", authorization: null}]);
        });

        it('should skip the refresh when the token expires further than the threshold', async () => {
            setAccessTokenCookie(600);
            const {ensureFreshSession} = await import('./ensure-fresh-session');

            await ensureFreshSession();

            expect(api.calls).toEqual([]);
        });

        it('should refresh when the token is within the freshness threshold', async () => {
            setAccessTokenCookie(30);
            api.queue(new Response(null, {status: 204}));
            const {ensureFreshSession} = await import('./ensure-fresh-session');

            await ensureFreshSession();

            expect(api.calls).toHaveLength(1);
            expect(api.calls[0].url).toContain("/auth/refresh");
        });

        it('should trigger a single refresh when called twice concurrently', async () => {
            setAccessTokenCookie(10);
            api.queue(new Response(null, {status: 204}));
            api.holdRefreshes();
            const {ensureFreshSession} = await import('./ensure-fresh-session');

            const first = ensureFreshSession();
            const second = ensureFreshSession();

            expect(api.calls).toHaveLength(1);
            api.releaseRefreshes();
            await Promise.all([first, second]);
            expect(api.calls).toHaveLength(1);
        });

        it('should resolve silently when the refresh request throws', async () => {
            vi.stubGlobal('fetch', (() => Promise.reject(new Error('network boom'))) as typeof fetch);
            const {ensureFreshSession} = await import('./ensure-fresh-session');

            await expect(ensureFreshSession()).resolves.toBeUndefined();
        });
    });

    describe('getAccessToken', () => {
        it('should return the current cookie value without hitting the network when the token is fresh', async () => {
            const token = setAccessTokenCookie(600);
            const {getAccessToken} = await import('./authenticated-fetch');

            const result = await getAccessToken();

            expect(result).toBe(token);
            expect(api.calls).toEqual([]);
        });

        it('should refresh the session, then return the rotated token from the cookie', async () => {
            setAccessTokenCookie(10, "stale");
            api.queue(new Response(null, {status: 204}));
            let rotatedToken = "";
            api.onRefresh(() => {
                rotatedToken = setAccessTokenCookie(3600, "rotated");
            });
            const {getAccessToken} = await import('./authenticated-fetch');

            const result = await getAccessToken();

            expect(result).toBe(rotatedToken);
            expect(api.calls).toHaveLength(1);
            expect(api.calls[0].url).toContain("/auth/refresh");
        });

        it('should return undefined when no cookie is present and the refresh does not set one', async () => {
            api.queue(new Response(null, {status: 401}));
            const {getAccessToken} = await import('./authenticated-fetch');

            const result = await getAccessToken();

            expect(result).toBeUndefined();
        });
    });

    describe('authenticatedFetch', () => {
        it('should pass a 200 response through unchanged, sending credentials but not touching Authorization', async () => {
            setAccessTokenCookie(600, "fresh");
            api.queue(new Response('{"ok":true}', {status: 200}));
            const {authenticatedFetch} = await import('./authenticated-fetch');

            const response = await authenticatedFetch('/api/v1/albums', {headers: {Authorization: 'Bearer fresh'}});

            expect(response.status).toBe(200);
            expect(api.calls).toEqual([{url: '/api/v1/albums', method: 'GET', authorization: 'Bearer fresh'}]);
        });

        it('should force a refresh and retry once with the rotated bearer token on 401', async () => {
            const staleToken = setAccessTokenCookie(600, "stale");
            api.queue(
                new Response(null, {status: 401}),
                new Response(null, {status: 204}),
                new Response('{"ok":true}', {status: 200}),
            );
            let rotatedToken = "";
            api.onRefresh(() => {
                rotatedToken = setAccessTokenCookie(3600, "rotated");
            });
            const {authenticatedFetch} = await import('./authenticated-fetch');

            const response = await authenticatedFetch('/api/v1/albums', {headers: {Authorization: `Bearer ${staleToken}`}});

            expect(response.status).toBe(200);
            expect(api.calls).toHaveLength(3);
            expect(api.calls[0]).toMatchObject({url: '/api/v1/albums', authorization: `Bearer ${staleToken}`});
            expect(api.calls[1].url).toContain('/auth/refresh');
            expect(api.calls[2]).toMatchObject({url: '/api/v1/albums', authorization: `Bearer ${rotatedToken}`});
        });

        it('should return the retried 401 unchanged when the forced refresh also fails to recover', async () => {
            setAccessTokenCookie(600);
            api.queue(
                new Response(null, {status: 401}),
                new Response(null, {status: 401}),
                new Response(null, {status: 401}),
            );
            const {authenticatedFetch} = await import('./authenticated-fetch');

            const response = await authenticatedFetch('/api/v1/albums');

            expect(response.status).toBe(401);
            expect(api.calls).toHaveLength(3);
        });

        it('should not retry on non-401 errors', async () => {
            setAccessTokenCookie(600);
            api.queue(new Response(null, {status: 500}));
            const {authenticatedFetch} = await import('./authenticated-fetch');

            const response = await authenticatedFetch('/api/v1/albums');

            expect(response.status).toBe(500);
            expect(api.calls).toHaveLength(1);
        });
    });
});
