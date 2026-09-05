# ACL

Controls authentication and access. Determines who a user is (via an external identity provider), which owner they are associated with, and what operations they are permitted to perform on catalog resources.

## Language

**User**:
A person who can authenticate on the web UI. Identified by a `UserId` (their email from the identity provider). There is no concept of a user in the CLI. _Avoid_: account, member, principal

**Owner**:
A role and the root resource to which all albums and medias are attached. A user becomes an owner by holding a `MainOwnerScope`. One user maps to at most one owner. _Avoid_: admin, tenant, account

**Scope**:
The fundamental permission record: a `ScopeType` granted to a `UserId` on a specific resource (identified by `ResourceOwner` + `ResourceId`). _Avoid_: permission, grant, right, role

**ScopeType**:
The kind of access a scope confers: `api` (admin endpoints), `owner:main` (full ownership), `album:visitor` (read-only album access), `album:contributor` (read + upload), `media:visitor` (read a specific media). _Avoid_: permission type, access level, role name

**Identity**:
The user's profile as asserted by an external identity provider (Google): email, display name, and avatar URL. Obtained during login; not stored as the authoritative record. _Avoid_: profile, user info, Google account

**Authentication**:
The result of a successful login: a DPhoto-issued access token, an optional refresh token, and expiry metadata. _Avoid_: session, login result, token response

**Claims**:
The decoded content of a DPhoto access token: the subject (`UserId`), the granted scopes map, and the derived `Owner`. _Avoid_: JWT payload, token claims, decoded token

**RefreshToken**:
A long-lived credential that allows a user to obtain a new access token without re-authenticating. Scoped to a purpose (e.g. `web`) and carries an absolute expiry. _Avoid_: session token, long-lived token
