export const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
export const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';
export const OAUTH_STATE_COOKIE = 'dphoto-oauth-state';
export const OAUTH_CODE_VERIFIER_COOKIE = 'dphoto-oauth-code-verifier';
export const OAUTH_NONCE_COOKIE = 'dphoto-oauth-nonce';
export const REDIRECT_AFTER_LOGIN_COOKIE = 'dphoto-redirect-after-login';

export interface AuthenticatedUser {
    name: string;
    email: string;
    picture?: string;
    isOwner: boolean;
}

export interface AccessToken {
    accessToken: string;
    expiresAt: Date;
}

export interface BackendSession {
    type: 'authenticated';
    accessToken: AccessToken;
    authenticatedUser: AuthenticatedUser;
}

// Session refresh status constants for testing
export const ACTIVE_SESSION = { status: 'active' as const };
export const ANONYMOUS_SESSION = { status: 'anonymous' as const };
export const EXPIRED_SESSION = { status: 'expired' as const };
export const USER_INFO_COOKIE = 'dphoto-user-info';