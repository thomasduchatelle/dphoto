import {AccessToken, AuthenticatedUser, LogoutListener} from "./security-state";
import {ApplicationContextType} from "../application";


export type AuthenticatedAction = {
    type: 'authenticated'
    accessToken: AccessToken
    user: AuthenticatedUser
    refreshTimeoutId: NodeJS.Timeout
    logoutListener?: LogoutListener
}

export type RefreshedTokenAction = {
    type: 'refreshed-token'
    accessToken: AccessToken
    currentTimeoutId: NodeJS.Timeout
    nextTimeoutId: NodeJS.Timeout
}

export type LoggedOutAction = {
    type: 'logged-out'
}

export type TimedOutAction = {
    type: 'timed-out'
}

export type SecurityAction = AuthenticatedAction | RefreshedTokenAction | LoggedOutAction | TimedOutAction

export function securityContextReducerSupports(action: any) {
    return action && action.type && ['authenticated', 'refreshed-token', 'logged-out', 'timed-out'].includes(action.type)
}

export function securityContextReducer(current: ApplicationContextType, action: SecurityAction): ApplicationContextType {
    switch (action.type) {
        case "authenticated":
            current.application.authenticationTimeoutIds.push(action.refreshTimeoutId)
            current.application.renewRefreshToken(action.accessToken)
            current.application.logoutListeners = action.logoutListener ? [action.logoutListener] : []

            return {
                ...current,
                security: {
                    ...current.security,
                    authenticatedUser: action.user,
                    hasTimedOut: false,
                },
            }

        case "refreshed-token":
            current.application.renewRefreshToken(action.accessToken)
            current.application.authenticationTimeoutIds = [...current.application.authenticationTimeoutIds.filter(id => id !== action.currentTimeoutId), action.nextTimeoutId]
            break

        case "logged-out":
            current.application.revokeAccessToken()
            current.application.authenticationTimeoutIds = []
            current.application.logoutListeners = []

            return {
                ...current,
                security: {
                    ...current.security,
                    authenticatedUser: undefined,
                    hasTimedOut: false,
                },
            }

        case 'timed-out':
            current.application.authenticationTimeoutIds = []
            current.application.logoutListeners = []
            current.application.revokeAccessToken()

            return {
                ...current,
                security: {
                    hasTimedOut: true,
                }
            }
    }

    return current
}