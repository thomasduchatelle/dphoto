/** Backend session is only available on the server side (SSR) */
export type BackendSession = AuthenticatedSession | AnonymousSession

export interface AuthenticatedSession {
    type: 'authenticated'
    accessToken: AccessToken
    refreshToken: string
    authenticatedUser: AuthenticatedUser
}

export interface AnonymousSession {
    type: 'anonymous'
}

export function newAnonymousSession(): AnonymousSession {
    return {
        type: 'anonymous'
    }
}

/** Session which is available in the client side, as an atom */
export interface ClientSession {
    accessToken: AccessToken
    authenticatedUser: AuthenticatedUser
}

export interface AccessToken {
    accessToken: string
    /** Expiration date of the access token, to be refreshed before it expires */
    expiresAt: Date
}

export interface AuthenticatedUser {
    name: string
    email: string
    picture?: string
    // Whether the authenticated user is OWNER or VISITOR, extracted fron JWT claims
    isOwner: boolean
}

export function toClientSession(backendSession: BackendSession): ClientSession | null {
    if (backendSession.type === 'authenticated') {
        return {
            accessToken: backendSession.accessToken,
            authenticatedUser: backendSession.authenticatedUser,
        }
    } else {
        return null
    }
}