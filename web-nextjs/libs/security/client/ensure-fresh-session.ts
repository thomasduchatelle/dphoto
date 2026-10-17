import {COOKIE_SESSION_ACCESS_TOKEN} from "@/libs/security/constants";
import {withBasePath} from "@/libs/requests/media-url";

const DEFAULT_FRESHNESS_THRESHOLD_SECONDS = 60;

export interface EnsureFreshSessionOptions {
    thresholdSeconds?: number;
}

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

function isFreshEnough(token: string | null, thresholdSeconds: number): boolean {
    if (!token) {
        return false;
    }
    const expiresAtSeconds = decodeJWTExpirySeconds(token);
    if (expiresAtSeconds === null) {
        return false;
    }
    const nowSeconds = Math.floor(Date.now() / 1000);
    return expiresAtSeconds - nowSeconds > thresholdSeconds;
}

async function performRefresh(): Promise<void> {
    try {
        await fetch(withBasePath("/auth/refresh"), {
            method: "POST",
            credentials: "include",
        });
    } catch (err) {
        console.warn("ensureFreshSession: refresh request failed", err);
    }
}

export function ensureFreshSession(opts?: EnsureFreshSessionOptions): Promise<void> {
    const thresholdSeconds = opts?.thresholdSeconds ?? DEFAULT_FRESHNESS_THRESHOLD_SECONDS;

    if (isFreshEnough(readAccessTokenFromCookie(), thresholdSeconds)) {
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
