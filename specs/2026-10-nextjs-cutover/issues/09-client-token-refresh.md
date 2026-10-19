# Client-side access-token refresh for long-lived SPA tabs

**Status**: done


## Description

This implements the ADR which is summarised as:

> we want the images 'scr' to be stable so the browser can cache it, and the authoriaation (AWS Lambda) to be cached. We will use exclusively a COOKIE for the get media endpoint, and EXCLUSIVELY the Authorization header for all the other call. It means that the access token must be accessible from the client (browser) as SPA style, and the session must be refreshable from the browser (access token in the cookies). So we decided to use a COOKIE readable from both server and client side to store the cookie, and an endpoint served by NextJS to refresh the access token cookie if it expired. The `ensureFreshSession` is designed to be called many times: a backend request will only occur when the token has expired.

Two new functions must be available on the client:

```typescript
/**
 * getAccessToken will return a valid access token.
 * @throws SessionExpiredException if no access token is accessible in the cookie, and it couldn't be refreshed.
 */
export async function getAccessToken(): string {
    // ...
}

/**
 * ensureFreshSession will refresh the session if the access token in the cookie has expired (or has been deleted).
 * The calls to this functions are debounced and will only make a backend requests when required: the function is safe to call on UI events like scrooling.
 * 
 * @throws SessionExpiredException if no access token is accessible in the cookie, and it couldn't be refreshed.
 */
export async function ensureFreshSession(): void {
    // ...
}
```

Each of these two functions must be tested on its own right, from the call to the API. Use API fake. Do not assert the interactions, only asserting the observable behaviour (returns, thrown exceptions, ...).

## Acceptance Criteria

The application is ready to support write operations (like create album) directly from the client as SPA (without full page refresh).

If the session has expired, the user is redirected to an expired session page with a button to reauthenticate (redirect to `/login` I would think, isn't it?).

## What does that mean?

In other term:

* `getAccessToken` and `ensureFreshSession` are implemented and tested
* they are integrated in @web-nextjs/domains/catalog/adapters/client-adapter-factory.ts and @web-nextjs/domains/catalog/adapters/api/FetchCatalogAdapter.ts
* a backend `POST /auth/refresh` is refreshing the session
* a new page for expired session, and a way to catch the exception and use it to redirect the user.

