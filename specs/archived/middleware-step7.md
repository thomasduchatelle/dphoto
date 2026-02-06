Implement caching for the OIDC configuration to avoid fetching it on every request. The configuration should be fetched once and cached in memory.

The files using it includes, and are not limited to:

* `web-nextjs/app/auth/callback/route.ts`
* `web-nextjs/app/auth/login/route.ts`
* `web-nextjs/proxy.ts`

The configuration must be cached without limit of time. Place the utilities in `web-nextjs/app/libs/security`.
