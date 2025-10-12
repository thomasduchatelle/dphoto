# Feature Summary
Migrate authentication and authorization to AWS Cognito using the Hosted UI with Google as the sole IdP. Tokens (access and refresh) will be issued by Cognito and stored in HttpOnly cookies, validated by SSR middleware (Waku) for page rendering, and enforced on APIs via API Gateway authorizers aligned with user groups (admins, owners, visitors).

# Ubiquity Language
- DPHOTO Cookies: The two cookies used by the app for auth: dphoto-access-token (JWT access token) and dphoto-refresh-token (refresh token). Always HttpOnly and Secure; scoped domain and SameSite decided during design.
- SSR Gatekeeper: The Waku SSR middleware that validates/refreshes tokens, triggers redirects to Cognito Hosted UI, and preserves the original URL via OAuth state.
- Login Return URL: The exact URL the user initially requested, encoded into the OAuth state parameter and restored after successful authentication.
- API Shield Authorizer: The API Gateway authorizer dedicated to a specific group (admins, owners, visitors) that permits/denies requests based on token claims.
- Auth Session: The pair of Cognito-issued tokens (access + refresh) represented by the DPHOTO Cookies and their lifecycle/rotation rules within the app.
- User-Unknown Error Route: The canonical error endpoint /errors/user-must-exists used when a Google identity has no matching Cognito user.

# Scenarios
TBD with lead developer during topic-by-topic decisions. We will add 5–8 complete journeys covering happy-path SSR, token refresh, unknown user error, API calls with each role, and logout/invalidation flows.

# Target architecture and decisions
This section will be filled as we converge on each topic. Initial constraints from the brief:
- Use Cognito Hosted UI with Google SSO (only). Users must pre-exist in Cognito and belong to admins, owners, or visitors (any combination).
- SSR (Waku) checks cookies, refreshes access token when valid refresh exists, otherwise redirects to Cognito and then restores the original URL.
- API Gateway enforces access via authorizers aligned with groups; token may be provided via Authorization header or cookie.

Out of scope (initial):
- Any non-Google identity providers; password-based sign-in; custom login pages beyond Hosted UI branding.
- Self-service signup flows (user must already exist in Cognito).
- Issuing or verifying legacy tokens from the previous auth system.
- RBAC/ABAC beyond the three groups listed unless later expanded.

# Topics to discuss
Please pick the first topic to cover; we’ll document decisions here and iterate.

- [ ] Cognito user pool structure and environments
  - Single pool per environment vs single multi-env pool; user pool domain naming; regional selection.
- [ ] Hosted UI and OAuth/OIDC flow details
  - Authorization Code with PKCE vs confidential client; callback paths; scopes (openid, email, profile); state/nonce handling.
- [ ] Google IdP configuration
  - OIDC vs SAML setup; required Google project credentials; allowed domains; email verification; linking behavior.
- [ ] User provisioning and group assignment
  - Source of truth for group membership (Terraform/IaC, admin console, SCIM); handling unknown users; sync strategy from existing system.
- [ ] Token model and lifetimes
  - Access token TTL (e.g., 15 min), refresh token TTL (e.g., 30 days), rotation policy, replay detection, revocation strategy, clock skew handling.
- [ ] Cookie strategy for DPHOTO Cookies
  - Domain, path, Secure, HttpOnly, SameSite (Lax/None), max-age, rotation, size constraints, cross-subdomain behavior, encryption vs signed values.
- [ ] Waku SSR Gatekeeper implementation
  - JWKS retrieval and caching, local signature verification vs token introspection, refresh flow with Cognito token endpoint, preserving original URL, error mapping.
- [ ] Callback endpoint(s) and redirect handling
  - /auth/callback responsibilities, code exchange, setting cookies, CSRF/state validation, handling partial failures.
- [ ] API Gateway authorizers per group
  - Lambda Authorizer vs JWT Authorizer trade-offs; one-per-group design; evaluating cognito:groups claim; policy caching, performance and cold starts.
- [ ] Request token source precedence
  - Authorization header vs cookie precedence; normalizing inbound requests at API Gateway or at an edge/middleware layer.
- [ ] Route protection matrix
  - Mapping site routes and API routes to required groups; public routes; mixed-mode pages (SSR + client fetches).
- [ ] Logging, monitoring, and alerting
  - CloudWatch logs/metrics for authorizers, Hosted UI events, SSR; alarms for spikes in auth failures; audit trails for group changes.
- [ ] Security considerations
  - PKCE enforcement, SameSite/CSRF strategy, refresh token theft mitigation, JWKS rotation cadence, minimal IAM for IaC and runtime roles.
- [ ] IaC and deployment
  - Terraform/CloudFormation/CDK choice, modules for Cognito, API Gateway authorizers, Waku config; per-env variables and secrets management.
- [ ] Migration and cutover plan
  - Parallel run vs big bang, user migration/matching, rollback strategy, subject/identifier continuity, comms plan for admins/owners/visitors.
- [ ] Testing strategy
  - E2E browser flows, negative tests (expired/invalid tokens), API authz tests per group, load testing for authorizers, chaos tests for JWKS/key rotations.
- [ ] Logout and session invalidation
  - Hosted UI sign-out endpoints, global sign-out, cookie clearing, refresh token revocation and propagation lag.
- [ ] CORS and cross-origin behavior
  - API and SSR cookie behavior across domains/subdomains; CORS headers; preflight caching and implications when using cookies.

---
Meta
- Date basis for this spec filename assumed 2025-10-12. Please confirm current date/timezone; we will rename if needed.
