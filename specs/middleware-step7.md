# Middleware Step 7: OIDC Configuration Caching

## Overview

Implement caching for the OIDC configuration to avoid fetching it on every request. The configuration should be fetched once and cached in memory.

## Files to Modify

- `web/src/middleware/authentication.tsx` - implement config caching
- `web/src/middleware/authentication.test.ts` - add test cases

## Test Cases

### Test 28: it should fetch OIDC configuration only once across multiple requests

**Input:**

- Multiple requests in sequence (simulate 3 requests):
    1. Path: `/`, no auth (redirect to login)
    2. Path: `/auth/callback`, with valid OAuth code
    3. Path: `/albums`, no auth (redirect to login)

**Expected Output:**

- OIDC discovery should be called only once
- All three requests should use the cached configuration

**Implementation Notes:**

- Track number of calls to `client.discovery()`
- Use module-level cache variable:
  ```typescript
  let oidcConfigCache: client.Configuration | null = null;
  ```
- First request fetches and caches
- Subsequent requests use cached value

---

### Test 29: it should handle OIDC configuration fetch failure gracefully

**Input:**

- Path: `/`, no auth
- Mock `client.discovery()` to throw network error

**Expected Output:**

- Status: `500 Internal Server Error` or `503 Service Unavailable`
- Error response body indicating configuration failure
- Error should be logged

**Implementation Notes:**

- Wrap `oidcConfig()` call in try-catch
- Don't cache failed results
- Return appropriate error response to user
- Log error for monitoring

---

### Test 30: it should allow cache invalidation for testing or reconfiguration

**Input:**

- Programmatic cache reset (for testing or hot reload scenarios)
- Followed by request that needs config

**Expected Output:**

- Config should be fetched again after cache clear

**Implementation Notes:**

- This is primarily for test isolation
- Export a cache-clear function:
  ```typescript
  export function clearOidcConfigCache() {
    oidcConfigCache = null;
  }
  ```
- Tests should call this in beforeEach/afterEach

---

## Implementation Pattern

```typescript
let oidcConfigCache: client.Configuration | null = null;

async function oidcConfig(): Promise<client.Configuration> {
  if (oidcConfigCache) {
    return oidcConfigCache;
  }
  
  try {
    const config = await client.discovery(
      new URL(getEnv("COGNITO_ISSUER")),
      getEnv("COGNITO_CLIENT_ID"),
      getEnv("COGNITO_CLIENT_SECRET"),
    );
    
    oidcConfigCache = config;
    return config;
  } catch (error) {
    console.error('Failed to fetch OIDC configuration:', error);
    throw new Error('Authentication service configuration unavailable');
  }
}

export function clearOidcConfigCache(): void {
  oidcConfigCache = null;
}
```

## Performance Impact

Without caching:

- Every unauthenticated request fetches OIDC config (network call)
- Every OAuth callback fetches config
- Adds 50-200ms latency per request

With caching:

- First request: ~100ms for config fetch
- Subsequent requests: 0ms (memory read)
- Significantly improves response time for authentication flow

