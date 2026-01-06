interface JWTPayload {
    sub?: string;
    Scopes?: string;  // Backend token - space-separated scopes with "owner:" prefix
    scope?: string;   // Cognito token - space-separated OAuth scopes
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