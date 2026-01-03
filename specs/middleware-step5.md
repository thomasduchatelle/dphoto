# Middleware Step 5: ID Token Processing and User Details

## Overview

Extract user information (name, email, picture, isOwner) from the ID token returned by Cognito and populate the backendSession correctly.

## Files to Modify

- `web-nextjs/middleware-authentication.tsx` - parse ID token and extract user claims
- `web/src/core/security/jwt-utils.ts` - add JWT decoding utilities if not present
- `web-nextjs/middleware-authentication.test.ts` - add test cases

## Test Cases

### Test 18: it should extract user details from ID token after successful callback

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=EXPECTED_STATE`
- Method: `GET`
- Cookies: (valid OAuth cookies)

**Mock ID Token Claims:**

```json
{
  "sub": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "picture": "https://lh3.googleusercontent.com/a/avatar.jpg",
  "cognito:groups": ["owners"],
  "iss": "https://cognito-idp.region.amazonaws.com/poolId",
  "aud": "client_id",
  "exp": 1735200000,
  "iat": 1735196400
}
```

**Expected Output:**

- Status: `302 Found`
- Cookies set with access and refresh tokens
- (Internal) User details should be captured for DynamoDB storage (see implementation notes)

**Implementation Notes:**

- Parse ID token JWT without verification (middleware runs after Cognito validates it)
- Extract claims: `email`, `name`, `picture`, `cognito:groups`
- Determine `isOwner` from `cognito:groups` array (contains "owners" or "admins")
- Store identity in DynamoDB via backend API call (see TODO in code)
- This aligns with Go backend: `pkg/acl/aclcore/authenticate_sso.go`

---

### Test 19: it should populate backendSession with user details from access token

**Input:**

- Path: `/albums`
- Method: `GET`
- Cookies: `dphoto-access-token=ENCODED_JWT`

**Mock Access Token Claims:**

```json
{
  "sub": "123e4567-e89b-12d3-a456-426614174000",
  "email": "jane.visitor@example.com",
  "name": "Jane Visitor",
  "picture": "https://example.com/avatar2.jpg",
  "cognito:groups": ["visitors"],
  "exp": 1735200000
}
```

**Expected Output:**

- Calls `next()` to proceed
- `ctx.data.backendSession` populated with:
  ```typescript
  {
    type: "authenticated",
    accessToken: {
      accessToken: "ENCODED_JWT",
      expiresAt: new Date(1735200000 * 1000)
    },
    refreshToken: "REFRESH_TOKEN_IF_PRESENT",
    authenticatedUser: {
      name: "Jane Visitor",
      email: "jane.visitor@example.com",
      picture: "https://example.com/avatar2.jpg",
      isOwner: false
    }
  }
  ```

**Implementation Notes:**

- Parse access token JWT to extract user details
- Extract `exp` claim for token expiration
- `isOwner = true` if `cognito:groups` contains "owners" or "admins"
- `isOwner = false` if only "visitors" group

---

### Test 20: it should handle ID token without picture claim

**Input:**

- Path: `/auth/callback?code=AUTH_CODE&state=EXPECTED_STATE`

**Mock ID Token Claims:**

```json
{
  "sub": "user-id",
  "email": "user@example.com",
  "name": "User Name",
  "cognito:groups": ["visitors"]
}
```

**Expected Output:**

- Successful redirect
- User details captured with `picture: undefined`

**Implementation Notes:**

- Picture is optional field in `AuthenticatedUser` interface
- Handle missing claims gracefully

---

### Test 21: it should handle missing cognito:groups claim as visitor

**Input:**

- Path: `/albums`
- Cookies: `dphoto-access-token=JWT_WITHOUT_GROUPS`

**Mock Access Token Claims:**

```json
{
  "sub": "user-id",
  "email": "unknown@example.com",
  "name": "Unknown User",
  "exp": 1735200000
}
```

**Expected Output:**

- `ctx.data.backendSession.authenticatedUser.isOwner = false`

**Implementation Notes:**

- If `cognito:groups` is missing or empty, default to visitor (isOwner: false)
- Security: err on side of least privilege

---

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
- Set-Cookie for `dphoto-refresh-token` with longer maxAge (decode refresh token or use default)

**Implementation Notes:**

- Use `tokens.expires_in` for access token cookie maxAge
- For refresh token, parse JWT to get exp claim, or use 30 days default
- This implements TODO: "use the real expiration time of the JWT token"

---

## Additional Implementation Details

### JWT Parsing Utility

Add to `web/src/core/security/jwt-utils.ts`:

```typescript
export interface CognitoTokenClaims {
  sub: string;
  email: string;
  name?: string;
  picture?: string;
  'cognito:groups'?: string[];
  exp: number;
  iat: number;
  iss: string;
}

export function parseJwtWithoutVerification(token: string): CognitoTokenClaims {
  const parts = token.split('.');
  if (parts.length !== 3) {
    throw new Error('Invalid JWT format');
  }
  
  const payload = Buffer.from(parts[1], 'base64url').toString('utf-8');
  return JSON.parse(payload);
}

export function isOwnerRole(groups?: string[]): boolean {
  if (!groups || groups.length === 0) {
    return false;
  }
  return groups.includes('owners') || groups.includes('admins');
}
```

### DynamoDB Storage Integration

This is referenced in TODO but may be implemented in a later step:

- Backend API endpoint to store identity details
- Called from middleware after successful authentication
- Matches Go implementation in `pkg/acl/aclcore/authenticate_sso.go`
- Can be a fire-and-forget call (don't block redirect on it)

