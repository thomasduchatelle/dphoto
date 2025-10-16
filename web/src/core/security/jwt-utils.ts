export interface JWTPayload {
    sub: string // Subject (user email)
    Scopes: string // Space-separated scopes
    iss: string // Issuer
    aud: string[] // Audience
    exp: number // Expiration time
    iat: number // Issued at
    jti: string // JWT ID
}

/**
 * Decode JWT payload without verification.
 * Note: This is safe because the JWT is verified on the backend.
 * We only need to read the payload to extract user information.
 */
export function decodeJWTPayload(token: string): JWTPayload | null {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) {
            return null;
        }

        const payload = parts[1];
        const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
        return decoded as JWTPayload;
    } catch (error) {
        console.error('Failed to decode JWT:', error);
        return null;
    }
}

/**
 * Check if the user is an owner based on JWT scopes.
 * An owner has scopes starting with "owner:", while a visitor has only "visitor" scope.
 */
export function isOwnerFromJWT(token: string): boolean {
    const payload = decodeJWTPayload(token);
    if (!payload || !payload.Scopes) {
        return false;
    }

    const scopes = payload.Scopes.split(' ');
    return scopes.some(scope => scope.startsWith('owner:'));
}
