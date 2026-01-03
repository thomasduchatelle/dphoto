# Authentication Middleware for web-nextjs

This middleware implements OAuth 2.0 / OpenID Connect authentication flow for the DPhoto NextJS application.

## Features

- Automatic redirect to authentication provider when not authenticated
- OAuth 2.0 Authorization Code flow with PKCE
- Session management with secure HTTP-only cookies
- JWT token parsing for user information
- Backend session data passed to application via custom headers

## Files

### Core Implementation

- `middleware.ts` - Main NextJS middleware implementing authentication flows
- `lib/security/constants.ts` - Cookie names and type definitions
- `lib/security/jwt-utils.ts` - JWT parsing utilities

### Testing

- `middleware.test.ts` - Comprehensive test suite covering all authentication flows
- `__tests__/helpers/fake-oidc-server.ts` - Mock OIDC server for testing
- `__tests__/helpers/test-helper-oidc.ts` - Test utilities and fixtures

## Configuration

The middleware requires the following environment variables:

- `COGNITO_ISSUER` - URL of the OIDC issuer (e.g., AWS Cognito domain)
- `COGNITO_CLIENT_ID` - OAuth client ID
- `COGNITO_CLIENT_SECRET` - OAuth client secret

## Authentication Flow

### 1. Unauthenticated Access

When a user accesses any page without an access token cookie:

1. Middleware generates PKCE code verifier and challenge
2. Redirects user to authorization endpoint with:
   - Client ID
   - Redirect URI pointing to `/auth/callback`
   - PKCE parameters
   - State parameter for CSRF protection
3. Sets temporary cookies for state and code verifier (5 min expiry)

### 2. OAuth Callback (`/auth/callback`)

When the OAuth provider redirects back:

1. Middleware validates the state parameter
2. Exchanges authorization code for tokens using PKCE code verifier
3. Extracts user information from ID token
4. Sets secure cookies:
   - `dphoto-access-token` - JWT access token
   - `dphoto-refresh-token` - Refresh token
   - `dphoto-user-info` - User profile information
5. Clears temporary auth cookies
6. Redirects to home page

### 3. Authenticated Requests

For requests with valid access token:

1. Middleware parses JWT to extract user information
2. Creates BackendSession object with:
   - Access token and expiry
   - Refresh token
   - User details (name, email, picture, isOwner flag)
3. Passes session via `x-backend-session` header
4. Allows request to proceed

### 4. Explicit Login (`/auth/login`)

Allows re-authentication even when already logged in.

## Testing

Run tests with:

```bash
npm test
```

The test suite covers:
- Redirect to authorization when not authenticated
- Explicit login flow
- OAuth callback handling
- Authenticated request processing with session data

## Security Features

- HTTP-only cookies prevent XSS attacks
- Secure flag ensures cookies only sent over HTTPS
- SameSite attribute prevents CSRF attacks
- PKCE (Proof Key for Code Exchange) protects authorization code
- State parameter validates OAuth responses

## Implementation Notes

### NextJS Integration

The middleware uses NextJS's native middleware system. For authenticated requests, it:

1. Sets custom request headers that can be read by server components
2. Sets response headers for easier testing
3. Uses NextJS's `NextResponse` API for redirects and header manipulation

### JWT Security

JWT tokens are decoded but **not verified** in the middleware. This is safe because:
1. Tokens are verified by the backend API
2. Middleware only reads user information for display purposes
3. All security decisions are made by the backend based on token validation

## Future Enhancements

- Token refresh before expiry
- Error page for OAuth failures
- Redirect to originally requested URL after login
- Token caching/memoization for OIDC configuration
