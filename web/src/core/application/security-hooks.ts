import {AuthenticateCase, AuthenticatedUser, SecurityState} from "../security";
import {useContext, useMemo, useRef} from "react";
import {ApplicationContext} from "./application-context";
import {AuthenticationAPIAdapter} from "../../apis/oauthapi/AuthenticationAPIAdapter";
import {LogoutCase} from "../security/LogoutCase";
import {AxiosInstance} from "axios";
import {AccessForbiddenError} from "./errors";
import {useAxios} from "./application-hooks";

export const useSecurityState = (): SecurityState => {
    return useContext(ApplicationContext).context.security
}

export const useIsAuthenticated = (): boolean => {
    return !!useSecurityState().authenticatedUser
}

export const useAuthenticationCase = (): AuthenticateCase => {
    const {dispatch} = useContext(ApplicationContext);
    const authenticate = useRef(new AuthenticateCase(dispatch, new AuthenticationAPIAdapter()));
    return authenticate.current
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
        throw new AccessForbiddenError("user must be authenticated")
    }

    return useMemo(() => ({
        accessToken: application.getAccessToken(),
        authenticatedAxios: axiosInstance,
        loggedUser: authenticatedUser,
        signOutCase: logoutCase
    }), [authenticatedUser, application, logoutCase, axiosInstance])
}