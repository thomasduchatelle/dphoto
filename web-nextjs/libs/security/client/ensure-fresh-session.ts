import {COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT} from "@/libs/security/constants";
import {withBasePath} from "@/libs/requests/media-url";

const DEFAULT_FRESHNESS_THRESHOLD_SECONDS = 60;

export interface EnsureFreshSessionOptions {
    thresholdSeconds?: number;
}

let inflightRefresh: Promise<void> | null = null;

function readExpiresAtSeconds(): number | null {
    if (typeof document === "undefined") {
        return null;
    }
    const prefix = `${COOKIE_SESSION_ACCESS_TOKEN_EXPIRES_AT}=`;
    const cookies = document.cookie ? document.cookie.split("; ") : [];
    for (const entry of cookies) {
        if (entry.startsWith(prefix)) {
            const raw = entry.substring(prefix.length);
            const parsed = parseInt(raw, 10);
            return Number.isFinite(parsed) ? parsed : null;
        }
    }
    return null;
}

function isFreshEnough(expiresAtSeconds: number | null, thresholdSeconds: number): boolean {
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

    if (isFreshEnough(readExpiresAtSeconds(), thresholdSeconds)) {
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
