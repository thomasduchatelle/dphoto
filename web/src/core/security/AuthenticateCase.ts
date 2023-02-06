import {AxiosError} from "axios";
import {Dispatch} from "react";
import {ApplicationAction} from "../application";
import {AuthenticatedUser, LogoutListener} from "./security-state";

interface ErrorBody {
    code: string
    error: string
}

export interface SuccessfulAuthenticationResponse extends AuthenticatedUser {
    accessToken: string
    expiresIn: number
}

export interface AuthenticateAPI {
    authenticateWithIdentityToken(identityToken: string): Promise<SuccessfulAuthenticationResponse>
}

export class AppAuthenticationError extends Error {
    constructor(public readonly code?: string,
                message ?: string) {
        super(message);
    }
}

export class AuthenticateCase {
    constructor(
        private dispatch: Dispatch<ApplicationAction>,
        private authenticateAPI: AuthenticateAPI,
    ) {
    }

    public authenticate = (identityToken: string, logoutListener: LogoutListener | undefined = undefined): Promise<AuthenticatedUser> => {
        return this.authenticateAPI.authenticateWithIdentityToken(identityToken)
            .then(user => {
                const timeoutId = this.refreshToken(identityToken, user.expiresIn)

                this.dispatch({
                    accessToken: {
                        accessToken: user.accessToken,
                        expiryTime: user.expiresIn,
                    },
                    logoutListener: logoutListener,
                    refreshTimeoutId: timeoutId,
                    user: user,
                    type: 'authenticated'
                })
                return user
            })
            .catch((err: AxiosError<ErrorBody>) => {
                console.log(`ERROR: authentication failed ${JSON.stringify(err)}`)
                return Promise.reject(new AppAuthenticationError(err.response?.data?.code, err.response?.data?.error))
            })
    }

    private refreshToken = (identityToken: string, expiresIn: number): NodeJS.Timeout => {
        const timeoutId = setTimeout(() => {
            this.authenticateAPI.authenticateWithIdentityToken(identityToken)
                .then(user => {
                    const nextTimeoutId = this.refreshToken(identityToken, user.expiresIn)

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
