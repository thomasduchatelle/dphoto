import {Dispatch} from "react";

export interface OAuth2LogoutApi {
    logout(): Promise<void>;
}

export type LoggedOutAction = {
    type: 'logged-out'
}

export class LogoutCase {
    constructor(readonly dispatch: Dispatch<LoggedOutAction>,
                readonly oauthApi: OAuth2LogoutApi,
    ) {
    }

    public logout = () => {
        this.oauthApi
            .logout()
            .then(() => this.dispatch({type: 'logged-out'}))
    }
}

// TODO AGENT - generate the test: it should dispatch logged-out action after successful logout API call.
// TODO AGENT - implement the feature, and the test: it should dispatch a generic error if logout API call fails.
// TODO AGENT - convert the LogoutCase into a thunk as described in the coding conventions.