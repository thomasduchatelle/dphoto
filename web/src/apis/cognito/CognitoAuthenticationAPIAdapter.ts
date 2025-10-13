import axios from "axios";
import {AuthenticateAPI, SuccessfulAuthenticationResponse} from "../../core/security/AuthenticateCase";

/**
 * Cognito authentication adapter.
 * Note: With Cognito and SSR, most authentication happens server-side.
 * Tokens are stored in HttpOnly cookies and managed by the SSR lambda.
 * This adapter is primarily for compatibility with the existing authentication flow.
 */
export class CognitoAuthenticationAPIAdapter implements AuthenticateAPI {

    /**
     * For Cognito flow, this method is not used as authentication happens via OAuth redirect.
     * The user is redirected to Cognito login, then back to /auth/callback.
     * This is kept for interface compatibility.
     */
    public authenticateWithIdentityToken(identityToken: string): Promise<SuccessfulAuthenticationResponse> {
        // In Cognito flow, we don't use identity tokens directly
        // Instead, users are redirected to Cognito hosted UI
        return Promise.reject(new Error("Cognito uses OAuth redirect flow, not identity tokens"));
    }

    /**
     * Refresh tokens using the Cognito refresh token.
     * With HttpOnly cookies, this is typically handled server-side by SSR.
     */
    public refreshTokens(refreshToken: string): Promise<SuccessfulAuthenticationResponse> {
        return axios.post<CognitoTokenResponse>("/api/auth/refresh", {
            refreshToken,
        }, {
            headers: {
                "Content-Type": "application/json",
            },
            withCredentials: true, // Important for cookies
        }).then(resp => {
            return {
                details: {
                    name: resp.data.name,
                    email: resp.data.email,
                    picture: resp.data.picture || '',
                },
                accessToken: resp.data.accessToken,
                refreshToken: resp.data.refreshToken,
                expiresIn: resp.data.expiresIn
            };
        });
    }

    /**
     * Logout from Cognito.
     * Clears server-side session and HttpOnly cookies.
     */
    public logout(refreshToken: string): Promise<void> {
        return axios.post<void>("/api/auth/logout", {
            refreshToken,
        }, {
            headers: {
                "Content-Type": "application/json",
            },
            withCredentials: true, // Important for cookies
        }).then();
    }

    /**
     * Initiate Cognito login flow by redirecting to Cognito hosted UI.
     */
    public initiateLogin(returnUrl?: string): void {
        const params = new URLSearchParams();
        if (returnUrl) {
            params.set('returnUrl', returnUrl);
        }
        window.location.href = `/api/auth/login${params.toString() ? '?' + params.toString() : ''}`;
    }
}

interface CognitoTokenResponse {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    name: string;
    email: string;
    picture?: string;
}
