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
        const cognitoIssuer = process.env.SST_COGNITO_ISSUER;
        const cognitoClientId = process.env.SST_COGNITO_CLIENT_ID;
        const cognitoClientSecret = process.env.SST_COGNITO_CLIENT_SECRET;

        console.log(`SST_CLOUD_FRONT_DOMAIN=${domainName}`);
        console.log(`SST_COGNITO_ISSUER=${cognitoIssuer}`);
        console.log(`SST_COGNITO_CLIENT_ID=${cognitoClientId}`);
        console.log(
            `SST_COGNITO_CLIENT_SECRET=${cognitoClientSecret ? "****" : "undefined"}`,
        );

        if (!domainName) {
            throw new Error("SST_CLOUD_FRONT_DOMAIN is not defined");
        }

        const router = new sst.aws.Router("Router", {
            domain: {
                name: domainName,
            }
        });

        new sst.aws.Nextjs("DPhotoWEB", {
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
                NEXT_PUBLIC_COGNITO_ISSUER: cognitoIssuer!,
                NEXT_PUBLIC_COGNITO_CLIENT_ID: cognitoClientId!,
                COGNITO_CLIENT_SECRET: cognitoClientSecret!,
            }
        });
    },
});
