import axios from "axios";
import {AuthenticateAPI, SuccessfulAuthenticationResponse} from "../../core/security";

interface IdentityResponse {
    email: string
    name: string
    picture: string
}

interface TokenResponse {
    access_token: string
    refresh_token: string
    identity: IdentityResponse
    expires_in: number
}

export class AuthenticationAPIAdapter implements AuthenticateAPI {

    public authenticateWithIdentityToken(identityToken: string): Promise<SuccessfulAuthenticationResponse> {
        return axios.post<TokenResponse>("/oauth/token", new URLSearchParams({
            "grant_type": "identity",
            "identity_token": identityToken,
        }), {
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            }
        }).then(resp => {
            return {
                details: {
                    name: resp.data.identity.name,
                    email: resp.data.identity.email,
                    picture: resp.data.identity.picture,
                },
                accessToken: resp.data.access_token,
                refreshToken: resp.data.refresh_token,
                expiresIn: resp.data.expires_in
            }
        })
    }

    public refreshTokens(refreshToken: string): Promise<SuccessfulAuthenticationResponse> {
        return axios.post<TokenResponse>("/oauth/token", new URLSearchParams({
            "grant_type": "refresh_token",
            "refresh_token": refreshToken,
        }), {
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            }
        }).then(resp => {
            return {
                details: {
                    name: resp.data.identity.name,
                    email: resp.data.identity.email,
                    picture: resp.data.identity.picture,
                },
                accessToken: resp.data.access_token,
                refreshToken: resp.data.refresh_token,
                expiresIn: resp.data.expires_in
            }
        })
    }

    public logout(refreshToken: string): Promise<void> {
        return axios.post<void>("/oauth/logout", {
            refreshToken,
        }, {
            headers: {
                "Content-Type": "application/json",
            }
        }).then()
    }

}
