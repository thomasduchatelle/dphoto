import {COOKIE_SESSION_ACCESS_TOKEN} from "@/libs/security/constants";
import {SessionExpiredException} from "./session-expired-exception";

const DEFAULT_FRESHNESS_THRESHOLD_SECONDS = 60;

let inflightRefresh: Promise<void> | null = null;

export function readAccessTokenFromCookie(): string | null {
    if (typeof document === "undefined") {
        return null;
    }
    const prefix = `${COOKIE_SESSION_ACCESS_TOKEN}=`;
    const cookies = document.cookie ? document.cookie.split("; ") : [];
    for (const entry of cookies) {
        if (entry.startsWith(prefix)) {
            const raw = entry.substring(prefix.length);
            return raw.length > 0 ? decodeURIComponent(raw) : null;
        }
    }
    return null;
}

function decodeJWTExpirySeconds(token: string): number | null {
    const parts = token.split(".");
    if (parts.length !== 3) {
        return null;
    }
    try {
        const padded = parts[1].replace(/-/g, "+").replace(/_/g, "/");
        const payload = JSON.parse(atob(padded));
        return typeof payload.exp === "number" ? payload.exp : null;
    } catch {
        return null;
    }
}

function isFreshEnough(token: string | null): boolean {
    if (!token) {
        return false;
    }
    const expiresAtSeconds = decodeJWTExpirySeconds(token);
    if (expiresAtSeconds === null) {
        return false;
    }
    const nowSeconds = Math.floor(Date.now() / 1000);
    return expiresAtSeconds - nowSeconds > DEFAULT_FRESHNESS_THRESHOLD_SECONDS;
}

async function performRefresh(): Promise<void> {
    let response: Response;
    try {
        response = await fetch("/auth/refresh", {
            method: "POST",
            credentials: "include",
        });
    } catch (err) {
        throw new SessionExpiredException(`Session refresh request failed: ${err instanceof Error ? err.message : String(err)}`);
    }

    if (!response.ok) {
        throw new SessionExpiredException(`Session refresh rejected with status ${response.status}`);
    }
}

export function ensureFreshSession(): Promise<void> {
    if (isFreshEnough(readAccessTokenFromCookie())) {
        return Promise.resolve();
    }

    if (inflightRefresh) {
        return inflightRefresh;
    }

    inflightRefresh = performRefresh().finally(() => {
        inflightRefresh = null;
    });

    return inflightRefresh;
}
