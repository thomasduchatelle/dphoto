import {AccessToken, LogoutListener} from "../security";
import axios, {AxiosInstance, InternalAxiosRequestConfig} from "axios";

export class DPhotoApplication {
    private accessToken?: AccessToken
    public readonly axiosInstance: AxiosInstance = axios.create({})
    private axiosInterceptorId ?: number

    constructor(public logoutListeners: LogoutListener[] = [],
                public authenticationTimeoutIds: NodeJS.Timeout[] = [],
    ) {
    }

    public renewRefreshToken(accessToken: AccessToken) {
        this.accessToken = accessToken
        if (!this.axiosInterceptorId) {
            this.axiosInterceptorId = this.axiosInstance.interceptors.request.use(this.axiosRequestInterceptor, error => Promise.reject(error));
        }
    }

    public revokeAccessToken() {
        if (this.axiosInterceptorId) {
            this.axiosInstance.interceptors.request.eject(this.axiosInterceptorId)
        }
    }

    public getAccessToken(): string {
        return this.accessToken?.accessToken ?? ''
    }

    private axiosRequestInterceptor = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
        if (this.accessToken) {
            config.headers['Authorization'] = `Bearer ${this.accessToken}`
        }

        return config
    };
}