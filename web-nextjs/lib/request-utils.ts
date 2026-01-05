import {NextRequest} from 'next/server';

/**
 * Parses the RFC 7239 Forwarded header.
 * Format: "by=<identifier>;for=<identifier>;host=<host>;proto=<protocol>"
 * Note: Values can be quoted according to RFC 7239
 * 
 * @param forwardedHeader - The Forwarded header value
 * @returns Object with extracted proto and host, or null if parsing fails
 */
function parseForwardedHeader(forwardedHeader: string): { proto: string; host: string } | null {
    try {
        const parts = forwardedHeader.split(';');
        let proto: string | undefined;
        let host: string | undefined;

        for (const part of parts) {
            const equalsIndex = part.indexOf('=');
            if (equalsIndex === -1) continue;
            
            const key = part.substring(0, equalsIndex).trim();
            let value = part.substring(equalsIndex + 1).trim();
            
            // Remove quotes if present (RFC 7239 allows quoted values)
            if (value.startsWith('"') && value.endsWith('"')) {
                value = value.substring(1, value.length - 1);
            }
            
            if (key === 'proto') {
                proto = value;
            } else if (key === 'host') {
                host = value;
            }
        }

        if (proto && host) {
            return { proto, host };
        }
    } catch (e) {
        // Parsing failed, return null
    }
    return null;
}

/**
 * Extracts the original origin (protocol + host + optional port) from a NextRequest.
 * 
 * When the application is deployed behind AWS API Gateway or CloudFront, the original
 * domain information is passed through the RFC 7239 Forwarded header. This function attempts
 * to reconstruct the original URL from this header, falling back to request.nextUrl.origin
 * if the header is not present or invalid.
 * 
 * @param request - The NextRequest object
 * @returns The original origin (e.g., "https://example.com" or "https://example.com:3000")
 */
export function getOriginalOrigin(request: NextRequest): string {
    const forwardedHeader = request.headers.get('forwarded');

    if (forwardedHeader) {
        const parsed = parseForwardedHeader(forwardedHeader);
        
        if (parsed) {
            const { proto, host } = parsed;
            
            // Extract and validate protocol - must be http or https
            const protocol = proto.toLowerCase();
            if (protocol !== 'http' && protocol !== 'https') {
                // Invalid protocol, fallback to original URL
                return request.nextUrl.origin;
            }
            
            // Extract host and optional port
            const normalizedHost = host.toLowerCase();
            
            // Basic host validation:
            // - Must not be empty
            // - Must start and end with alphanumeric (allowing port)
            // - Can contain alphanumeric, dots, hyphens, and colons in between
            // - Must not have consecutive dots
            if (!normalizedHost || normalizedHost.includes('..')) {
                // Invalid host, fallback to original URL
                return request.nextUrl.origin;
            }
            
            // Check if host contains a port
            const hostParts = normalizedHost.split(':');
            const hostWithoutPort = hostParts[0];
            const portStr = hostParts[1];
            
            // Validate host part (without port)
            if (!/^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$/i.test(hostWithoutPort)) {
                // Invalid host, fallback to original URL
                return request.nextUrl.origin;
            }
            
            // Validate port if present
            let port: number | undefined;
            if (portStr) {
                const portNum = parseInt(portStr, 10);
                if (!isNaN(portNum) && Number.isInteger(portNum) && portNum >= 1 && portNum <= 65535) {
                    port = portNum;
                } else {
                    // Invalid port, fallback to original URL
                    return request.nextUrl.origin;
                }
            }
            
            // Only include port if it's valid and non-standard for the protocol
            const isStandardPort = 
                (protocol === 'https' && (!port || port === 443)) ||
                (protocol === 'http' && (!port || port === 80));
            
            if (!port || isStandardPort) {
                return `${protocol}://${hostWithoutPort}`;
            } else {
                return `${protocol}://${hostWithoutPort}:${port}`;
            }
        }
    }

    // Fallback to request URL if Forwarded header is not present or invalid
    return request.nextUrl.origin;
}
