import {ensureFreshSession} from "@/libs/security/client/ensure-fresh-session";

export async function authenticatedFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    await ensureFreshSession();

    const requestInit: RequestInit = {credentials: "include", ...init};
    const response = await fetch(input, requestInit);

    if (response.status !== 401) {
        return response;
    }

    await ensureFreshSession({thresholdSeconds: Number.POSITIVE_INFINITY});
    return fetch(input, requestInit);
}
