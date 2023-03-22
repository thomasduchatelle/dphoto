import {PageState} from "./login-model";

type OnUnsuccessfulAutoLoginAttemptAction = {
    type: 'OnUnsuccessfulAutoLoginAttempt'
}

type OnExpiredSessionAction = {
    type: 'OnExpiredSession'
}

type ErrorAction = {
    type: 'error'
    message: string
}

type UpdateLoadingAction = {
    type: 'update-loading'
    message: string
}

type OnSuccessfulAuthenticationAction = {
    type: 'on-successful-authentication'
}

export type PageAction =
    OnUnsuccessfulAutoLoginAttemptAction
    | ErrorAction
    | UpdateLoadingAction
    | OnSuccessfulAuthenticationAction
    | OnExpiredSessionAction

export const initialLoginPageState: PageState = {
    error: "",
    loading: true,
    stage: "Please wait, authenticating...",
    timeout: false,
    promptForLogin: false,
}

export function reduce(current: PageState, action: PageAction): PageState {
    switch (action.type) {
        case "error":
            return {
                error: action.message, loading: false, stage: "", timeout: false, promptForLogin: true,
            }
        case "update-loading":
            return {
                ...current, loading: true, stage: action.message, error: "",
            }
        case 'OnUnsuccessfulAutoLoginAttempt':
            return {
                ...current, loading: false, stage: "", promptForLogin: true,
            }
        case 'OnExpiredSession':
            return {
                error: "",
                timeout: true,
                loading: false,
                stage: "",
                promptForLogin: true,
            }
    }

    return current
}

