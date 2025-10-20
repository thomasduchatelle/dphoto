import { describe, it, expect } from 'vitest';
import { extractTokensFromCookies, getTokenExpiration, isTokenValid } from './server-token-utils';

describe('server-token-utils', () => {
  describe('extractTokensFromCookies', () => {
    it('should extract tokens from cookie header', () => {
      const cookieHeader = 'dphoto-access-token=access-value; dphoto-refresh-token=refresh-value';
      const tokens = extractTokensFromCookies(cookieHeader);

      expect(tokens.accessToken).toBe('access-value');
      expect(tokens.refreshToken).toBe('refresh-value');
    });

    it('should return empty object for undefined header', () => {
      const tokens = extractTokensFromCookies(undefined);

      expect(tokens).toEqual({});
    });
  });

  describe('getTokenExpiration', () => {
    it('should decode and return expiration time', () => {
      // Create a mock JWT with exp claim
      const payload = { exp: 1234567890 };
      const encodedPayload = Buffer.from(JSON.stringify(payload)).toString('base64');
      const mockToken = `header.${encodedPayload}.signature`;

      const expiration = getTokenExpiration(mockToken);

      expect(expiration).toBe(1234567890 * 1000); // Converted to milliseconds
    });

    it('should return null for invalid token', () => {
      const invalidToken = 'not-a-valid-jwt';
      const expiration = getTokenExpiration(invalidToken);

      expect(expiration).toBeNull();
    });
  });

  describe('isTokenValid', () => {
    it('should return true for non-expired token', () => {
      const futureTime = Math.floor(Date.now() / 1000) + 3600; // 1 hour from now
      const payload = { exp: futureTime };
      const encodedPayload = Buffer.from(JSON.stringify(payload)).toString('base64');
      const mockToken = `header.${encodedPayload}.signature`;

      const valid = isTokenValid(mockToken);

      expect(valid).toBe(true);
    });

    it('should return false for expired token', () => {
      const pastTime = Math.floor(Date.now() / 1000) - 3600; // 1 hour ago
      const payload = { exp: pastTime };
      const encodedPayload = Buffer.from(JSON.stringify(payload)).toString('base64');
      const mockToken = `header.${encodedPayload}.signature`;

      const valid = isTokenValid(mockToken);

      expect(valid).toBe(false);
    });

    it('should return false for invalid token', () => {
      const invalid = isTokenValid('invalid-token');

      expect(invalid).toBe(false);
    });
  });
});
