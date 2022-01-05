import axios from 'axios'
import {GoogleLoginResponse, GoogleLoginResponseOffline} from "react-google-login";
import {SecurityContextType} from "./security.model";

interface IdentityResponse {
  email: string
  name: string
  picture: string
}

interface TokenResponse {
  access_token: string
  identity: IdentityResponse
}

export let owner = ""
export let accessToken = ""

function isValidResponse(value: GoogleLoginResponse | GoogleLoginResponseOffline): value is GoogleLoginResponse {
  return value.hasOwnProperty('profileObj');
}

export const authenticateWithGoogle = (googleAnswer: (GoogleLoginResponse | GoogleLoginResponseOffline)): Promise<SecurityContextType> => {
  if (isValidResponse(googleAnswer)) {
    return axios.post<TokenResponse>("/api/oauth/token", {}, {
      headers: {
        'Authorization': `Bearer ${googleAnswer.tokenId}`
      }
    }).then(resp => {
      console.log("access_token = " + resp.data.access_token)
      owner = resp.data.identity.email
      accessToken = resp.data.access_token
      return {
        loggedUser: {
          name: resp.data.identity.name,
          email: resp.data.identity.email,
          picture: resp.data.identity.picture,
        }
      }
    })
  }

  return Promise.reject("Not implemented.")
}

export const logoutFromGoogle = (): Promise<SecurityContextType> => {
  return Promise.resolve({})
}