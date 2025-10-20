// Client-side cookie utilities
// The access token is stored in two cookies:
// 1. dphoto-access-token (HttpOnly) - for automatic inclusion in requests
// 2. dphoto-access-token-client (readable by JS) - for manual Authorization header inclusion

export function getAccessTokenFromCookie(): string | null {
  if (typeof document === 'undefined') {
    return null;
  }

  const cookies = document.cookie.split(';');
  
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=');
    // Read from the client-accessible cookie
    if (name === 'dphoto-access-token-client') {
      return decodeURIComponent(value);
    }
  }
  
  return null;
}

export function getTokenExpirationFromCookie(): number | null {
  const token = getAccessTokenFromCookie();
  if (!token) {
    return null;
  }

  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000; // Convert to milliseconds
  } catch {
    return null;
  }
}
