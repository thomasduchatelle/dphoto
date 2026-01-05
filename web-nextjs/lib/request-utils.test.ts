// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {NextRequest} from 'next/server';
import {getOriginalOrigin} from './request-utils';

describe('getOriginalOrigin', () => {
    it('should extract origin from RFC 7239 Forwarded header', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com;proto=https',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should handle Forwarded header with different order', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'proto=https;host=example.com;for=83.106.145.60;by=3.248.245.105',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should extract host with port from Forwarded header', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com:8443;proto=https',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com:8443');
    });

    it('should not include standard HTTPS port 443', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com:443;proto=https',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should not include standard HTTP port 80', () => {
        const request = new NextRequest('http://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com:80;proto=http',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('http://example.com');
    });

    it('should include non-standard port for HTTP', () => {
        const request = new NextRequest('http://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com:8080;proto=http',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('http://example.com:8080');
    });

    it('should fallback to request.url when Forwarded header is not present', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should fallback when Forwarded header is missing proto', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;host=example.com',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should fallback when Forwarded header is missing host', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'forwarded': 'by=3.248.245.105;for=83.106.145.60;proto=https',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    describe('security validations', () => {
        it('should fallback when proto contains invalid protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=malicious.com;proto=javascript',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when proto contains data protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=malicious.com;proto=data',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when proto contains ftp protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=malicious.com;proto=ftp',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when host is empty', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when host contains invalid characters', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=evil@malicious.com;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when host contains consecutive dots', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=evil..com;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when host starts with a dot', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=.evil.com;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when host ends with a dot', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=evil.com.;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when port is invalid', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=example.com:abc;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when port is negative', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=example.com:-1;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when port is too large', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=example.com:999999;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should normalize protocol to lowercase', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=example.com;proto=HTTPS',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should normalize host to lowercase', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=Example.COM;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should accept IPv4 addresses as valid hosts', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host=192.168.1.1;proto=https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://192.168.1.1');
        });

        it('should handle malformed Forwarded header gracefully', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'invalid-format',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should handle Forwarded header with whitespace around values', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'forwarded': 'host = example.com ; proto = https',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });
    });
});

