import {Dispatch} from "react";
import {ApplicationAction, DPhotoApplication} from "../application";
import {AuthenticateAPI} from "./AuthenticateCase";
import {REFRESH_TOKEN_KEY} from "./security-state";

export class LogoutCase {
    constructor(readonly dispatch: Dispatch<ApplicationAction>,
                readonly application: DPhotoApplication,
                readonly oauthApi: AuthenticateAPI,
    ) {
    }

    public logout = () => {
        this.application.authenticationTimeoutIds.forEach(timeoutId => {
            clearTimeout(timeoutId)
        })
        this.application.logoutListeners.forEach(listener => {
            listener.onLogout()
        })

        const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
        if (refreshToken) {
            localStorage.removeItem(REFRESH_TOKEN_KEY)

            this.oauthApi
                .logout(refreshToken)
                .then(() => this.dispatch({type: 'logged-out'}))
        } else {
            this.dispatch({type: 'logged-out'})
        }
    }
}