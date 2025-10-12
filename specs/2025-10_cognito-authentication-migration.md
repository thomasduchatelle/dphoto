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

### Scenario 6: Admin Creates New Owner
1. Admin user (member of `admins` group) calls API endpoint to create a new owner
2. API validates admin's access token and group membership
3. New Cognito user is created in the `owners` group with provided email
4. Owner can now authenticate via Google SSO using that email address

### Scenario 7: Album Sharing Creates Visitor
1. Owner shares an album with an email address not in the system
2. API automatically creates a new Cognito user in the `visitors` group
3. New visitor receives invitation and can authenticate via Google SSO
4. Visitor gains access only to the shared album(s)

### Scenario 8: Token Validation Failures During SSR
1. User navigates to a protected page
2. Waku SSR encounters one of the following token issues:
   - Access token expired and refresh token is also expired/invalid
   - Malformed or tampered tokens with invalid signatures
   - Missing or corrupted token data
3. SSR logs security event (for signature failures) and clears both cookies
4. User is redirected to `/errors/session-timed-out`
5. Error page displays "Your session has timed out. Please log in again."
6. User clicks login button to restart authentication flow

### Scenario 9: Network Error During Authentication
1. User completes Google OAuth flow successfully
2. Network error occurs while exchanging authorization code for tokens
3. Authentication Lambda logs error and returns 500
4. User sees technical error page: "A technical error occurred. Please try again later."
5. User can retry the authentication process

### Scenario 10: API Authorization Failure
1. User with `visitors` group makes API request requiring `owners` access
2. API Gateway authorizer validates token but denies access due to insufficient permissions
3. API returns 403 Forbidden with error message
4. Frontend displays user-friendly "Access denied" message

### Scenario 11: API Technical Error
1. User makes valid API request with proper authorization
2. API Gateway authorizer Lambda fails due to technical issue
3. API Gateway returns 500 Internal Server Error
4. Error is logged in CloudWatch for 30-day retention
5. Frontend displays technical error message to user

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

### User Matching and Group Assignment
- **Pre-provisioning Required**: All users must exist in Cognito before authentication (matches existing system behavior)
- **Identity Matching**: Google email address must exactly match Cognito username (email field)
- **No Email Changes**: System does not support users changing their Google email address
- **Group Assignment Rules**:
  - `admins`: Users who can create other owners and perform administrative functions
  - `owners`: Users with full control over their content (MainOwnerScope equivalent)
  - `visitors`: Users with read access to specific shared albums (AlbumVisitorScope equivalent)
- **User Creation Flows**:
  - **Admin creates owner**: API endpoint protected by `admins` group creates new Cognito user in `owners` group
  - **Album sharing auto-creation**: When owner shares album with unknown email, automatically create Cognito user in `visitors` group
  - **Migration strategy**: Manually recreate existing 6 users in appropriate Cognito groups (no complex migration needed)
- **CLI Authentication**: Current CLI will continue using direct AWS resource access via IAM (no device authentication migration)

### Token Management Strategy
- **Cookie Configuration**:
  - Names: `dphoto-access-token` and `dphoto-refresh-token`
  - Security: HttpOnly, Secure (HTTPS only), SameSite=Strict
  - Domain: Set to application domain for cross-subdomain access if needed
  - Path: `/` for application-wide access
- **Token Storage**: Cookies only - no server-side token storage for stateless design
- **Refresh Mechanism**: 
  - Proactive refresh when access token expires in < 5 minutes
  - Handle refresh token rotation (Cognito rotates refresh tokens on use)
  - Clear cookies and redirect on refresh failure
- **Security Considerations**: 
  - Tokens never exposed to client-side JavaScript (HttpOnly)
  - HTTPS required for Secure cookie flag
  - Short access token lifetime (1 hour) limits exposure window

### Migration Strategy
- **Simplicity First**: With only 6 existing users, manual recreation is simpler than complex migration
- **User Recreation Process**:
  1. Create Cognito User Pool with Google IdP configuration
  2. Manually create each user in appropriate Cognito group based on current permissions
  3. Test authentication with each user before go-live
  4. Deploy new authentication system and decommission old system
- **Rollback Plan**: Keep old authentication system deployable until new system is fully validated
- **User Communication**: Inform users of the change and verify their Google email addresses match expectations

### API Gateway Authorizers
- **Single Lambda Authorizer**: One unified Lambda authorizer function for all API Gateway routes
- **Token Extraction Logic**:
  1. Extract token from `Authorization: Bearer {token}` header (priority)
  2. If not found, extract from `dphoto-access-token` cookie (fallback)
  3. Return 401 Unauthorized if no token found
- **Token Validation Process**:
  1. Validate JWT signature against Cognito JWKS
  2. Check token expiration and issuer
  3. Extract user groups from token claims (`cognito:groups`)
