import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {tokenHolder} from "../security/client-utils";

export class DPhotoApplication {
    private axiosInterceptorId ?: number

    constructor(
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

    private axiosRequestInterceptor = (config: AxiosRequestConfig): AxiosRequestConfig => {
        // TODO AGENT - trigger a refresh if the token is expired or about to expire (<5 min)
        if (tokenHolder.accessToken) {
            config.headers = {
                ...config.headers,
                'Authorization': `Bearer ${tokenHolder.accessToken.accessToken}`,
            }
        }

        return config
    };
}