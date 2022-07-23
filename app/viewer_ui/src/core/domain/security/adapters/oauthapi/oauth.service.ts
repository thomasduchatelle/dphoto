import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {AuthenticatedUser, OAuthService} from "../../security.domain";

const instance = axios.create({})

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

export class OauthServiceImpl implements OAuthService {
  constructor(private dphotoAccessToken?: string,
              private axiosInterceptorId?: number) {
  }

  public clearTokens = (): void => {
    if (this.axiosInterceptorId) {
      instance.interceptors.request.eject(this.axiosInterceptorId);
    }
    this.dphotoAccessToken = undefined
  }

  public dispatchAccessToken = (accessToken: string): void => {
    this.dphotoAccessToken = accessToken
    if (!this.axiosInterceptorId) {
      this.axiosInterceptorId = instance.interceptors.request.use(this.axiosRequestInterceptor);
    }
  }

  public authenticateWithGoogleId = (googleIdToken: string): Promise<AuthenticatedUser> => {
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
        expiresIn: resp.data.expires_in
      }
    })
  }

  public getDPhotoAccessToken(): string {
    return this.dphotoAccessToken ?? ""
  }

  private axiosRequestInterceptor = (config: AxiosRequestConfig): Promise<AxiosRequestConfig> => {
    if (!this.dphotoAccessToken) {
      // safeguard - interceptor should have been ejected before
      return Promise.resolve(config)
    }

    return Promise.resolve({
      ...config,
      headers: {
        ...config.headers,
        'Authorization': `Bearer ${this.dphotoAccessToken}`
      }
    })
  };
}

export function authenticatedAxios(): AxiosInstance {
  return instance
}
