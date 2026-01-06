# Middleware Step 2: Redirect Path Preservation

## Overview

Implement the ability to remember the originally requested URL before redirecting to login, and restore it after successful authentication. 

## Files to Modify

- `web-nextjs/proxy.tsx` - add redirect path cookie logic
- `web-nextjs/proxy.test.ts` - add test cases

```typescript
export const REDIRECT_AFTER_LOGIN_COOKIE = 'dphoto-redirect-after-login';
```

## Test Cases

### Test 5: it should store the original URL in a cookie before redirecting to login

**Input:**

- Path: `/albums/2024-summer`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: none (no access token)

**Expected Output:**

- Status: `302 Found`
- Header `Location`: points to Cognito authorization URL
- Set-Cookie headers include:
    - `dphoto-redirect-after-login=/albums/2024-summer` (maxAge: 300s, sameSite: lax)
    - Other OAuth cookies (state, code_verifier)

**Implementation Notes:**

- Only store the path, not the full URL (security: avoid open redirect)
- Validate the path starts with `/` to prevent open redirects

---

### Test 6: it should redirect to stored URL after successful OAuth callback

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=EXPECTED_STATE`
- Method: `GET`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER`
    - `dphoto-redirect-after-login=/albums/2024-summer`

**Expected Output:**

- Status: `302 Found`
- Header `Location`: `https://example.com/albums/2024-summer` (restored path)
- Set-Cookie headers:
    - Access and refresh token cookies set
    - OAuth cookies cleared
    - `dphoto-redirect-after-login` cleared (maxAge: 0)

**Implementation Notes:**

- After successful token exchange, read redirect cookie
- Default to `/` if cookie is missing or invalid
- Clear the redirect cookie after use

---

### Test 7: it should redirect to home if redirect cookie contains external URL

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=EXPECTED_STATE`
- Method: `GET`
- Cookies:
    - OAuth cookies (valid)
    - `dphoto-redirect-after-login=https://evil.com/phishing`

**Expected Output:**

- Status: `302 Found`
- Header `Location`: `https://example.com/` (defaults to home, ignores external URL)

**Implementation Notes:**

- Security measure: validate redirect path starts with `/` and doesn't contain `://`
- This prevents open redirect vulnerabilities

---

### Test 8: it should not store redirect path for /auth/login explicit requests

**Input:**

- Path: `/auth/login`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: none

**Expected Output:**

- Status: `302 Found`
- Header `Location`: points to Cognito authorization URL
- Set-Cookie headers should NOT include `dphoto-redirect-after-login`
- Only OAuth cookies (state, code_verifier) are set

**Implementation Notes:**

- Explicit login requests should redirect to home after completion
- Don't store `/auth/login` as the redirect target

