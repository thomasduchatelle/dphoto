import {LoginPageActions, PageState} from "./login-model";
import {useAuthenticationCase, useSecurityState} from "../../../core/application";
import {useEffect, useMemo, useReducer} from "react";
import {LoginController} from "./login-controller";
import {initialLoginPageState, reduce} from "./login-reducer";
import {useCatalogLoader} from "../../../core/catalog";

export interface LoginPageState extends PageState, LoginPageActions {
}

export const useLoginPageCase = (onSuccessfulAuthentication: () => void): LoginPageState => {
    const {hasTimedOut} = useSecurityState();
    const [state, dispatch] = useReducer(reduce, {...initialLoginPageState, timeout: hasTimedOut})
    const authenticationCase = useAuthenticationCase()
    const catalogLoader = useCatalogLoader()

    const loginController = useMemo<LoginPageActions>(() => new LoginController(
        action => {
            if (action.type === "on-successful-authentication") {
                onSuccessfulAuthentication()
            } else {
                dispatch(action)
            }
        },
        authenticationCase,
        {
            warmupApplication: catalogLoader
        },
    ), [dispatch, onSuccessfulAuthentication, authenticationCase, catalogLoader])

    useEffect(() => {
        loginController.attemptToAutoAuthenticate()
    }, [loginController])

    return {
        ...state,
        ...loginController,
    }
}