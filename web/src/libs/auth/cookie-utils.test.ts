import { describe, it, expect } from 'vitest';
import { setTokenCookies, clearTokenCookies, parseTokenCookies } from './cookie-utils';

describe('cookie-utils', () => {
  describe('setTokenCookies', () => {
    it('should create cookie strings with correct attributes', () => {
      const cookies = setTokenCookies('access-token-value', 'refresh-token-value');

      expect(cookies).toHaveLength(2);
      expect(cookies[0]).toContain('dphoto-access-token=access-token-value');
      expect(cookies[0]).toContain('HttpOnly');
      expect(cookies[0]).toContain('SameSite=Strict');
      expect(cookies[0]).toContain('Path=/');
      expect(cookies[0]).toContain('Max-Age=3600'); // 1 hour

      expect(cookies[1]).toContain('dphoto-refresh-token=refresh-token-value');
      expect(cookies[1]).toContain('HttpOnly');
      expect(cookies[1]).toContain('SameSite=Strict');
      expect(cookies[1]).toContain('Path=/');
      expect(cookies[1]).toContain('Max-Age=2592000'); // 30 days
    });

    it('should include Secure flag in production', () => {
      const originalEnv = process.env.NODE_ENV;
      process.env.NODE_ENV = 'production';

      const cookies = setTokenCookies('access-token', 'refresh-token');

      expect(cookies[0]).toContain('Secure');
      expect(cookies[1]).toContain('Secure');

      process.env.NODE_ENV = originalEnv;
    });
  });

  describe('clearTokenCookies', () => {
    it('should create cookie strings with Max-Age=0', () => {
      const cookies = clearTokenCookies();

      expect(cookies).toHaveLength(2);
      expect(cookies[0]).toContain('dphoto-access-token=');
      expect(cookies[0]).toContain('Max-Age=0');
      expect(cookies[1]).toContain('dphoto-refresh-token=');
      expect(cookies[1]).toContain('Max-Age=0');
    });
  });

  describe('parseTokenCookies', () => {
    it('should parse cookies from header string', () => {
      const cookieHeader = 'dphoto-access-token=access-value; dphoto-refresh-token=refresh-value';
      const tokens = parseTokenCookies(cookieHeader);

      expect(tokens.accessToken).toBe('access-value');
      expect(tokens.refreshToken).toBe('refresh-value');
    });

    it('should handle missing cookies', () => {
      const cookieHeader = 'other-cookie=value';
      const tokens = parseTokenCookies(cookieHeader);

      expect(tokens.accessToken).toBeUndefined();
      expect(tokens.refreshToken).toBeUndefined();
    });

    it('should return empty object for undefined header', () => {
      const tokens = parseTokenCookies(undefined);

      expect(tokens).toEqual({});
    });
  });
});
