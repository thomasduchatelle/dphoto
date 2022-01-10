import axios from "axios";
import {AuthenticatedUser, OAuthService} from "../../security.domain";

interface IdentityResponse {
  email: string
  name: string
  picture: string
}

interface TokenResponse {
  access_token: string
  identity: IdentityResponse
}

export const OauthServiceImpl: OAuthService = {
  authenticateWithGoogleId(googleIdToken: string): Promise<AuthenticatedUser> {
    return axios.post<TokenResponse>("/api/oauth/token", {}, {
      headers: {
        'Authorization': `Bearer ${googleIdToken}`
      }
    }).then(resp => {
      return {
        name: resp.data.identity.name,
        email: resp.data.identity.email,
        picture: resp.data.identity.picture,
        accessToken: resp.data.access_token,
      }
    })
  }

}