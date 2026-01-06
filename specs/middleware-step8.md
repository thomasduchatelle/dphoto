Objective: implements session refresh and session expiry.

Acceptance criteria:

1. When a user load a page and its access token has expired, then a new access token is generated with the refresh token, and the requested page is loaded.
2. When a user load a page and its access token has expired, and the refresh token couldn't be used, then the user is redirected to the /auth/login page.

Implementation details:

* use the `web-nextjs/proxy.ts` to check the access token, and try to refresh it.