// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {NextRequest} from 'next/server';
import {getOriginalOrigin} from './request-utils';

describe('getOriginalOrigin', () => {
    it('should extract origin from RFC 7239 Forwarded header with https', () => {
        const request = new NextRequest('https://internal-api-gateway.example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com;proto=https',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should extract origin from RFC 7239 Forwarded header with http', () => {
        const request = new NextRequest('https://internal-api-gateway.example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'proto=http;host=example.com;for=83.106.145.60;by=3.248.245.105',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('http://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should accept forwarded host with port when hostname matches', () => {
        const request = new NextRequest('https://internal-api-gateway.example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com:8443;proto=https',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com:8443');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should handle quoted values in Forwarded header', () => {
        const request = new NextRequest('https://internal-api-gateway.example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'host="example.com";proto="https"',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should accept quoted host with port when hostname matches', () => {
        const request = new NextRequest('https://internal-api-gateway.example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'host="example.com:8443";proto="https"',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com:8443');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should fallback to request.url when Forwarded header is not present', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should fallback when Forwarded header is missing proto', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    it('should fallback when Forwarded header is missing host', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;proto=https',
            },
        });

        const result = getOriginalOrigin(request);

        expect(result.origin).toBe('https://example.com');
        expect(result.pathname).toBe('/nextjs/path');
    });

    describe('security validations', () => {
        it('should reject forwarded host that is not a subdomain of server host', () => {
            const request = new NextRequest('https://api-gateway.example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=malicious.com;proto=https',
                },
            });

            const result = getOriginalOrigin(request);

            expect(result.origin).toBe('https://api-gateway.example.com');
            expect(result.pathname).toBe('/nextjs/path');
        });

        it('should reject invalid protocol like ftp', () => {
            const request = new NextRequest('https://api-gateway.example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=api-gateway.example.com;proto=ftp',
                },
            });

            const result = getOriginalOrigin(request);

            expect(result.origin).toBe('https://api-gateway.example.com');
            expect(result.pathname).toBe('/nextjs/path');
        });

        it('should reject invalid protocol like javascript', () => {
            const request = new NextRequest('https://api-gateway.example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=api-gateway.example.com;proto=javascript',
                },
            });

            const result = getOriginalOrigin(request);

            expect(result.origin).toBe('https://api-gateway.example.com');
            expect(result.pathname).toBe('/nextjs/path');
        });

        it('should handle malformed Forwarded header gracefully', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'invalid-format',
                },
            });

            const result = getOriginalOrigin(request);

            expect(result.origin).toBe('https://example.com');
            expect(result.pathname).toBe('/nextjs/path');
        });
    });

    describe('basePath handling', () => {
        it('should add basePath to pathname', () => {
            const request = new NextRequest('https://example.com/albums', {
                method: 'GET',
            });

            const result = getOriginalOrigin(request);

            expect(result.pathname).toBe('/nextjs/albums');
        });

        it('should remove trailing slash after adding basePath', () => {
            const request = new NextRequest('https://example.com/albums/', {
                method: 'GET',
            });

            const result = getOriginalOrigin(request);

            expect(result.pathname).toBe('/nextjs/albums');
        });

        it('should handle root path correctly', () => {
            const request = new NextRequest('https://example.com/', {
                method: 'GET',
            });

            const result = getOriginalOrigin(request);

            expect(result.pathname).toBe('/nextjs');
        });
    });
});

