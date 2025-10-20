// Client-side cookie utilities
// Note: HttpOnly cookies cannot be accessed by JavaScript, but they are automatically
// sent with requests. This utility helps read non-HttpOnly tokens if needed.

export function getAccessTokenFromCookie(): string | null {
  if (typeof document === 'undefined') {
    return null;
  }

  const cookies = document.cookie.split(';');
  
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=');
    if (name === 'dphoto-access-token') {
      return decodeURIComponent(value);
    }
  }
  
  return null;
}

export function getRefreshTokenFromCookie(): string | null {
  if (typeof document === 'undefined') {
    return null;
  }

  const cookies = document.cookie.split(';');
  
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=');
    if (name === 'dphoto-refresh-token') {
      return decodeURIComponent(value);
    }
  }
  
  return null;
}

// Note: Since we're using HttpOnly cookies, these functions won't be able to read
// the auth tokens. The tokens will automatically be sent with requests via the
// Cookie header. The axios interceptor doesn't need to manually add them.
// However, we keep these functions for potential non-HttpOnly use cases.
