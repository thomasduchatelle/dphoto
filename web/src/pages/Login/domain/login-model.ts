import {AuthenticatedUser, LogoutListener} from "../../../core/security";

export class IdentityProviderError extends Error {
}

export interface LoginPageActions {
    attemptToAutoAuthenticate(): void

    loginWithIdentityToken(identityToken: string): void

    onError(err: Error): void
}

export interface PageState {
    loading: boolean

    stage: string

    promptForLogin: boolean

    error: string

    timeout: boolean
}

export interface AuthenticationPort {

    authenticate(identityToken: string, logoutListener: LogoutListener | undefined): Promise<AuthenticatedUser>
}

export interface LoadingPort {

    warmupApplication(user: AuthenticatedUser): Promise<void>
}

