import {ensureFreshSession, readAccessTokenFromCookie} from "./ensure-fresh-session";
import {SessionExpiredException} from "./session-expired-exception";

export async function getAccessToken(): Promise<string> {
    await ensureFreshSession();
    const token = readAccessTokenFromCookie();
    if (!token) {
        throw new SessionExpiredException("Access token cookie is not available after session refresh");
    }
    return token;
}
