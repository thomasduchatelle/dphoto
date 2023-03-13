import {AccessToken, LogoutListener} from "../security";
import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {AccessTokenHolder} from "./application-model";

export class DPhotoApplication implements AccessTokenHolder {
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

    private axiosRequestInterceptor = (config: AxiosRequestConfig): AxiosRequestConfig => {
        if (this.accessToken) {
            config.headers = {
                ...config.headers,
                'Authorization': `Bearer ${this.accessToken.accessToken}`,
            }
        }

        return config
    };
}