# Authentication Middleware - Production Readiness Plan

## Overview

This document outlines the comprehensive test-driven plan to make the authentication middleware in `web/src/middleware/authentication.tsx` production-ready. The
plan is divided into 8 steps, each focusing on a specific aspect of the middleware functionality.

## Current State Analysis

The middleware currently has several TODO comments indicating incomplete features:

1. ❌ OIDC configuration is fetched on every request (not cached)
2. ❌ OAuth callback errors are not handled properly
3. ❌ Token exchange failures crash the middleware
4. ❌ ID token is not parsed to extract user details
5. ❌ JWT expiration times are not used for cookie maxAge
6. ❌ Original URL is not preserved during login flow
7. ❌ No logout functionality implemented
8. ❌ No distinction between navigation and API requests
9. ❌ No path whitelisting for public resources

## Implementation Steps

### Step 1: Basic Authentication Flows

**File:** `specs/middleware-step1.md`  
**Tests:** 4 test cases  
**Focus:** Core authentication flows

- Redirect unauthenticated users to Cognito
- Handle OAuth callback successfully
- Process tokens and set cookies
- Allow authenticated requests to proceed

**Key Features:**

- PKCE flow implementation
- Token exchange
- Cookie management
- Basic session establishment

---

### Step 2: Redirect Path Preservation

**File:** `specs/middleware-step2.md`  
**Tests:** 4 test cases  
**Focus:** User experience during authentication

- Store original URL before redirect
- Restore URL after successful login
- Prevent open redirect vulnerabilities
- Handle explicit login requests

**Key Features:**

- Redirect cookie management
- Security validation
- Default redirect behavior

---

### Step 3: Non-Navigation Request Handling

**File:** `specs/middleware-step3.md`  
**Tests:** 4 test cases  
**Focus:** API vs browser request handling

- Return 401 for unauthenticated API requests
- Distinguish between navigation and AJAX requests
- Handle Accept header properly
- Prevent redirecting API calls to login page

**Key Features:**

- Request type detection
- Appropriate status codes
- Client-side redirect capability

---

### Step 4: OAuth Callback Error Handling

**File:** `specs/middleware-step4.md`  
**Tests:** 5 test cases  
**Focus:** Robust error handling

- Handle OAuth errors in query parameters
- Catch token exchange failures
- Manage state mismatch
- Handle missing cookies
- Clean up OAuth state on errors

**Key Features:**

- Error detection and categorization
- User-friendly error messages
- OAuth state cleanup
- Error page creation

**New Files:**

- `web/src/pages/auth/error.tsx` - error display page

---

### Step 5: ID Token Processing and User Details

**File:** `specs/middleware-step5.md`  
**Tests:** 5 test cases  
**Focus:** User identity management

- Extract user details from ID token
- Parse JWT claims (name, email, picture)
- Determine user role from Cognito groups
- Use real token expiration times
- Populate backendSession correctly

**Key Features:**

- JWT parsing utilities
- Role determination logic
- Token expiration handling
- User detail extraction

**Enhanced Files:**

- `web/src/core/security/jwt-utils.ts` - JWT parsing functions

---

### Step 6: Logout Flow

**File:** `specs/middleware-step6.md`  
**Tests:** 5 test cases  
**Focus:** Secure logout

- Sign out from Cognito
- Clear all authentication cookies
- Redirect appropriately
- Whitelist public paths
- Handle post-logout landing

**Key Features:**

- Cognito logout integration
- Cookie cleanup
- Path whitelisting
- Public resource access

**New Files:**

- `web/src/pages/auth/logout-complete.tsx` - post-logout page (optional)

---

### Step 7: OIDC Configuration Caching

**File:** `specs/middleware-step7.md`  
**Tests:** 3 test cases  
**Focus:** Performance optimization

- Cache OIDC configuration in memory
- Fetch configuration only once
- Handle fetch failures gracefully
- Provide cache invalidation for testing

**Key Features:**

- In-memory configuration cache
- Error handling
- Test utilities

---

### Step 8: Edge Cases and Security Hardening

**File:** `specs/middleware-step8.md`  
**Tests:** 10 test cases  
**Focus:** Production robustness

