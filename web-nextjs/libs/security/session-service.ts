import {getValidAccessToken} from "@/libs/security/access-token-service";
import {clearAuthSession, clearFullSession, loadAuthSession, storeAuthSession, storeSession} from "@/libs/security/backend-store";
import {getLogoutUrl} from "@/libs/security/logout-utils";
import {getOidcConfigFromEnv, oidcConfig} from "@/libs/security/oidc-config";
import * as client from "openid-client";
import {redirectUrl, requestUrlWithBaseBath} from "@/libs/requests";
import {decodeJWTPayload} from "@/libs/security/jwt-utils";

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

export async function getValidAuthentication(): Promise<BackendSession> {
    const validSession = await getValidAccessToken();
    if (!validSession || !validSession.accessToken) {
        return {status: "anonymous"}
    }

    const {accessToken, idToken} = validSession;

    return {
        status: 'authenticated',
        authenticatedUser: {
            ...readIdToken(idToken),
            isOwner: accessToken.isOwner,
        },
        logoutUrl: await getLogoutUrl(),
    };
}


/**
 * To logout, the user must be first redirected to the identity provider's logout endpoint. Then this function can be called.
 */
export async function completeLogout() {
    await clearFullSession()
}

export interface Redirection {
    redirectTo: URL;
}

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

    await storeAuthSession({
        nonce: parameters.nonce,
        state: parameters.state,
        codeVerifier: codeVerifier,
        redirectAfterLogin: path,
    })

    const redirectTo: URL = client.buildAuthorizationUrl(config, parameters);
    return {redirectTo};
}

async function redirectToErrorPage(error: string, errorDescription?: string): Promise<Redirection> {
    // await clearAuthSession();

    const errorUrl = await redirectUrl("/auth/error")
    errorUrl.searchParams.set('error', error);
    if (errorDescription) {
        errorUrl.searchParams.set('error_description', errorDescription);
    }

    return {redirectTo: errorUrl};
}

export async function authenticate(requestUrl: URL): Promise<Redirection> {
    const searchParams = requestUrl.searchParams;
    const authenticationFlowState = await loadAuthSession();

    await clearAuthSession()

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

    let currentUrlWithBasePath = await requestUrlWithBaseBath(requestUrl);
    console.log("authenticate > currentUrlWithBasePath:", currentUrlWithBasePath.toString());

    try {
        const tokens: client.TokenEndpointResponse = await client.authorizationCodeGrant(
            config,
            currentUrlWithBasePath, // This is the URL of this route
            {
                pkceCodeVerifier: authenticationFlowState.codeVerifier,
                expectedState: authenticationFlowState.state,
                expectedNonce: authenticationFlowState.nonce,
            }
        );

        await storeSession({
            accessToken: tokens.access_token,
            accessTokenExpiresIn: tokens.expires_in ?? 3600,
            refreshToken: tokens.refresh_token ?? '',
            idToken: tokens.id_token ?? '',
        })

        return {redirectTo: (await redirectUrl(authenticationFlowState.redirectAfterLogin ?? "/"))}

    } catch (error) {
        console.error('OAuth callback error:', error);
        return redirectToErrorPage('token-exchange-failed');
    }
}