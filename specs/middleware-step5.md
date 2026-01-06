# Middleware Step 5: ID Token Processing and User Details

## Overview

Improve how tokens are stored in the cookies.

## Files to Modify

* `web-nextjs/app/auth/callback/route.ts`
* `web-nextjs/app/auth/callback/route.test.ts`

## Test Cases

### Test 22: it should use real token expiration for cookie maxAge

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=EXPECTED_STATE`

**Mock Token Response:**

```typescript
{
  access_token: "JWT_TOKEN",
  refresh_token: "REFRESH_TOKEN",
  expires_in: 3600,  // 1 hour
  id_token: "ID_JWT"
}
```

**Expected Output:**

- Set-Cookie for `dphoto-access-token` with `maxAge=3600`
- Set-Cookie for `dphoto-refresh-token` with 30 days

**Implementation Notes:**

- Use `tokens.expires_in` for access token cookie maxAge
