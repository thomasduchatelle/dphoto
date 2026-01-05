import {NextRequest} from 'next/server';

/**
 * Extracts the original origin (protocol + host + optional port) from a NextRequest.
 * 
 * When the application is deployed behind AWS API Gateway or CloudFront, the original
 * domain information is passed through X-Forwarded-* headers. This function attempts
 * to reconstruct the original URL from these headers, falling back to request.nextUrl.origin
 * if the headers are not present or invalid.
 * 
 * @param request - The NextRequest object
 * @returns The original origin (e.g., "https://example.com" or "https://example.com:3000")
 */
export function getOriginalOrigin(request: NextRequest): string {
    const forwardedProto = request.headers.get('x-forwarded-proto');
    const forwardedHost = request.headers.get('x-forwarded-host');
    const forwardedPort = request.headers.get('x-forwarded-port');

    if (forwardedProto && forwardedHost) {
        // Extract and validate protocol - must be http or https
        const rawProtocol = forwardedProto.split(',')[0].trim().toLowerCase();
        if (rawProtocol !== 'http' && rawProtocol !== 'https') {
            // Invalid protocol, fallback to original URL
            return request.nextUrl.origin;
        }
        const protocol = rawProtocol;
        
        // Extract and normalize host
        const host = forwardedHost.split(',')[0].trim().toLowerCase();
        
        // Basic host validation - must not be empty and contain valid characters
        if (!host || !/^[a-z0-9.-]+$/i.test(host)) {
            // Invalid host, fallback to original URL
            return request.nextUrl.origin;
        }
        
        // Extract and validate port
        let port: number | undefined;
        if (forwardedPort) {
            const portStr = forwardedPort.split(',')[0].trim();
            const portNum = Number(portStr);
            
            // Validate port is a valid integer between 1 and 65535
            if (Number.isInteger(portNum) && portNum >= 1 && portNum <= 65535 && portStr === portNum.toString()) {
                port = portNum;
            }
        }
        
        // Only include port if it's valid and non-standard for the protocol
        const isStandardPort = 
            (protocol === 'https' && (!port || port === 443)) ||
            (protocol === 'http' && (!port || port === 80));
        
        if (!port || isStandardPort) {
            return `${protocol}://${host}`;
        } else {
            return `${protocol}://${host}:${port}`;
        }
    }

    // Fallback to request URL if X-Forwarded headers are not present
    return request.nextUrl.origin;
}
