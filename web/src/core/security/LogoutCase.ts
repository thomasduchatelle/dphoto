import {Dispatch} from "react";
import {ApplicationAction, DPhotoApplication} from "../application";

export class LogoutCase {
    constructor(readonly dispatch: Dispatch<ApplicationAction>,
                readonly application: DPhotoApplication) {
    }

    public logout = () => {
        this.application.authenticationTimeoutIds.forEach(timeoutId => {
            clearTimeout(timeoutId)
        })
        this.application.logoutListeners.forEach(listener => {
            listener.onLogout()
        })

        this.dispatch({type: 'logged-out'})
    }
}