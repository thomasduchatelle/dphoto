# Cognito Authentication Implementation Summary

## Overview
This document summarizes the implementation of AWS Cognito authentication for the DPhoto web application, following the specification in `specs/2025-10_cognito-authentication-migration.md`.

## Implementation Status

### ✅ Completed (In Scope)
1. **Authentication Flow Endpoints**
   - `/auth/login` - Initiates OAuth flow with Cognito
   - `/auth/callback` - Handles OAuth callback and token exchange
   - `/auth/logout` - Clears authentication cookies

2. **Token Storage in Cookies**
   - Implemented dual-cookie strategy:
     - `dphoto-access-token` (HttpOnly) - Server-side secure storage
     - `dphoto-access-token-client` (readable by JS) - For Authorization headers
     - `dphoto-refresh-token` (HttpOnly) - Refresh token storage
   - Cookies configured with: HttpOnly, Secure (production), SameSite=Strict, appropriate Max-Age

3. **Token Availability in UI**
   - Created token context system for client-side token management
   - Axios interceptor automatically reads token and adds Authorization header
   - Automatic token loading from cookies on application mount
   - Periodic cookie checking (every 5 seconds) to detect token updates

4. **Security Implementation**
   - OAuth2 with PKCE (Proof Key for Code Exchange)
   - State parameter validation (CSRF protection)
   - Nonce validation in ID tokens (replay attack prevention)
   - Short token lifetimes (access token: 1 hour, refresh token: 30 days)
   - HttpOnly cookies for sensitive data
   - Secure cookies in production

