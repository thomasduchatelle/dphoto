import {AuthenticateCase, AuthenticatedUser, SecurityState} from "../security";
import {useContext, useMemo, useRef} from "react";
import {ApplicationContext} from "./application-context";
import {AuthenticationAPIAdapter} from "../../apis/oauthapi/AuthenticationAPIAdapter";
import {LogoutCase} from "../security/LogoutCase";
import {AxiosInstance} from "axios";
import {AccessForbiddenError} from "./application-errors";
import {useAxios} from "./application-hooks";

export const useSecurityState = (): SecurityState => {
    return useContext(ApplicationContext).context.security
}

export const useIsAuthenticated = (): boolean => {
    return !!useSecurityState().authenticatedUser
}

export const useAuthenticatedUser = (): AuthenticatedUser => {
    const securityState = useSecurityState();
    if (!securityState.authenticatedUser) {
        throw new AccessForbiddenError("user is not authenticated")
    }
    return securityState.authenticatedUser;
}

export const useAuthenticationCase = (): AuthenticateCase => {
    const {dispatch} = useContext(ApplicationContext);
    const authenticate = useMemo(() => new AuthenticateCase(dispatch, new AuthenticationAPIAdapter()), [dispatch]);
    return authenticate
}

export const useLogoutCase = (): LogoutCase => {
    const {dispatch, context: {application}} = useContext(ApplicationContext);
    const logout = useRef(new LogoutCase(dispatch, application));
    return logout.current
}

export interface MustBeAuthenticated {
    loggedUser: AuthenticatedUser
    signOutCase: LogoutCase
    authenticatedAxios: AxiosInstance
    accessToken: string
}

export const useMustBeAuthenticated = (): MustBeAuthenticated => {
    const {context: {security: {authenticatedUser}, application}} = useContext(ApplicationContext);
    const axiosInstance = useAxios();
    const logoutCase = useLogoutCase();
    if (!authenticatedUser) {
        throw new AccessForbiddenError("user is not authenticated")
    }

    return useMemo(() => ({
        accessToken: application.getAccessToken(),
        authenticatedAxios: axiosInstance,
        loggedUser: authenticatedUser,
        signOutCase: logoutCase
    }), [authenticatedUser, application, logoutCase, axiosInstance])
}