export const initialSecurityState: SecurityState = {
    hasTimedOut: false,
}

export interface AuthenticatedUser {
    name: string
    email: string
    picture?: string
}

export interface AccessToken {
    accessToken: string
    expiryTime: number
}

export interface LogoutListener {
    onLogout(): void
}

export interface SecurityState {

    authenticatedUser?: AuthenticatedUser

    hasTimedOut: boolean
}

export const REFRESH_TOKEN_KEY = "refreshToken";