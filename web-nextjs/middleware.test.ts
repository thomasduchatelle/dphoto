// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {NextRequest} from 'next/server';
import {middleware, skipProxyForPageMatching} from './middleware';
import {ACCESS_TOKEN_COOKIE} from './lib/security/constants';

describe('authentication middleware', () => {
    it('should redirect to login page when requesting home page without access token', async () => {
        const request = new NextRequest('https://example.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await middleware(request);

        expect(response.status).toBe(307);
        expect(response.headers.get('Location')).toBe('https://example.com/auth/login');
    });

    it('should allow authenticated request to proceed with backendSession', async () => {
        const accessToken = 'VALID_ACCESS_TOKEN';

        const request = new NextRequest('https://example.com/albums', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
                Cookie: `${ACCESS_TOKEN_COOKIE}=${accessToken}`,
            },
        });

        const response = await middleware(request);

        expect(response.status).toBe(200);
    });
});

describe('skipProxyForPageMatching regex', () => {
    it('should NOT match static image paths', () => {
        expect(skipProxyForPageMatching.test('images/photo.png')).toBe(false);
        expect(skipProxyForPageMatching.test('images/photo.jpg')).toBe(false);
        expect(skipProxyForPageMatching.test('images/photo.gif')).toBe(false);
        expect(skipProxyForPageMatching.test('static/image.svg')).toBe(false);
        expect(skipProxyForPageMatching.test('photo.png')).toBe(false);
    });

    it('should NOT match /auth/* paths', () => {
        expect(skipProxyForPageMatching.test('auth')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/login')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/callback')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/logout')).toBe(false);
        expect(skipProxyForPageMatching.test('auth/refresh')).toBe(false);
    });

    it('should NOT match /api/* paths', () => {
        expect(skipProxyForPageMatching.test('api/albums')).toBe(false);
        expect(skipProxyForPageMatching.test('api/photos/123')).toBe(false);
    });

    it('should NOT match favicon.ico', () => {
        expect(skipProxyForPageMatching.test('favicon.ico')).toBe(false);
    });

    it('should NOT match Next.js internal paths', () => {
        expect(skipProxyForPageMatching.test('_next/static/chunks/main.js')).toBe(false);
        expect(skipProxyForPageMatching.test('_next/image')).toBe(false);
    });

    it('should NOT match JavaScript files', () => {
        expect(skipProxyForPageMatching.test('scripts/main.js')).toBe(false);
        expect(skipProxyForPageMatching.test('bundle.js')).toBe(false);
    });

    it('should match application pages that require authentication', () => {
        expect(skipProxyForPageMatching.test('')).toBe(true);
        expect(skipProxyForPageMatching.test('albums')).toBe(true);
        expect(skipProxyForPageMatching.test('albums/123')).toBe(true);
        expect(skipProxyForPageMatching.test('photos')).toBe(true);
        expect(skipProxyForPageMatching.test('settings')).toBe(true);
    });
});
