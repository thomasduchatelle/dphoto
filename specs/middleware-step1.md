# Middleware Step 1: Basic Authentication Flows

## Overview

Implement and test the core authentication flows for the middleware, including redirect to login, handling OAuth callback, and allowing authenticated requests.

## Files to Modify

- `web/src/middleware/authentication.tsx` - main middleware implementation
- `web/src/middleware/authentication.test.ts` - new test file to create

## Test Cases

### Test 1: it should redirect to authorization authority when requesting home page without access token

**Input:**

- Path: `/`
- Method: `GET`
- Headers: `Accept: text/html` (indicates browser navigation)
- Cookies: none (no access token)

**Expected Output:**

- Status: `302 Found`
- Header `Location`: points to Cognito authorization URL with PKCE parameters
- Set-Cookie headers:
    - `dphoto-oauth-state` with generated state value (maxAge: 300s, sameSite: lax)
    - `dphoto-oauth-code-verifier` with generated code verifier (maxAge: 300s, sameSite: lax)
- Authorization URL should contain:
    - `client_id`
    - `redirect_uri` pointing to `/auth/callback`
    - `scope=openid profile email`
    - `code_challenge` and `code_challenge_method=S256`
    - `state` parameter

**Implementation Notes:**

- Use a fake OIDC config provider to avoid network calls
- Mock `client.randomPKCECodeVerifier()` and `client.randomState()` for deterministic testing
- Mock `client.buildAuthorizationUrl()` to verify parameters

---

### Test 2: it should redirect to authorization authority when explicitly requesting /auth/login

**Input:**

- Path: `/auth/login`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: valid access token present

**Expected Output:**

- Status: `302 Found`
- Header `Location`: points to Cognito authorization URL
- Set-Cookie headers for state and code verifier
- Should redirect even when already authenticated (allows re-authentication)

**Implementation Notes:**

- This tests the explicit login path that should trigger re-authentication

---

### Test 3: it should handle OAuth callback with valid authorization code

**Input:**

- Path: `/auth/callback?code=AUTH_CODE_123&state=EXPECTED_STATE`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER_123`

**Expected Output:**

- Status: `302 Found`
- Header `Location`: `https://example.com/` (redirects to home)
- Set-Cookie headers:
    - `dphoto-access-token` with access token value (httpOnly, secure, sameSite: true)
    - `dphoto-refresh-token` with refresh token value (httpOnly, secure, sameSite: true)
    - `dphoto-oauth-state` cleared (maxAge: 0)
    - `dphoto-oauth-code-verifier` cleared (maxAge: 0)

**Implementation Notes:**

- Mock `client.authorizationCodeGrant()` to return successful token response:
  ```typescript
  {
    access_token: "ACCESS_TOKEN_VALUE",
    refresh_token: "REFRESH_TOKEN_VALUE",
    expires_in: 3600,
    id_token: "ID_TOKEN_JWT"
  }
  ```
- Mock JWT decoding to extract user info from id_token

---

### Test 4: it should allow authenticated request to proceed with backendSession

**Input:**

- Path: `/albums`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: `dphoto-access-token=VALID_ACCESS_TOKEN`

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.backendSession` should be set with:
  ```typescript
  {
    type: "authenticated",
    accessToken: {
      accessToken: "VALID_ACCESS_TOKEN",
      expiresAt: Date
    },
    refreshToken: "", // or actual refresh token if present
    authenticatedUser: {
      name: "John Doe",
      email: "john@example.com",
      picture: "https://example.com/avatar.jpg",
      isOwner: true
    }
  }
  ```

**Implementation Notes:**

- Mock JWT parsing to extract user claims (name, email, picture, role)
- Verify `next()` is called
- This requires parsing the access token (JWT) to get user details

