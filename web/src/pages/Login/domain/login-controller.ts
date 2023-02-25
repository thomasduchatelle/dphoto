import {AuthenticationPort, IdentityProviderError, LoadingPort, LoginPageActions} from "./login-model";
import {AppAuthenticationError, AuthenticatedUser, LogoutListener} from "../../../core/security";
import {Dispatch} from "react";
import {PageAction} from "./login-reducer";

const IDENTITY_TOKEN_KEY = "identityToken";

export class LoginController implements LoginPageActions {
    constructor(
        private dispatch: Dispatch<PageAction>,
        private authenticationPort: AuthenticationPort,
        private loadingPort: LoadingPort) {
    }

    public attemptToAutoAuthenticate = () => {
        const identityToken = localStorage.getItem(IDENTITY_TOKEN_KEY)
        if (identityToken) {
            this.authenticationPort
                .authenticate(identityToken, {
                    onLogout() {
                        localStorage.removeItem(IDENTITY_TOKEN_KEY)
                    }
                })
                .then(this.onSuccessfulAuthentication)
                .catch(err => {
                    console.log(`WARN: couldn't restore the session from identity token: ${err.message}`)
                    localStorage.removeItem(IDENTITY_TOKEN_KEY)

                    this.dispatch({type: "OnUnsuccessfulAutoLoginAttempt"})
                })

        } else {
            this.dispatch({type: "OnUnsuccessfulAutoLoginAttempt"})
        }
    }

    public loginWithIdentityToken = (identityToken: string, logoutListener?: LogoutListener): void => {
        this.authenticationPort
            .authenticate(identityToken, {
                onLogout() {
                    localStorage.removeItem(IDENTITY_TOKEN_KEY)
                    return logoutListener?.onLogout()
                }
            })
            .then(user => {
                localStorage.setItem(IDENTITY_TOKEN_KEY, identityToken)
                return user
            })
            .then(this.onSuccessfulAuthentication)
            .catch(this.onError)
    }

    private onSuccessfulAuthentication = (user: AuthenticatedUser): void => {
        this.dispatch({type: 'update-loading', message: "Please wait, loading your catalog..."})
        this.loadingPort.warmupApplication(user)
            .then(() => {
                this.dispatch({type: 'on-successful-authentication'})
            })
            .catch(this.onError)
    }

    public onError = (err: Error): void => {
        console.log(`ERROR ${err.name}: ${err.message}`)

        let message = "An unexpected error occurred, please report it to the maintainer."
        if (err instanceof IdentityProviderError) {
            message = "Authentication failed, please clear your cookies and retry."
        } else if (err instanceof AppAuthenticationError) {
            switch (err.code) {
                case "oauth.user-not-preregistered":
                    message = "You must be pre-registered to use this service."
                    break

                default:
                    message = "Sorry, DPhoto is not able to authenticate you at the moment. Please retry later."
            }
        }

        this.dispatch({type: 'error', message})
    }
}
