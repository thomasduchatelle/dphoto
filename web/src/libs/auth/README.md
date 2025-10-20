# Authentication Library

This library implements AWS Cognito authentication for the DPhoto web application following the design specified in `specs/2025-10_cognito-authentication-migration.md`.

## Architecture Overview

### Components

1. **Lambda Handler Wrapper** (`auth-handler-wrapper.mjs`)
   - Intercepts `/auth/login`, `/auth/callback`, and `/auth/logout` routes
   - Implements OAuth2 authorization code flow with PKCE
   - Manages session state for OAuth flow
   - Sets authentication cookies

2. **Client-Side Token Management**
   - `token-context.ts` - Token storage and lifecycle management
   - `client-cookie-utils.ts` - Cookie reading utilities
   - `TokenInitializer.tsx` - React component to initialize tokens on mount
   - `AuthProvider.tsx` - Wrapper component for the application

3. **Server-Side Utilities** (used in Lambda handler)
   - `cognito-client.ts` - OpenID Connect client configuration
   - `cookie-utils.ts` - Cookie serialization utilities
   - `session-store.ts` - OAuth session management
   - `token-utils.ts` - Token validation and exchange

## Authentication Flow

### Login Flow

1. User visits `/auth/login?returnUrl=/albums`
2. Handler generates OAuth state, nonce, and PKCE code verifier
3. Session stored in memory with 10-minute TTL
4. User redirected to Cognito login page
5. User authenticates with Google OAuth via Cognito
6. Cognito redirects to `/auth/callback?code=...&state=...`
7. Handler validates state, exchanges code for tokens
8. Tokens stored in cookies:
   - `dphoto-access-token` (HttpOnly, 1 hour)
   - `dphoto-access-token-client` (readable by JS, 1 hour)
   - `dphoto-refresh-token` (HttpOnly, 30 days)
9. User redirected to original URL

### Token Usage

On the client side:
```typescript
import { getClientAccessToken } from 'src/libs/auth';

// Get current access token (automatically reads from cookie if needed)
const token = getClientAccessToken();

// Axios interceptor automatically adds Authorization header
// No manual token management needed
```

The `DPhotoApplication` class automatically:
1. Reads tokens from the cookie via `token-context`
2. Adds `Authorization: Bearer <token>` header to all API requests
3. Periodically checks for token updates (every 5 seconds)

### Logout Flow

1. User navigates to `/auth/logout`
2. Handler clears all authentication cookies
3. User redirected to `/auth/login`

## Security Features

- **OAuth2 with PKCE**: Prevents authorization code interception attacks
- **HttpOnly Cookies**: Refresh tokens cannot be accessed by JavaScript
- **Dual Cookie Strategy**: 
  - HttpOnly for server-side security
  - Client-accessible for JavaScript Authorization header
- **Short Token Lifetime**: Access tokens expire after 1 hour
- **SameSite=Strict**: Prevents CSRF attacks
- **Secure Flag**: Cookies only sent over HTTPS in production

## Cookie Structure

| Cookie Name | Purpose | HttpOnly | Max-Age | Accessible |
|-------------|---------|----------|---------|------------|
| `dphoto-access-token` | Server-side validation | Yes | 1 hour | Server only |
| `dphoto-access-token-client` | Client-side API calls | No | 1 hour | JavaScript |
| `dphoto-refresh-token` | Token refresh | Yes | 30 days | Server only |

## Environment Variables

Required environment variables (set in Lambda):
- `COGNITO_USER_POOL_ID` - Cognito User Pool ID
- `COGNITO_CLIENT_ID` - Cognito App Client ID
- `COGNITO_CLIENT_SECRET` - Cognito App Client Secret
- `COGNITO_DOMAIN` - Cognito domain (e.g., `https://your-domain.auth.region.amazoncognito.com`)
- `COGNITO_ISSUER` - Cognito issuer URL (e.g., `https://cognito-idp.region.amazonaws.com/user-pool-id`)
- `APP_URL` - Application base URL (optional, derived from request if not set)

## Testing

Run tests:
```bash
npm run test:unit
```

Test files:
- `token-context.test.ts` - Token management tests
- `cookie-utils.test.ts` - Cookie serialization tests
- `server-token-utils.test.ts` - Token validation tests

## Integration

The authentication is integrated into the application through:

1. **Build Process**: `package.json` build script copies the handler wrapper
2. **Application Providers**: `AuthProvider` wraps the app in `Providers.tsx`
3. **Token Initialization**: `TokenInitializer` loads tokens on mount
4. **Axios Integration**: `DPhotoApplication` automatically adds tokens to requests

## Future Enhancements

Out of scope for this implementation:
- Token refresh flow (handled by backend)
- Cleanup of legacy login page
- DynamoDB session store (currently using in-memory store)
- Token revocation

## Troubleshooting

### Tokens not available in client
- Check browser cookies for `dphoto-access-token-client`
- Verify Cognito callback was successful
- Check Lambda logs for authentication errors

### Authorization header not added
- Verify `TokenInitializer` is mounted
- Check `getClientAccessToken()` returns a valid token
- Verify token is not expired

### OAuth flow fails
- Check Cognito configuration (callback URLs, client secret)
- Verify session exists (10-minute TTL)
- Check Lambda logs for detailed error messages
