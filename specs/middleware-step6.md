# Middleware Step 6: Logout Flow

## Overview

Implement the logout functionality that signs the user out from Cognito, clears all authentication cookies, and redirects appropriately.

## Test Cases

### Test 23: it should clear all cookies and redirect to OIDC logout URL

**Input:**

- Path: `/auth/logout`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies:
    - `dphoto-access-token=VALID_TOKEN`
    - `dphoto-refresh-token=REFRESH_TOKEN`

**Expected Output:**

- Set-Cookie headers clearing all auth cookies (maxAge: 0):
    - `dphoto-access-token`
    - `dphoto-refresh-token`
    - `dphoto-oauth-state`
    - `dphoto-oauth-code-verifier`
    - `dphoto-redirect-after-login`
- 307 code with `Location=${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`
  - logoutUri is `/auth/logout-success`

---

### Test 24: it should redirect to cognito logout URL even if there is no cookies

**Input:**

- Path: `/auth/logout`
- Method: `GET`
- Headers: `Accept: text/html`
- Cookies: none

**Expected Output:**

- 307 code with `Location=${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`
  - logoutUri is `/auth/logout-success`
- No need to have a specific handler: Set-Cookie headers clearing all auth cookies (maxAge: 0):
  - `dphoto-access-token`
  - `dphoto-refresh-token`
  - `dphoto-oauth-state`
  - `dphoto-oauth-code-verifier`
  - `dphoto-redirect-after-login`

---

## Create the logout success page

Page path: `/auth/logout-success`

A logout success page must be created. Use a simple layout to notify the success of the operation.

A link to re-authenticate (link to `/auth/login`) should appear.
