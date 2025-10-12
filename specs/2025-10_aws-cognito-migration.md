# AWS Cognito Migration — Requirements
Date: 2025-10-12

1. Feature Summary
- Migrate authentication and authorization to AWS Cognito using the Hosted UI with Google as the sole IdP, replacing our existing token generation and login pages while preserving current access semantics.
- SSR pages (Waku) will read tokens from cookies, refresh when needed, or redirect to Cognito for login; APIs behind API Gateway will be protected by authorizers that validate Cognito access tokens from Authorization header or cookies.

2. Ubiquity Language
- Hosted UI: Cognito-managed OAuth/OIDC login pages and flows.
- User Pool: Cognito user directory containing users and groups.
- App Client: Cognito application configuration used by our web app (PKCE-enabled).
- IdP (Google): External identity provider (Google OAuth) linked to the User Pool.
- Access Token: JWT used to authorize API requests; short-lived.
- Refresh Token: Opaque token used to obtain new access tokens; longer-lived.
- ID Token: JWT conveying user identity/claims; not used for API authorization.
- Authorizer: API Gateway component (Cognito or Lambda/JWT) validating access tokens.
- Waku SSR: Our server-side rendering layer that runs before sending HTML to clients.

3. Scenarios
(TBD after we align on topics; will include 5–8 full user journeys covering happy-path, permission boundaries, and error handling.)

4. Target architecture and decisions
(Decisions will be logged here as we cover topics below. Also list out-of-scope items.)

5. Topics to discuss
- [ ] Authentication flow selection (Authorization Code with PKCE via Cognito Hosted UI with Google IdP)
- [ ] Cognito resources and configuration
  - [ ] User Pool setup (attributes, policies, advanced security)
  - [ ] App Client settings (PKCE only, callback/logout URLs, allowed OAuth flows/scopes)
  - [ ] Cognito domain/Hosted UI (custom domain vs AWS subdomain)
  - [ ] Resource Servers and custom scopes (if needed)
- [ ] Google IdP configuration
  - [ ] Google Cloud project, OAuth consent screen, client ID/secret (stored in AWS Secrets Manager)
  - [ ] Scopes (openid, email, profile), domain/workspace restrictions if any
- [ ] User lifecycle and group management
  - [ ] Pre-provisioning users vs. just-in-time (require “user must exist”)
  - [ ] Group assignment flow for admins, owners, visitors
  - [ ] Mapping from Google identity (sub/email) to Cognito user and groups
  - [ ] Operational process to add/remove users and audit changes
- [ ] Token strategy and cookie management
  - [ ] Access/refresh token lifetimes, rotation, revocation, clock skew tolerance
  - [ ] Cookie names: dphoto-access-token, dphoto-refresh-token
  - [ ] Cookie attributes: HttpOnly, Secure, SameSite (Lax vs None), Domain and Path strategy
  - [ ] Where cookies are set/cleared (SSR/backend only; never via client JS)
- [ ] SSR integration (Waku)
  - [ ] Middleware to read/validate cookies, refresh access token if invalid, else redirect to Hosted UI
  - [ ] Preserving original URL across redirects (state param), and handling fragments
  - [ ] Handling callback on /auth/callback (code exchange, PKCE verification, cookie set)
  - [ ] Error routing to /errors/user-must-exists and other failure pages
- [ ] OAuth/OIDC security controls
  - [ ] CSRF protection via state, replay protection via PKCE, nonce for ID token
  - [ ] TLS requirements, CSP updates, cookie downgrade prevention
  - [ ] JWKS retrieval and caching, key rotation handling, leeway for exp/nbf
- [ ] API Gateway protection
  - [ ] Cognito Authorizer vs Lambda/JWT Authorizer (per group or single authorizer + policy mapping)
  - [ ] Validating token from Authorization header or from cookie (normalizing at edge)
  - [ ] Authorizer result caching (TTL), policy granularity, mapping to routes
- [ ] Error handling and edge cases
  - [ ] Unknown user (redirect to /errors/user-must-exists)
  - [ ] Group changes mid-session, disabled users, token expiry mid-request
  - [ ] Refresh failures and invalid/tempered cookies; 401 vs redirect behavior
- [ ] Logout and session termination
  - [ ] App logout endpoint, Cognito logout endpoint, cookie clearing across subdomains
  - [ ] Global sign-out and refresh token revocation
- [ ] Infrastructure as Code and environments
  - [ ] Tooling choice (Terraform/CDK/CloudFormation), per-environment config
  - [ ] Secrets management (AWS Secrets Manager/SSM), CI/CD integration
- [ ] Observability
  - [ ] Structured logs for auth flows, metrics (login success/fail), CloudWatch dashboards/alarms
  - [ ] Tracing boundaries around SSR/auth callbacks/authorizers
- [ ] Local development and testing
  - [ ] Using Cognito in dev vs local OIDC mock; seeded test users/groups
  - [ ] End-to-end tests for SSR redirects, cookie handling, authorizers
- [ ] Migration and rollout plan
  - [ ] Dual-stack period with feature flag, cutover plan, rollback strategy
  - [ ] User/group backfill/import, communication to stakeholders
- [ ] Performance and cost
  - [ ] JWKS cache, authorizer latency and cache TTL, API Gateway cost impacts
- [ ] Frontend UX implications
  - [ ] Login/logout links, error pages, messages; preserve deep links and query params
- [ ] Access control model
  - [ ] Exact permissions per group (admins, owners, visitors), route/API mapping
- [ ] Out of scope
  - [ ] Anything explicitly not addressed (to be finalized)
