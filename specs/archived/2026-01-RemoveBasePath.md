You must remove the basepath from the NextJS application.

You'll need to modify the following places. This is not exhaustive, search if anything has been missed.

1. `deployments/cdk`: in CDK, use the `$default` route for NextJS URL redirection (instead of proxy), and remove any reference to Waku.
2. `deployments/cdk`: in SST, remove the base path configuration
3. `web-nextjs`:
    1. remove the basepath from the nextjs config
    2. in `web-nextjs/libs/requests/`, remove the basepath value
    3. in `web-nextjs/proxy.ts`, use the recommended approach to filter out the request for which the proxy apply instead of the current approach.

You need to wipe clean the concept of basepath from anywhere except `web-nextjs/libs/requests/` which can keep an empty string value. Anywhere else, CDK,
SST, ... there should be no trace that a basepath have ever been used, no comment, no function, no dead code, ...