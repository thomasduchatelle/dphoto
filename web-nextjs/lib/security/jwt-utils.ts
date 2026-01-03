interface JWTPayload {
    sub?: string;
    Scopes?: string;
    name?: string;
    email?: string;
    picture?: string;
    exp?: number;
    iat?: number;
    [key: string]: any;
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
        // Handle both base64 and base64url encoding
        const decoded = JSON.parse(
            Buffer.from(payload.replace(/-/g, '+').replace(/_/g, '/'), 'base64').toString('utf-8')
        );
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
    return scopes.some((scope: string) => scope.startsWith('owner:'));
}
