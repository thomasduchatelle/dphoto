import {LoginPageActions, PageState} from "./login-model";
import {useAuthenticationCase, useSecurityState} from "../../../core/application";
import {useReducer, useRef} from "react";
import {LoginController} from "./login-controller";
import {initialPageState, reduce} from "./login-reducer";

export interface LoginPageState extends PageState, LoginPageActions {
}

export const useLoginPageCase = (onSuccessfulAuthentication: () => void): LoginPageState => {
    const {hasTimedOut} = useSecurityState();
    const [state, dispatch] = useReducer(reduce, {...initialPageState, timeout: hasTimedOut})
    const authenticationCase = useAuthenticationCase()

    const loginController = useRef<LoginPageActions>(new LoginController(
        action => {
            if (action.type === "on-successful-authentication") {
                onSuccessfulAuthentication()
            } else {
                dispatch(action)
            }
        },
        authenticationCase,
        {
            warmupApplication: () => Promise.resolve()
        },
    ))

    return {
        ...state,
        ...loginController.current,
    }
}