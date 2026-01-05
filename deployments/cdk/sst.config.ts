// eslint-disable-next-line @typescript-eslint/triple-slash-reference
/// <reference path="./.sst/platform/config.d.ts" />
export default $config({
    app(input) {
        return {
            name: "web-nextjs",
            removal: "remove",
            protect: false,
            home: "aws",
        };
    },
    async run() {
        const domainName = process.env.SST_CLOUD_FRONT_DOMAIN;
        const env = {
            OAUTH_ISSUER_URL: process.env.OAUTH_ISSUER_URL,
            OAUTH_CLIENT_ID: process.env.OAUTH_CLIENT_ID,
            OAUTH_CLIENT_SECRET: process.env.OAUTH_CLIENT_SECRET,
            DPHOTO_DOMAIN_NAME: process.env.DPHOTO_DOMAIN_NAME,
        }

        console.log(`SST_CLOUD_FRONT_DOMAIN=${domainName}`);
        console.log(`Environment variables:`, JSON.stringify({...env, OAUTH_CLIENT_SECRET: env.OAUTH_CLIENT_SECRET ? "****" : "undefined"}, null, 2));

        if (!domainName) {
            throw new Error("SST_CLOUD_FRONT_DOMAIN is not defined");
        }

        const router = new sst.aws.Router("Router", {
            domain: {
                name: domainName,
            }
        });

        new sst.aws.Nextjs("DPhotoWEB", {
            path: "../../web-nextjs",
            buildCommand: "exit 0", // skip the build (done is previous workflow step)
            router: {
                instance: router,
                domain: domainName,
                path: "/nextjs",
            },
            server: {
                memory: "512 MB",
            },
            environment: {
                ...env,
            }
        });
    },
});
