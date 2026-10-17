import {ensureFreshSession, readAccessTokenFromCookie} from "./ensure-fresh-session";

export async function getAccessToken(): Promise<string | undefined> {
    await ensureFreshSession();
    return readAccessTokenFromCookie() ?? undefined;
}

export async function authenticatedFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    const response = await fetch(input, {credentials: "include", ...init});

    if (response.status !== 401) {
        return response;
    }

    await ensureFreshSession({thresholdSeconds: Number.POSITIVE_INFINITY});
    const refreshedToken = readAccessTokenFromCookie();
    const retryHeaders = new Headers(init?.headers);
    if (refreshedToken) {
        retryHeaders.set("Authorization", `Bearer ${refreshedToken}`);
    }
    return fetch(input, {credentials: "include", ...init, headers: retryHeaders});
}
