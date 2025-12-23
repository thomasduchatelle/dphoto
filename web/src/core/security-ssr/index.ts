import {getContextData} from "waku/middleware/context";
import {BackendSession, newAnonymousSession} from "../security";

export function readBackendSession() {
    const {backendSession} = getContextData() as { backendSession: BackendSession };
    return backendSession ?? newAnonymousSession();
}