- **Authorization Logic**:
  - **Route Configuration**: Each API Gateway route specifies required group via authorizer configuration
  - **Multi-group Support**: Users can belong to multiple groups (`admins` + `owners`)
  - **Group Access Rules**:
    - `admins` routes: Require `admins` group membership
    - `owners` routes: Require `owners` group membership
    - `visitors` routes: Any group membership grants access (hierarchical)
- **Backend Integration**:
  - **No Authorization Caching**: Full token validation on every request
  - **Token Pass-through**: Original access token forwarded to backend services in `Authorization` header
  - **Backend Re-validation**: Each backend service independently validates the token for security

### Error Handling and Edge Cases
- **Token Validation Failures**:
  - Invalid, malformed, or expired tokens: Clear cookies and redirect to session timeout error page
  - Failed token refresh: Clear cookies and redirect with "session timed out" message requiring re-authentication
  - Security violations (signature failures): Log event and treat as session timeout
- **Network and Service Errors**:
  - Network failures during authentication: Return 500 error with technical error message
  - Service unavailability: Fail with 500 and user-friendly error message, log all errors to CloudWatch (30-day retention)
- **OAuth Flow Errors**:
  - State mismatch, session timeouts, or flow failures: Clear session data and redirect to try-again error page
  - Only secure authorization code flow implemented - no fallback mechanisms
- **API Authorization**:
  - Insufficient permissions: Return 403 Forbidden with appropriate error message
  - Authorizer technical failures: Return 500 with error logged to CloudWatch
  - SSR authorization failures: Render user-friendly error pages
- **Stateless Token Design**:
  - Access tokens remain valid for full lifetime (1 hour) regardless of group changes
  - No server-side token revocation - rely on short token lifetime
  - Cognito prevents refresh of revoked refresh tokens automatically
- **Error Pages Required**:
  - Technical error page (500 scenarios)
  - Forbidden access page (403 scenarios)  
  - Session timed out page (token expiration/invalid)
  - User must exist page (authentication without valid Cognito user)
- **No Retry Logic**: All failures are final - users must manually retry operations
- **Logging Strategy**: All Lambda functions log errors to CloudWatch with 30-day retention, no additional alerting required

### Security and Compliance
- **CORS Configuration**:
  - API Gateway CORS configured to allow requests only from application domains
  - Restrict allowed origins to specific application URLs (no wildcard origins)
  - Standard CORS headers: `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`
  - No credentials sharing with external domains
- **Token Security**:
  - Tokens stored exclusively in HttpOnly cookies (not accessible to JavaScript)
  - Secure flag enforced (HTTPS-only transmission)
  - SameSite=Strict prevents CSRF attacks
  - Short access token lifetime (1 hour) limits exposure window
  - Refresh token rotation on each use prevents replay attacks
- **Transport Security**:
  - HTTPS mandatory for all authentication endpoints and cookie transmission
  - TLS 1.2+ required for all external communication
  - No fallback to HTTP for authentication flows
- **OAuth Security**:
  - Authorization code flow with PKCE (Proof Key for Code Exchange)
  - State parameter validation prevents CSRF
  - Nonce validation in ID tokens prevents replay attacks
  - Callback URL validation against registered application domains
- **Session Security**:
  - OAuth sessions stored in DynamoDB with 10-minute TTL
  - Session IDs cryptographically secure and unpredictable
  - No sensitive data stored in browser storage (localStorage/sessionStorage)
- **Compliance Requirements**: None - this is a personal application with no regulatory compliance requirements
- **Security Logging**:
  - Failed authentication attempts logged to CloudWatch
  - Token signature validation failures logged as security events
  - No PII (Personally Identifiable Information) logged in security events
  - 30-day log retention for security audit trail

## Topics to Discuss

- [X] **Cognito User Pool Configuration** - How to structure the user pool, groups (admins, owners, visitors), and Google SSO integration
- [X] **SSR Authentication Flow** - How Waku will handle token validation during server-side rendering and the redirect logic
- [X] **User Matching and Group Assignment** - How Google SSO users will be mapped to existing Cognito users and assigned to appropriate groups
- [X] **Token Management Strategy** - Cookie configuration, token refresh mechanisms, and security considerations (HttpOnly, Secure, SameSite attributes)
- [X] **Migration Strategy** - How to transition from the current authentication system to Cognito without disrupting existing users
- [X] **API Gateway Authorizers** - Implementation details for the unified authorizer, token validation logic, and group-based authorization
- [X] **Error Handling and Edge Cases** - Token expiration scenarios, network failures, invalid tokens, and user access denied flows
- [X] **Security and Compliance** - CORS configuration, token storage security, and any compliance requirements
- [ ] **Testing and Monitoring** - How to validate the authentication flow and monitor token usage/failures
- [ ] **Performance Considerations** - Caching strategies for token validation and potential impact on page load times
- [ ] **Device Authentication for CLI** - Future consideration for migrating CLI from direct AWS access to API-based authentication
- [ ] **Amazon Verified Permissions** - Evaluate Amazon Verified Permissions service for fine-grained authorization policies and permissions management
