// @vitest-environment node

import {describe, expect, it} from 'vitest';
import {getLogoutUrl} from './oidc-config';

describe('getLogoutUrl', () => {
    it('should generate correct logout URL with client_id and logout_uri parameters', () => {
        const issuerUrl = 'https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_EXAMPLE';
        const clientId = 'test-client-id';
        const logoutUri = 'https://example.com/auth/logout-callback';

        const logoutUrl = getLogoutUrl(issuerUrl, clientId, logoutUri);

        expect(logoutUrl).toBe(
            'https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_EXAMPLE/logout?client_id=test-client-id&logout_uri=https%3A%2F%2Fexample.com%2Fauth%2Flogout-callback'
        );
    });

    it('should properly encode logout_uri parameter', () => {
        const issuerUrl = 'https://auth.example.com';
        const clientId = 'my-client';
        const logoutUri = 'https://app.example.com/logout?redirect=/home';

        const logoutUrl = getLogoutUrl(issuerUrl, clientId, logoutUri);

        expect(logoutUrl).toContain('logout_uri=https%3A%2F%2Fapp.example.com%2Flogout%3Fredirect%3D%2Fhome');
    });
});
