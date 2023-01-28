import {PageState} from "./login-model";

type StopLoadingAction = {
    type: 'on-waiting-for-user-input'
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

export type PageAction = StopLoadingAction | ErrorAction | UpdateLoadingAction | OnSuccessfulAuthenticationAction

export const initialPageState: PageState = {
    error: "", loading: true, stage: "Please wait, authenticating...", timeout: false
}

export function reduce(current: PageState, action: PageAction): PageState {
    switch (action.type) {
        case "error":
            return {
                error: action.message, loading: false, stage: "", timeout: false
            }
        case "update-loading":
            return {
                ...current, loading: true, stage: action.message, error: "",
            }
        case 'on-waiting-for-user-input':
            return {
                ...current, loading: false, stage: "",
            }
    }

    return current
}

