import {
    COOKIE_AUTH_CODE_VERIFIER,
    COOKIE_AUTH_NONCE,
    COOKIE_AUTH_REDIRECT_AFTER_LOGIN,
    COOKIE_AUTH_STATE,
    COOKIE_SESSION_ACCESS_TOKEN,
    COOKIE_SESSION_REFRESH_TOKEN,
    COOKIE_SESSION_USER_INFO
} from './constants';
import {COOKIE_VALUE_DELETE} from "@/libs/security/session-service";
import {ReadCookieStore, SetCookies} from "@/libs/nextjs-cookies";

export interface StoredSession {
    accessToken?: string,
    refreshToken?: string,
    idToken?: string,
}

export interface StoredAuthenticationFlow {
    codeVerifier: string;
    nonce: string;
    redirectAfterLogin: string;
    state: string;
}

export function clearFullSession(): SetCookies {
    return {
        [COOKIE_SESSION_ACCESS_TOKEN]: COOKIE_VALUE_DELETE,
        [COOKIE_SESSION_REFRESH_TOKEN]: COOKIE_VALUE_DELETE,
        [COOKIE_SESSION_USER_INFO]: COOKIE_VALUE_DELETE,

        [COOKIE_AUTH_STATE]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_CODE_VERIFIER]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_NONCE]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_REDIRECT_AFTER_LOGIN]: COOKIE_VALUE_DELETE,
    }
}

export  function clearAuthSession(): SetCookies {
    return {
        [COOKIE_AUTH_STATE]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_CODE_VERIFIER]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_NONCE]: COOKIE_VALUE_DELETE,
        [COOKIE_AUTH_REDIRECT_AFTER_LOGIN]: COOKIE_VALUE_DELETE,
    }
}

export function storeAuthSession(parameters: StoredAuthenticationFlow): SetCookies {
    return {
        [COOKIE_AUTH_NONCE]: {value: parameters.nonce, maxAge: 600},
        [COOKIE_AUTH_STATE]: {value: parameters.state, maxAge: 600},
        [COOKIE_AUTH_CODE_VERIFIER]: {value: parameters.codeVerifier, maxAge: 600},
        [COOKIE_AUTH_REDIRECT_AFTER_LOGIN]: {value: parameters.redirectAfterLogin, maxAge: 600},
    }
}

export function loadAuthSession(cookieStore: ReadCookieStore): StoredAuthenticationFlow {
    return {
        codeVerifier: cookieStore.get(COOKIE_AUTH_CODE_VERIFIER) || '',
        nonce: cookieStore.get(COOKIE_AUTH_NONCE) || '',
        redirectAfterLogin: cookieStore.get(COOKIE_AUTH_REDIRECT_AFTER_LOGIN) || '',
        state: cookieStore.get(COOKIE_AUTH_STATE) || '',
    }
}

export function storeSession(session: {
    accessToken: string;
    accessTokenExpiresIn: number;
    refreshToken: string;
    idToken: string;
}): SetCookies {
    return {
        [COOKIE_SESSION_ACCESS_TOKEN]: {value: session.accessToken, maxAge: session.accessTokenExpiresIn},
        [COOKIE_SESSION_REFRESH_TOKEN]: {value: session.refreshToken, maxAge: 30 * 24 * 3600},
        [COOKIE_SESSION_USER_INFO]: {value: session.idToken, maxAge: 30 * 24 * 3600},
    };
}

export function loadSession(cookieStore: ReadCookieStore): StoredSession {
    return {
        accessToken: cookieStore.get(COOKIE_SESSION_ACCESS_TOKEN),
        refreshToken: cookieStore.get(COOKIE_SESSION_REFRESH_TOKEN),
        idToken: cookieStore.get(COOKIE_SESSION_USER_INFO),
    };

}
