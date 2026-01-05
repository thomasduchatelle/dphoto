// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {NextRequest} from 'next/server';
import {getOriginalOrigin} from './request-utils';

describe('getOriginalOrigin', () => {
    it('should extract origin from X-Forwarded-Proto and X-Forwarded-Host headers', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https',
                'x-forwarded-host': 'example.com',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should handle X-Forwarded-Host with comma-separated values (taking first)', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https',
                'x-forwarded-host': 'example.com, proxy1.internal, proxy2.internal',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should handle X-Forwarded-Proto with comma-separated values (taking first)', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https, http',
                'x-forwarded-host': 'example.com',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should include non-standard port for HTTPS', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https',
                'x-forwarded-host': 'example.com',
                'x-forwarded-port': '8443',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com:8443');
    });

    it('should include non-standard port for HTTP', () => {
        const request = new NextRequest('http://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'http',
                'x-forwarded-host': 'example.com',
                'x-forwarded-port': '8080',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('http://example.com:8080');
    });

    it('should not include standard HTTPS port 443', () => {
        const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https',
                'x-forwarded-host': 'example.com',
                'x-forwarded-port': '443',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should not include standard HTTP port 80', () => {
        const request = new NextRequest('http://internal-api-gateway.amazonaws.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'http',
                'x-forwarded-host': 'example.com',
                'x-forwarded-port': '80',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('http://example.com');
    });

    it('should fallback to request.url when X-Forwarded headers are not present', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should fallback to request.url when only X-Forwarded-Proto is present', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-proto': 'https',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    it('should fallback to request.url when only X-Forwarded-Host is present', () => {
        const request = new NextRequest('https://example.com/path', {
            method: 'GET',
            headers: {
                'x-forwarded-host': 'other.com',
            },
        });

        const origin = getOriginalOrigin(request);

        expect(origin).toBe('https://example.com');
    });

    describe('security validations', () => {
        it('should fallback when X-Forwarded-Proto contains invalid protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'javascript',
                    'x-forwarded-host': 'malicious.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when X-Forwarded-Proto contains data protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'data',
                    'x-forwarded-host': 'malicious.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when X-Forwarded-Proto contains ftp protocol', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'ftp',
                    'x-forwarded-host': 'malicious.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when X-Forwarded-Host is empty', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': '',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should fallback when X-Forwarded-Host contains invalid characters', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'evil@malicious.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should ignore invalid port and omit it from URL', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': 'abc',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should ignore negative port number', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': '-1',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should ignore port number that is too large', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': '999999',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should ignore empty port value', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': '',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should handle comma-separated X-Forwarded-Port values (taking first)', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': '8443, 443',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com:8443');
        });

        it('should normalize protocol to lowercase', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'HTTPS',
                    'x-forwarded-host': 'example.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should normalize host to lowercase', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'Example.COM',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should accept IPv4 addresses as valid hosts', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': '192.168.1.1',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://192.168.1.1');
        });

        it('should reject hosts with consecutive dots', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'evil..com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should reject hosts ending with a dot', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'evil.com.',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should reject hosts starting with a dot', () => {
            const request = new NextRequest('https://example.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': '.evil.com',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com');
        });

        it('should accept port with leading zeros', () => {
            const request = new NextRequest('https://internal-api-gateway.amazonaws.com/path', {
                method: 'GET',
                headers: {
                    'x-forwarded-proto': 'https',
                    'x-forwarded-host': 'example.com',
                    'x-forwarded-port': '08443',
                },
            });

            const origin = getOriginalOrigin(request);

            expect(origin).toBe('https://example.com:8443');
        });
    });
});

