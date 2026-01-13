import {cookies} from 'next/headers';
import {
    COOKIE_AUTH_CODE_VERIFIER,
    COOKIE_AUTH_NONCE,
    COOKIE_AUTH_REDIRECT_AFTER_LOGIN,
    COOKIE_AUTH_STATE,
    COOKIE_SESSION_ACCESS_TOKEN,
    COOKIE_SESSION_REFRESH_TOKEN,
    COOKIE_SESSION_USER_INFO
} from './constants';

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

const deleteCookieOpt = {maxAge: 0, path: '/'};

const baseCookieOptions: any = {
    httpOnly: true,
    secure: true,
    sameSite: 'lax', // lax is required when the Referer is a different site (which happens during OAuth flow when user is not already authenticated on Cognito: user comes from the Social login)
    path: '/',
};

const authenticationFlowCookieOptions = {
    ...baseCookieOptions,
    maxAge: 10 * 60, // 10 minutes
}

const tokensCookieOptions = {
    ...baseCookieOptions,
    maxAge: 30 * 24 * 3600, // 30 days in seconds
}

export async function clearAuthSession(): Promise<void> {
    const cookieStore = await cookies()

    cookieStore.set(COOKIE_SESSION_ACCESS_TOKEN, '', deleteCookieOpt);
    cookieStore.set(COOKIE_SESSION_REFRESH_TOKEN, '', deleteCookieOpt);
}

export async function clearFullSession(): Promise<void> {
    const cookieStore = await cookies()
    cookieStore.set(COOKIE_SESSION_ACCESS_TOKEN, '', deleteCookieOpt);
    cookieStore.set(COOKIE_SESSION_REFRESH_TOKEN, '', deleteCookieOpt);

    cookieStore.set(COOKIE_AUTH_STATE, '', deleteCookieOpt);
    cookieStore.set(COOKIE_AUTH_CODE_VERIFIER, '', deleteCookieOpt);
    cookieStore.set(COOKIE_AUTH_NONCE, '', deleteCookieOpt);
    cookieStore.set(COOKIE_AUTH_REDIRECT_AFTER_LOGIN, '', deleteCookieOpt);
}

export async function storeAuthSession(authSession: StoredAuthenticationFlow): Promise<void> {
    const cookieStore = await cookies();

    cookieStore.set(COOKIE_AUTH_CODE_VERIFIER, authSession.codeVerifier, authenticationFlowCookieOptions);
    cookieStore.set(COOKIE_AUTH_NONCE, authSession.nonce, authenticationFlowCookieOptions);
    cookieStore.set(COOKIE_AUTH_REDIRECT_AFTER_LOGIN, authSession.redirectAfterLogin, authenticationFlowCookieOptions);
    cookieStore.set(COOKIE_AUTH_STATE, authSession.state, authenticationFlowCookieOptions);
}

export async function loadAuthSession(): Promise<StoredAuthenticationFlow> {
    const cookieStore = await cookies();
    return {
        codeVerifier: cookieStore.get(COOKIE_AUTH_CODE_VERIFIER)?.value || '',
        nonce: cookieStore.get(COOKIE_AUTH_NONCE)?.value || '',
        redirectAfterLogin: cookieStore.get(COOKIE_AUTH_REDIRECT_AFTER_LOGIN)?.value || '',
        state: cookieStore.get(COOKIE_AUTH_STATE)?.value || '',
    }
}

export async function storeSession(session: {
    accessToken: string;
    accessTokenExpiresIn: number;
    refreshToken: string;
    idToken: string;
}) {
    const cookieStore = await cookies();

    cookieStore.set(COOKIE_SESSION_ACCESS_TOKEN, session.accessToken, {...tokensCookieOptions, maxAge: session.accessTokenExpiresIn});
    cookieStore.set(COOKIE_SESSION_REFRESH_TOKEN, session.refreshToken, tokensCookieOptions);
    cookieStore.set(COOKIE_SESSION_USER_INFO, session.idToken, tokensCookieOptions);
}

export async function loadSession(): Promise<StoredSession> {
    const cookieStore = await cookies();
    return {
        accessToken: cookieStore.get(COOKIE_SESSION_ACCESS_TOKEN)?.value,
        refreshToken: cookieStore.get(COOKIE_SESSION_REFRESH_TOKEN)?.value,
        idToken: cookieStore.get(COOKIE_SESSION_USER_INFO)?.value,
    };

}