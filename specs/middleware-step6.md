# Middleware Step 6: Logout Flow

## Overview

Implement the logout functionality that signs the user out from Cognito, clears all authentication cookies, and redirects appropriately.

## Files to Modify

- `web-nextjs/middleware-authentication.tsx` - add logout path handler
- `web/src/core/security/consts.ts` - may need logout-related constants
- `web/src/pages/logout.tsx` - create logout confirmation page (optional)
- `web-nextjs/middleware-authentication.test.ts` - add test cases

## Test Cases

### Test 23: it should handle logout request and clear all cookies

**Input:**

- Path: `/auth/logout`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies:
    - `dphoto-access-token=VALID_TOKEN`
    - `dphoto-refresh-token=REFRESH_TOKEN`

**Expected Output:**

- Status: `302 Found`
- Header `Location`: Cognito logout endpoint URL with:
    - `client_id` parameter
    - `logout_uri` parameter pointing to `/` or a post-logout page
- Set-Cookie headers clearing all auth cookies (maxAge: 0):
    - `dphoto-access-token`
    - `dphoto-refresh-token`
    - `dphoto-oauth-state`
    - `dphoto-oauth-code-verifier`
    - `dphoto-redirect-after-login`

**Implementation Notes:**

- Cognito logout URL: `https://{domain}.auth.{region}.amazoncognito.com/logout`
- Query parameters: `?client_id={clientId}&logout_uri={encodedUri}`
- Clear all cookies to prevent any stale auth state
- Redirect to Cognito to invalidate session there

---

### Test 24: it should handle logout when already unauthenticated

**Input:**

- Path: `/auth/logout`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: none

**Expected Output:**

- Status: `302 Found`
- Header `Location`: redirects to home `/` (no need to contact Cognito)
- Set-Cookie headers clearing any potential stale cookies

**Implementation Notes:**

- If no access token present, skip Cognito logout
- Still clear cookies as defensive measure
- Redirect directly to home page

---

### Test 25: it should allow unauthenticated access to post-logout page

**Input:**

- Path: `/auth/logout-complete` (or similar post-logout landing page)
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: none

**Expected Output:**

- Calls `next()` to render the page
- No redirect to login
- No backendSession set (or set as anonymous)

**Implementation Notes:**

- Whitelist `/auth/logout-complete` from authentication requirement
- Can show "You have been logged out" message
- Provide link to login again

---

### Test 26: it should whitelist static assets and public paths from authentication

**Input:**

- Path: `/assets/logo.png` (or `/public/*`, `/_next/*`, etc.)
- Method: `GET`
- Headers: `Accept: image/png`
- Cookies: none

**Expected Output:**

- Calls `next()` without authentication
- No redirect

**Implementation Notes:**

- Define public paths that don't require authentication:
    - `/assets/*`
    - `/public/*`
    - `/favicon.ico`
    - `/auth/logout-complete`
    - `/auth/error`
    - Any other public resources
- Check path prefix before enforcing authentication

---

### Test 27: it should handle logout with custom redirect parameter

**Input:**

- Path: `/auth/logout?redirect=/welcome`
- Method: `GET`
- Cookies: (authenticated)

**Expected Output:**

- Status: `302 Found`
- Header `Location`: Cognito logout with `logout_uri` pointing to `/welcome`
- Cookies cleared

**Implementation Notes:**

- Optional: allow specifying post-logout destination
- Validate redirect parameter (must start with `/`, no external URLs)
- Default to `/` if not specified or invalid

---

## Cognito Logout Endpoint Details

The Cognito logout endpoint format:

```
https://{cognito-domain}.auth.{region}.amazoncognito.com/logout?
  client_id={client_id}&
  logout_uri={url-encoded-redirect-uri}
```

Configuration needed:

- Cognito domain: from `getEnv("COGNITO_DOMAIN")` or derive from issuer
- The `logout_uri` must be in Cognito's allowed logout URLs list

Alternative: If full Cognito logout is not required (just clear cookies):

- Simplified logout can just clear cookies and redirect
- No need to contact Cognito if session doesn't need server-side invalidation
- Decision depends on security requirements

## Path Whitelist

Paths that should bypass authentication:

```typescript
const PUBLIC_PATHS = [
  '/auth/callback',
  '/auth/error', 
  '/auth/logout-complete',
  '/assets',
  '/public',
  '/favicon.ico',
  '/__waku',  // Waku internal paths
];

function isPublicPath(path: string): boolean {
  return PUBLIC_PATHS.some(publicPath => 
    path === publicPath || path.startsWith(publicPath + '/')
  );
}
```