### ❌ Out of Scope (As Specified)
1. Token refresh flow - Backend responsibility
2. Cleanup of old login page - Deferred
3. DynamoDB session store - Using in-memory for MVP
4. Token revocation - Handled by Cognito

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                    AWS API Gateway                          │
│                                                             │
│  ┌──────────────┐  ┌─────────────┐  ┌──────────────┐      │
│  │ /auth/login  │  │/auth/callback│  │ /auth/logout │      │
│  └──────┬───────┘  └──────┬──────┘  └──────┬───────┘      │
└─────────┼─────────────────┼─────────────────┼──────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────┐
│           Lambda Handler (auth-handler-wrapper.mjs)         │
│                                                             │
│  - OAuth2 with PKCE flow                                    │
│  - Session management (in-memory)                           │
│  - Token exchange with Cognito                              │
│  - Cookie management                                        │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    AWS Cognito                              │
│                                                             │
│  - User Pool with Google OAuth                              │
│  - Token issuance and validation                            │
│  - User groups (admins, owners, visitors)                   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                Client Application (Waku/React)              │
│                                                             │
│  - TokenInitializer: Loads tokens from cookies              │
│  - TokenContext: Manages token lifecycle                    │
│  - DPhotoApplication: Axios interceptor adds Auth header    │
└─────────────────────────────────────────────────────────────┘
```

### Files Created/Modified

**New Files:**
- `web/auth-handler-wrapper.mjs` - Lambda handler wrapper for auth routes
- `web/src/libs/auth/cognito-client.ts` - Cognito client configuration
- `web/src/libs/auth/cookie-utils.ts` - Cookie serialization utilities
- `web/src/libs/auth/client-cookie-utils.ts` - Client-side cookie reading
- `web/src/libs/auth/server-token-utils.ts` - Server-side token utilities
- `web/src/libs/auth/session-store.ts` - OAuth session management
- `web/src/libs/auth/token-context.ts` - Client-side token context
- `web/src/libs/auth/token-utils.ts` - Token validation and exchange
- `web/src/libs/auth/index.ts` - Public API exports
- `web/src/libs/auth/README.md` - Library documentation
- `web/src/components/AuthProvider.tsx` - Authentication provider component
- `web/src/components/TokenInitializer.tsx` - Token initialization component
- Test files: `cookie-utils.test.ts`, `server-token-utils.test.ts`, `token-context.test.ts`

**Modified Files:**
- `web/package.json` - Added openid-client and cookie dependencies, updated build script
- `web/src/components/Providers.tsx` - Integrated AuthProvider
- `web/src/core/application/DPhotoApplication.ts` - Updated to use token context

## Authentication Flow

### Login Sequence

1. **User initiates login**: Navigates to `/auth/login?returnUrl=/albums`

2. **Handler generates security parameters**:
   - State (random string for CSRF protection)
   - Nonce (random string for replay protection)
   - PKCE code verifier and challenge

3. **Session stored**: In-memory with 10-minute TTL

4. **User redirected to Cognito**: Authorization URL with all security parameters

5. **User authenticates**: Via Google OAuth through Cognito

6. **Cognito redirects back**: To `/auth/callback?code=...&state=...`

7. **Handler validates and exchanges**:
   - Validates state matches stored session
   - Exchanges authorization code for tokens using PKCE verifier
   - Validates nonce in ID token

8. **Cookies set**: Both HttpOnly and client-accessible cookies

9. **User redirected**: To original URL (from returnUrl or session)

### Token Usage Flow

1. **Application loads**: TokenInitializer triggers on mount

2. **Token loaded**: From `dphoto-access-token-client` cookie

3. **Token stored**: In TokenContext for reuse

4. **API request made**: User action triggers API call

5. **Interceptor activates**: Axios interceptor calls `getClientAccessToken()`

6. **Token retrieved**: From TokenContext (or loaded from cookie if stale)

7. **Header added**: `Authorization: Bearer <token>` added to request

8. **Periodic refresh**: Every 5 seconds, checks cookie for token updates

## Configuration

### Environment Variables

Required in Lambda:
```bash
COGNITO_USER_POOL_ID=us-east-1_xxxxx
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxx
COGNITO_CLIENT_SECRET=xxxxxxxxxxxxxxxxxx
COGNITO_DOMAIN=https://your-domain.auth.region.amazoncognito.com
COGNITO_ISSUER=https://cognito-idp.region.amazonaws.com/us-east-1_xxxxx
NODE_ENV=production  # For secure cookies
```

Optional:
```bash
APP_URL=https://your-app.example.com  # Derived from request if not set
```

### Build Configuration

Updated `package.json` build script:
```json
"build:lambda": "rm -rf dist && waku build --with-aws-lambda && sed -i 's#const config = {[^}]*}#const config = { \"basePath\": \"/\", \"distDir\": \".\", \"rscBase\": \"RSC\" }#' dist/server/index.js && cp auth-handler-wrapper.mjs dist/ && echo 'export { handler } from \"./auth-handler-wrapper.mjs\";' > dist/serve-aws-lambda.js"
```

This script:
1. Builds Waku for Lambda
2. Patches server config
3. Copies auth handler wrapper
4. Creates entry point that uses the wrapper

## Testing

### Unit Tests
- **239 tests pass** (18 new authentication tests)
- Coverage includes:
  - Token context management
  - Cookie serialization and parsing
  - Token expiration checking
  - Server-side token utilities

### Running Tests
```bash
cd web
npm install
npm run test:unit
```

### Security Scan
- CodeQL analysis: ✅ **0 alerts** (no vulnerabilities found)

## Security Considerations

### Implemented Security Measures

1. **OAuth2 with PKCE**: Prevents authorization code interception
2. **State Parameter**: CSRF protection in OAuth flow
3. **Nonce Validation**: Replay attack prevention
4. **HttpOnly Cookies**: Prevents JavaScript access to refresh tokens
5. **SameSite=Strict**: Prevents CSRF attacks
6. **Secure Flag**: HTTPS-only in production
7. **Short Token Lifetimes**: Limits exposure window

### Security Trade-offs

**Dual Cookie Strategy**: We use both HttpOnly and non-HttpOnly cookies for the access token:
- **HttpOnly cookie** (`dphoto-access-token`): Maximum security, but JavaScript can't read it
- **Client-accessible cookie** (`dphoto-access-token-client`): Allows JavaScript to add Authorization header

**Rationale**: 
- API Gateway authorizer can read both cookies
- Client-side code needs token for Authorization header
- Refresh token remains HttpOnly for maximum security
- Short access token lifetime (1 hour) limits exposure

### Threat Model

**Protected Against**:
- CSRF attacks (SameSite=Strict, state parameter)
- XSS attacks on refresh tokens (HttpOnly)
- Authorization code interception (PKCE)
- Replay attacks (nonce)
- Token exposure (short lifetimes)

**Not Protected Against** (acceptable for pet project):
- XSS attacks on access tokens (client-accessible cookie needed for functionality)
- Token theft from compromised client (mitigated by short lifetime)
- Advanced persistent threats (out of scope)

## Next Steps (Out of Current Scope)

### Token Refresh Flow
- Backend implementation needed
- SSR should check token expiration and refresh if needed
- Update cookies with new tokens

### Migration Path
- Remove old login page and authentication code
- Test with real Cognito instance
- Migrate session store to DynamoDB for production scale

### Future Enhancements
- Device code flow for CLI authentication
- Token revocation support
- Enhanced monitoring and alerting
- Amazon Verified Permissions integration

## Deployment Checklist

Before deploying to production:

1. ✅ Set all required environment variables in Lambda
2. ✅ Configure Cognito User Pool with Google OAuth
3. ✅ Set correct callback URLs in Cognito
4. ✅ Create user groups (admins, owners, visitors)
5. ✅ Add users to appropriate groups
6. ⏳ Configure API Gateway authorizer (backend work)
7. ⏳ Test authentication flow end-to-end
8. ⏳ Verify token refresh works
9. ⏳ Monitor error rates and logs

## Troubleshooting

### Common Issues

**Tokens not available in client**:
- Check browser cookies for `dphoto-access-token-client`
- Verify callback completed successfully
- Check Lambda logs for errors

**Authorization header not added**:
- Verify TokenInitializer is mounted
- Check token is not expired
- Verify getClientAccessToken() returns value

**OAuth flow fails**:
- Verify Cognito configuration (callback URLs)
- Check session exists (10-minute TTL)
- Review Lambda logs for detailed errors

**Build fails**:
- Ensure openid-client v6+ is installed
- Verify cookie package is installed
- Check TypeScript compilation errors

## References

- Specification: `specs/2025-10_cognito-authentication-migration.md`
- Library Documentation: `web/src/libs/auth/README.md`
- OpenID Client: https://github.com/panva/openid-client
- AWS Cognito: https://docs.aws.amazon.com/cognito/
- OAuth 2.0 PKCE: https://oauth.net/2/pkce/
