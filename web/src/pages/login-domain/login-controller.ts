import {AuthenticationPort, IdentityProviderError, LoadingPort, LoginPageActions} from "./login-model";
import {AppAuthenticationError, ExpiredSessionError, LogoutListener, SuccessfulAuthenticationResponse} from "../../core/security";
import {Dispatch} from "react";
import {initialLoginPageState, PageAction} from "./login-reducer";

export class LoginController implements LoginPageActions {
    constructor(
        private dispatch: Dispatch<PageAction>,
        private authenticationPort: AuthenticationPort,
        private loadingPort: LoadingPort) {
    }

    public attemptToAutoAuthenticate = () => {
        this.authenticationPort
            .restoreSession()
            .then(this.onSuccessfulAuthentication)
            .catch(err => {
                console.log(`WARN: couldn't restore the session from refresh token: ${err.message}`)

                if (err instanceof ExpiredSessionError) {
                    this.dispatch({type: "OnExpiredSession"});

                } else {
                    this.dispatch({type: "OnUnsuccessfulAutoLoginAttempt"});
                }
            })
    }

    public loginWithIdentityToken = (identityToken: string, logoutListener?: LogoutListener): void => {
        this.dispatch({type: "update-loading", message: initialLoginPageState.stage})
        this.authenticationPort
            .authenticate(identityToken, {
                onLogout() {
                    return logoutListener?.onLogout()
                }
            })
            .then(this.onSuccessfulAuthentication)
            .catch(this.onError)
    }

    private onSuccessfulAuthentication = (user: SuccessfulAuthenticationResponse): void => {
        this.dispatch({type: 'update-loading', message: "Please wait, loading your state..."})
        this.loadingPort.warmupApplication(user.details)
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
