# Authentication Middleware - Production Readiness Plan

## ✅ Plan Complete - Ready for Implementation

I've created a comprehensive test-driven plan to make the authentication middleware production-ready. The plan addresses all the TODO comments in the code plus
additional requirements you specified.

---

## 📁 Files Created

### Summary Document

- **`specs/middleware-SUMMARY.md`** (307 lines) - Complete overview and implementation guide

### Step-by-Step Specifications

1. **`specs/middleware-step1.md`** (124 lines) - Basic Authentication Flows
2. **`specs/middleware-step2.md`** (101 lines) - Redirect Path Preservation
3. **`specs/middleware-step3.md`** (94 lines) - Non-Navigation Request Handling
4. **`specs/middleware-step4.md`** (171 lines) - OAuth Callback Error Handling
5. **`specs/middleware-step5.md`** (213 lines) - ID Token Processing & User Details
6. **`specs/middleware-step6.md`** (165 lines) - Logout Flow
7. **`specs/middleware-step7.md`** (116 lines) - OIDC Configuration Caching
8. **`specs/middleware-step8.md`** (255 lines) - Edge Cases & Security Hardening

**Total:** 1,546 lines of detailed specifications

---

## 🎯 Test Coverage: 40 Test Cases

Each test case follows the principle: **1 test = 1 middleware call**

| Step | Tests | Focus Area                                                  |
|------|-------|-------------------------------------------------------------|
| 1    | 4     | Basic auth flows (redirect, callback, authenticated access) |
| 2    | 4     | Preserve & restore original URL during login                |
| 3    | 4     | API requests get 401, not redirects                         |
| 4    | 5     | Handle OAuth errors gracefully                              |
| 5    | 5     | Extract user details from JWT tokens                        |
| 6    | 5     | Complete logout flow + path whitelisting                    |
| 7    | 3     | Cache OIDC config for performance                           |
| 8    | 10    | Edge cases, security, malformed data                        |

---

## 🔑 Key Features Covered

### ✅ Happy Paths

- User visits page → redirects to Cognito login
- User authenticates → returns to original page
- Authenticated user accesses protected resources
- User logs out → clears cookies, signs out from Cognito

### ✅ Error Handling

- OAuth callback errors (invalid_grant, state_mismatch, etc.)
- Malformed JWT tokens
- Expired tokens
- Missing environment configuration
- Token exchange failures

### ✅ Security

- Open redirect prevention
- CSRF protection via OAuth state
- Cookie security (httpOnly, secure, sameSite)
- JWT validation
- Role-based access (isOwner from cognito:groups)
- Path whitelisting for public resources

### ✅ User Experience

- Original URL preservation during login
- Proper handling of API vs navigation requests
- User-friendly error pages
- Seamless authenticated navigation

### ✅ Performance

- OIDC configuration caching (fetch once, use forever)
- Efficient JWT parsing
- Minimal overhead for authenticated requests

---

## 📋 Requirements Coverage

### Your Specified Cases

1. ✅ **Happy path redirect flow** - Steps 1 & 2
2. ✅ **Do not redirect AJAX requests** - Step 3
3. ✅ **Logout happy path** - Step 6

### Additional TODO Comments from Code

4. ✅ **Cache OIDC configuration** - Step 7
5. ✅ **Handle OAuth callback errors** - Step 4
6. ✅ **Handle token exchange failures** - Step 4
7. ✅ **Extract user details from ID token** - Step 5
8. ✅ **Use real JWT expiration times** - Step 5
9. ✅ **Redirect to original URL** - Step 2

### Production Readiness Enhancements

10. ✅ **Edge case handling** - Step 8
11. ✅ **Security hardening** - Step 8
12. ✅ **Comprehensive error handling** - Steps 4 & 8
13. ✅ **Public path whitelisting** - Step 6
14. ✅ **CORS preflight support** - Step 8

---

## 🎨 Test Structure

Each test specification includes:

- **Explicit title**: "it should redirect to authorization authority when requesting home page without access token"
- **Input details**: Path, method, headers, cookies
- **Expected output**: Status code, headers, cookies, context data
- **Implementation notes**: Mocking strategy, security considerations

Example:

```typescript
it('should redirect to authorization authority when requesting home page without access token', async () => {
  // Input: GET / with no cookies
  // Output: 302 to Cognito, OAuth cookies set
  // See specs/middleware-step1.md for full details
});
```

---

## 🚀 Implementation Guidance

### Recommended Order

1. **Step 1** - Foundation (basic auth works)
2. **Step 7** - Performance (caching in place)
3. **Step 5** - Identity (user details correct)
4. **Step 2** - UX (redirect preservation)
5. **Step 3** - API support (proper status codes)
6. **Step 4** - Robustness (error handling)
7. **Step 6** - Complete flow (logout)
8. **Step 8** - Production ready (hardening)

### Files to Create/Modify

**Modify:**

- `web-nextjs/middleware-authentication.tsx` - main implementation
- `web/src/core/security/consts.ts` - add cookie constants
- `web/src/core/security/jwt-utils.ts` - JWT parsing utilities

**Create:**

- `web-nextjs/middleware-authentication.test.ts` - test suite
- `web/src/pages/auth/error.tsx` - OAuth error page
- `web/src/pages/auth/logout-complete.tsx` - post-logout page (optional)

---

## 📊 Success Criteria

The middleware is production-ready when:

- ✅ All 40 test cases pass
- ✅ No TODO comments remain
- ✅ Security checklist completed
- ✅ Error handling covers all code paths
- ✅ Performance optimized (cached config)
- ✅ UX is seamless (preserved redirects)
- ✅ Follows DPhoto coding standards (no comments in code)

---

## 💡 Next Steps

To implement:

1. Read `specs/middleware-SUMMARY.md` for complete overview
2. Start with `specs/middleware-step1.md`
3. Create test file with first 4 tests
4. Implement features to make tests pass
5. Continue through steps 2-8

Each step is self-contained with all context needed for an agent (human or AI) to implement it independently.

---

## 📝 Notes

- Tests use **TDD approach**: write test first, implement feature
- **No network calls** in tests: all external dependencies mocked
- **Deterministic tests**: fixed values for state, verifier, tokens
- **Clean code**: follow DPhoto standards (no comments, simple, well-tested)
- **Security first**: every step includes security considerations

---

**Total Specification Size:** 1,546 lines across 9 files
**Ready for implementation by any developer or AI agent with full context**
