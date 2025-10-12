# Authentication with AWS Cognito

1. Feature Summary
Implement authentication for WAKU (DPhoto) using AWS Cognito Hosted UI with Google SSO and Authorization Code + PKCE. The UI stores Cognito access and refresh tokens in secure httpOnly cookies and performs SSR-driven redirects/refresh. API Gateway protects /api/** with Lambda authorizers that require specific Cognito group membership; fine-grained permissions remain enforced in the application using the existing DynamoDB ACL.

2. Ubiquity Language
- Cognito User Pool: AWS-managed identity store where users are pre-created and federated to Google.
- Google IdP: External identity provider configured in Cognito for SSO.
- Hosted UI: Cognito’s OAuth/OIDC UI used for sign-in.
- PKCE: Code challenge/verification used with the Authorization Code flow.
- dphoto-access-token: httpOnly cookie that stores the Cognito access token (JWT).
- dphoto-refresh-token: httpOnly cookie that stores the Cognito refresh token.
- dphoto-pkce-verifier: short-lived httpOnly cookie carrying the PKCE code_verifier during the sign-in redirect flow.
- dphoto-auth-state: short-lived httpOnly cookie carrying the OAuth state and original URL to return to.
- AdminsAuthorizer / OwnersAuthorizer / VisitorsAuthorizer: API Gateway Lambda authorizers that require the caller to be in the respective Cognito group; strict membership (no implicit admin bypass).
- Groups: Cognito user-pool groups named admins, owners, visitors. Users may belong to multiple groups; users with no group are denied.
- ACL DB: Existing DynamoDB-based authorization model holding fine-grained resource grants (owner, album, media) that the app enforces server-side.
- Pre-Signup Lambda: Cognito trigger that blocks sign-ups unless the user is pre-created.
- Token Configuration: Cognito client configuration adding email to the access token; built-in cognito:groups is present in access tokens.

3. Scenarios

3.1 First-time sign-in (pre-provisioned user) via SSR
- Before: User navigates to a protected SSR page without cookies or with invalid/expired tokens.
- Steps:
  1) SSR detects missing/invalid dphoto-access-token or dphoto-refresh-token; generates PKCE code_verifier/challenge and a random state containing the original URL.
  2) SSR sets dphoto-pkce-verifier and dphoto-auth-state cookies (httpOnly, secure, SameSite=Lax, TTL ~10 min).
  3) SSR redirects the browser to Cognito Hosted UI (authorization endpoint) with client_id, redirect_uri, scope=openid email, response_type=code, code_challenge(+method=S256), and state.
  4) User selects Google, completes SSO; Cognito redirects back to /auth/callback with code and state.
  5) Backend validates state against dphoto-auth-state and reconstructs the original URL; validates PKCE using dphoto-pkce-verifier; exchanges code for tokens server-side.
  6) Backend sets dphoto-access-token and dphoto-refresh-token (httpOnly, secure, SameSite=Lax, Path=/), clears PKCE/state cookies, and 302 redirects to the original URL.
  7) SSR re-renders the page, now authenticated.
- Validation and errors:
  - Missing/mismatched state or PKCE → reject, clear cookies, redirect to /auth/access-denied.
  - OAuth error from Cognito → redirect to /auth/access-denied.
- After: User is authenticated; subsequent SSR pages skip re-auth until refresh or expiry is needed.

3.2 Unapproved user attempts SSO
- Before: User not pre-created in Cognito attempts sign-in via Google SSO.
- Steps:
  1) Cognito Pre-Signup trigger denies the sign-up attempt (no matching user in pool).
  2) Hosted UI returns an error to the callback.
  3) /auth/callback detects the error, clears PKCE/state cookies, does not set token cookies, and redirects to /auth/access-denied.
- After: User sees an access denied page with guidance to contact an administrator.

3.3 SSR token refresh when access token is near expiry
- Before: User is authenticated, access token will expire in < 5 minutes.
- Steps:
  1) On SSR, the server parses dphoto-access-token claims to read exp.
  2) If exp - now < 5 minutes, SSR calls POST /auth/refresh (server-to-server call).
  3) Backend uses dphoto-refresh-token to obtain a new access token (and refresh token if provided), updates cookies, returns 204.
  4) SSR proceeds to render the page.
- Errors:
  - Refresh fails (expired/revoked) → backend clears token cookies and returns 401; SSR responds by redirecting the user to Cognito (Scenario 3.1).
- After: User continues seamlessly if refresh succeeds; otherwise re-authenticates.

3.4 Admin-only API request
- Before: User has admins group and valid access token; wants to call an admin-protected endpoint (e.g., POST /api/users).
- Steps:
  1) API Gateway attaches AdminsAuthorizer to the route.
  2) Authorizer extracts dphoto-access-token from Cookie header, verifies JWT (issuer/audience/signature/exp) using Cognito JWKS.
  3) Authorizer reads cognito:groups from the access token and checks membership in admins.
  4) If allowed, request forwards to the backend; backend applies domain logic (including ACL as needed).
- Errors:
  - Not in admins or missing groups → 403 with JSON error via Gateway Responses.
  - Invalid/expired token → 401/403 as appropriate.
- After: Admin operations succeed only for admins.

3.5 Owners-authorized API request (strict group) with app-level ACL
- Before: User belongs to owners group; calling an owners-protected route (e.g., share/unshare album).
- Steps:
  1) API Gateway uses OwnersAuthorizer; validates token and requires owners group.
  2) Backend, after authorizer passes, enforces fine-grained ACL via DynamoDB (e.g., owner:main or album-level grants).
- Errors:
  - User in admins but not in owners → 403 (strict authorizer).
  - Fails ACL check in backend → 403 with domain-specific error.
- After: Only owners (that also pass domain ACL) can perform the action.

3.6 Logout (local-only)
- Before: User is authenticated.
- Steps:
  1) User triggers /auth/logout.
  2) Backend clears dphoto-access-token, dphoto-refresh-token, PKCE/state cookies (if present) and redirects to /.
  3) Cognito session remains active; next SSR on a protected page restarts Scenario 3.1 and may result in an immediate SSO if the Hosted UI session is still valid.
- Errors: None specific; clearing non-existent cookies is a no-op.
- After: Local session ends; user may appear re-signed if Cognito Hosted UI still holds a session.

4. Technical Context

Provided by external services or other domains:
- AWS Cognito
  - User Pool with Google as an external IdP; Hosted UI domain configured.
  - Authorization Code + PKCE flow; client app has redirect URL /auth/callback and allowed logout URL /.
  - Token lifetimes: default access token (~60 min) and refresh token (e.g., 30 days).
  - Token Configuration: include email in access token; built-in cognito:groups included.
  - Pre-Signup Lambda: denies sign-ups unless user pre-exists; auto-confirms/linking as needed.
  - Users are pre-created by the CLI and assigned to groups admins/owners/visitors.

- AWS CDK
  - Provisions Cognito User Pool, App Client, Hosted UI domain, Google IdP configuration.
  - Injects Google client ID/secret from AWS KMS-encrypted parameters.
  - Outputs necessary env/config (pool id, client id, issuer URL, JWKS URL, callback URL).

- API Gateway
  - All /api/** routes are protected; requests must include a valid dphoto-access-token cookie.
  - Three Lambda authorizers (AdminsAuthorizer, OwnersAuthorizer, VisitorsAuthorizer) each enforce strict membership in their respective group.
  - Gateway Responses configured to return JSON error bodies for 401/403 to align with API error format.

- Lambda Authorizers (custom)
  - Extract dphoto-access-token from Cookie, validate JWT (iss, aud, kid, exp) using Cognito JWKS.
  - Decode claims and require presence of email and cognito:groups; enforce exact group membership.
  - Return allow/deny policy and, on allow, inject authorizer context (e.g., sub, email, groups) to downstream integration.

- Backend application (SSR server)
  - On SSR, parses token claims to decide redirect or refresh; not responsible for enforcing groups (that’s at Gateway).
  - Endpoints:
    - /auth/callback: validate state/PKCE, exchange code for tokens, set cookies, redirect back.
    - /auth/refresh: use refresh token to mint new access token when exp < 5 min; returns 204; on failure clears cookies and 401.
    - /auth/logout: clear cookies and redirect to /.
    - /auth/access-denied: simple page for unauthorized users.
  - Cookie attributes:
    - dphoto-access-token, dphoto-refresh-token: httpOnly, secure, SameSite=Lax, Path=/, Domain=app host, reasonable Max-Age.
    - dphoto-pkce-verifier, dphoto-auth-state: httpOnly, secure, SameSite=Lax, Path=/, TTL ~10 minutes.

- DynamoDB ACL (existing domain)
  - Source of truth for fine-grained permissions (owner/main, album visitors/contributors, media visitors).
  - Continues to be enforced by backend services; not duplicated in Cognito.

- CLI (existing, Go)
  - Updated to manage Cognito users and group memberships (admins, owners, visitors) and to continue managing the ACL DB entries.
  - Self sign-up disabled; all users are provisioned by CLI.

Out of scope
- Mapping specific /api routes to each authorizer (owned by the integration agent).
- Replacing the existing ACL model with Cognito constructs (groups are coarse only).
- Centralized authorization using AVP or ABAC policy engines.
- Cross-subdomain cookies, mobile apps, and non-browser clients.
- Global logout from Cognito Hosted UI; we only implement local logout.
- Admin UI for user/group management (handled via CLI).
- Handling identity changes (email changes) and historical ACL migrations.

5. Explorations

Open design questions impacting behavior:
- API Gateway flavor and responses:
  - Will we use REST API or HTTP API? This affects authorizer interfaces and the ability to customize 401/403 JSON responses (Gateway Responses availability).
  - Confirm desired JSON shape for 401/403 produced by Gateway; if required, define Gateway Responses accordingly.

- Cookie domain and lifetime details:
  - Confirm exact domain to scope cookies (e.g., app.example.com). Any need to share across subdomains in the future?
  - Desired Max-Age for access/refresh cookies (match token TTLs or shorter?), and clock skew tolerance policy.

- CSRF posture for /auth/refresh:
  - Should /auth/refresh require POST and include an anti-CSRF token (double-submit or same-origin header check), or is SameSite=Lax sufficient given SSR-only invocation?

- Group changes and session consistency:
  - If group membership changes during a session, authorizer decisions are based on the access token’s cognito:groups at issuance time. Define expectations for when changes take effect (e.g., after next refresh or forced re-login).

- Email as identifier:
  - Email is added to the access token via Token Configuration. Must email be verified? Define behavior if the IdP returns an unverified email or if email changes in IdP.

- Pre-Signup behavior:
  - Exact rejection messaging and localization for /auth/access-denied. Should we expose a support contact or request ID?

- Rate limits and protection:
  - Define limits for /auth/refresh and /auth/callback to mitigate abuse. Consider WAF rules on OAuth endpoints.

- Observability:
  - Required logs/metrics/traces for the authorizers and /auth endpoints (e.g., group membership decisions, refresh failures) and PII redaction policy.

- Future admin route gating:
  - Any endpoints requiring composite roles (e.g., both admin and owner) or tenant-aware admin? If so, specify policy to avoid drift.
