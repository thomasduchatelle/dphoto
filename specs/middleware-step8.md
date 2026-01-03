# Middleware Step 8: Edge Cases and Security Hardening

## Overview

Handle various edge cases, security concerns, and unexpected scenarios to make the middleware production-ready.

## Files to Modify

- `web-nextjs/middleware-authentication.tsx` - add edge case handling
- `web-nextjs/middleware-authentication.test.ts` - add test cases

## Test Cases

### Test 31: it should handle malformed JWT in access token cookie

**Input:**

- Path: `/albums`
- Method: `GET`
- Cookies: `dphoto-access-token=not.a.valid.jwt`

**Expected Output:**

- Status: `302 Found` (redirect to login for browser navigation)
- OR Status: `401 Unauthorized` (for API requests)
- Clear the invalid access token cookie

**Implementation Notes:**

- Wrap JWT parsing in try-catch
- Treat malformed token same as missing token
- Clear the bad cookie to prevent infinite loops

---

### Test 32: it should handle expired access token in cookie

**Input:**

- Path: `/albums`
- Method: `GET`
- Cookies:
    - `dphoto-access-token=EXPIRED_JWT` (exp claim is in the past)
    - `dphoto-refresh-token=VALID_REFRESH`

**Expected Output:**

- Option A (Simple): Redirect to login, let refresh happen client-side
- Option B (Advanced): Use refresh token to get new access token

For production-ready Step 8, implement Option A (simpler):

- Status: `302 Found`
- Redirect to login
- Clear expired access token cookie

**Implementation Notes:**

- Check JWT `exp` claim against current time
- If expired and refresh token present, could refresh (future enhancement)
- For now, treat as unauthenticated
- Client-side can handle token refresh proactively

---

### Test 33: it should handle missing required environment variables gracefully

**Input:**

- Path: `/` (any request)
- Environment: `COGNITO_ISSUER` is undefined

**Expected Output:**

- Status: `500 Internal Server Error`
- Error message indicating configuration issue
- Logged error with details

**Implementation Notes:**

- Validate env vars at startup or first use
- Provide clear error messages for misconfiguration
- Don't expose sensitive details in error response to client

---

### Test 34: it should prevent open redirect via malicious redirect cookie

**Input:**

- Path: `/auth/callback?code=CODE&state=STATE`
- Cookies:
    - Valid OAuth cookies
    - `dphoto-redirect-after-login=//evil.com/phishing`

**Expected Output:**

- Status: `302 Found`
- Header `Location`: `https://example.com/` (ignores malicious redirect)

**Implementation Notes:**

- Validate redirect path thoroughly:
  ```typescript
  function isValidRedirectPath(path: string): boolean {
    if (!path || typeof path !== 'string') return false;
    if (!path.startsWith('/')) return false;
    if (path.startsWith('//')) return false;  // Protocol-relative URL
    if (path.includes('://')) return false;    // Absolute URL
    return true;
  }
  ```

---

### Test 35: it should handle concurrent requests during authentication

**Input:**

- Two simultaneous requests:
    1. Path: `/albums` (no auth)
    2. Path: `/photos` (no auth)
- Both arrive while OIDC config is being fetched

**Expected Output:**

- Both requests should complete successfully
- OIDC config fetched only once (not twice)
- Both redirected to login with their own OAuth state

**Implementation Notes:**

- If using async cache pattern, ensure thread-safety
- Each request gets unique state/verifier
- Config can be shared once fetched

---

### Test 36: it should handle very long cookie values safely

**Input:**

- Path: `/auth/callback`
- Cookies with extremely long values (near 4KB browser limit)

**Expected Output:**

- Should handle gracefully without crashing
- May truncate or reject if over limits

**Implementation Notes:**

- Browsers limit cookies to ~4KB
- JWT tokens can be large (1-2KB)
- Multiple cookies together could exceed limits
- Monitor cookie sizes in production

---

### Test 37: it should handle missing state parameter in callback

**Input:**

- Path: `/auth/callback?code=AUTH_CODE`
- No `state` query parameter
- Cookies: (OAuth cookies present)

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.oidcError` set:
  ```typescript
  {
    error: "invalid_request",
    errorDescription: "Missing required state parameter"
  }
  ```

**Implementation Notes:**

- Validate required OAuth parameters before token exchange
- State parameter is required for CSRF protection

---

### Test 38: it should handle callback with both error and code parameters

**Input:**

- Path: `/auth/callback?code=CODE&error=invalid_request&state=STATE`
- Malformed or malicious request

**Expected Output:**

- Prefer error parameter (indicates OAuth provider rejection)
- `ctx.data.oidcError` set with the error details

**Implementation Notes:**

- If `error` parameter exists, process it first
- Don't attempt token exchange when error is present

---

### Test 39: it should handle OPTIONS requests for CORS preflight

**Input:**

- Path: `/api/albums`
- Method: `OPTIONS`
- Headers: CORS headers

**Expected Output:**

- Calls `next()` without authentication
- No redirect, no 401

**Implementation Notes:**

- OPTIONS requests should bypass authentication
- CORS preflight must complete before authentication
- Add to public request types:
  ```typescript
  if (ctx.req.method === 'OPTIONS') {
    return next();
  }
  ```

---

### Test 40: it should rate-limit authentication attempts

**Input:**

- Multiple rapid requests to `/auth/login` from same IP
- Example: 10 requests in 10 seconds

**Expected Output:**

- First 5 requests: normal redirect to Cognito
- Subsequent requests: `429 Too Many Requests`
- Rate limit resets after cooldown period

**Implementation Notes:**

- This is advanced and may need separate rate-limiting middleware
- For production: use Cognito's built-in rate limiting
- For additional protection: implement IP-based rate limit
- Can use in-memory store for simple case, or Redis for distributed
- Consider this optional for initial production release

---

## Security Checklist

Before production deployment, verify:

- ✅ All cookies use `httpOnly: true` (prevent XSS)
- ✅ All cookies use `secure: true` (HTTPS only)
- ✅ OAuth state/verifier cookies use `sameSite: 'lax'` for OIDC flow
- ✅ Auth cookies use `sameSite: true` (strict) to prevent CSRF
- ✅ Redirect validation prevents open redirect attacks
- ✅ JWT parsing handles malformed tokens gracefully
- ✅ Expired tokens are detected and handled
- ✅ Error messages don't leak sensitive information
- ✅ Environment variables are validated
- ✅ OIDC config failures don't crash the server
- ✅ Public paths are properly whitelisted
- ✅ CORS preflight requests bypass authentication
- ✅ Rate limiting prevents abuse (or rely on Cognito)

## Logging and Monitoring

Add appropriate logging:

- ✅ Log authentication attempts (success/failure)
- ✅ Log OAuth callback errors
- ✅ Log OIDC config fetch failures
- ✅ Don't log sensitive data (tokens, passwords)
- ✅ Include request IDs for tracing

Example log points:

```typescript
console.log('Auth: redirecting to login', {path, hasToken: !!cookies.accessToken});
console.log('Auth: OAuth callback success', {email: user.email});
console.error('Auth: OAuth callback failed', {error: err.message, code: err.code});
console.warn('Auth: invalid JWT in cookie', {path});
```

