import axios from "axios";
import {AuthenticateAPI, SuccessfulAuthenticationResponse} from "../../core/security";

interface IdentityResponse {
    email: string
    name: string
    picture: string
}

interface TokenResponse {
    access_token: string
    identity: IdentityResponse
    expires_in: number
}

export class AuthenticationAPIAdapter implements AuthenticateAPI {

    public authenticateWithIdentityToken(identityToken: string): Promise<SuccessfulAuthenticationResponse> {
        return axios.post<TokenResponse>("/api/oauth/token", {}, {
            headers: {
                'Authorization': `Bearer ${identityToken}`
            }
        }).then(resp => {
            return {
                details: {
                    name: resp.data.identity.name,
                    email: resp.data.identity.email,
                    picture: resp.data.identity.picture,
                },
                accessToken: resp.data.access_token,
                expiresIn: resp.data.expires_in
            }
        })
    }

}
