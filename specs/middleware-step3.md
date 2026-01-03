# Middleware Step 3: Non-Navigation Request Handling

## Overview

Implement logic to distinguish between browser navigation requests and AJAX/fetch requests. Non-navigation requests should return 401 Unauthorized instead of
redirecting to login.

## Files to Modify

- `web/src/middleware/authentication.tsx` - add request type detection
- `web/src/middleware/authentication.test.ts` - add test cases

## Test Cases

### Test 9: it should return 401 for unauthenticated API request instead of redirecting

**Input:**

- Path: `/api/albums`
- Method: `GET`
- Headers:
    - `Accept: application/json`
    - `X-Requested-With: XMLHttpRequest` (optional indicator)
- Cookies: none (no access token)

**Expected Output:**

- Status: `401 Unauthorized`
- Body: JSON error response
  ```json
  {
    "error": "unauthorized",
    "message": "Authentication required"
  }
  ```
- No redirect, no cookies set

**Implementation Notes:**

- Detect non-navigation requests by checking `Accept` header
- API requests typically have `Accept: application/json`
- Navigation requests have `Accept: text/html`
- Do not redirect AJAX requests to login page

---

### Test 10: it should return 401 for fetch request without credentials

**Input:**

- Path: `/albums`
- Method: `GET`
- Headers: `Accept: application/json`
- Cookies: none

**Expected Output:**

- Status: `401 Unauthorized`
- No redirect

**Implementation Notes:**

- Even paths that look like pages should return 401 if Accept header indicates JSON
- Frontend can handle 401 by redirecting to login client-side

---

### Test 11: it should allow authenticated API request to proceed

**Input:**

- Path: `/api/albums`
- Method: `GET`
- Headers: `Accept: application/json`
- Cookies: `dphoto-access-token=VALID_TOKEN`

**Expected Output:**

- Calls `next()` to proceed with request
- `ctx.data.backendSession` populated with authenticated user
- No response modification

**Implementation Notes:**

- Authenticated requests proceed regardless of Accept header
- Only unauthenticated requests are affected by request type detection

---

### Test 12: it should redirect browser navigation even for API-like paths

**Input:**

- Path: `/api/albums`
- Method: `GET`
- Headers: `Accept: text/html,application/xhtml+xml` (browser navigation)
- Cookies: none

**Expected Output:**

- Status: `302 Found`
- Header `Location`: points to Cognito authorization URL
- OAuth cookies set

**Implementation Notes:**

- If a browser navigates directly to API path, redirect to login
- Accept header is the primary indicator of request type

