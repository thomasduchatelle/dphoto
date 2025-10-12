# AWS Cognito Authentication Migration

## Feature Summary

Migrate the existing authentication and authorization system to AWS Cognito while maintaining the same user experience and permissions model. Users will authenticate via Google SSO through Cognito, with tokens stored in cookies for seamless SSR integration with Waku, and API access controlled by Cognito-based authorizers.

## Ubiquity Language

*No project-specific terms defined yet.*

## Scenarios

### Scenario 1: Successful Authentication Flow (Happy Path)
1. User navigates to any protected page in the application
2. Waku SSR checks for `dphoto-access-token` and `dphoto-refresh-token` cookies
3. Both tokens are missing, so user is redirected to Cognito login page
4. User clicks "Sign in with Google" on Cognito login page
5. Google OAuth flow completes successfully
6. Cognito matches the Google identity to an existing user in the `owners` group
7. Cognito issues access and refresh tokens
8. Tokens are set as HttpOnly cookies (`dphoto-access-token`, `dphoto-refresh-token`)
9. User is redirected back to the original protected page
10. Waku SSR validates the access token and renders the page successfully

### Scenario 2: Token Refresh During SSR
1. User navigates to a protected page
2. Waku SSR finds both tokens in cookies
3. Access token has expired, but refresh token is still valid
4. SSR makes a request to Cognito to refresh the access token using the refresh token
5. New access token is received and updated in the cookie
6. Page is rendered successfully with the refreshed token

### Scenario 3: Unknown User Authentication Failure
1. User navigates to a protected page and is redirected to Cognito login
2. User completes Google OAuth successfully
3. Cognito cannot find a matching user in any of the configured groups (admins, owners, visitors)
4. Authentication fails and user is redirected to `/errors/user-must-exists`
5. Error page displays message explaining that access must be granted by an administrator

### Scenario 4: API Access with Valid Token
1. Authenticated user makes an API request to a protected endpoint
2. Frontend includes the access token in the `Authorization` header
3. API Gateway authorizer validates the token against Cognito
4. User belongs to the `admins` group, which has access to this endpoint
5. Request is authorized and forwarded to the backend service
6. API response is returned successfully

### Scenario 5: API Access with Cookie-based Token
1. User's browser makes an API request (e.g., from JavaScript) without explicit Authorization header
2. Request includes the `dphoto-access-token` cookie automatically
3. API Gateway authorizer extracts and validates the token from the cookie
4. Token is valid and user has appropriate group membership
5. Request is authorized and processed successfully

## Target Architecture and Decisions

### Cognito User Pool Configuration
- **Single User Pool**: One Cognito User Pool will manage all users with three groups: `admins`, `owners`, and `visitors`
- **Google SSO Integration**: 
  - Use Cognito's built-in Google identity provider with OAuth 2.0
  - No domain restrictions - any Google account will be accepted for authentication
  - Map Google attributes (email, given_name, family_name) to Cognito user attributes
- **User Pool Settings**:
  - Username configuration: email as username
  - Required attributes: email, given_name, family_name
  - Email verification: disabled (Google already verifies email)
  - Password policy: not applicable (Google SSO only)
- **App Client Configuration**:
  - Authentication flows: `ALLOW_USER_SRP_AUTH` and `ALLOW_REFRESH_TOKEN_AUTH` only
  - Token expiration: Access token (1 hour), Refresh token (30 days), ID token (1 hour)
  - Callback URLs: application domains + `/auth/callback`
  - Logout URLs: application domains + `/auth/logout`
- **Group Strategy**: Users can belong to multiple groups simultaneously, with group membership determining API access permissions through token scopes

### SSR Authentication Flow
- **Library**: Use `openid-client` library for proper OIDC/OAuth2 flow implementation with automatic JWKS handling and refresh token rotation support
- **Token Validation**: Local JWT validation using `openid-client` with JWKS cached in Lambda memory
- **OAuth State Management**: DynamoDB session store for secure OAuth flow state
    - Session format: `AUTH_SESSION#{sessionId}` storing `{originalUrl, nonce, codeVerifier}`
  - TTL: 10 minutes for OAuth sessions
  - Session ID passed in Cognito's state parameter for security
- **Token Refresh Strategy**: 
    - Attempt refresh when access token is about to expire (< 5 min)
  - If refresh succeeds: render page with updated cookies (including rotated refresh token)
  - If refresh fails: clear cookies and redirect to Cognito login
- **Internal API Calls**: SSR passes validated access token in `Authorization: Bearer {token}` header to internal APIs
- **Stateless Design**: No server-side token storage - tokens read from cookies and validated on each request

## Topics to Discuss

- [X] **Cognito User Pool Configuration** - How to structure the user pool, groups (admins, owners, visitors), and Google SSO integration
- [X] **SSR Authentication Flow** - How Waku will handle token validation during server-side rendering and the redirect logic
- [ ] **Token Management Strategy** - Cookie configuration, token refresh mechanisms, and security considerations (HttpOnly, Secure, SameSite attributes)
- [ ] **API Gateway Authorizers** - Implementation details for the three authorizers (one per group) and token validation logic
- [ ] **User Matching and Group Assignment** - How Google SSO users will be mapped to existing Cognito users and assigned to appropriate groups
- [ ] **Error Handling and Edge Cases** - Token expiration scenarios, network failures, invalid tokens, and user access denied flows
- [ ] **Migration Strategy** - How to transition from the current authentication system to Cognito without disrupting existing users
- [ ] **Security and Compliance** - CORS configuration, token storage security, and any compliance requirements
- [ ] **Testing and Monitoring** - How to validate the authentication flow and monitor token usage/failures
- [ ] **Performance Considerations** - Caching strategies for token validation and potential impact on page load times
