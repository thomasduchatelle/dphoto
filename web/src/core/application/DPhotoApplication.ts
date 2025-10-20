import {AccessToken, LogoutListener} from "../security";
import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {AccessTokenHolder} from "./application-model";
import { getClientAccessToken } from "../../libs/auth/token-context";

export class DPhotoApplication implements AccessTokenHolder {
    private accessToken?: AccessToken
    private axiosInterceptorId ?: number

    constructor(
        public logoutListeners: LogoutListener[] = [],
        public authenticationTimeoutIds: NodeJS.Timeout[] = [],
        public readonly axiosInstance: AxiosInstance = axios.create({}),
    ) {
        // Set up axios interceptor immediately to use tokens from cookie/context
        if (!this.axiosInterceptorId) {
            this.axiosInterceptorId = this.axiosInstance.interceptors.request.use(this.axiosRequestInterceptor, error => Promise.reject(error));
        }
    }

    public renewRefreshToken(accessToken: AccessToken) {
        this.accessToken = accessToken
    }

    public revokeAccessToken() {
        if (this.axiosInterceptorId) {
            this.axiosInstance.interceptors.request.eject(this.axiosInterceptorId)
            this.axiosInterceptorId = undefined
        }
        this.accessToken = undefined
    }

    public getAccessToken(): string {
        // Try to get token from context (Cognito tokens) first
        const cognitoToken = getClientAccessToken();
        if (cognitoToken) {
            return cognitoToken;
        }
        
        // Fall back to legacy token
        return this.accessToken?.accessToken ?? ''
    }

    private axiosRequestInterceptor = (config: AxiosRequestConfig): AxiosRequestConfig => {
        // Get token from context (Cognito) or legacy source
        const token = this.getAccessToken();
        
        if (token) {
            config.headers = {
                ...config.headers,
                'Authorization': `Bearer ${token}`,
            }
        }

        return config
    };
}