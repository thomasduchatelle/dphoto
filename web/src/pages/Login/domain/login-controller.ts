import {AuthenticationPort, IdentityProviderError, LoadingPort, LoginPageActions} from "./login-model";
import {AppAuthenticationError, LogoutListener} from "../../../core/security";
import {Dispatch} from "react";
import {PageAction} from "./login-reducer";

export class LoginController implements LoginPageActions {
    constructor(
        private dispatch: Dispatch<PageAction>,
        private authenticationPort: AuthenticationPort,
        private loadingPort: LoadingPort) {
    }

    public onWaitingForUserInput = (): void => {
        this.dispatch({type: 'on-waiting-for-user-input'})
    }

    public loginWithIdentityToken = (identityToken: string, logoutListener?: LogoutListener): void => {
        this.authenticationPort.authenticate(identityToken, logoutListener)
            .then(user => {
                this.dispatch({type: 'update-loading', message: "Please wait, loading your catalog..."})
                return this.loadingPort.warmupApplication(user)
            })
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
