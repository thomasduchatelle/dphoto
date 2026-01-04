import {NextRequest} from 'next/server';
import {basePath} from "./basepath";

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
            return {proto, host};
        }
    } catch (e) {
        // Parsing failed, return null
    }
    return null;
}

function overloadWithForwardedUrl(serverUrl: URL, forwardedHeader: string): URL {
    const forwardedFrom = parseForwardedHeader(forwardedHeader);
    if (forwardedFrom
        && (forwardedFrom.proto === 'http' || forwardedFrom.proto === 'https') // only HTTP and HTTPS protocols are allowed
        && isParentDomain(forwardedFrom.host, serverUrl.host)  // only parent domains are allowed
    ) {
        let url = new URL(serverUrl.toString());

        url.protocol = forwardedFrom.proto + ':';
        url.host = forwardedFrom.host;

        return url
    }

    console.warn("WARNING: Ignoring invalid Forwarded header:", forwardedHeader);
    return serverUrl;
}

function isParentDomain(parentHost: string, childHost: string): boolean {
    // Extract hostname without port
    const parent = parentHost.split(':')[0];
    const child = childHost.split(':')[0];

    // Same host is allowed
    if (parent === child) {
        return true;
    }

    // Check if child is a subdomain of parent
    // e.g., internal-gateway.my-domain.com should match parent my-domain.com
    return child.endsWith('.' + parent);
}

export function getOriginalOrigin(request: NextRequest): URL {
    let url = new URL(request.url)

    const forwardedHeader = request.headers.get('forwarded');
    if (forwardedHeader) {
        url = overloadWithForwardedUrl(url, forwardedHeader);
    }

    if (basePath) {
        url.pathname = basePath + url.pathname

        if (url.pathname.endsWith("/")) {
            url.pathname = url.pathname.slice(0, -1);
        }
    }

    return url
}
