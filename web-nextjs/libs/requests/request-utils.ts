import {NextRequest} from 'next/server';
import {basePath} from "./basepath";
import {headers} from "next/headers";

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

function overloadWithForwardedUrl(serverUrl: URL, forwardedHeader: string | null): URL {
    if (!forwardedHeader) {
        return serverUrl
    }

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

function getRequestURL(request: NextRequest): URL {
    let url = new URL(request.url)

    const forwardedHeader = request.headers.get('forwarded');
    if (forwardedHeader) {
        url = overloadWithForwardedUrl(url, forwardedHeader);
    }

    if (basePath && !url.pathname.startsWith(basePath)) { // the proxy gets the base path in its request URL.
        url.pathname = basePath + url.pathname

        if (url.pathname.endsWith("/")) {
            url.pathname = url.pathname.slice(0, -1);
        }
    }

    return url
}

export interface Origin {
    getCurrentUrl(): Promise<URL>
}

export function newOriginFromRequest(request: NextRequest): Origin {
    return {
        async getCurrentUrl(): Promise<URL> {
            return getRequestURL(request);
        }
    }
}

export function newOriginFromHeaders(): Origin {
    return {
        async getCurrentUrl(): Promise<URL> {
            const h = await headers()

            const host = h.get("host")
            const proto = host?.startsWith("localhost") || host?.startsWith("127.0.0.1") ? "http" : "https"
            return overloadWithForwardedUrl(new URL(`${proto || 'http'}://${host}`), h.get("forwarded"));
        }
    }
}

/**
 * Generates a redirect URL from the original HOST and PROTO (if the request got forwarded), and which includes the basepath (if set)
 */
export async function redirectUrl(path: string, origin: Origin = newOriginFromHeaders()): Promise<URL> {
    return new URL(`${basePath}${path}`, await origin.getCurrentUrl());
}


export async function requestUrlWithBaseBath(request: URL): Promise<URL> {
    const h = await headers()
    let updated = new URL(request)

    if (basePath && !request.pathname.startsWith(basePath)) {
        updated.pathname = `${basePath}${request.pathname}`;
    }

    return overloadWithForwardedUrl(updated, h.get("forwarded"));
}
