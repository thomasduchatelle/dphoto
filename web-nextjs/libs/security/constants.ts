export const ACCESS_TOKEN_COOKIE = 'dphoto-access-token';
export const REFRESH_TOKEN_COOKIE = 'dphoto-refresh-token';
export const OAUTH_STATE_COOKIE = 'dphoto-oauth-state';
export const OAUTH_CODE_VERIFIER_COOKIE = 'dphoto-oauth-code-verifier';
export const OAUTH_NONCE_COOKIE = 'dphoto-oauth-nonce';

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
