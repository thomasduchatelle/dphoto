import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {AccessTokenHolder} from "./application-model";
import { loadClientSession } from "../../libs/auth/client-token-utils";

export class DPhotoApplication implements AccessTokenHolder {
    private axiosInterceptorId ?: number

    constructor(
        public logoutListeners: any[] = [],
        public authenticationTimeoutIds: NodeJS.Timeout[] = [],
        public readonly axiosInstance: AxiosInstance = axios.create({}),
    ) {
        // Set up axios interceptor to use tokens from cookies
        if (!this.axiosInterceptorId) {
            this.axiosInterceptorId = this.axiosInstance.interceptors.request.use(this.axiosRequestInterceptor, error => Promise.reject(error));
        }
    }

    public renewRefreshToken(accessToken: any) {
        // No-op for backwards compatibility with legacy code
    }

    public revokeAccessToken() {
        if (this.axiosInterceptorId) {
            this.axiosInstance.interceptors.request.eject(this.axiosInterceptorId)
            this.axiosInterceptorId = undefined
        }
    }

    public getAccessToken(): string {
        // Get token from Cognito cookies
        const session = loadClientSession();
        return session?.accessToken.value ?? '';
    }

    private axiosRequestInterceptor = (config: AxiosRequestConfig): AxiosRequestConfig => {
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