- Handle malformed JWTs
- Manage expired tokens
- Validate environment configuration
- Prevent open redirects
- Handle concurrent requests
- Manage cookie size limits
- Validate OAuth parameters
- Support CORS preflight
- Rate limiting considerations
- Security hardening

**Key Features:**

- Comprehensive error handling
- Security validations
- Edge case coverage
- Logging and monitoring

---

## Test Implementation Strategy

Each step includes:

1. **Unit tests** for the middleware function
2. **One test = one function call** principle
3. **Explicit test titles** that describe the scenario
4. **Detailed specifications** including:
    - Input (path, headers, cookies)
    - Expected output (status, headers, cookies, context data)
    - Implementation notes

### Test File Structure

All tests will be in: `web/src/middleware/authentication.test.ts`

```typescript
describe('Authentication Middleware', () => {
  describe('Basic Authentication Flows', () => {
    it('should redirect to authorization authority when requesting home page without access token', async () => {
      // Test implementation
    });
    // ... more tests
  });
  
  describe('Redirect Path Preservation', () => {
    // Tests from step 2
  });
  
  // ... other test groups
});
```

### Mock Strategy

Tests will use:

- **Fake HandlerContext**: Mock request/response objects
- **Fake OIDC client**: Mock `openid-client` library functions
- **Deterministic values**: Fixed state, verifier, tokens for predictable tests
- **No network calls**: All external dependencies mocked

---

## Implementation Order

Recommended implementation sequence:

1. **Step 1** - Get basic auth working (foundation)
2. **Step 7** - Add caching early (affects all other features)
3. **Step 5** - User details (needed for proper session)
4. **Step 2** - Redirect preservation (UX improvement)
5. **Step 3** - API request handling (critical for frontend)
6. **Step 4** - Error handling (robustness)
7. **Step 6** - Logout (complete the auth lifecycle)
8. **Step 8** - Hardening (production readiness)

## Success Criteria

The middleware is production-ready when:

- ✅ All 40 test cases pass
- ✅ No TODO comments remain
- ✅ Security checklist completed
- ✅ Error handling covers all paths
- ✅ Performance is optimized (cached config)
- ✅ User experience is smooth (preserved redirects)
- ✅ API and browser requests handled correctly
- ✅ Logout flow complete
- ✅ Edge cases covered
- ✅ Code follows project standards (no comments, clean code)

## Test Metrics

Total test cases planned: **40**

Breakdown by step:

- Step 1: 4 tests (Basic flows)
- Step 2: 4 tests (Redirects)
- Step 3: 4 tests (API handling)
- Step 4: 5 tests (Error handling)
- Step 5: 5 tests (User details)
- Step 6: 5 tests (Logout)
- Step 7: 3 tests (Caching)
- Step 8: 10 tests (Edge cases)

## Development Workflow

For each step:

1. **Read the step specification** (`specs/middleware-step{n}.md`)
2. **Write the tests** based on specifications
3. **Run tests** (they should fail - TDD)
4. **Implement the feature** to make tests pass
5. **Refactor** for clean code
6. **Verify** no other tests broke
7. **Move to next step**

## Files Modified/Created

### Modified:

- `web/src/middleware/authentication.tsx` - main implementation
- `web/src/core/security/consts.ts` - new cookie constants
- `web/src/core/security/jwt-utils.ts` - JWT parsing utilities

### Created:

- `web/src/middleware/authentication.test.ts` - comprehensive test suite
- `web/src/pages/auth/error.tsx` - OAuth error display page
- `web/src/pages/auth/logout-complete.tsx` - post-logout page (optional)

### Documentation:

- `specs/middleware-step1.md` through `specs/middleware-step8.md`

---

## Next Actions

To begin implementation:

1. Start with Step 1: Read `specs/middleware-step1.md`
2. Create `web/src/middleware/authentication.test.ts`
3. Implement the first 4 test cases
4. Update middleware to make tests pass
5. Proceed to Step 7 (caching) for better performance foundation
6. Continue through remaining steps

Each step is self-contained with clear inputs, outputs, and implementation guidance.

