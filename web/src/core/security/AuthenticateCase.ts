import {AxiosError} from "axios";
import {Dispatch} from "react";
import {ApplicationAction} from "../application";
import {AuthenticatedUser, LogoutListener} from "./security-state";
import {AuthenticationPort} from "../../pages/Login/domain";
import {REFRESH_TOKEN_KEY} from "../../pages/Login/domain/login-controller";

interface ErrorBody {
    code: string
    error: string
}

export interface SuccessfulAuthenticationResponse {
    details: AuthenticatedUser
    accessToken: string
    refreshToken: string
    expiresIn: number
}

export interface AuthenticateAPI {
    authenticateWithIdentityToken(identityToken: string): Promise<SuccessfulAuthenticationResponse>

    refreshTokens(refreshToken: string): Promise<SuccessfulAuthenticationResponse>
}

export class AppAuthenticationError extends Error {
    constructor(public readonly code?: string,
                message ?: string) {
        super(message);
    }
}

export class AuthenticateCase implements AuthenticationPort {
    constructor(
        private dispatch: Dispatch<ApplicationAction>,
        private authenticateAPI: AuthenticateAPI,
    ) {
    }

    public authenticate = (identityToken: string, logoutListener: LogoutListener | undefined = undefined): Promise<SuccessfulAuthenticationResponse> => {
        return this.authenticateAPI.authenticateWithIdentityToken(identityToken)
            .then(user => {
                const timeoutId = this.scheduleTokensRefresh(user.refreshToken, user.expiresIn)

                this.dispatch({
                    accessToken: {
                        accessToken: user.accessToken,
                        expiryTime: user.expiresIn,
                    },
                    logoutListener: logoutListener,
                    refreshTimeoutId: timeoutId,
                    user: user.details,
                    type: 'authenticated'
                })
                return user
            })
            .catch((err: AxiosError<ErrorBody>) => {
                console.log(`ERROR: authentication failed ${JSON.stringify(err)}`)
                return Promise.reject(new AppAuthenticationError(err.response?.data?.code, err.response?.data?.error))
            })
    }

    public restoreSession = (refreshToken: string, logoutListener: LogoutListener | undefined): Promise<SuccessfulAuthenticationResponse> => {
        return this.authenticateAPI.refreshTokens(refreshToken)
            .then(user => {
                const timeoutId = this.scheduleTokensRefresh(user.refreshToken, user.expiresIn)

                this.dispatch({
                    accessToken: {
                        accessToken: user.accessToken,
                        expiryTime: user.expiresIn,
                    },
                    logoutListener: logoutListener,
                    refreshTimeoutId: timeoutId,
                    user: user.details,
                    type: 'authenticated'
                })
                return user
            })
            .catch((err: AxiosError<ErrorBody>) => {
                return Promise.reject(new AppAuthenticationError(err.response?.data?.code, err.response?.data?.error))
            })
    }

    private scheduleTokensRefresh = (refreshToken: string, expiresIn: number): NodeJS.Timeout => {
        const timeoutId = setTimeout(() => {
            this.authenticateAPI.refreshTokens(refreshToken)
                .then(user => {
                    const nextTimeoutId = this.scheduleTokensRefresh(user.refreshToken, user.expiresIn)
                    localStorage.setItem(REFRESH_TOKEN_KEY, user.refreshToken)

                    this.dispatch({
                        accessToken: {
                            accessToken: user.accessToken,
                            expiryTime: user.expiresIn,
                        },
                        nextTimeoutId: nextTimeoutId,
                        currentTimeoutId: timeoutId,
                        type: "refreshed-token",
                    })
                })
                .catch(err => {
                    console.log(`ERROR: refresh failed ${JSON.stringify(err)}`)
                    this.dispatch({type: 'timed-out'})
                })

        }, (expiresIn - 60) * 1000)

        return timeoutId
    }

}
