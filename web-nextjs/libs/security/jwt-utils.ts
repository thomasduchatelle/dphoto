interface JWTPayload {
    sub?: string;
    scope?: string;              // Cognito token - space-separated OAuth scopes
    "dphoto:email"?: string;     // Injected by the pre-token-generation trigger - the DPhoto UserId
    "dphoto:owner"?: string;     // Injected by the pre-token-generation trigger - empty/absent for non-owners
    "dphoto:scopes"?: string;    // Injected by the pre-token-generation trigger - space-separated DPhoto scopes
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