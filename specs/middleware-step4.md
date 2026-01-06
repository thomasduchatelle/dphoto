# Middleware Step 4: OAuth Callback Error Handling

## Overview

Implement proper error handling for OAuth callback failures, including Cognito errors in query parameters and token exchange failures.

## Files to Modify

* `web-nextjs/app/auth/login/route.ts`
* `web-nextjs/app/auth/callback/route.ts`

## Test Cases

### Test 13: it should redirect to the error page with the same parameters when an error occurred

**Input:**

- Path: `/auth/callback?error=invalid_request&error_description=user.email%3A+Attribute+cannot+be+updated.
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies:
    - `dphoto-oauth-state=EXPECTED_STATE`
    - `dphoto-oauth-code-verifier=CODE_VERIFIER`

**Expected Output:**

- 307 code with location to `/auth/error?error=invalid_request&error_description=user.email%3A+Attribute+cannot+be+updated.`
- Cookies used for authentication are deleted

---

### Test 14: it should redirect to error page when the state mismatch

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=WRONG_STATE`
- Method: `GET`
- Cookies:
  - `dphoto-oauth-state=EXPECTED_STATE`
  - `dphoto-oauth-code-verifier=CODE_VERIFIER`

**Expected Output:**

- 307 code with location to `/auth/error?error=state-mismatch`
- Cookies used for authentication are deleted


### Test 15: it should redirect when authentication cookies are not present

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=SOME_STATE`
- Method: `GET`
- Cookies: none (cookies expired or deleted)

**Expected Output:**


**Expected Output:**

- 307 code with location to `/auth/error?error=missing-authentication-cookies`
- Cookies used for authentication are deleted

---

## New Page to Create

File: `web-nextjs/app/auth/error.tsx`

This page needs to show a user-friendly description of the errors: invalid_request, state-mismatch, missing-authentication-cookies.

It should also handle the most common error types that can be used in OIDC protocol by Cognito.

It needs to have a button to invite the user to retry: it's a link to `/auth/login`.
