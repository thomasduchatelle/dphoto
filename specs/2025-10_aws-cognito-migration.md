# AWS Cognito Authentication and Authorization Migration

1. Feature Summary
Migrate authentication and authorization to AWS Cognito using the Hosted UI with Google SSO. Maintain existing application behavior, but delegate token issuance and login UI to Cognito; store tokens in cookies (dphoto-access-token, dphoto-refresh-token); integrate with SSR (Waku) and API Gateway authorizers.

2. Ubiquity Language
No terms defined yet. We will add only project-specific, uncommon terms when identified by the lead developer.

3. Scenarios
To be authored after decisions on flows and integrations. Will include success and failure paths for SSR rendering, token refresh, API access, and unknown-user handling.

4. Target architecture and decisions
To be populated as we discuss each topic, including explicit service settings (e.g., Cognito App Client settings, callback URLs, API Gateway authorizer configuration) and out-of-scope items.

5. Topics to discuss
- [ ] Cognito User Pool design
  - Single User Pool vs multiple; App Client(s) per environment
  - Domain for Hosted UI (Cognito domain vs custom domain)
  - Regions and tenancy implications
- [ ] Identity Provider (IdP): Google SSO
  - OIDC configuration (Client ID/Secret), allowed domains (if any), scopes (openid, email, profile)
  - Authorization Code + PKCE vs other flows (recommend: Authorization Code + PKCE)
- [ ] Cognito Groups and authorization model
  - Groups: admins, owners, visitors; mapping Google users to these groups
  - How unknown users are handled (redirect to /errors/user-must-exists)
  - Whether group membership is managed manually or via automated synchronization
- [ ] Callback and logout endpoints
  - Application callback route(s), preserving and validating state and nonce, original URL redirection
  - Post-logout redirect URIs and logout propagation across apps
- [ ] Token model and lifetimes
  - Access/refresh token TTLs, rotation policy, max session duration
  - Claims required by the app (e.g., sub, email, groups) and custom claims if needed
- [ ] Cookie strategy
  - Cookie names: dphoto-access-token, dphoto-refresh-token
  - Attributes: HttpOnly, Secure, SameSite (Lax vs None), path, domain
  - CSRF protections for browser flows
- [ ] SSR (Waku) integration
  - Middleware/handler to parse cookies, validate/refresh access token, redirect to Cognito as needed
  - Preserving original URL across redirects; error handling on callback
- [ ] Token refresh logic
  - Where refresh occurs (SSR only vs client too), retry strategies, handling refresh failure
  - Clock skew tolerance
- [ ] API Gateway authorizers
  - One authorizer per group vs single authorizer with RBAC evaluation (trade-offs: manageability, latency, cache footprint)
  - JWT (Cognito) authorizer vs Lambda authorizer (requirement to support tokens in Authorization header or Cookie)
  - Policy caching, TTL, and failure modes (401 vs 403)
- [ ] Passing tokens to APIs
  - Enforce Authorization: Bearer header vs permitting Cookie; if Cookie is allowed, consistent name and gateway mapping
  - Frontend/API client changes, including SSR or fetch interceptors
- [ ] Error handling and UX
  - /errors/user-must-exists behavior and copy
  - Handling IdP/Cognito errors (account disabled, consent denied, invalid callback)
- [ ] Infrastructure as Code
  - Tooling (Terraform/CDK/SAM), module structure, environment separation (dev/stage/prod)
  - Secrets management (Google client secret), SSM/Secrets Manager
- [ ] DNS and certificates
  - Hosted UI domain and TLS certs (ACM), application domain cookies scope
- [ ] Observability and audit
  - CloudWatch logs/metrics, structured logging, alarms
  - Audit trail: login attempts, group changes, admin actions
- [ ] Security posture
  - PKCE, nonce/state handling, replay protection, token revocation strategy
  - Minimum TLS version, HSTS, WAF considerations
- [ ] Migration and rollout plan
  - Parallel run/feature flag, user import/linking strategy, rollback plan
  - Communication to users/admins about new login flow
- [ ] Rate limiting and quotas
  - API Gateway throttling, Cognito limits, protection against token refresh storms
- [ ] Testing strategy
  - Dev/staging Cognito User Pools, test Google project, end-to-end tests
  - Local development approach (using real Cognito vs mocks)
- [ ] Documentation and runbooks
  - On-call procedures for auth outages, token invalidation, user onboarding/offboarding
- [ ] Out of scope (to confirm)
  - Non-Google IdPs
  - Password-based sign-in
  - Authorization changes beyond group-based access currently in place
