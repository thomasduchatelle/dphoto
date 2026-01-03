# Middleware Step 4: OAuth Callback Error Handling

## Overview

Implement proper error handling for OAuth callback failures, including Cognito errors in query parameters and token exchange failures.

## Files to Modify

- `web-nextjs/middleware-authentication.tsx` - add error handling in callback path
- `web/src/pages/auth/error.tsx` - new Waku page to create for displaying errors
- `web-nextjs/middleware-authentication.test.ts` - add test cases

## Test Cases

### Test 13: it should handle OAuth error in callback query parameters

**Input:**

- Path: `/auth/callback?error=invalid_request&error_description=user.email%3A+Attribute+cannot+be+updated.&state=EXPECTED_STATE`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER`

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.oidcError` should be set with:
  ```typescript
  {
    error: "invalid_request",
    errorDescription: "user.email: Attribute cannot be updated."
  }
  ```
- No cookies cleared yet (page will handle display)

**Implementation Notes:**

- Check for `error` query parameter before attempting token exchange
- URL decode `error_description`
- Set error in context data for the error page to display
- Call `next()` to allow Waku to render `/auth/error` page

---

### Test 14: it should handle token exchange failure with invalid_grant

**Input:**

- Path: `/auth/callback?code=INVALID_CODE&state=EXPECTED_STATE`
- Method: `GET`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER`

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.oidcError` should be set with:
  ```typescript
  {
    error: "invalid_grant",
    errorDescription: "The provided authorization grant is invalid, expired, or revoked."
  }
  ```

**Implementation Notes:**

- Mock `client.authorizationCodeGrant()` to throw error:
  ```typescript
  {
    cause: { error: 'invalid_grant' },
    code: 'OAUTH_RESPONSE_BODY_ERROR',
    error: 'invalid_grant',
    status: 400
  }
  ```
- Catch the exception and set `ctx.data.oidcError`
- Call `next()` to render error page

---

### Test 15: it should handle state mismatch in OAuth callback

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=WRONG_STATE`
- Method: `GET`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER`

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.oidcError` should contain:
  ```typescript
  {
    error: "state_mismatch",
    errorDescription: "OAuth state validation failed. Please try again."
  }
  ```

**Implementation Notes:**

- `client.authorizationCodeGrant()` should throw when state doesn't match
- Catch and convert to user-friendly error

---

### Test 16: it should handle missing OAuth cookies in callback

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=SOME_STATE`
- Method: `GET`
- Cookies: none (cookies expired or deleted)

**Expected Output:**

- Status: not modified (calls `next()`)
- `ctx.data.oidcError` should contain:
  ```typescript
  {
    error: "missing_credentials",
    errorDescription: "Authentication session expired. Please try logging in again."
  }
  ```

**Implementation Notes:**

- Check if state or codeVerifier cookies are missing
- Set appropriate error before attempting token exchange

---

### Test 17: it should clear OAuth cookies after error in callback

**Input:**

- Path: `/auth/callback?error=access_denied&error_description=User+cancelled`
- Method: `GET`
- Cookies: (OAuth cookies present)

**Expected Output:**

- `ctx.data.oidcError` set
- Should still clear OAuth cookies (set maxAge: 0):
    - `dphoto-oauth-state`
    - `dphoto-oauth-code-verifier`
    - `dphoto-redirect-after-login`

**Implementation Notes:**

- Always clean up OAuth flow cookies, even on error
- Prevents stale state from interfering with retry

---

## New Page to Create

File: `web/src/pages/auth/error.tsx`

This page should:

1. Check `ctx.data.oidcError` from the middleware
2. Display user-friendly error message
3. Provide "Try Again" button that redirects to `/auth/login`
4. Handle case where error data is missing (shouldn't happen, but defensive)

Example structure:

```typescript
export default async function AuthErrorPage({ ctx }) {
  const error = ctx.data.oidcError || {
    error: 'unknown_error',
    errorDescription: 'An unexpected error occurred during authentication.'
  };
  
  return (
    <div>
      <h1>Authentication Error</h1>
      <p>{error.errorDescription}</p>
      <a href="/auth/login">Try Again</a>
    </div>
  );
}
```

