import {LogoutListener} from "../../../core/security";

export class IdentityProviderError extends Error {
}

export interface LoginPageActions {
    onWaitingForUserInput(): void

    loginWithIdentityToken(identityToken: string): void

    onError(err: Error): void
}

export interface PageState {
    loading: boolean

    stage: string

    error: string

    timeout: boolean
}

export interface AuthenticationPort {

    authenticate(identityToken: string, logoutListener: LogoutListener | undefined): Promise<void>
}

export interface LoadingPort {

    warmupApplication(): Promise<void>
}

