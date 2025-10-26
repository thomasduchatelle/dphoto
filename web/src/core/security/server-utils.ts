import {getContextData} from "waku/middleware/context";
import {BackendSession, newAnonymousSession} from "./security-model";

export function readBackendSession() {
    const {backendSession} = getContextData() as { backendSession: BackendSession };
    return backendSession ?? newAnonymousSession();
}