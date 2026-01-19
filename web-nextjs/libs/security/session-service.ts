import {clearAuthSession, clearFullSession, loadAuthSession, loadSession, storeAuthSession, storeSession} from "@/libs/security/backend-store";
import {getLogoutUrl} from "@/libs/security/logout-utils";
import {getOidcConfigFromEnv, oidcConfig} from "@/libs/security/oidc-config";
import * as client from "openid-client";
import {redirectUrl, requestUrlWithBaseBath} from "@/libs/requests";
import {decodeJWTPayload} from "@/libs/security/jwt-utils";
import {parseCurrentAccessToken} from "@/libs/security/access-token-service";
import {CookieValue, ReadCookieStore, Redirection} from "@/libs/nextjs-cookies";

export interface AuthenticatedUser {
    name: string;
    email: string;
    picture?: string;
    isOwner: boolean;
}

export type BackendSession = AuthenticatedSession | AnonymousSession

export interface AuthenticatedSession {
    status: 'authenticated';
    authenticatedUser: AuthenticatedUser;
    logoutUrl: string;
    aboutToExpire: boolean;
}

export interface AnonymousSession {
    status: 'anonymous';
}

interface IDTokenPayload {
    given_name?: string;
    family_name?: string;
    email?: string;
    picture?: string;
    exp?: number;

    [key: string]: any;
}

function readIdToken(idToken: string): {
    name: string;
    email: string;
    picture?: string;
} {

    const idTokenPayload = decodeJWTPayload(idToken) as IDTokenPayload | null;
    if (idTokenPayload) {
        const firstName = idTokenPayload.given_name || '';
        const lastName = idTokenPayload.family_name || '';
        const fullName = [firstName, lastName].filter(Boolean).join(' ');

        return {
            name: fullName,
            email: idTokenPayload.email || '',
            picture: idTokenPayload.picture,
        };
    }

    return {
        name: 'Monsieur, Madame',
        email: '',
    };
}

export async function getCurrentAuthentication(cookieStore: ReadCookieStore): Promise<BackendSession> {
    const session = loadSession(cookieStore);
    if (!session.accessToken || !session.accessToken || !session.idToken) {
        return {status: "anonymous"}
    }

    const accessToken = await parseCurrentAccessToken(session.accessToken);
    if (!accessToken) {
        return {status: "anonymous"}
    }

    return {
        status: 'authenticated',
        authenticatedUser: {
            ...readIdToken(session.idToken),
            isOwner: accessToken.isOwner,
        },
        logoutUrl: await getLogoutUrl(),
        aboutToExpire: accessToken.aboutToExpire,
    };
}


/**
 * To logout, the user must be first redirected to the identity provider's logout endpoint. Then this function can be called.
 */
export async function completeLogout() {
    await clearFullSession()
}

export const COOKIE_VALUE_DELETE: CookieValue = {value: '', maxAge: 0};

export async function initiateAuthenticationFlow(path: string = "/"): Promise<Redirection> {
    const config = await oidcConfig(getOidcConfigFromEnv());

    const codeVerifier: string = client.randomPKCECodeVerifier();
    const code_challenge: string = await client.calculatePKCECodeChallenge(codeVerifier);

    let originalUrl = (await redirectUrl("/auth/callback")).toString();
    const parameters: Record<string, string> = {
        redirect_uri: originalUrl,
        scope: 'openid profile email',
        code_challenge,
        code_challenge_method: 'S256',
        state: client.randomState(),
        nonce: client.randomNonce(),
    };

    const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);
    return {
        redirectTo,
        cookies: storeAuthSession({
            nonce: parameters.nonce,
            state: parameters.state,
            codeVerifier: codeVerifier,
            redirectAfterLogin: path,
        })
    };
}

async function redirectToErrorPage(error: string, errorDescription?: string): Promise<Redirection> {
    const errorUrl = await redirectUrl("/auth/error")
    errorUrl.searchParams.set('error', error);
    if (errorDescription) {
        errorUrl.searchParams.set('error_description', errorDescription);
    }

    return {
        redirectTo: errorUrl,
        cookies: clearFullSession(),
    };
}

export async function authenticate(requestUrl: URL, cookiesStore: ReadCookieStore): Promise<Redirection> {
    const searchParams = requestUrl.searchParams;
    const authenticationFlowState = loadAuthSession(cookiesStore);

    const errorParam = searchParams.get('error');
    if (errorParam) {
        const errorDescription = searchParams.get('error_description');
        return redirectToErrorPage(errorParam, errorDescription ?? undefined);
    }

    if (!authenticationFlowState.state || !authenticationFlowState.codeVerifier || !authenticationFlowState.nonce) {
        console.log("Invalid authenticationFlowState:", authenticationFlowState);
        return redirectToErrorPage('missing-authentication-cookies');
    }

    const config = await oidcConfig(getOidcConfigFromEnv());

    try {
        const tokens: client.TokenEndpointResponse = await client.authorizationCodeGrant(
            config,
            await requestUrlWithBaseBath(requestUrl), // This is the URL of this route
            {
                pkceCodeVerifier: authenticationFlowState.codeVerifier,
                expectedState: authenticationFlowState.state,
                expectedNonce: authenticationFlowState.nonce,
            }
        );

        return {
            redirectTo: (await redirectUrl(authenticationFlowState.redirectAfterLogin ?? "/")),
            cookies: {
                ...clearAuthSession(),
                ...storeSession({
                    accessToken: tokens.access_token,
                    accessTokenExpiresIn: tokens.expires_in ?? 3600,
                    refreshToken: tokens.refresh_token ?? '',
                    idToken: tokens.id_token ?? '',
                }),
            },
        }

    } catch (error) {
        console.error('OAuth callback error:', error);
        return redirectToErrorPage('token-exchange-failed');
    }
}