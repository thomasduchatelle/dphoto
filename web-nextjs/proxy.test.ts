// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {NextRequest} from 'next/server';
import {proxy} from './proxy';
import {ACCESS_TOKEN_COOKIE} from './lib/security/constants';

describe('authentication middleware', () => {
    it('should redirect to login page when requesting home page without access token', async () => {
        const request = new NextRequest('https://example.com/', {
            method: 'GET',
            headers: {
                Accept: 'text/html',
            },
        });

        const response = await proxy(request);

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

        const response = await proxy(request);

        expect(response.status).toBe(200);
    });
});
