import {NextRequest} from 'next/server';

/**
 * Extracts the original origin (protocol + host + optional port) from a NextRequest.
 * 
 * When the application is deployed behind AWS API Gateway or CloudFront, the original
 * domain information is passed through X-Forwarded-* headers. This function attempts
 * to reconstruct the original URL from these headers, falling back to request.url
 * if the headers are not present.
 * 
 * @param request - The NextRequest object
 * @returns The original origin (e.g., "https://example.com" or "https://example.com:3000")
 */
export function getOriginalOrigin(request: NextRequest): string {
    const forwardedProto = request.headers.get('x-forwarded-proto');
    const forwardedHost = request.headers.get('x-forwarded-host');
    const forwardedPort = request.headers.get('x-forwarded-port');

    if (forwardedProto && forwardedHost) {
        const protocol = forwardedProto.split(',')[0].trim();
        const host = forwardedHost.split(',')[0].trim();
        
        // Only include port if it's non-standard for the protocol
        const isStandardPort = 
            (protocol === 'https' && (!forwardedPort || forwardedPort === '443')) ||
            (protocol === 'http' && (!forwardedPort || forwardedPort === '80'));
        
        if (isStandardPort || !forwardedPort) {
            return `${protocol}://${host}`;
        } else {
            return `${protocol}://${host}:${forwardedPort}`;
        }
    }

    // Fallback to request URL if X-Forwarded headers are not present
    return request.nextUrl.origin;
}
