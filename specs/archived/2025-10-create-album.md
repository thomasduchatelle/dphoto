The button "create album" from `AlbumListActions` needs to be disabled if the user is a simple visitor, not a owner.

You can find the user role from its JWT that can be decoded on the frontend side. I suggest you add a new property `isOwner: boolean` on the
`web/src/core/security/security-state.ts > AuthenticatedUser`. It is populated by an action dispatched from `web/src/core/security/AuthenticateCase.ts`.

You should find what to look for on the JWT from the package `pkg/acl/aclcore`.