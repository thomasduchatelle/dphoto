# Feature Summary
Migrate authentication and authorization to AWS Cognito using the Hosted UI with Google as the sole IdP. Tokens (access and refresh) will be issued by Cognito and stored in HttpOnly cookies, validated by SSR middleware (Waku) for page rendering, and enforced on APIs via API Gateway authorizers aligned with user groups (admins, owners, visitors).

# Ubiquity Language
- Turnstile: The SSR authentication gate that runs before page rendering. It validates local auth state, attempts token refresh when possible, or redirects the browser to the Hosted Login carrying a Return Stub.
- Entry Keycards: The two HttpOnly cookies used by this app to represent the signed-in session: dphoto-access-token (Access Keycard) and dphoto-refresh-token (Refresh Keycard).
- Return Stub: The URL-safe encoded original request URL we attach to the OAuth state parameter so the user returns to the exact page after login.
- Stage Pass: The API authorization check bound to a specific group (admins, owners, visitors). Implemented as a per-endpoint policy that evaluates the group claim in the access token.
- Unknown Guest: The failure state when a signed-in Google identity does not correspond to a known user in our system; we redirect to /errors/user-must-exists.

# Scenarios
## Scenario 1: Happy path SSR + Cognito login with Google SSO
1. user load a page (SSR rendered with Waku)
2. SSR page check the cookies to find refresh and access token
    * if all valid -> render the page
    * if access token invalid -> refresh it and render the page
    * if both are invalid or are not present -> redirect to cognito login page, then move to step 3
3. on cognito login page, user can authenticate with Google SSO (only)
4. once its Google identity received, it will be matched to users in cognito and its group
    * if either part of the group `admins`, `owners`, or `visitors` (or a combination of then), issue a set of tokens, then move to step 5
    * if user is unknown, fail the authentication and redirect to `/errors/user-must-exists`
5. tokens are recorded in the cookies (`dphoto-access-token` and `dphoto-refresh-token`) and the user is redirected to the original page
6. the page is rendered

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
