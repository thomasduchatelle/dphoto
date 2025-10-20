import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getClientAccessToken, setClientAccessToken, clearClientAccessToken } from './token-context';
import * as clientCookieUtils from './client-cookie-utils';

describe('token-context', () => {
  beforeEach(() => {
    clearClientAccessToken();
    vi.clearAllMocks();
  });

  describe('setClientAccessToken and getClientAccessToken', () => {
    it('should store and retrieve access token', () => {
      const token = {
        accessToken: 'test-token',
        expiresAt: Date.now() + 3600000, // 1 hour from now
      };

      setClientAccessToken(token);
      const retrieved = getClientAccessToken();

      expect(retrieved).toBe('test-token');
    });

    it('should return null for expired token', () => {
      const token = {
        accessToken: 'test-token',
        expiresAt: Date.now() - 1000, // Expired 1 second ago
      };

      setClientAccessToken(token);
      const retrieved = getClientAccessToken();

      expect(retrieved).toBeNull();
    });
  });

  describe('clearClientAccessToken', () => {
    it('should clear stored token', () => {
      const token = {
        accessToken: 'test-token',
        expiresAt: Date.now() + 3600000,
      };

      setClientAccessToken(token);
      clearClientAccessToken();
      const retrieved = getClientAccessToken();

      expect(retrieved).toBeNull();
    });
  });

  describe('cookie integration', () => {
    it('should load token from cookie if not in memory', () => {
      const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZXhwIjoxOTk5OTk5OTk5fQ.x';
      
      vi.spyOn(clientCookieUtils, 'getAccessTokenFromCookie').mockReturnValue(mockToken);
      vi.spyOn(clientCookieUtils, 'getTokenExpirationFromCookie').mockReturnValue(Date.now() + 3600000);

      const retrieved = getClientAccessToken();

      expect(retrieved).toBe(mockToken);
    });

    it('should return null if no token in cookie', () => {
      vi.spyOn(clientCookieUtils, 'getAccessTokenFromCookie').mockReturnValue(null);
      vi.spyOn(clientCookieUtils, 'getTokenExpirationFromCookie').mockReturnValue(null);

      const retrieved = getClientAccessToken();

      expect(retrieved).toBeNull();
    });
  });
});
